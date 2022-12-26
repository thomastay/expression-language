package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	vm := vm.New(vm.Params{})
	result, err := vm.Eval(comp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("VM result:", result.Val)
}
