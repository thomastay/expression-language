package parser

// Parser based on the Pratt Parser by Matklad

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ParseString(s string) (Expr, error) {
	lexer := FromString(s)
	sexpr, err := exprBP(&lexer, 0)
	if err != nil {
		errWithTrace, ok := errors.Cause(err).(stackTracer)
		if !ok {
			return sexpr, err
		}
		return sexpr, fmt.Errorf(`%w
 | %s
 | %s^--- Parser stopped here

%+v
`, err, s, strings.Repeat(" ", lexer.Pos()), errWithTrace.StackTrace())
	}

	// Possible that exprBP doesn't parse all the string, so check that we've fully consumed everything
	switch t := lexer.Peek().(type) {
	case TokEOF:
		return sexpr, nil
	default:
		return sexpr, errors.New(fmt.Sprintf(`Unparsed character %s at end of Parse.
 | %s
 | %s^--- Parser stopped here
`, t, s, strings.Repeat(" ", lexer.Pos())))
	}
}

// Parsing
func exprBP(lexer *Lexer, minBP int) (Expr, error) {
	var lhs Expr
	var err error
	switch nextVal := lexer.Next().(type) {
	case TokOp:
		if nextVal == '(' {
			// Handle parenthesis
			lhs, err = exprBP(lexer, 0)
			if err != nil {
				return nil, err
			}
			switch end := lexer.Next().(type) {
			case TokEndExpr:
				if end != ')' {
					return lhs, errors.New("Unmatched (")
				}
			default:
				return lhs, errors.New("Unmatched (")
			}
		} else {
			// general operator

			rp := prefixBP[nextVal]
			rhs, err := exprBP(lexer, rp)
			if err != nil {
				return nil, err
			}
			lhs = &EUnOp{
				op:  nextVal,
				val: rhs,
			}
		}
	case TokIdent:
		lhs = &EValue{val: TIdent(nextVal)}
	default:
		return nil, errors.Errorf("Unrecognized token %s %T", nextVal, nextVal)
	}

Loop:
	for {
		var op TokOp
		switch nextOp := lexer.Peek().(type) {
		case TokOp:
			op = nextOp
		case TokEOF:
			break Loop
		case TokEndExpr:
			break Loop
		default:
			return nil, errors.Errorf("Unrecognized token after expr %s, %T", nextOp, nextOp)
		}
		// optional postfix op
		if lp, ok := postFixBP[op]; ok {
			if lp < minBP {
				break
			}
			lexer.Next() // skip op
			switch op {
			case '[':
				// Array indexing
				inner, err := exprBP(lexer, 0)
				if err != nil {
					return nil, err
				}
				switch end := lexer.Next().(type) {
				case TokEndExpr:
					if end != ']' {
						return lhs, errors.New("Unmatched [")
					}
				default:
					return lhs, errors.New("Unmatched [")
				}
				lhs = &EIdxAccess{
					base:  lhs,
					index: inner,
				}
			case '.':
				// Call operator with a base
				lhs, err = parseCallWithBase(lhs, lexer)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.Errorf("No other postfix operators %s", op)
			}
			continue
		}

		// infix ops
		if infixPowers, ok := infixBP[op]; ok {
			lp, rp := infixPowers.l, infixPowers.r
			if lp < minBP {
				break
			}
			// Skip the operator token
			lexer.Next()
			if op == '?' {
				// special case ternaries
				inner, err := exprBP(lexer, 0)
				if err != nil {
					return nil, err
				}
				switch end := lexer.Next().(type) {
				case TokOp:
					if end != ':' {
						return lhs, errors.New("Unmatched ?")
					}
				default:
					return lhs, errors.New("Unmatched ?")
				}
				rhs, err := exprBP(lexer, rp)
				if err != nil {
					return nil, err
				}
				lhs = &ECond{
					cond:   lhs,
					first:  inner,
					second: rhs,
				}
			} else {
				rhs, err := exprBP(lexer, rp)
				if err != nil {
					return nil, err
				}
				lhs = &EBinOp{
					op:    op,
					left:  lhs,
					right: rhs,
				}
			}
			continue
		}
		break
	}
	return lhs, nil
}

func parseCallWithBase(base Expr, lexer *Lexer) (*Call, error) {
	var expr *Call

	switch ident := lexer.Next().(type) {
	case TokIdent:
		// A method call is a base.ident, then followed by possible expression list.
		var exprList ExprList
		switch next := lexer.Peek().(type) {
		case TokOp:
			if next == '(' {
				// It is an expression list. Start to parse.
				lexer.Next()
			ExprLoop:
				for {
					param, err := exprBP(lexer, 0)
					if err != nil {
						return nil, err
					}
					exprList = append(exprList, param)
					switch op := lexer.Next().(type) {
					case TokEndExpr:
						switch op {
						case ',':
							continue
						case ')':
							break ExprLoop
						default:
							return nil, errors.Errorf("Unrecognized end of expression in param list: %s", op)
						}
					default:
						return nil, errors.Errorf("Unrecognized token in parsing param list: %s", op)
					}
				}
			}
			// else, fallthrough
		default:
			// fallthrough, do nothing here.
		}
		expr = &Call{
			base:   base,
			method: TIdent(ident),
			exprs:  exprList,
		}
	default:
		return nil, errors.Errorf("Only identifiers can be used for a method call, found %s", ident)
	}
	return expr, nil
}

type InfixBP struct {
	l, r int
}

var infixBP = map[TokOp]InfixBP{
	'?': {2, 1},
	'+': {3, 4},
	'-': {3, 4},
	'*': {5, 6},
	'/': {5, 6},
}

var prefixBP = map[TokOp]int{
	'+': 7,
	'-': 7,
}

var postFixBP = map[TokOp]int{
	// This is a Call operator on a base class
	'.': 11,
	// This is the indexing operator
	'[': 9,
}
