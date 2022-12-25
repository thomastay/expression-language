// Code generated by Participle. DO NOT EDIT.
package parser

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"unicode/utf8"
	"regexp/syntax"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var _ syntax.Op
var _ fmt.State
const _ = utf8.RuneError

var BackRefCache sync.Map
var Lexer lexer.Definition = lexerDefinitionImpl{}

type lexerDefinitionImpl struct {}

func (lexerDefinitionImpl) Symbols() map[string]lexer.TokenType {
	return map[string]lexer.TokenType{
      "Bool": -15,
      "DoubleString": -12,
      "DoubleStringEnd": -2,
      "EOF": -1,
      "EndExpr": -17,
      "Float": -19,
      "Ident": -18,
      "Int": -20,
      "Op": -16,
      "SingleString": -13,
      "whitespace": -14,
	}
}

func (lexerDefinitionImpl) LexString(filename string, s string) (lexer.Lexer, error) {
	return &lexerImpl{
		s: s,
		pos: lexer.Position{
			Filename: filename,
			Line:     1,
			Column:   1,
		},
		states: []lexerState{ {name: "Root"} },
	}, nil
}

func (d lexerDefinitionImpl) LexBytes(filename string, b []byte) (lexer.Lexer, error) {
	return d.LexString(filename, string(b))
}

func (d lexerDefinitionImpl) Lex(filename string, r io.Reader) (lexer.Lexer, error) {
	s := &strings.Builder{}
	_, err := io.Copy(s, r)
	if err != nil {
		return nil, err
	}
	return d.LexString(filename, s.String())
}

type lexerState struct {
	name    string
	groups  []string
}

type lexerImpl struct {
	s       string
	p       int
	pos     lexer.Position
	states  []lexerState
}

func (l *lexerImpl) Next() (lexer.Token, error) {
	if l.p == len(l.s) {
		return lexer.EOFToken(l.pos), nil
	}
	var (
		state = l.states[len(l.states)-1]
		groups []int
		sym lexer.TokenType
	)
	switch state.name {
	case "DoubleString":if match := matchDoubleStringEnd(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -2
			groups = match[:]
			l.states = l.states[:len(l.states)-1]
		}
	case "Expr":if match := matchDoubleString(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -12
			groups = match[:]
			l.states = append(l.states, lexerState{name: "DoubleString"})
		} else if match := matchSingleString(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -13
			groups = match[:]
		} else if match := matchwhitespace(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -14
			groups = match[:]
		} else if match := matchBool(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -15
			groups = match[:]
		} else if match := matchOp(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -16
			groups = match[:]
		} else if match := matchEndExpr(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -17
			groups = match[:]
		} else if match := matchIdent(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -18
			groups = match[:]
		} else if match := matchFloat(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -19
			groups = match[:]
		} else if match := matchInt(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -20
			groups = match[:]
		}
	case "Root":if match := matchDoubleString(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -12
			groups = match[:]
			l.states = append(l.states, lexerState{name: "DoubleString"})
		} else if match := matchSingleString(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -13
			groups = match[:]
		} else if match := matchwhitespace(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -14
			groups = match[:]
		} else if match := matchBool(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -15
			groups = match[:]
		} else if match := matchOp(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -16
			groups = match[:]
		} else if match := matchEndExpr(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -17
			groups = match[:]
		} else if match := matchIdent(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -18
			groups = match[:]
		} else if match := matchFloat(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -19
			groups = match[:]
		} else if match := matchInt(l.s, l.p, l.states[len(l.states)-1].groups); match[1] != 0 {
			sym = -20
			groups = match[:]
		}
	}
	if groups == nil {
		sample := []rune(l.s[l.p:])
		if len(sample) > 16 {
			sample = append(sample[:16], []rune("...")...)
		}
		return lexer.Token{}, participle.Errorf(l.pos, "invalid input text %q", string(sample))
	}
	pos := l.pos
	span := l.s[groups[0]:groups[1]]
	l.p = groups[1]
	l.pos.Advance(span)
	return lexer.Token{
		Type:  sym,
		Value: span,
		Pos:   pos,
	}, nil
}

func (l *lexerImpl) sgroups(match []int) []string {
	sgroups := make([]string, len(match)/2)
	for i := 0; i < len(match)-1; i += 2 {
		sgroups[i/2] = l.s[l.p+match[i]:l.p+match[i+1]]
	}
	return sgroups
}
// "
func matchDoubleStringEnd(s string, p int, backrefs []string) (groups [2]int) {
if p < len(s) && s[p] == '"' {
groups[0] = p
groups[1] = p + 1
}
return
}

// "
func matchDoubleString(s string, p int, backrefs []string) (groups [2]int) {
if p < len(s) && s[p] == '"' {
groups[0] = p
groups[1] = p + 1
}
return
}

// '[^']*'
func matchSingleString(s string, p int, backrefs []string) (groups [2]int) {
// ' (Literal)
l0 := func(s string, p int) int {
if p < len(s) && s[p] == '\'' { return p+1 }
return -1
}
// [^'] (CharClass)
l1 := func(s string, p int) int {
if len(s) <= p { return -1 }
var (rn rune; n int)
if s[p] < utf8.RuneSelf {
  rn, n = rune(s[p]), 1
} else {
  rn, n = utf8.DecodeRuneInString(s[p:])
}
switch {
case rn >= '\x00' && rn <= '&': return p+1
case rn >= '(' && rn <= '\U0010ffff': return p+n
}
return -1
}
// [^']* (Star)
l2 := func(s string, p int) int {
for len(s) > p {
if np := l1(s, p); np == -1 { return p } else { p = np }
}
return p
}
// '[^']*' (Concat)
l3 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
if p = l2(s, p); p == -1 { return -1 }
if p = l0(s, p); p == -1 { return -1 }
return p
}
np := l3(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [\t-\n\f-\r ]+
func matchwhitespace(s string, p int, backrefs []string) (groups [2]int) {
// [\t-\n\f-\r ] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '\t' && rn <= '\n': return p+1
case rn >= '\f' && rn <= '\r': return p+1
case rn == ' ': return p+1
}
return -1
}
// [\t-\n\f-\r ]+ (Plus)
l1 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
np := l1(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// true|false
func matchBool(s string, p int, backrefs []string) (groups [2]int) {
// true (Literal)
l0 := func(s string, p int) int {
if p+4 <= len(s) && s[p:p+4] == "true" { return p+4 }
return -1
}
// false (Literal)
l1 := func(s string, p int) int {
if p+5 <= len(s) && s[p:p+5] == "false" { return p+5 }
return -1
}
// true|false (Alternate)
l2 := func(s string, p int) int {
if np := l0(s, p); np != -1 { return np }
if np := l1(s, p); np != -1 { return np }
return -1
}
np := l2(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// and|not|or|\+=|-=|\*=|/=|[\(\*-\+\--/:\?\[]
func matchOp(s string, p int, backrefs []string) (groups [2]int) {
// and (Literal)
l0 := func(s string, p int) int {
if p+3 <= len(s) && s[p:p+3] == "and" { return p+3 }
return -1
}
// not (Literal)
l1 := func(s string, p int) int {
if p+3 <= len(s) && s[p:p+3] == "not" { return p+3 }
return -1
}
// or (Literal)
l2 := func(s string, p int) int {
if p+2 <= len(s) && s[p:p+2] == "or" { return p+2 }
return -1
}
// \+= (Literal)
l3 := func(s string, p int) int {
if p+2 <= len(s) && s[p:p+2] == "+=" { return p+2 }
return -1
}
// -= (Literal)
l4 := func(s string, p int) int {
if p+2 <= len(s) && s[p:p+2] == "-=" { return p+2 }
return -1
}
// \*= (Literal)
l5 := func(s string, p int) int {
if p+2 <= len(s) && s[p:p+2] == "*=" { return p+2 }
return -1
}
// /= (Literal)
l6 := func(s string, p int) int {
if p+2 <= len(s) && s[p:p+2] == "/=" { return p+2 }
return -1
}
// [\(\*-\+\--/:\?\[] (CharClass)
l7 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn == '(': return p+1
case rn >= '*' && rn <= '+': return p+1
case rn >= '-' && rn <= '/': return p+1
case rn == ':': return p+1
case rn == '?': return p+1
case rn == '[': return p+1
}
return -1
}
// and|not|or|\+=|-=|\*=|/=|[\(\*-\+\--/:\?\[] (Alternate)
l8 := func(s string, p int) int {
if np := l0(s, p); np != -1 { return np }
if np := l1(s, p); np != -1 { return np }
if np := l2(s, p); np != -1 { return np }
if np := l3(s, p); np != -1 { return np }
if np := l4(s, p); np != -1 { return np }
if np := l5(s, p); np != -1 { return np }
if np := l6(s, p); np != -1 { return np }
if np := l7(s, p); np != -1 { return np }
return -1
}
np := l8(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [\),\]]
func matchEndExpr(s string, p int, backrefs []string) (groups [2]int) {
// [\),\]] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch rn {
case ')',',',']': return p+1
}
return -1
}
np := l0(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [A-Za-z][0-9A-Z_a-z]*
func matchIdent(s string, p int, backrefs []string) (groups [2]int) {
// [A-Za-z] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= 'A' && rn <= 'Z': return p+1
case rn >= 'a' && rn <= 'z': return p+1
}
return -1
}
// [0-9A-Z_a-z] (CharClass)
l1 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '0' && rn <= '9': return p+1
case rn >= 'A' && rn <= 'Z': return p+1
case rn == '_': return p+1
case rn >= 'a' && rn <= 'z': return p+1
}
return -1
}
// [0-9A-Z_a-z]* (Star)
l2 := func(s string, p int) int {
for len(s) > p {
if np := l1(s, p); np == -1 { return p } else { p = np }
}
return p
}
// [A-Za-z][0-9A-Z_a-z]* (Concat)
l3 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
if p = l2(s, p); p == -1 { return -1 }
return p
}
np := l3(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [0-9]*\.[0-9]+(e[0-9]+)?
func matchFloat(s string, p int, backrefs []string) (groups [4]int) {
// [0-9] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '0' && rn <= '9': return p+1
}
return -1
}
// [0-9]* (Star)
l1 := func(s string, p int) int {
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
// \. (Literal)
l2 := func(s string, p int) int {
if p < len(s) && s[p] == '.' { return p+1 }
return -1
}
// [0-9]+ (Plus)
l3 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
for len(s) > p {
if np := l0(s, p); np == -1 { return p } else { p = np }
}
return p
}
// e (Literal)
l4 := func(s string, p int) int {
if p < len(s) && s[p] == 'e' { return p+1 }
return -1
}
// e[0-9]+ (Concat)
l5 := func(s string, p int) int {
if p = l4(s, p); p == -1 { return -1 }
if p = l3(s, p); p == -1 { return -1 }
return p
}
// (e[0-9]+) (Capture)
l6 := func(s string, p int) int {
np := l5(s, p)
if np != -1 {
  groups[2] = p
  groups[3] = np
}
return np}
// (e[0-9]+)? (Quest)
l7 := func(s string, p int) int {
if np := l6(s, p); np != -1 { return np }
return p
}
// [0-9]*\.[0-9]+(e[0-9]+)? (Concat)
l8 := func(s string, p int) int {
if p = l1(s, p); p == -1 { return -1 }
if p = l2(s, p); p == -1 { return -1 }
if p = l3(s, p); p == -1 { return -1 }
if p = l7(s, p); p == -1 { return -1 }
return p
}
np := l8(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}

// [1-9][0-9]*
func matchInt(s string, p int, backrefs []string) (groups [2]int) {
// [1-9] (CharClass)
l0 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '1' && rn <= '9': return p+1
}
return -1
}
// [0-9] (CharClass)
l1 := func(s string, p int) int {
if len(s) <= p { return -1 }
rn := s[p]
switch {
case rn >= '0' && rn <= '9': return p+1
}
return -1
}
// [0-9]* (Star)
l2 := func(s string, p int) int {
for len(s) > p {
if np := l1(s, p); np == -1 { return p } else { p = np }
}
return p
}
// [1-9][0-9]* (Concat)
l3 := func(s string, p int) int {
if p = l0(s, p); p == -1 { return -1 }
if p = l2(s, p); p == -1 { return -1 }
return p
}
np := l3(s, p)
if np == -1 {
  return
}
groups[0] = p
groups[1] = np
return
}
