// Package compile implements semantic analysis (aka Parsing of ints, type checking, constant folding, etc)
// as well as compiling the instructions down into bytecode
package compiler

import (
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
	// We are going to build this up slowly.
	// Let's start by being able to interpret Add operations
	c := Compilation{}
	var compileRec func(parser.Expr)
	compileRec = func(expr parser.Expr) {
		if expr == nil {
			return
		}
		switch node := expr.(type) {
		case *parser.EValue:
			switch node.Val.Type {
			case parser.TokInt:
				tok := node.Val
				val, err := strconv.ParseInt(tok.Value, 10, 63)
				if err != nil {
					c.Errors = append(c.Errors, CompileError{
						Err:   err,
						Start: tok,
						End:   tok,
					})
					return
				}
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: instructions.OpConstant,
					Val:  val,
				})
			default:
				panic("Not implemented")
			}
		case *parser.EBinOp:
			// parse left
			compileRec(node.Left)
			compileRec(node.Right)
			// parse right
			switch node.Op.Value {
			case "+":
				c.Bytecode = append(c.Bytecode, Bytecode{
					Inst: instructions.OpAdd,
				})
			default:
				panic("Not implemented")
			}
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

type CompileError struct {
	Err   error
	Start *lexer.Token
	End   *lexer.Token
}