package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Helper function, don't call this.
func walkInner(ptrToExpr *Expr, visit visitor, preorder bool) walkError {
	var compileErrors []CompileError
	walkAndAdd := func(p *Expr) {
		errs := walk(p, visit)
		compileErrors = append(compileErrors, errs...)
	}
	if preorder {
		errs := visit(ptrToExpr)
		compileErrors = append(compileErrors, errs...)
	}
	switch node := (*ptrToExpr).(type) {
	case *EValue:
	case *EUnOp:
		walkAndAdd(&node.Val)
	case *EBinOp:
		walkAndAdd(&node.Left)
		walkAndAdd(&node.Right)
	case *EFieldAccess:
		walkAndAdd(&node.Base)
	case *EIdxAccess:
		walkAndAdd(&node.Base)
		walkAndAdd(&node.Index)
	case *ECond:
		walkAndAdd(&node.Cond)
		walkAndAdd(&node.First)
		walkAndAdd(&node.Second)
	case *ECall:
		if node.Base != nil {
			walkAndAdd(&node.Base)
		}
		for i := range node.Exprs {
			walkAndAdd(&node.Exprs[i])
		}
	case *EArray:
		for i := range *node {
			walkAndAdd(&(*node)[i])
		}
	}
	if !preorder {
		errs := visit(ptrToExpr)
		compileErrors = append(compileErrors, errs...)
	}
	return compileErrors
}

// Walks the tree in a post order traversal
func walk(ptrToExpr *Expr, visit visitor) walkError {
	return walkInner(ptrToExpr, visit, false)
}

// Walks the tree in a pre order traversal
func walkTopDown(ptrToExpr *Expr, visit visitor) walkError {
	return walkInner(ptrToExpr, visit, true)
}

type visitor func(ptrToExpr *Expr) walkError

type walkError []CompileError
