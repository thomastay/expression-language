package compiler

import (
	. "github.com/thomastay/expression_language/pkg/ast"
)

// Walks the tree
func Walk(expr Expr, visit Visitor) error {
	switch node := expr.(type) {
	case *EValue, *EUnOp:
		err := Walk(node, visit)
		if err != nil {
			return err
		}
	case *EBinOp:
		err := Walk(node.Left, visit)
		if err != nil {
			return err
		}
		err = Walk(node.Right, visit)
		if err != nil {
			return err
		}
	case *EFieldAccess:
		err := Walk(node.Base, visit)
		if err != nil {
			return err
		}
	case *EIdxAccess:
		err := Walk(node.Base, visit)
		if err != nil {
			return err
		}
		err = Walk(node.Index, visit)
		if err != nil {
			return err
		}
	case *ECond:
		err := Walk(node.Cond, visit)
		if err != nil {
			return err
		}
		err = Walk(node.First, visit)
		if err != nil {
			return err
		}
		err = Walk(node.Second, visit)
		if err != nil {
			return err
		}
	case *ECall:
		err := Walk(node.Base, visit)
		if err != nil {
			return err
		}
		for _, expr := range node.Exprs {
			err = Walk(expr, visit)
			if err != nil {
				return err
			}
		}
	case *EArray:
		var err error
		if err != nil {
			return err
		}
		for _, expr := range []Expr(*node) {
			err = Walk(expr, visit)
			if err != nil {
				return err
			}
		}
	}
	return visit(expr)
}

type Visitor func(expr Expr) error
