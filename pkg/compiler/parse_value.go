package compiler

import (
	"log"
	"strconv"

	. "github.com/thomastay/expression_language/pkg/ast"
	"github.com/thomastay/expression_language/pkg/parser"
)

func ParseValue(ptrToExpr *Expr) walkError {
	var errs []CompileError
	expr := *ptrToExpr
	switch node := expr.(type) {
	case *EValue:
		switch node.Val.Type {
		case parser.TokInt, parser.TokHexInt, parser.TokOctInt, parser.TokBinInt:
			tok := node.Val
			val, err := strconv.ParseInt(tok.Value, 0, 64)
			if err != nil {
				errs = append(errs, CompileError{
					Err:   err,
					Start: tok,
					End:   tok,
				})
			}
			// override node
			*ptrToExpr = (*EInt)(&val)
		case parser.TokFloat:
			tok := node.Val
			val, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				errs = append(errs, CompileError{
					Err:   err,
					Start: tok,
					End:   tok,
				})
			}
			// override node
			*ptrToExpr = (*EFloat)(&val)
		case parser.TokSingleString:
			val := node.Val.Value
			val = val[1 : len(val)-1]
			// override node
			*ptrToExpr = (*EStr)(&val)
		case parser.TokIdent:
			val := node.Val.Value
			// override node
			*ptrToExpr = (*EIdent)(&val)
		case parser.TokBool:
			val := node.Val.Value
			boolVal := false
			if val == "true" {
				boolVal = true
			}
			*ptrToExpr = (*EBool)(&boolVal)
		default:
			log.Panicf("Token %s type %d not implemented", node.Val.Value, node.Val.Type)
		}

	// else do nothing
	case *EInt:
	case *EFloat:
	case *EStr:
	case *EIdent:
	case *EBool:
	case *EUnOp:
	case *EBinOp:
	case *EFieldAccess:
	case *EIdxAccess:
	case *ECond:
	case *ECall:
	case *EArray:
	default:
		log.Panicf("AST type %T is not impl", expr)
	}

	return errs
}
