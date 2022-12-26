package vm_test

import (
	"fmt"
	"testing"

	"github.com/thomastay/expression_language/pkg/vm"
)

func TestValidStrings(t *testing.T) {
	var tests = []string{
		// Floats
		"1.1",
		"1.1e10",
		// Ints
		"1", // base 10
		"0x10",
		"0b10",
		"0o70",
		// Calculator
		"1 * 10",
		"1 + 10",
		"1 - 100 - 3",
		"1 / 10",
		"100 / 10 * 3",
		"((10 * 3.0) ? 3 : 10) * 5.0e10",
		"((10 * 3.0) ? 3 : 10) * 5.0",
		// strings
		"'a'",
		"'a' * 2",
		"'a' + 'b'",
		// conditionals
		"0.7 or 9",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			_, err := vm.EvalString(tt)
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
		"((10 * 3.0) ? 3 : 1000000000000000000000000000000000000) * 5.0e1000000000",
		// div 0
		"1 / 0",
		// strings
		"'a' + 2",
		"2 + 'a'",
		"'a' + 2.0",
		"2.0 + 'a'",
		"'a' / 2",
		"2 / 'a'",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			_, err := vm.EvalString(tt)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
		})
	}
}
