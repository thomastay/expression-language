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
	m.AddInt("a", 43)
	m.AddInt("b", 2)
	m.AddInt("c", 15)
	m.AddFloat("foo", 10.5)
	m.AddStr("bar", "I am a string")
	m.AddStr("fizzbuzz", "fizzbuzz")
	m.AddStr("fizz", "fizz")
	m.AddStr("buzz", "buzz")
	m.AddFunc("foobar", vm.Wrap1(func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		return bytecode.BNull{}
	}))
	// TODO make wrap1 return a BFunc
	bazFn := vm.Wrap1(func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		xx := x.(bytecode.BInt)
		return bytecode.BFloat(float64(xx) * 43.4)
	})
	fooObjVal := map[string]bytecode.BVal{
		"bar": bytecode.BInt(10),
		"baz": bytecode.BFunc{
			Fn:      bazFn.Fn,
			NumArgs: bazFn.NumArgs,
			Name:    "baz",
		},
	}
	m.AddObject("fooObj", fooObjVal)
	m.AddObject("emptyObj", nil)
	result, err := m.Eval(comp, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("VM result:", result.Val)
}
