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
			case TokOp:
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
		return nil, errors.Errorf("Bad token %s", nextVal)
	}

Loop:
	for {
		var op TokOp
		switch nextOp := lexer.Peek().(type) {
		case TokEOF:
			break Loop
		case TokOp:
			op = nextOp
		default:
			return nil, errors.Errorf("Bad op token %s", nextOp)
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
				case TokOp:
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
				switch ident := lexer.Next().(type) {
				case TokIdent:
					// TODO parse expression lists. Assume no params
					lhs = &Call{
						base:   lhs,
						method: TIdent(ident),
						// exprs: nil, // TODO implement!
					}
				default:
					return nil, errors.Errorf("Only identifiers can be used for a method call, found %s", ident)
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

func isOp(r rune) bool {
	tr := TokOp(r)
	_, ok := infixBP[tr]
	if ok {
		return true
	}
	_, ok = prefixBP[tr]
	if ok {
		return true
	}
	_, ok = postFixBP[tr]
	if ok {
		return true
	}
	if r == '(' || r == ')' || r == ']' || r == ':' {
		return true
	}
	return false
}
