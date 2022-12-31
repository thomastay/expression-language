package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Walks the tree in a post order traversal
func walk(expr Expr, visit visitor) walkError {
	var compileErrors []CompileError
	walkAndAdd := func(e Expr) {
		errs := walk(e, visit)
		compileErrors = append(compileErrors, errs...)
	}
	switch node := expr.(type) {
	case *EValue, *EUnOp:
		walkAndAdd(node)
	case *EBinOp:
		walkAndAdd(node.Left)
		walkAndAdd(node.Right)
	case *EFieldAccess:
		walkAndAdd(node.Base)
	case *EIdxAccess:
		walkAndAdd(node.Base)
		walkAndAdd(node.Index)
	case *ECond:
		walkAndAdd(node.Cond)
		walkAndAdd(node.First)
		walkAndAdd(node.Second)
	case *ECall:
		walkAndAdd(node.Base)
		for _, expr := range node.Exprs {
			walkAndAdd(expr)
		}
	case *EArray:
		for _, expr := range *node {
			walkAndAdd(expr)
		}
	}
	errs := visit(expr)
	compileErrors = append(compileErrors, errs...)
	return compileErrors
}

type visitor func(expr Expr) walkError

type walkError []CompileError
