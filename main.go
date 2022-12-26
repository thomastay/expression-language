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
	var bytecodeStrs []string
	for _, b := range comp.Bytecode {
		bytecodeStrs = append(bytecodeStrs, b.String())
	}
	fmt.Println("Bytecode:", strings.Join(bytecodeStrs, ", "))
	vm := vm.New(vm.Params{})
	result := vm.Eval(comp)
	if result.Err != nil {
		log.Fatal(err)
	}
	fmt.Println("VM result:", result.Val)
}
