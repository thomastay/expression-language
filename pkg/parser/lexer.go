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
	if len(txt) == 1 && isOp(rune(txt[0])) {
		return TokOp(txt[0])
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
	if len(txt) == 1 && isOp(rune(txt[0])) {
		tok := TokOp(txt[0])
		l.peekedTok = tok
		return tok
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

type Token interface {
	String() string
	Len() int
}

type TokIdent string
type TokOp rune
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

func (t TokIdent) Len() int {
	return len(t)
}
func (t TokOp) Len() int {
	return 1
}
func (t TokEOF) Len() int {
	return 0
}
