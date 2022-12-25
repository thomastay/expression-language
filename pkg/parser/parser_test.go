package parser_test

import (
	"fmt"
	"testing"

	"github.com/thomastay/expression_language/pkg/parser"
)

func TestValidStrings(t *testing.T) {
	var tests = []string{
		// Floats
		"1.1",
		"1.1e10",
		// Ints
		"1", // base 10
		// "0x10", // hex (not implemented)
		// "0o30", // octal (not implemented)
		// Regular binary expressions
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
		// Method calls on identifier
		"bar()",
		"bar(20)",
		"bar(20, 30, 40)",
		// method calls on integers
		"20.to_int",
		// Arrays (not impl)
		// "[1, 2, 3]",
		// // We accept mixed arrays in the parser, can reject in sema
		// "[1.1, 2.2, potato]",
		// // Arrays have methods
		// "[1.1, 2.2, potato].len",
		// // Arrays are expressions
		// "[1.1, 2.2, potato] * 2",
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
