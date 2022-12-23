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
		// Method calls
		"foo.bar",
		"foo.bar(3*3, 2/2*(4+xoo))",
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
