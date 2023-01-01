// Package compile implements semantic analysis (aka Parsing of ints, type checking, constant folding, etc)
// as well as compiling the instructions down into bytecode
package compiler

import (
	"fmt"
	"log"

	"github.com/alecthomas/participle/v2/lexer"
	. "github.com/thomastay/expression_language/pkg/ast"
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
	return Compile(expr, Params{})
}

// Compiles the parse tree down to bytecode, represented by the Compilation object
// Note that a compilation may have errors. Users are expected to check the Compilation.Errors
// object to report them. Note that a Compilation may have errors but still be ok to interpret
// because this package will attempt a best effort compile
func Compile(expr Expr, params Params) Compilation {
	// This function merely orchestrates the steps, then compileToBytecode does the actual compiling
	// Stage 1: Parsing of values and reporting overflow errors or simple errors
	c := Compilation{}
	c.Errors = walk(&expr, ParseValue)
	if len(c.Errors) > 0 {
		return c
	}

	// Stage 2: Optimization
	errs := walk(&expr, ConstFold)
	c.Errors = append(c.Errors, []CompileError(errs)...)
	if len(c.Errors) > 0 {
		return c
	}
	if params.Debug {
		fmt.Println("Expression after const fold:", expr.String())
	}
	errs = walk(&expr, ConstPushDown)
	c.Errors = append(c.Errors, []CompileError(errs)...)
	if len(c.Errors) > 0 {
		return c
	}
	if params.Debug {
		fmt.Println("Expression after const Pushdown:", expr.String())
	}

	c.compileToBytecode(expr)
	return c
}

// This function does the actual compiling to bytecode
func (c *Compilation) compileToBytecode(expr Expr) {
	seen := newSeenConstants()
	// ?: why doesn't this defer work?
	// defer func() { c.Constants = seen.constants }()

	var compileRec func(Expr)
	compileRec = func(expr Expr) {
		if expr == nil {
			panic("No nil expressions!")
		}
		switch node := expr.(type) {
		case *EValue:
			panic("No more EValues at this point")
		case *EInt:
			val := int64(*node)
			pos := seen.AddInt(val)
			c.Bytecode.Push(Bytecode{
				Inst: OpConst,
				Val:  pos,
			})
		case *EFloat:
			val := float64(*node)
			pos := seen.AddFloat(val)
			c.Bytecode.Push(Bytecode{
				Inst: OpConst,
				Val:  pos,
			})
		case *EStr:
			val := string(*node)
			pos := seen.AddStr(val)
			c.Bytecode.Push(Bytecode{
				Inst: OpConst,
				Val:  pos,
			})
		case *EIdent:
			val := string(*node)
			pos := seen.AddStr(val)
			c.Bytecode.Push(Bytecode{
				Inst: OpLoad,
				Val:  pos,
			})
		case *EBool:
			val := bool(*node)
			pos := falseConstPos
			if val {
				pos = trueConstPos
			}
			c.Bytecode.Push(Bytecode{
				Inst: OpConst,
				Val:  pos,
			})
		case *EBinOp:
			if inst, ok := simpleBinaryOps[node.Op.Value]; ok {
				compileRec(node.Left)
				compileRec(node.Right)
				c.Bytecode.Push(Bytecode{Inst: inst})
			} else {
				switch node.Op.Value {
				case "and":
					// Bytecode:
					// | 0        First expr
					// | 1   |    BR_IF_FALSE_OR_POP
					// | 2   |    Second Expr
					// | 3   ---> ...
					compileRec(node.Left)
					c.Bytecode.Push(Bytecode{
						Inst: OpBrIfFalseOrPop,
						// patch the val later on
					})
					jumpIdx := c.Bytecode.Len() - 1
					compileRec(node.Right)
					c.Bytecode.IntData[jumpIdx] = c.Bytecode.Len()
				case "or":
					// Bytecode:
					// | 0        First expr
					// | 1   |    BR_IF_OR_POP
					// | 2   |    Second Expr
					// | 3   ---> ...
					compileRec(node.Left)
					c.Bytecode.Push(Bytecode{
						Inst: OpBrIfOrPop,
						// patch the val later on
					})
					jumpIdx := c.Bytecode.Len() - 1
					compileRec(node.Right)
					c.Bytecode.IntData[jumpIdx] = c.Bytecode.Len()
				default:
					log.Panicf("Not implemented %v", node)
				}
			}
		case *EUnOp:
			if inst, ok := unaryOps[node.Op.Value]; ok {
				compileRec(node.Val)
				c.Bytecode.Push(Bytecode{Inst: inst})
			} else {
				log.Panicf("Not implemented %v", node)
			}
		case *ECond:
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
			c.Bytecode.Push(Bytecode{
				Inst: OpBrIf,
				// patch the val later on
			})
			firstJumpIdx := c.Bytecode.Len() - 1

			compileRec(node.Second)
			c.Bytecode.Push(Bytecode{
				Inst: OpBr,
			})
			secondJumpIdx := c.Bytecode.Len() - 1
			// Patch the first jump
			c.Bytecode.IntData[firstJumpIdx] = c.Bytecode.Len()

			compileRec(node.First)
			// Patch the second jump
			c.Bytecode.IntData[secondJumpIdx] = c.Bytecode.Len()
		case *ECall:
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
			val := node.Method.Value
			if node.Base != nil {
				compileRec(node.Base)
				pos := seen.AddStr(val)
				c.Bytecode.Push(Bytecode{
					Inst: OpConst,
					Val:  pos,
				})
				c.Bytecode.Push(Bytecode{
					Inst: OpLoadAttr,
				})
			} else {
				pos := seen.AddStr(val)
				c.Bytecode.Push(Bytecode{
					Inst: OpLoad,
					Val:  pos,
				},
				)
			}
			c.Bytecode.Push(
				Bytecode{
					Inst: OpCall,
					Val:  numParams,
				},
			)
		case *EFieldAccess:
			// Load Base, then Field
			compileRec(node.Base)
			val := node.Field.Value
			pos := seen.AddStr(val)
			c.Bytecode.Push(
				Bytecode{
					Inst: OpConst,
					Val:  pos,
				},
			)
			c.Bytecode.Push(
				Bytecode{Inst: OpLoadAttr},
			)
		// ------------------- Arrays ------------------------------
		case *EArray:
			// For simplicity, we're just going to dump them all on the stack in reverse order and evaluate them for now
			// Future optimization may make it such that integers are not dumped on stack
			n := len(*node)
			for i := n - 1; i >= 0; i-- {
				expr := (*node)[i]
				compileRec(expr)
			}
			c.Bytecode.Push(
				Bytecode{
					Inst: OpNewArray,
					Val:  n,
				},
			)
		case *EIdxAccess:
			compileRec(node.Base)
			compileRec(node.Index)
			c.Bytecode.Push(Bytecode{Inst: OpLoadSubscript})
		default:
			log.Panicf("Not implemented %v", node)
		}
	}
	compileRec(expr)
	c.Constants = seen.constants
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

// represents the result of compiling an AST
type Compilation struct {
	Bytecode  ByteCodes
	Constants []BVal
	Errors    []CompileError
}

// Helper class to aid in constructing the constant table
type seenConstants struct {
	constants   []BVal
	seenInts    map[int64]int
	seenFloats  map[float64]int
	seenStrings map[string]int
}

const (
	falseConstPos = 0
	trueConstPos  = 1
)

func newSeenConstants() seenConstants {
	seen := seenConstants{}
	seen.seenInts = make(map[int64]int)
	seen.seenFloats = make(map[float64]int)
	seen.seenStrings = make(map[string]int)

	// We initialize common constants here. true and false are always 1 and 0 respectively
	// Also initialize 0, 1, 2, 3
	seen.constants = append(seen.constants,
		BBool(false), BBool(true),
		BInt(0), BInt(1), BInt(2), BInt(3),
	)
	seen.seenInts[0] = 2 // position 2 cos 0 and 1 are bools
	seen.seenInts[1] = 3
	seen.seenInts[2] = 4
	seen.seenInts[3] = 5

	return seen
}

func (seen *seenConstants) AddInt(x int64) int {
	if pos, ok := seen.seenInts[x]; ok {
		return pos
	}
	pos := len(seen.constants)
	seen.constants = append(seen.constants, BInt(x))
	seen.seenInts[x] = pos
	return pos
}
func (seen *seenConstants) AddFloat(x float64) int {
	if pos, ok := seen.seenFloats[x]; ok {
		return pos
	}
	pos := len(seen.constants)
	seen.constants = append(seen.constants, BFloat(x))
	seen.seenFloats[x] = pos
	return pos
}
func (seen *seenConstants) AddStr(x string) int {
	if pos, ok := seen.seenStrings[x]; ok {
		return pos
	}
	pos := len(seen.constants)
	seen.constants = append(seen.constants, BStr(x))
	seen.seenStrings[x] = pos
	return pos
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

type Params struct {
	Debug bool
}
