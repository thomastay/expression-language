package vm_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/vm"
)

var validStrings = []string{
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
	// unary operators
	"+foo",
	"-10",
	"not 10",
	// strings
	"'a'",
	"'a' * 2",
	"'a' + 'b'",
	// conditionals
	"0.7 or 9",
	// variables
	"a * 30",
	"buzz * 30",
	"a % 3 ? fizz : buzz",
	// Comparison ops
	"10 < 30 ? 20 : 40",
	"10 > 30 ? 20 : 40",
	"10 <= 10 ? 20 : 40",
	"10 >= 10 ? 20 : 40",
	"10.0 < 30 ? 20 : 40",
	"10 > 30.3 ? 20 : 40",
	"10.3 <= 10 ? 20 : 40",
	"10 >= 10.5 ? 20 : 40",
	"'asd' < buzz ? fizz : 'bar'",
	// comparing weird values
	"1 == 4",
	"2 == foo",
	"3 != 'false'",
	"4 == foobar(10)",
	"true == false",
	"true == foo",
	"true != 'false'",
	"true == foobar(10)",
	// Fizzbuzz!
	"a % 3 ? a % 5 ? a : 'buzz' : a % 5 ? fizz : fizzbuzz",
	// Collatz
	"a % 2 ? 3 * a + 1 : a // 2",
	// functions!
	"foobar(123)",
	"ba(123)",
	"vv(123)",
	// objects
	"fooObj.bar * 10",
	"fooObj.baz(30) * 10",
	"fooObj.baz(40)",
	// Arrays
	"[1, 2, 3]",
	"[]",
	"[1, 2, 3] + [4, 5, 6]",
	"[1, 2, 3] * 3",
	"[1, 2, 3][0]",
}

func TestValidStrings(t *testing.T) {
	var tests = validStrings

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			seedVM(vm)
			_, err := vm.EvalString(tt, nil)
			if err != nil {
				t.Error(err)
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
		"[1, 2, 3][5]",
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			seedVM(vm)
			_, err := vm.EvalString(tt, nil)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
		})
	}
}

func TestFizzBuzz(t *testing.T) {
	m := vm.New(vm.Params{})
	s := "i % 3 ? i % 5 ? i : 'buzz' : i % 5 ? fizz : 'fizzbuzz'"
	m.AddStr("fizz", "fizz")
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		t.Fatal("Found compile errors")
	}
	for i := 0; i < 100; i++ {
		m.AddInt("i", int64(i))
		result, err := m.Eval(compilation, nil)
		if err != nil {
			t.Error(err)
		}
		switch fizzStr := result.Val.(type) {
		case bytecode.BStr:
			if string(fizzStr) != fizzBuzz(i) {
				t.Errorf("Wanted %s, got %s", fizzBuzz(i), fizzStr)
			}
		case bytecode.BInt:
			// temp until we can do int->string in the VM
			if fmt.Sprint(fizzStr) != fizzBuzz(i) {
				t.Errorf("Wanted %s, got %s", fizzBuzz(i), fizzStr)
			}
		}
	}
}

func TestCollatz(t *testing.T) {
	m := vm.New(vm.Params{})
	s := "i % 2 == 0 ? i//2 : 3*i + 1"
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		t.Fatal("Found compile errors")
	}
	for i := 1000; i > 1; {
		m.AddInt("i", int64(i))
		result, err := m.Eval(compilation, nil)
		if err != nil {
			t.Error(err)
		}
		switch cltz := result.Val.(type) {
		case bytecode.BInt:
			if int(cltz) != collatz(i) {
				t.Fatalf("Wanted %d, got %d", collatz(i), cltz)
			}
			i = int(cltz)
		default:
			t.Fatal("Bad return")
		}
	}
}

func BenchmarkCollatz(b *testing.B) {
	// Just for fun
	m := vm.New(vm.Params{})
	s := "i % 2 ? 3*i + 1 : i//2 "
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		b.Fatal("Found compile errors")
	}
	for numRuns := 0; numRuns < b.N; numRuns++ {
		for i := 100000; i > 1; {
			m.AddInt("i", int64(i))
			result, err := m.Eval(compilation, nil)
			if err != nil {
				b.Error(err)
			}
			switch cltz := result.Val.(type) {
			case bytecode.BInt:
				i = int(cltz)
			default:
				b.Fatal("Bad return")
			}
		}
	}
}

func BenchmarkCollatzRegular(b *testing.B) {
	// Just for fun
	for numRuns := 0; numRuns < b.N; numRuns++ {
		for i := 100000; i > 1; {
			i = collatz(i)
		}
	}
}

func fizzBuzz(i int) string {
	if i%3 == 0 {
		if i%5 == 0 {
			return "fizzbuzz"
		}
		return "fizz"
	} else if i%5 == 0 {
		return "buzz"
	}
	return fmt.Sprint(i)
}

func collatz(i int) int {
	if i%2 == 0 {
		return i / 2
	}
	return 3*i + 1
}

func seedVM(m vm.VMState) {
	// seed the VM with some useful variables
	m.AddInt("a", 43)
	m.AddInt("b", 2)
	m.AddFloat("foo", 10.5)
	m.AddStr("bar", "I am a string")
	m.AddStr("fizzbuzz", "fizzbuzz")
	m.AddStr("fizz", "fizz")
	m.AddStr("buzz", "buzz")
	m.AddFunc("foobar", func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		return bytecode.BNull{}
	})
	m.AddFunc("ba", func(x bytecode.BVal) (bytecode.BVal, error) {
		log.Println(x)
		xx := x.(bytecode.BInt)
		return bytecode.BFloat(float64(xx) * 43.4), nil
	})
	m.AddFunc("vv",
		func(x bytecode.BVal) {
			log.Println(x)
		})
	fooObjVal := map[string]bytecode.BVal{
		"bar": bytecode.BInt(10),
		"baz": vm.WrapFn("baz", func(x bytecode.BVal) bytecode.BVal {
			log.Println(x)
			xx := x.(bytecode.BInt)
			return bytecode.BFloat(float64(xx) * 43.4)
		}),
		"ba": vm.WrapFn("baz", func(x bytecode.BVal) (bytecode.BVal, error) {
			log.Println(x)
			xx := x.(bytecode.BInt)
			return bytecode.BFloat(float64(xx) * 43.4), nil
		}),
		"vv": vm.WrapFn("baz", func(x bytecode.BVal) {
			log.Println(x)
		}),
	}
	m.AddObject("fooObj", fooObjVal)
}

// ---------------------- Fuzzing ---------------------------------

var testingVMParams = vm.Params{
	MaxMemory:       32768,
	MaxInstructions: 10000,
}

// VM should never fail
func FuzzVM(f *testing.F) {
	testcases := validStrings
	for _, tc := range testcases {
		f.Add(tc)
	}
	vm := vm.New(testingVMParams)
	seedVM(vm)
	f.Fuzz(func(t *testing.T, orig string) {
		vm.EvalString(orig, nil)
	})
}
