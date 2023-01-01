package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/runtime"
)

// Constant folding pass does two things:
//  1. If node is a binary op with two constants, it replaces them
//  1. If node is binary op with one EValue, it rotates the constant to the RHS
//     ^^ this trick is to enable constant pushdown, which will be implemented next
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
			switch node.Val.(type) {
			case *EInt, *EFloat, *EBool, *EStr:
				bVal := toBVal(node.Val)
				result, err := runtime.Negate(bVal)
				if err != nil {
					errs = append(errs, CompileError{
						Err:   err,
						Start: node.Op,
						End:   node.Op,
					})
				} else {
					*ptrToExpr = bValToNode(result)
				}
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
			case *EInt, *EFloat, *EBool, *EStr:
				bVal := toBVal(node.Val)
				flip := !bVal.IsTruthy()
				*ptrToExpr = (*EBool)(&flip)
			case *EArray:
				// special case this for now until const arrays
				flip := len(*inner) == 0
				*ptrToExpr = (*EBool)(&flip)
			default:
				// do nothing, fallthrough
			}
		}
	case *EBinOp:
		if isConst(node.Left) && isConst(node.Right) {
			newExpr, err := foldBinaryOpBothConst(node)
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
		// Some other optimizations
		// 1. If LHS is const and the op is OR, then we can immediately fold
		//    e.g.  10 or x and y becomes just 10
		if isConst(node.Left) && node.Op.Value == "or" {
			left := toBVal(node.Left)
			if left.IsTruthy() {
				*ptrToExpr = node.Left
			} else {
				*ptrToExpr = node.Right
			}
		}
		// 2. If either is const and falsey, and the op is AND, then we can immediately fold
		if node.Op.Value == "and" {
			if isConst(node.Right) {
				right := toBVal(node.Right)
				if !right.IsTruthy() {
					*ptrToExpr = node.Left
				}
			}
			if isConst(node.Left) {
				left := toBVal(node.Left)
				if !left.IsTruthy() {
					*ptrToExpr = node.Right
				}
			}
		}
		// else, swap
		if _, ok := commutativeOps[node.Op.Value]; ok && isConst(node.Left) {
			node.Left, node.Right = node.Right, node.Left
		}

	case *ECond:
		if isConst(node.Cond) {
			val := toBVal(node.Cond)
			if val.IsTruthy() {
				*ptrToExpr = node.First
			} else {
				*ptrToExpr = node.Second
			}
		}
		// else, don't optimize cond

	// Do nothing for now (not impl)
	case *EFieldAccess:
	case *EIdxAccess:
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
func foldBinaryOpBothConst(node *EBinOp) (Expr, error) {
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
	case "%":
		result, err = runtime.Modulo(left, right)
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
		if right.IsTruthy() {
			return bValToNode(right), nil
		}
		return bValToNode(left), nil
	case "or":
		if left.IsTruthy() {
			return bValToNode(left), nil
		}
		return bValToNode(right), nil
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

var commutativeOps = map[string]struct{}{
	// TODO - and /
	"+":  struct{}{},
	"*":  struct{}{},
	"==": struct{}{},
}
