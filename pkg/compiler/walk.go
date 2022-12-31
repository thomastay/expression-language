package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Walks the tree in a post order traversal
func walk(ptrToExpr *Expr, visit visitor) walkError {
	var compileErrors []CompileError
	walkAndAdd := func(p *Expr) {
		errs := walk(p, visit)
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
		walkAndAdd(&node.Base)
		for _, expr := range node.Exprs {
			walkAndAdd(&expr)
		}
	case *EArray:
		for _, expr := range *node {
			walkAndAdd(&expr)
		}
	}
	errs := visit(ptrToExpr)
	compileErrors = append(compileErrors, errs...)
	return compileErrors
}

type visitor func(ptrToExpr *Expr) walkError

type walkError []CompileError
