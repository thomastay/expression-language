package parser_test

import (
	"fmt"
	"testing"

	"github.com/thomastay/expression_language/pkg/parser"
)

func TestValidStrings(t *testing.T) {
	var tests = []string{
		// Regular binary expressions
		"1",
		"1 + 1",
		"1 - 1",
		"500 - foo",
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
		"foo.bar()",
		"foo.bar(3*3, 2/2*(4+xoo))",
		// Method calls on identifier (not implemented)
		// "bar(20)",
		// method calls on integers
		"20.to_int",
		// Regressions
		"not temp[i] ? 5 / -2 : 10 * f.foo",
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

func TestInvalidStrings(t *testing.T) {
	var tests = []string{
		// Exponentiation is not allowed
		"1 ** 2",
		// Method call with number as ident
		"foo.3",
		// unmatched
		"[a.i",
		// unmatched
		"(a.i",
		// unmatched
		"a.[i",
		// Use of keywords as identifiers
		// Note: can probably improve error message for this. It currently points to the * as the error,
		// but the error is `not`
		"not * 3",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			_, err := parser.ParseString(tt)
			if err == nil {
				t.Errorf("Should have got an error, but received nothing")
			}
		})
	}
}
