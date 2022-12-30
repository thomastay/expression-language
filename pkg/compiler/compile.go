// Package compile implements semantic analysis (aka Parsing of ints, type checking, constant folding, etc)
// as well as compiling the instructions down into bytecode
package compiler

import (
	"log"
	"strconv"

	"github.com/alecthomas/participle/v2/lexer"
	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/parser"
)

// The main entry point for apps who want to cache the bytecode across several runs
// If you don't, then just do vm.EvalString(s)
func CompileString(s string) Compilation {
	expr, err := parser.ParseString(s)
	if err != nil {
		return Compilation{
			Errors: []CompileError{{Err: err}},
		}
	}
	return Compile(expr)
}

// Compiles the parse tree down to IR, represented by the Compilation object
// Note that a compilation may have errors. Users are expected to check the Compilation.Errors
// object to report them. Note that a Compilation may have errors but still be ok to interpret
// because this package will attempt a best effort compile
func Compile(expr parser.Expr) Compilation {
	c := Compilation{}
	var compileRec func(parser.Expr)
	compileRec = func(expr parser.Expr) {
		if expr == nil {
			panic("No nil expressions!")
			// log.Println("Expressions shouldn't be nil.")
			// return
		}
		switch node := expr.(type) {
		case *parser.EValue:
			switch node.Val.Type {
			case parser.TokInt, parser.TokHexInt, parser.TokOctInt, parser.TokBinInt:
				tok := node.Val
				val, err := strconv.ParseInt(tok.Value, 0, 64)
				if err != nil {
					c.Errors = append(c.Errors, CompileError{
						Err:   err,
						Start: tok,
						End:   tok,
					})
					return
				}
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: OpConst,
					Val:  BInt(val),
				})
			case parser.TokFloat:
				tok := node.Val
				val, err := strconv.ParseFloat(tok.Value, 64)
				if err != nil {
					c.Errors = append(c.Errors, CompileError{
						Err:   err,
						Start: tok,
						End:   tok,
					})
					return
				}
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: OpConst,
					Val:  BFloat(val),
				})
			case parser.TokSingleString:
				val := node.Val.Value
				val = val[1 : len(val)-1]
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: OpConst,
					Val:  BStr(val),
				})
			// Parse and load identifier
			case parser.TokIdent:
				val := node.Val.Value
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: OpLoad,
					Val:  BStr(val),
				})
			case parser.TokBool:
				val := node.Val.Value
				iVal := 0
				if val == "true" {
					iVal = 1
				}
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: OpConst,
					Val:  BInt(iVal),
				})
			default:
				log.Panicf("Not implemented %v", node)
			}
		case *parser.EBinOp:
			if inst, ok := simpleBinaryOps[node.Op.Value]; ok {
				compileRec(node.Left)
				compileRec(node.Right)
				c.Bytecode = append(c.Bytecode, Bytecode{Inst: inst})
			} else {
				switch node.Op.Value {
				case "and":
					// Bytecode:
					// | 0        First expr
					// | 1   |    BR_IF_FALSE_OR_POP
					// | 2   |    Second Expr
					// | 3   ---> ...
					compileRec(node.Left)
					c.Bytecode = append(c.Bytecode, Bytecode{
						Inst: OpBrIfFalseOrPop,
						// patch the val later on
					})
					jumpIdx := len(c.Bytecode) - 1
					compileRec(node.Right)
					c.Bytecode[jumpIdx].IntVal = len(c.Bytecode)
				case "or":
					// Bytecode:
					// | 0        First expr
					// | 1   |    BR_IF_OR_POP
					// | 2   |    Second Expr
					// | 3   ---> ...
					compileRec(node.Left)
					c.Bytecode = append(c.Bytecode, Bytecode{
						Inst: OpBrIfOrPop,
						// patch the val later on
					})
					jumpIdx := len(c.Bytecode) - 1
					compileRec(node.Right)
					c.Bytecode[jumpIdx].IntVal = len(c.Bytecode)
				default:
					log.Panicf("Not implemented %v", node)
				}
			}
		case *parser.EUnOp:
			if inst, ok := unaryOps[node.Op.Value]; ok {
				compileRec(node.Val)
				c.Bytecode = append(c.Bytecode, Bytecode{Inst: inst})
			} else {
				log.Panicf("Not implemented %v", node)
			}
		case *parser.ECond:
			// This places the condition val onto the stack.
			compileRec(node.Cond)
			// Next, we want to add a branch instruction to branch if true
			// Bytecode:
			// | 0       Condition Expression
			// | 1   |    BR_IF
			// | 2   |    Else clause
			// | 3   | |  BR
			// | 4   ---> Then clause
			// | 5     --> ....
			// Thus, we put the else clause first, and branch to the then clause if true
			c.Bytecode = append(c.Bytecode, Bytecode{
				Inst: OpBrIf,
				// patch the val later on
			})
			firstJumpIdx := len(c.Bytecode) - 1

			compileRec(node.Second)
			c.Bytecode = append(c.Bytecode, Bytecode{
				Inst: OpBr,
			})
			secondJumpIdx := len(c.Bytecode) - 1
			// Patch the first jump
			c.Bytecode[firstJumpIdx].IntVal = len(c.Bytecode)

			compileRec(node.First)
			// Patch the second jump
			c.Bytecode[secondJumpIdx].IntVal = len(c.Bytecode)
		case *parser.ECall:
			// Just a plain old function call
			// Bytecode:
			// | 0        Load params in reverse order onto stack
			// | 1        Load Base (if needed)
			// | 2        Load Field
			// | 3        Load Base.Field (if needed)
			// | 4       	Call (num params embedded in bytecode)
			numParams := len(node.Exprs)
			for i := numParams - 1; i >= 0; i-- {
				param := node.Exprs[i]
				compileRec(param)
			}
			if node.Base != nil {
				compileRec(node.Base)
				c.Bytecode = append(c.Bytecode,
					Bytecode{
						Inst: OpConst,
						Val:  BStr(node.Method.Value),
					},
					Bytecode{Inst: OpLoadAttr})
			} else {
				c.Bytecode = append(c.Bytecode,
					Bytecode{
						Inst: OpLoad,
						Val:  BStr(node.Method.Value),
					},
				)
			}
			c.Bytecode = append(c.Bytecode,
				Bytecode{
					Inst:   OpCall,
					IntVal: numParams,
				},
			)
		case *parser.EFieldAccess:
			// Load Base, then Field
			compileRec(node.Base)
			c.Bytecode = append(c.Bytecode,
				Bytecode{
					Inst: OpConst,
					Val:  BStr(node.Field.Value),
				},
				Bytecode{Inst: OpLoadAttr},
			)
		// ------------------- Arrays ------------------------------
		case *parser.EArray:
			// For simplicity, we're just going to dump them all on the stack in reverse order and evaluate them for now
			// Future optimization may make it such that integers are not dumped on stack
			n := len(*node)
			for i := n - 1; i >= 0; i-- {
				expr := (*node)[i]
				compileRec(expr)
			}
			c.Bytecode = append(c.Bytecode,
				Bytecode{
					Inst:   OpNewArray,
					IntVal: n,
				},
			)
		case *parser.EIdxAccess:
			compileRec(node.Base)
			compileRec(node.Index)
			c.Bytecode = append(c.Bytecode, Bytecode{Inst: OpLoadSubscript})
		default:
			log.Panicf("Not implemented %v", node)
		}
	}
	compileRec(expr)
	return c
}

var simpleBinaryOps = map[string]Instruction{
	"+":  OpAdd,
	"-":  OpMinus,
	"*":  OpMul,
	"/":  OpDiv,
	"//": OpFloorDiv,
	"%":  OpMod,
	"**": OpPow,
	"<":  OpLt,
	">":  OpGt,
	">=": OpGe,
	"<=": OpLe,
	"==": OpEq,
	"!=": OpNe,
}

var unaryOps = map[string]Instruction{
	"+":   OpUnaryPlus,
	"-":   OpUnaryMinus,
	"not": OpUnaryNot,
}

// represents the
type Compilation struct {
	Bytecode []Bytecode
	Errors   []CompileError
}

type CompileError struct {
	Err   error
	Start *lexer.Token
	End   *lexer.Token
}

func (c CompileError) Error() string {
	return c.Err.Error()
}

func (c CompileError) IsErr() bool {
	return c.Err != nil
}
