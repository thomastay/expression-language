package vm_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
	"github.com/thomastay/expression_language/pkg/runtime"
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
	// Booleans are treated as integers
	"true == 1",
	"true == 1.1",
	"1.1 == false",
	"1 == true",
	"true + 1",
	"true + 1.1",
	"1.1 + false",
	"1 - true",
	"true - 1",
	"true - 1.1",
	"1.1 - false",
	"1 * true",
	"true * 1",
	"true * 1.1",
	"1.1 * false",
	"1 * true",
	"[1, 2, 3] * true",
	"false * [1, 2, 3]",
	"'1, 2, 3' * true",
	"false * '1, 2, 3'",
	"1 / true",
	"true / 1",
	"true / 1.1",
	"1.1 / true",
	"1 // true",
	"true // 1",
	"true // 1.1",
	"1.1 // true",
	"1 % true",
	"true % 1",
	"true % 1.1",
	"1.1 % true",
	"1 < true",
	"true < 1",
	"true < 1.1",
	"1.1 < false",
	"1 <= true",
	"true <= 1",
	"true <= 1.1",
	"1.1 <= false",
	"1 > true",
	"true > 1",
	"true > 1.1",
	"1.1 > false",
	"1 >= true",
	"true >= 1",
	"true >= 1.1",
	"1.1 >= false",
	"[1, 2,3][true]",
	"[1, 2,3][false]",
	"+true",
	"+false",
	"-true",
	"-false",
	"not true",
	"not false",
	"10 or unknown1 and unknown2", // folding should work!
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
	// regressions
	"fooObj != 1",
}

type InputOutput struct {
	in       string
	expected bytecode.BVal
}

// This list is basically autogenerated (with some human changes)
// with the following command:
// (Powershell)
// go test -v -count=1 ./pkg/vm -run=TestValidStrings  | sls "^\{"  > out.txt
// (Unix Shell)
// go test -v -count=1 ./pkg/vm -run=TestValidStrings  | grep "^\{"  > out.txt
var validStringsInOut = []InputOutput{
	{"1.1", bytecode.BFloat(1.100000)},
	{"1.1e10", bytecode.BFloat(11000000000.000000)},
	{"1", bytecode.BInt(1)},
	{"0x10", bytecode.BInt(16)},
	{"0b10", bytecode.BInt(2)},
	{"0o70", bytecode.BInt(56)},
	{"1 * 10", bytecode.BInt(10)},
	{"1 + 10", bytecode.BInt(11)},
	{"1 - 100 - 3", bytecode.BInt(-96)},
	{"1 / 10", bytecode.BFloat(0.100000)},
	{"100 / 10 * 3", bytecode.BFloat(30.000000)},
	{"((10 * 3.0) ? 3 : 10) * 5.0e10", bytecode.BFloat(150000000000.000000)},
	{"((10 * 3.0) ? 3 : 10) * 5.0", bytecode.BFloat(15.000000)},
	{"+foo", bytecode.BFloat(10.500000)},
	{"-10", bytecode.BInt(-10)},
	{"not 10", bytecode.BBool(false)},
	{"'a'", bytecode.BStr("a")},
	{"'a' * 2", bytecode.BStr("aa")},
	{"'a' + 'b'", bytecode.BStr("ab")},
	{"0.7 or 9", bytecode.BFloat(0.700000)},
	{"a * 30", bytecode.BInt(1290)},
	{"buzz * 30", bytecode.BStr("buzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzzbuzz")},
	{"a % 3 ? fizz : buzz", bytecode.BStr("fizz")},
	{"10 < 30 ? 20 : 40", bytecode.BInt(20)},
	{"10 > 30 ? 20 : 40", bytecode.BInt(40)},
	{"10 <= 10 ? 20 : 40", bytecode.BInt(20)},
	{"10 >= 10 ? 20 : 40", bytecode.BInt(20)},
	{"10.0 < 30 ? 20 : 40", bytecode.BInt(20)},
	{"10 > 30.3 ? 20 : 40", bytecode.BInt(40)},
	{"10.3 <= 10 ? 20 : 40", bytecode.BInt(40)},
	{"10 >= 10.5 ? 20 : 40", bytecode.BInt(40)},
	{"'asd' < buzz ? fizz : 'bar'", bytecode.BStr("fizz")},
	{"1 == 4", bytecode.BBool(false)},
	{"2 == foo", bytecode.BBool(false)},
	{"3 != 'false'", bytecode.BBool(true)},
	{"4 == foobar(10)", bytecode.BBool(false)},
	{"true == false", bytecode.BBool(false)},
	{"true == foo", bytecode.BBool(false)},
	{"true != 'false'", bytecode.BBool(true)},
	{"true == foobar(10)", bytecode.BBool(false)},
	{"true == 1", bytecode.BBool(true)},
	{"true == 1.1", bytecode.BBool(false)},
	{"1.1 == false", bytecode.BBool(false)},
	{"1 == true", bytecode.BBool(true)},
	{"true + 1", bytecode.BInt(2)},
	{"true + 1.1", bytecode.BFloat(2.100000)},
	{"1.1 + false", bytecode.BFloat(1.100000)},
	{"1 - true", bytecode.BInt(0)},
	{"true - 1", bytecode.BInt(0)},
	// {"true - 1.1", bytecode.BFloat(-0.100000)},
	{"1.1 - false", bytecode.BFloat(1.100000)},
	{"1 * true", bytecode.BInt(1)},
	{"true * 1", bytecode.BInt(1)},
	{"true * 1.1", bytecode.BFloat(1.100000)},
	{"1.1 * false", bytecode.BFloat(0.000000)},
	{"1 * true", bytecode.BInt(1)},
	// { "[1, 2, 3] * true", bytecode.BArray([1, 2, 3]) },
	// { "false * [1, 2, 3]", bytecode.BArray([]) },
	{"'1, 2, 3' * true", bytecode.BStr("1, 2, 3")},
	{"false * '1, 2, 3'", bytecode.BStr("")},
	{"1 / true", bytecode.BFloat(1.000000)},
	{"true / 1", bytecode.BFloat(1.000000)},
	// {"true / 1.1", bytecode.BFloat(0.909091)},
	{"1.1 / true", bytecode.BFloat(1.100000)},
	{"1 // true", bytecode.BInt(1)},
	{"true // 1", bytecode.BInt(1)},
	// {"true // 1.1", bytecode.BFloat(0.909091)},
	{"1.1 // true", bytecode.BFloat(1.100000)},
	{"1 % true", bytecode.BInt(0)},
	{"true % 1", bytecode.BInt(0)},
	{"true % 1.1", bytecode.BFloat(1.000000)},
	// {"1.1 % true", bytecode.BFloat(0.100000)},
	{"1 < true", bytecode.BBool(false)},
	{"true < 1", bytecode.BBool(false)},
	{"true < 1.1", bytecode.BBool(true)},
	{"1.1 < false", bytecode.BBool(false)},
	{"1 <= true", bytecode.BBool(true)},
	{"true <= 1", bytecode.BBool(true)},
	{"true <= 1.1", bytecode.BBool(true)},
	{"1.1 <= false", bytecode.BBool(false)},
	{"1 > true", bytecode.BBool(false)},
	{"true > 1", bytecode.BBool(false)},
	{"true > 1.1", bytecode.BBool(false)},
	{"1.1 > false", bytecode.BBool(true)},
	{"1 >= true", bytecode.BBool(true)},
	{"true >= 1", bytecode.BBool(true)},
	{"true >= 1.1", bytecode.BBool(false)},
	{"1.1 >= false", bytecode.BBool(true)},
	{"[1, 2,3][true]", bytecode.BInt(2)},
	{"[1, 2,3][false]", bytecode.BInt(1)},
	{"+true", bytecode.BInt(1)},
	{"+false", bytecode.BInt(0)},
	{"-true", bytecode.BInt(-1)},
	{"-false", bytecode.BInt(0)},
	{"not true", bytecode.BBool(false)},
	{"not false", bytecode.BBool(true)},
	{"a % 3 ? a % 5 ? a : 'buzz' : a % 5 ? fizz : fizzbuzz", bytecode.BInt(43)},
	{"a % 2 ? 3 * a + 1 : a // 2", bytecode.BInt(130)},
	{"foobar(123)", bytecode.BNull{}},
	{"ba(123)", bytecode.BFloat(5338.200000)},
	{"vv(123)", bytecode.BNull{}},
	{"fooObj.bar * 10", bytecode.BInt(100)},
	{"fooObj.baz(30) * 10", bytecode.BFloat(13020.000000)},
	{"fooObj.baz(40)", bytecode.BFloat(1736.000000)},
	// { "[1, 2, 3]", bytecode.BArray([1, 2, 3]) },
	// { "[]", bytecode.BArray([]) },
	// { "[1, 2, 3] + [4, 5, 6]", bytecode.BArray([1, 2, 3, 4, 5, 6]) },
	// { "[1, 2, 3] * 3", bytecode.BArray([1, 2, 3, 1, 2, 3, 1, 2, 3]) },
	{"[1, 2, 3][0]", bytecode.BInt(1)},
	{"fooObj != 1", bytecode.BBool(true)},
	{"10 or unknown1 and unknown2", bytecode.BInt(10)},
	{"10 ? b : unknownvariable", bytecode.BInt(2)},
	{"0 ? unknownvar : b", bytecode.BInt(2)},
}

func TestValidStrings(t *testing.T) {
	var tests = validStrings

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			_, err := vm.EvalString(tt, vmSeed)
			// result, err := vm.EvalString(tt, vmSeed)
			// val := result.Val
			// fmt.Printf("{ \"%s\", %T(%s) },\n", tt, val, val)
			// Convert the above with: , ([^']*)' --> , $1"
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestValidStringsWithOutput(t *testing.T) {
	for _, tt := range validStringsInOut {
		testname := fmt.Sprintf("%s", tt.in)
		t.Run(testname, func(t *testing.T) {
			m := vm.New(vm.Params{})
			result, err := m.EvalString(tt.in, vmSeed)
			if err != nil {
				t.Error(err)
			}
			val := result.Val
			if !runtime.Eq(tt.expected, val) {
				t.Errorf("Expected %s, got %s", tt.expected, val)
			}
		})
	}
}

var invalidStrings = []string{
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
	"[foo(), bar, fooObj.bar][5]",
	"[foo(), *2 * 3, fooObj.bar][5]",
	"[foo(), (2, 3), fooObj.bar][5]",
	"[1//2nasdijio2 * 5, (2, 3), 35]",
	"emptyObj + 2",
}

func TestInvalidStrings(t *testing.T) {
	var tests = invalidStrings

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt)
		t.Run(testname, func(t *testing.T) {
			vm := vm.New(vm.Params{})
			_, err := vm.EvalString(tt, vmSeed)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
		})
	}
}

func TestFizzBuzz(t *testing.T) {
	env := vm.CloneEnv(vmSeed)
	m := vm.New(vm.Params{})
	s := "i % 3 ? i % 5 ? i : 'buzz' : i % 5 ? fizz : 'fizzbuzz'"
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		t.Fatal("Found compile errors")
	}
	for i := 0; i < 100; i++ {
		env["i"] = bytecode.BInt(int64(i))
		result, err := m.Eval(compilation, env)
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
	env := vm.CloneEnv(vmSeed)
	s := "i % 2 == 0 ? i//2 : 3*i + 1"
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		t.Fatal("Found compile errors")
	}
	for i := 1000; i > 1; {
		env["i"] = bytecode.BInt(int64(i))
		result, err := m.Eval(compilation, env)
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
	env := vm.CloneEnv(vmSeed)
	s := "i % 2 ? 3*i + 1 : i//2"
	compilation := compiler.CompileString(s)
	if len(compilation.Errors) > 0 {
		b.Fatal("Found compile errors")
	}
	for numRuns := 0; numRuns < b.N; numRuns++ {
		for i := 100000; i > 1; {
			env["i"] = bytecode.BInt(int64(i))
			result, err := m.Eval(compilation, env)
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

//go:noinline
func collatz(i int) int {
	// no inlining for a fairer comparison
	if i%2 == 0 {
		return i / 2
	}
	return 3*i + 1
}

var vmSeed = vm.VMEnv{
	"a":        bytecode.BInt(43),
	"b":        bytecode.BInt(2),
	"c":        bytecode.BInt(15),
	"d":        bytecode.BArray{bytecode.BInt(1), bytecode.BFloat(2.2)},
	"e":        bytecode.BStr("Echo location for dolphins"),
	"f":        bytecode.BFloat(3.14),
	"foo":      bytecode.BFloat(10.5),
	"s":        bytecode.BStr("I am a string!"),
	"fizz":     bytecode.BStr("fizz"),
	"buzz":     bytecode.BStr("buzz"),
	"fizzbuzz": bytecode.BStr("fizzbuzz"),
	"emptyObj": make(bytecode.BObj),
	"null":     bytecode.BNull{},
	"fooObj": bytecode.BObj(map[string]bytecode.BVal{
		"bar": bytecode.BInt(10),
		"baz": vm.WrapFn("baz", func(x bytecode.BVal) (bytecode.BVal, error) {
			log.Println(x)
			xx, ok := x.(bytecode.BInt)
			if !ok {
				return nil, fmt.Errorf("Was not passed in a integer")
			}
			return bytecode.BFloat(float64(xx) * 43.4), nil
		}),
	}),
	// functions
	"foobar": vm.WrapFn("foobar", func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		return bytecode.BNull{}
	}),
	"z": vm.WrapFn("z", func(x bytecode.BVal, y bytecode.BVal) (bytecode.BVal, error) {
		log.Println(x)
		xx, ok := x.(bytecode.BInt)
		if !ok {
			return nil, fmt.Errorf("Was not passed in a integer")
		}
		yy, ok := y.(bytecode.BFloat)
		if !ok {
			return nil, fmt.Errorf("Was not passed in a float")
		}
		return bytecode.BFloat(float64(xx) * float64(yy) * 200), nil
	}),
	"ba": vm.WrapFn("ba", func(x bytecode.BVal) (bytecode.BVal, error) {
		log.Println(x)
		xx, ok := x.(bytecode.BInt)
		if !ok {
			return nil, fmt.Errorf("Was not passed in a integer")
		}
		return bytecode.BFloat(float64(xx) * 43.4), nil
	}),
	"vv": vm.WrapFn("vv", func(x bytecode.BVal) {
		log.Println(x)
	}),
}

// ---------------------- Fuzzing ---------------------------------

var testingVMParams = vm.Params{
	MaxMemory:       32768,
	MaxInstructions: 10000,
}

// VM should never fail
func FuzzVM(f *testing.F) {
	testcases := append(validStrings, invalidStrings...)
	for _, tc := range testcases {
		f.Add(tc)
	}
	vm := vm.New(testingVMParams)
	f.Fuzz(func(t *testing.T, orig string) {
		vm.EvalString(orig, vmSeed)
	})
}

func FuzzVMRandAST(f *testing.F) {
	maxDepth := 20
	vm := vm.New(testingVMParams)
	// Don't make the test corpus too bug, Go runs it as part of every Go unit
	// test run
	for i := 0; i < 100; i++ {
		for j := 0; j < 10; j++ {
			f.Add(uint32(i), uint(j))
		}
	}
	f.Fuzz(func(t *testing.T, seed uint32, depth uint) {
		if depth >= uint(maxDepth) {
			return
		}
		ast := parser.GenRandomAST(seed, depth)
		s := ast.String()
		vm.EvalString(s, vmSeed)
	})
}

func TestRandAST(t *testing.T) {
	type Pair struct {
		seed  uint32
		depth uint
	}
	vm := vm.New(testingVMParams)
	knownRegressions := []Pair{
		{1, 20}, // at 29, becomes too big.
	}
	for _, regression := range knownRegressions {
		t.Run(fmt.Sprintf("%d", regression), (func(t *testing.T) {
			seed := regression.seed
			depth := regression.depth
			ast := parser.GenRandomAST(seed, depth)
			s := ast.String()
			// fmt.Println(s)
			vm.EvalString(s, vmSeed)
		}))
	}
}
