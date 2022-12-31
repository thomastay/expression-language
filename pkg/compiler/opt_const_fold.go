package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

// Constant folding pass does two things:
//  1. If node is a binary op with two constants, it replaces them
//  1. If node is binary op with one EValue, it rotates the constant to the LHS
//     ^^ this trick is to enable constant pushdown, which will be implemented next
//     ^^ Note: this is the opposite order from usual assembly ADDI, where constant is on the right.
//     This is because we have to do less tree rotations this way. We can always swap it back before
//     emitting bytecode, but i don't think we need to do that.
func ConstFold(ptrToExpr *Expr) walkError {
	var errs []CompileError
	switch node := (*ptrToExpr).(type) {
	case *EValue:
		panic("No more EValues after parsing")
	case *EInt:
	case *EFloat:
	case *EStr:
	case *EIdent:
	case *EBool:
		// do nothing
	case *EUnOp:
		if !isPossiblyConst(node.Val) {
			break
		}
		switch node.Op.Value {
		case "+":
			switch inner := node.Val.(type) {
			case *EInt:
				*ptrToExpr = node.Val
			case *EFloat:
				*ptrToExpr = node.Val
			case *EBool:
				// create an int node
				boolVal := int64(0)
				if *inner {
					boolVal = 1
				}
				*ptrToExpr = (*EInt)(&boolVal)
			case *EStr:
				errs = append(errs, CompileError{
					Err:   errUnaryType("+", "string"),
					Start: node.Op,
					End:   node.Op,
				})
			case *EArray:
				errs = append(errs, CompileError{
					Err:   errUnaryType("+", "array"),
					Start: node.Op,
					End:   node.Op,
				})
			default:
				// do nothing, fallthrough
			}
		case "-":
			switch inner := node.Val.(type) {
			case *EInt:
				val := -*inner
				*ptrToExpr = (*EInt)(&val)
			case *EFloat:
				val := -*inner
				*ptrToExpr = (*EFloat)(&val)
			case *EBool:
				// create an int node
				innerIntValNegated := int64(0)
				if *inner {
					innerIntValNegated = -1
				}
				*ptrToExpr = (*EInt)(&innerIntValNegated)
			case *EStr:
				errs = append(errs, CompileError{
					Err:   errUnaryType("-", "string"),
					Start: node.Op,
					End:   node.Op,
				})
			case *EArray:
				errs = append(errs, CompileError{
					Err:   errUnaryType("-", "array"),
					Start: node.Op,
					End:   node.Op,
				})
			default:
				// do nothing, fallthrough
			}
		case "not":
			switch inner := node.Val.(type) {
			case *EInt:
				flip := *inner == 0
				*ptrToExpr = (*EBool)(&flip)
			case *EFloat:
				flip := *inner == 0
				*ptrToExpr = (*EBool)(&flip)
			case *EBool:
				// create an int node
				flip := *inner != true
				*ptrToExpr = (*EBool)(&flip)
			case *EStr:
				flip := len(*inner) == 0
				*ptrToExpr = (*EBool)(&flip)
			case *EArray:
				flip := len(*inner) == 0
				*ptrToExpr = (*EBool)(&flip)
			default:
				// do nothing, fallthrough
			}

		}

	// Do nothing for now (not impl)
	case *EBinOp:
	case *EFieldAccess:
	case *EIdxAccess:
	case *ECond:
	case *ECall:
	case *EArray:
	default:
		panic("AST type is not impl")
	}
	return errs
}

func isPossiblyConst(expr Expr) bool {
	// TODO this function needs rework once we implement const array parsing
	switch expr.(type) {
	// arrays not here but MAY NOT BE CONST for now until we implement constant array parsing
	case *EInt, *EFloat, *EStr, *EBool, *EArray:
		return true
	}
	return false
}

func toBVal(expr Expr) BVal {
	switch inner := expr.(type) {
	case *EInt:
		return BInt(*inner)
	case *EFloat:
		return BFloat(*inner)
	case *EStr:
		return BStr(*inner)
	case *EBool:
		return BBool(*inner)
	// case *EArray:
	// 	return BArray(*inner)
	default:
		panic("Only call toBVal on const")
	}
}
