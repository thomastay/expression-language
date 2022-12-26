// Package compile implements semantic analysis (aka Parsing of ints, type checking, constant folding, etc)
// as well as compiling the instructions down into bytecode
package compiler

import (
	"fmt"
	"strconv"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/thomastay/expression_language/pkg/instructions"
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
			default:
				panic("Not implemented")
			}
			c.Bytecode = append(c.Bytecode, Bytecode{
				Inst: instructions.OpConst,
				Val:  val,
			})
		case *parser.EBinOp:
			// parse left
			compileRec(node.Left)
			compileRec(node.Right)
			// parse right
			var inst instructions.Instruction
			switch node.Op.Value {
			case "+":
				inst = instructions.OpAdd
			case "-":
				inst = instructions.OpMinus
			case "*":
				inst = instructions.OpMul
			case "/":
				inst = instructions.OpDiv
			default:
				panic("Not implemented")
			}
			c.Bytecode = append(c.Bytecode, Bytecode{
				Inst: inst,
			})
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
				Inst: instructions.OpBrIf,
				// patch the val later on
			})
			firstJumpIdx := len(c.Bytecode) - 1

			compileRec(node.Second)
			c.Bytecode = append(c.Bytecode, Bytecode{
				Inst: instructions.OpBr,
			})
			secondJumpIdx := len(c.Bytecode) - 1
			// Patch the first jump
			c.Bytecode[firstJumpIdx].Val = int64(len(c.Bytecode))

			compileRec(node.First)
			// Patch the second jump
			c.Bytecode[secondJumpIdx].Val = int64(len(c.Bytecode))
		default:
			panic("Not impl")
		}
	}
	compileRec(expr)
	return c
}

// represents the
type Compilation struct {
	Bytecode []Bytecode
	Errors   []CompileError
}

type Bytecode struct {
	Inst instructions.Instruction
	Val  int64 // TODO
}

func (b Bytecode) String() string {
	switch b.Inst {
	// No value
	case instructions.OpAdd:
		fallthrough
	case instructions.OpMul:
		fallthrough
	case instructions.OpDiv:
		fallthrough
	case instructions.OpMinus:
		return fmt.Sprintf("%s", b.Inst)
	default:
		return fmt.Sprintf("%s %d", b.Inst, b.Val)
	}
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
