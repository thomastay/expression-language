package parser

// Parser based on the Pratt Parser by Matklad

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ParseString(s string) (Expr, error) {
	genLexer, err := Lexer.Lex("memory", strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	peekLexer, err := lexer.Upgrade(genLexer, Lexer.Symbols()["whitespace"])
	if err != nil {
		return nil, err
	}
	sexpr, err := exprBP(peekLexer, 0)
	if err != nil {
		errWithTrace, ok := errors.Cause(err).(stackTracer)
		if !ok {
			return sexpr, err
		}
		return sexpr, fmt.Errorf(`%w
 | %s
 | %s^--- Parser stopped here

%+v
`, err, s, strings.Repeat(" ", peekLexer.Peek().Pos.Column-1), errWithTrace.StackTrace())
	}

	// Possible that exprBP doesn't parse all the string, so check that we've fully consumed everything
	t := peekLexer.Peek()
	switch t.Type {
	case lexer.EOF:
		return sexpr, nil
	default:
		return sexpr, errors.New(fmt.Sprintf(`Unparsed character %s at end of Parse.
 | %s
 | %s^--- Parser stopped here
`, t, s, strings.Repeat(" ", t.Pos.Column-1)))
	}
}

// Parsing
func exprBP(lex *lexer.PeekingLexer, minBP int) (Expr, error) {
	var lhs Expr
	var err error
	nextVal := lex.Next()
	switch nextVal.Type {
	case TokOp:
		if nextVal.Value == "(" {
			// Handle parenthesis
			lhs, err = exprBP(lex, 0)
			if err != nil {
				return nil, err
			}
			end := lex.Next()
			switch end.Type {
			case TokEndExpr:
				if end.Value != ")" {
					return lhs, errors.New("Unmatched (")
				}
			default:
				return lhs, errors.New("Unmatched (")
			}
		} else {
			// general operator

			rp := prefixBP[nextVal.Value]
			rhs, err := exprBP(lex, rp)
			if err != nil {
				return nil, err
			}
			lhs = &EUnOp{
				op:  nextVal,
				val: rhs,
			}
		}
	case TokIdent:
		lhs = &EValue{val: nextVal}
	case TokInt:
		lhs = &EValue{val: nextVal}
	default:
		return nil, errors.Errorf("Unrecognized token %s %d", nextVal, nextVal.Type)
	}

Loop:
	for {
		op := lex.Peek()
		switch op.Type {
		case TokOp:
			// do nothing, continue
		case lexer.EOF:
			break Loop
		case TokEndExpr:
			break Loop
		default:
			return nil, errors.Errorf("Unrecognized token after expr %s, %d", op, op.Type)
		}
		// optional postfix op
		if lp, ok := postFixBP[op.Value]; ok {
			if lp < minBP {
				break
			}
			lex.Next() // skip op
			switch op.Value {
			case "[":
				// Array indexing
				inner, err := exprBP(lex, 0)
				if err != nil {
					return nil, err
				}
				end := lex.Next()
				switch end.Type {
				case TokEndExpr:
					if end.Value != "]" {
						return lhs, errors.New("Unmatched [")
					}
				default:
					return lhs, errors.New("Unmatched [")
				}
				lhs = &EIdxAccess{
					base:  lhs,
					index: inner,
				}
			case ".":
				// Call operator with a base
				lhs, err = parseCallWithBase(lhs, lex)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.Errorf("No other postfix operators %s", op)
			}
			continue
		}

		// infix ops
		if infixPowers, ok := infixBP[op.Value]; ok {
			lp, rp := infixPowers.l, infixPowers.r
			if lp < minBP {
				break
			}
			// Skip the operator token
			lex.Next()
			if op.Value == "?" {
				// special case ternaries
				inner, err := exprBP(lex, 0)
				if err != nil {
					return nil, err
				}
				end := lex.Next()
				switch end.Type {
				case TokOp:
					if end.Value != ":" {
						return lhs, errors.New("Unmatched ?")
					}
				default:
					return lhs, errors.New("Unmatched ?")
				}
				rhs, err := exprBP(lex, rp)
				if err != nil {
					return nil, err
				}
				lhs = &ECond{
					cond:   lhs,
					first:  inner,
					second: rhs,
				}
			} else {
				rhs, err := exprBP(lex, rp)
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

func parseCallWithBase(base Expr, lex *lexer.PeekingLexer) (*Call, error) {
	var expr *Call

	ident := lex.Next()
	switch ident.Type {
	case TokIdent:
		// A method call is a base.ident, then followed by possible expression list.
		var exprList ExprList
		next := lex.Peek()
		switch next.Type {
		case TokOp:
			if next.Value == "(" {
				// It is an expression list. Start to parse.
				lex.Next()
			ExprLoop:
				for {
					param, err := exprBP(lex, 0)
					if err != nil {
						return nil, err
					}
					exprList = append(exprList, param)
					op := lex.Next()
					switch op.Type {
					case TokEndExpr:
						switch op.Value {
						case ",":
							continue
						case ")":
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
			method: ident,
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

var infixBP = map[string]InfixBP{
	"?": {6, 5},
	"+": {7, 8},
	"-": {7, 8},
	"*": {9, 10},
	"/": {9, 10},
	"and": {3, 4},
	"or": {2, 1},
}

var prefixBP = map[string]int{
	"+": 7,
	"-": 7,
	"not": 5,
}

var postFixBP = map[string]int{
	// This is a Call operator on a base class
	".": 13,
	// This is the indexing operator
	"[": 13,
}

var TokOp = Lexer.Symbols()["Op"]
var TokInt = Lexer.Symbols()["Int"]
var TokIdent = Lexer.Symbols()["Ident"]
var TokEndExpr = Lexer.Symbols()["EndExpr"]
