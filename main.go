package main

import (
	"bufio"
	"flag"
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
	shouldSeed := flag.Bool("seed", true, "Seed the VM")
	flag.Parse()
	m := vm.New(vm.Params{})
	var env vm.VMEnv
	if *shouldSeed {
		env = seedEnv
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			continue
		}
		if strings.HasPrefix(text, "var") {
			// var decl
			panic("vardecl")
		} else {
			runOnString(text, m, env)
		}
	}
}

func runOnString(s string, m vm.VMState, env vm.VMEnv) {
	expr, err := parser.ParseString(s)
	if err != nil {
		log.Println("Parse Error:", err)
		return
	}
	fmt.Println("Expression:", expr.String())
	comp := compiler.Compile(expr)
	if len(comp.Errors) > 0 {
		for _, compErr := range comp.Errors {
			log.Println(compErr)
		}
		return
	}
	fmt.Println("Bytecode:")
	for i, b := range comp.Bytecode {
		fmt.Println("  ", i, b.String())
	}
	result, err := m.Eval(comp, env)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("VM result:", result.Val)
}

// seed the VM with some useful variables
var fooObjVal = map[string]bytecode.BVal{
	"bar": bytecode.BInt(10),
	"baz": vm.WrapFn("baz", func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		xx := x.(bytecode.BInt)
		return bytecode.BFloat(float64(xx) * 43.4)
	}),
}
var seedEnv = vm.VMEnv{
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
