package optimizer

import "github.com/thomastay/expression_language/pkg/parser"

// Implements Visitor
func ConstFold(expr parser.Expr) error {
	var err error
	switch node := expr.(type) {
	case *parser.EValue:
		return nil
	case *parser.EUnOp:
		switch node.Val.(type) {
		case *parser.EValue:

		}

	case *parser.EBinOp:
	case *parser.EFieldAccess:
	case *parser.EIdxAccess:
	case *parser.ECond:
	case *parser.ECall:
	case *parser.EArray:
	default:
		panic("AST type is not impl")
	}
	return err
}
