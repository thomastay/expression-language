// Package compile implements semantic analysis (aka Parsing of ints, type checking, constant folding, etc)
// as well as compiling the instructions down into bytecode
package compiler

import (
	"strconv"

	"github.com/alecthomas/participle/v2/lexer"
	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/parser"
)

// Compiles the parse tree down to IR, represented by the Compilation object
// Note that a compilation may have errors. Users are expected to check the Compilation.Errors
// object to report them. Note that a Compilation may have errors but still be ok to interpret
// because this package will attempt a best effort compile
func Compile(expr parser.Expr) Compilation {
	c := Compilation{}
	var compileRec func(parser.Expr)
	compileRec = func(expr parser.Expr) {
		if expr == nil {
			return
		}
		switch node := expr.(type) {
		case *parser.EValue:
			var val int64
			var err error
			switch node.Val.Type {
			case parser.TokInt:
				tok := node.Val
				val, err = strconv.ParseInt(tok.Value, 10, 64)
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
			default:
				panic("Not implemented")
			}
		case *parser.EBinOp:
			if isSimpleBinOp(node.Op.Value) {
				compileRec(node.Left)
				compileRec(node.Right)
				var inst Instruction
				switch node.Op.Value {
				case "+":
					inst = OpAdd
				case "-":
					inst = OpMinus
				case "*":
					inst = OpMul
				case "/":
					inst = OpDiv
				default:
					panic("Not a simple binary op!")
				}
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
					panic("Not implemented")
				}
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
		default:
			panic("Not impl")
		}
	}
	compileRec(expr)
	return c
}

func isSimpleBinOp(op string) bool {
	switch op {
	case "+":
		return true
	case "-":
		return true
	case "*":
		return true
	case "/":
		return true
	}
	return false
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
