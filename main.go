package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
	"github.com/thomastay/expression_language/pkg/vm"

	"github.com/k0kubun/pp/v3"
)

var blurb = `Expression Lang repl v0.1.0. Press Ctrl-C to quit. Made by Thomas Tay.
Type .env to see all variables`

var varDeclRegex *regexp.Regexp

func init() {
	varDeclRegex, _ = regexp.Compile(`var ([a-zA-Z]\w*) (.*)`)
}

func main() {
	shouldSeed := flag.Bool("seed", true, "Seed the VM")
	flag.Parse()
	m := vm.New(vm.Params{Debug: true})
	var env vm.VMEnv
	if *shouldSeed {
		env = seedEnv
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(blurb)
	fmt.Print(">>> ")
	for scanner.Scan() {
		text := scanner.Text()
		runOnLine(text, m, env)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func runOnLine(text string, m vm.VMState, env vm.VMEnv) {
	defer func() { fmt.Print(">>> ") }()
	text = strings.Trim(text, " \t\r\n")
	if text == "" {
		return
	}
	if text == ".env" {
		pp.Println(env)
		return
	}
	if varDeclRegex.MatchString(text) {
		submatches := varDeclRegex.FindStringSubmatch(text)
		varName := submatches[1]
		text := submatches[2]
		val, err := runOnString(text, m, env)
		if err == nil {
			env[varName] = val
			fmt.Printf("Setting %s := %s \n", varName, val)
		}
	} else {
		runOnString(text, m, env)
	}
}

func runOnString(s string, m vm.VMState, env vm.VMEnv) (bytecode.BVal, error) {
	expr, err := parser.ParseString(s)
	if err != nil {
		log.Println("Parse Error:", err)
		return nil, fmt.Errorf("Parse Error")
	}
	fmt.Println("Expression:", expr.String())
	comp := compiler.Compile(expr)
	if len(comp.Errors) > 0 {
		for _, compErr := range comp.Errors {
			log.Println(compErr)
		}
		return nil, fmt.Errorf("Compile Error")
	}
	fmt.Println("Bytecode:")
	for i, b := range comp.Bytecode.Insts {
		bte := bytecode.Bytecode{
			Inst:   b,
			IntVal: comp.Bytecode.IntData[i],
		}
		fmt.Println("  ", i, bte.String())
	}
	result, err := m.Eval(comp, env)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Runtime Error")
	}
	fmt.Println("VM result:", result.Val)
	return result.Val, nil
}

//export runOnString
func wasmRunOnString(s string) (vm.Result, error) {
	// For wasm only! Don't use internally
	m := vm.New(vm.Params{})
	return m.EvalString(s, seedEnv)
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
	"emptyObj": make(bytecode.BObj),
	"fooObj":   bytecode.BObj(fooObjVal),
	// functions
	"foobar": vm.WrapFn("foobar", func(x bytecode.BVal) bytecode.BVal {
		log.Println(x)
		return bytecode.BNull{}
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
