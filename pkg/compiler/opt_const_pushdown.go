package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Pushes down constants in the AST
func ConstPushDown(ptrToExpr *Expr) walkError {
	var errs []CompileError
	// Only applies to binary ops
	if node, ok := (*ptrToExpr).(*EBinOp); ok {
		// Illustrative case: numHours * 24 * 60 * 60 * 1000
		//      	        *
		//              *   1000
		//            *  60
		//          *  60
		// numHours  24
		//
		// We go top down, so we first check if the two binOps are the same op, and if both are constant.
		op := node.Op.Value
		if _, ok := commutativeOps[op]; !ok {
			return errs
		}
		if isConst(node.Left) {
			panic("Should have been fixed before this")
		}
		if !isConst(node.Right) {
			return errs
		}
		if lNode, ok := node.Left.(*EBinOp); ok {
			if lNode.Op.Value != op {
				// Not the same op
				return errs
			}
			if isConst(lNode.Left) {
				panic("Should have been fixed before this")
			}
			if isConst(lNode.Right) {
				// jackpot! Fold
				newBinOp := EBinOp{
					Op:    node.Op,
					Left:  lNode.Right,
					Right: node.Right,
				}
				newExpr, err := foldBinaryOpBothConst(&newBinOp)
				if err != nil {
					errs = append(errs, CompileError{
						Err:   err,
						Start: lNode.Op,
						End:   node.Op,
					})
					return errs
				}
				// Rotate right
				lNode.Right = newExpr
				*ptrToExpr = lNode
				// Call again, in case we can push down again
				cErrs := ConstPushDown(ptrToExpr)
				errs = append(errs, cErrs...)
			} else {
				// Else, we PUSH the constant down (hence the name)
				lNode.Right, node.Right = node.Right, lNode.Right
			}
		}
	}
	return errs
}
