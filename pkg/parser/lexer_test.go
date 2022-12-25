package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/thomastay/expression_language/pkg/parser"
)

func TestValidTokens(t *testing.T) {
	syms := parser.Lexer.Symbols()
	var tests = []struct {
		input  string
		tokens []lexer.TokenType
	}{
		// Regular binary expressions
		{"1", []lexer.TokenType{syms["Int"]}},
		{"1 + 1", []lexer.TokenType{syms["Int"], syms["Op"], syms["Int"]}},
		{"500 - foo", []lexer.TokenType{syms["Int"], syms["Op"], syms["Ident"]}},
		// "baz * foo",
		// "baz * foo + 3",
		// "baz / xa + 3 -xxx",
		// // braced expressions
		// "(1 + 2) * 3",
		// // Ternary expressions
		// "baz ? 1 : 2",
		// "baz  ?(potato/2) : x*x",
		// // Indexing operator
		// "a[i]",
		// "a[5 * 2 * (4/3)]",
		// // Method calls on base
		// "foo.bar",
		// "foo.bar(3*3, 2/2*(4+xoo))",
		// Method calls on identifier (not implemented)
		// "bar(20)",
		// method calls on integers (disabled for now since the lexer is broken)
		// "20.to_int",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			genLexer, err := parser.Lexer.Lex("memory", strings.NewReader(tt.input))
			if err != nil {
				t.Errorf("got %v", err)
			}
			peekLexer, err := lexer.Upgrade(genLexer, syms["whitespace"])
			if err != nil {
				t.Errorf("got %v", err)
			}
			tokens := getTokens(peekLexer)
			if len(tt.tokens) != len(tokens) {
				t.Fatalf("wrong number of tokens, wanted %d, got %d", len(tt.tokens), len(tokens))
			}
			for i, tokType := range tt.tokens {
				if tokType != tokens[i].Type {
					t.Errorf("Invalid token %d, wanted %d, got %d", i, tokType, tokens[i].Type)
				}
			}
		})
	}
}

func getTokens(l *lexer.PeekingLexer) (result []*lexer.Token) {
	var tok *lexer.Token
	for tok = l.Next(); tok.Type != lexer.EOF; tok = l.Next() {
		result = append(result, tok)
	}
	return result
}
