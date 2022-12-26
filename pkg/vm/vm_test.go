package vm_test

import (
	"fmt"
	"testing"

	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
	"github.com/thomastay/expression_language/pkg/vm"
)

func TestValidStrings(t *testing.T) {
	var tests = []string{
		// Floats
		// "1.1",
		// "1.1e10",
		// Ints
		"1", // base 10
		// Calculator
		"1 * 10",
		"1 + 10",
		"1 - 100 - 3",
		"1 / 10",
		"100 / 10 * 3",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			expr, err := parser.ParseString(tt)
			if err != nil {
				t.Fatalf("got %v", err)
			}
			comp := compiler.Compile(expr)
			if len(comp.Errors) > 0 {
				t.Fatal(comp.Errors)
			}
			vm := vm.New(vm.Params{})
			_, err = vm.Eval(comp)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestInvalidStrings(t *testing.T) {
	var tests = []string{
		// Overflow
		"1 + 101000000000000000 * 20000000000000000",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			expr, err := parser.ParseString(tt)
			if err != nil {
				t.Fatalf("got %v", err)
			}
			comp := compiler.Compile(expr)
			if len(comp.Errors) > 0 {
				t.Fatal(comp.Errors)
			}
			vm := vm.New(vm.Params{})
			_, err = vm.Eval(comp)
			if err == nil {
				t.Errorf("Should have got an error, but received nothing")
			}
		})
	}
}
