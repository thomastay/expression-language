package parser_test

import (
	"fmt"
	"testing"

	"github.com/thomastay/kaitai-expr-lang/pkg/parser"
)

func TestValidStrings(t *testing.T) {
	var tests = []string{
		// Regular binary expressions
		"1",
		"1 + 1",
		"1 - 1",
		"1 - foo",
		"baz * foo",
		"baz * foo + 3",
		"baz / xa + 3 -xxx",
		// braced expressions
		"(1 + 2) * 3",
		// Ternary expressions
		"baz ? 1 : 2",
		"baz  ?(potato/2) : x*x",
		// Indexing operator
		"a[i]",
		"a[5 * 2 * (4/3)]",
		// Method calls on base
		"foo.bar",
		"foo.bar(3*3, 2/2*(4+xoo))",
		// Method calls on identifier (not implemented)
		// "bar(20)",
		// method calls on integers (disabled for now since the lexer is broken)
		// "20.to_int",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			_, err := parser.ParseString(tt)
			if err != nil {
				t.Errorf("got %v", err)
			}
		})
	}
}