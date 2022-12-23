package parser

import (
	"strings"
	"text/scanner"
)

type Lexer struct {
	sc scanner.Scanner
	// We need to use this since Scanners have no way to peek, so we just store the token here
	// to cache it for use later on
	peekedTok Token
}

func FromString(s string) Lexer {
	l := Lexer{}
	l.sc.Init(strings.NewReader(s))
	return l
}

func (l *Lexer) Next() Token {
	if l.peekedTok != nil {
		t := l.peekedTok
		l.peekedTok = nil
		return t
	}
	tok := l.sc.Scan()
	if tok == scanner.EOF {
		return TokEOF{}
	}
	txt := l.sc.TokenText()
	if len(txt) == 1 {
		r := rune(txt[0])
		var tok Token
		if isOp(r) {
			tok = TokOp(r)
		}
		if isEndExpr(r) {
			tok = TokEndExpr(r)
		}
		if tok != nil {
			return tok
		}
	}
	return TokIdent(txt)
}

func (l *Lexer) Peek() Token {
	if l.peekedTok != nil {
		return l.peekedTok
	}
	tok := l.sc.Scan()
	if tok == scanner.EOF {
		return TokEOF{}
	}
	txt := l.sc.TokenText()
	if len(txt) == 1 {
		r := rune(txt[0])
		var tok Token
		if isOp(r) {
			tok = TokOp(r)
		}
		if isEndExpr(r) {
			tok = TokEndExpr(r)
		}
		if tok != nil {
			l.peekedTok = tok
			return tok
		}
	}
	ident := TokIdent(txt)
	l.peekedTok = ident
	return ident
}

func (l *Lexer) Pos() int {
	if l.peekedTok != nil {
		return l.sc.Pos().Offset - l.peekedTok.Len()
	}
	return l.sc.Pos().Offset
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
	if r == '(' || r == ':' {
		return true
	}
	return false
}

func isEndExpr(r rune) bool {
	return r == ')' || r == ']' || r == ','
}

type Token interface {
	isToken()
	String() string
	Len() int
}

func (t TokIdent) isToken()   {}
func (t TokOp) isToken()      {}
func (t TokEndExpr) isToken() {}
func (t TokEOF) isToken()     {}

type TokIdent string
type TokOp rune
type TokEndExpr rune // end of expression
type TokEOF struct{}

func (t TokIdent) String() string {
	return string(t)
}
func (t TokOp) String() string {
	return string(t)
}
func (t TokEOF) String() string {
	return "EOF"
}
func (t TokEndExpr) String() string {
	return string(t)
}

func (t TokIdent) Len() int {
	return len(t)
}
func (t TokOp) Len() int {
	return 1
}
func (t TokEOF) Len() int {
	return 0
}
func (t TokEndExpr) Len() int {
	return 1
}
