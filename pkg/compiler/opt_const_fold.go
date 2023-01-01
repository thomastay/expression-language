package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/runtime"
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
	case *EBinOp:
		if isConst(node.Left) && isConst(node.Right) {
			newExpr, err := foldBinaryOp(node)
			if err != nil {
				errs = append(errs, CompileError{
					Err:   err,
					Start: node.Op,
					End:   node.Op,
				})
			} else {
				*ptrToExpr = newExpr
			}

		}
		// else, swap

	// Do nothing for now (not impl)
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

// TODO
var compilerMemoryLimit = 100000

// Helper function to fold a Binary operation with both children constant
func foldBinaryOp(node *EBinOp) (Expr, error) {
	left := toBVal(node.Left)
	right := toBVal(node.Right)
	var result BVal
	var err error
	// We just use the runtime module to do this
	switch node.Op.Value {
	// simple ops
	case "+":
		result, err = runtime.Add(left, right)
	case "-":
		result, err = runtime.Sub(left, right)
	case "*":
		result, err = runtime.Mul(left, right, compilerMemoryLimit)
	case "/":
		result, err = runtime.Div(left, right)
	case "//":
		result, err = runtime.FloorDiv(left, right)
	case "**":
		result, err = runtime.Pow(left, right)
	// Less simple ops
	case "<", ">", "<=", ">=":
		ord, err := runtime.Cmp(left, right, node.Op.Value)
		if err != nil {
			return nil, err
		}
		b := runtime.OrdToBool(node.Op.Value, ord)
		return (*EBool)(&b), nil
	case "==":
		result := runtime.Eq(left, right)
		return (*EBool)(&result), nil
	case "!=":
		result := !runtime.Eq(left, right)
		return (*EBool)(&result), nil
	// Conditionals
	case "and":
		result := true
		if left.IsTruthy() && right.IsTruthy() {
			return (*EBool)(&result), nil
		}
		result = false
		return (*EBool)(&result), nil
	case "or":
		result := true
		if left.IsTruthy() || right.IsTruthy() {
			return (*EBool)(&result), nil
		}
		result = false
		return (*EBool)(&result), nil
	default:
		panic("not impl")
	}
	if err != nil {
		return nil, err
	}
	return bValToNode(result), nil
}

func isConst(expr Expr) bool {
	// TODO this function needs rework once we implement const array parsing
	switch expr.(type) {
	// arrays MAY NOT BE CONST for now until we implement constant array parsing
	case *EInt, *EFloat, *EStr, *EBool: // , *EArray:
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

func bValToNode(val BVal) Expr {
	switch x := val.(type) {
	case BBool:
		return (*EBool)(&x)
	case BInt:
		return (*EInt)(&x)
	case BFloat:
		return (*EFloat)(&x)
	case BStr:
		return (*EStr)(&x)
	default:
		panic("no other bvals can be nodes (for now)")
	}
}
