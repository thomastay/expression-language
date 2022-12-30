package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
	"github.com/thomastay/expression_language/pkg/vm"
)

func main() {
	expr, err := parser.ParseString(strings.Join(os.Args[1:], " "))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Expression:", expr.String())
	comp := compiler.Compile(expr)
	if len(comp.Errors) > 0 {
		for _, compErr := range comp.Errors {
			log.Println(compErr)
		}
		os.Exit(1)
	}
	fmt.Println("Bytecode:")
	for i, b := range comp.Bytecode {
		fmt.Println("  ", i, b.String())
	}
	m := vm.New(vm.Params{})
	// seed the VM with some useful variables
	fooObjVal := map[string]bytecode.BVal{
		"bar": bytecode.BInt(10),
		"baz": vm.WrapFn("baz", func(x bytecode.BVal) bytecode.BVal {
			log.Println(x)
			xx := x.(bytecode.BInt)
			return bytecode.BFloat(float64(xx) * 43.4)
		}),
	}
	env := vm.VMEnv{
		"a":        bytecode.BInt(43),
		"b":        bytecode.BInt(2),
		"c":        bytecode.BInt(15),
		"foo":      bytecode.BFloat(10.5),
		"s":        bytecode.BStr("I am a string!"),
		"fizz":     bytecode.BStr("fizz"),
		"buzz":     bytecode.BStr("buzz"),
		"fizzbuzz": bytecode.BStr("fizzbuzz"),
		"emptyObj": nil,
		"fooObj":   bytecode.BObj(fooObjVal),
		// functions
		"foobar": vm.WrapFn("foobar", func(x bytecode.BVal) bytecode.BVal {
			log.Println(x)
			return bytecode.BNull{}
		}),
		"ba": vm.WrapFn("ba", func(x bytecode.BVal) (bytecode.BVal, error) {
			log.Println(x)
			xx := x.(bytecode.BInt)
			return bytecode.BFloat(float64(xx) * 43.4), nil
		}),
		"vv": vm.WrapFn("vv", func(x bytecode.BVal) {
			log.Println(x)
		}),
	}
	result, err := m.Eval(comp, env)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("VM result:", result.Val)
}
