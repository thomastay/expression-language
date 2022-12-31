package compiler

import "github.com/thomastay/expression_language/pkg/parser"

// Walks the tree
func Walk(expr parser.Expr, visit Visitor) error {
	switch node := expr.(type) {
	case *parser.EValue, *parser.EUnOp:
		err := Walk(node, visit)
		if err != nil {
			return err
		}
	case *parser.EBinOp:
		err := Walk(node.Left, visit)
		if err != nil {
			return err
		}
		err = Walk(node.Right, visit)
		if err != nil {
			return err
		}
	case *parser.EFieldAccess:
		err := Walk(node.Base, visit)
		if err != nil {
			return err
		}
	case *parser.EIdxAccess:
		err := Walk(node.Base, visit)
		if err != nil {
			return err
		}
		err = Walk(node.Index, visit)
		if err != nil {
			return err
		}
	case *parser.ECond:
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
	case *parser.ECall:
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
	case *parser.EArray:
		var err error
		if err != nil {
			return err
		}
		for _, expr := range []parser.Expr(*node) {
			err = Walk(expr, visit)
			if err != nil {
				return err
			}
		}
	}
	return visit(expr)
}

type Visitor func(expr parser.Expr) error
