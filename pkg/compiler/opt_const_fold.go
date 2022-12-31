package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Implements Visitor type
// Constant folding pass does two things:
//  1. If node is a binary op with two constants, it replaces them
//  1. If node is binary op with one constant, it rotates the constant to the LHS
//     ^^ this trick is to enable constant pushdown, which will be implemented next
//     ^^ Note: this is the opposite order from usual assembly ADDI, where constant is on the right.
//     This is because we have to do less tree rotations this way. We can always swap it back before
//     emitting bytecode, but i don't think we need to do that.
func ConstFold(expr Expr) walkError {
	var errs []CompileError
	switch node := expr.(type) {
	case *EValue:
		return nil
	case *EUnOp:
		switch node.Val.(type) {
		case *EValue:

		}

	// case *EBinOp:
	// case *EFieldAccess:
	// case *EIdxAccess:
	// case *ECond:
	// case *ECall:
	// case *EArray:
	default:
		panic("AST type is not impl")
	}
	return errs
}
