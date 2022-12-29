// Package vm implements the Virtual Machine that will interpret the Bytecode
// There are a few features that (will) be implemented for the VM
// 1. A VM is something that your program interfacts with, so you can do things like:
// vm.AddVariable("foo", IntType, 32) --> makes the foo variable an int in the vm which can be used by the expression language
// 1. VM state is serializable, so we can save and resume it from disk
// 1. VM state keeps track of the total number of instructions performed, so we can set a maximum # of instructions
package vm

import (
	"errors"
	"fmt"

	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
)

var defaultMaxInstructions = 1000

func New(params Params) VMState {
	if params.MaxInstructions == 0 {
		params.MaxInstructions = defaultMaxInstructions
	}
	vars := make(map[string]BVal)
	return VMState{params: params, variables: vars}
}

type VMState struct {
	// stack     Stack // we don't need a stack for now
	variables map[string]BVal
	params    Params
}

func (vm *VMState) AddInt(key string, val int64) {
	vm.variables[key] = BInt(val)
}
func (vm *VMState) AddFloat(key string, val float64) {
	vm.variables[key] = BFloat(val)
}
func (vm *VMState) AddStr(key string, val string) {
	vm.variables[key] = BStr(val)
}
func (vm *VMState) AddObject(key string, obj BObj) {
	if obj == nil {
		// No nil objects in the VM
		obj = make(map[string]BVal)
	}
	vm.variables[key] = obj
}

// Adds a function to the VM. Users must specify the number of arguments the function takes.
func (vm *VMState) AddFunc(name string, fnWithArgs VMFuncWithArgs) {
	vm.variables[name] = BFunc{
		Fn:      fnWithArgs.Fn,
		NumArgs: fnWithArgs.NumArgs,
		Name:    name,
	}
}

type VMFuncWithArgs struct {
	Fn      VMFunc
	NumArgs int
}

// Convenience method if you just want to evaluate a string. Concatenates all compile errors into one
func (vm *VMState) EvalString(s string) (Result, error) {
	expr, err := parser.ParseString(s)
	if err != nil {
		return Result{}, err
	}
	comp := compiler.Compile(expr)
	if len(comp.Errors) > 0 {
		var errString string
		for _, c := range comp.Errors {
			errString += c.Error()
		}
		return Result{}, errors.New(errString)
	}
	return vm.Eval(comp)
}

func (vm *VMState) Eval(compilation compiler.Compilation) (Result, error) {
	executedInsts := 0
	pc := 0
	// stack := vm.stack
	stack := make(Stack, 0, 64) // preallocate some space for items
	codes := compilation.Bytecode
InstLoop:
	for pc < len(codes) && executedInsts < vm.params.MaxInstructions {
		executedInsts++
		code := codes[pc]
		switch code.Inst {
		case OpReturn:
			break InstLoop
		case OpConst:
			stack.push(code.Val)
		case OpLoad:
			identName := code.Val.(BStr)
			val, ok := vm.variables[string(identName)]
			if !ok {
				return Result{}, fmt.Errorf("NameError: name '%s' is not defined", identName)
			}
			stack = append(stack, val)
		// ----------------Binary Operations------------------
		case OpAdd:
			b := stack.pop()
			a := stack.pop()
			result, err := add(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMinus:
			b := stack.pop()
			a := stack.pop()
			result, err := sub(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMul:
			b := stack.pop()
			a := stack.pop()
			result, err := mul(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpDiv:
			b := stack.pop()
			a := stack.pop()
			result, err := div(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpFloorDiv:
			b := stack.pop()
			a := stack.pop()
			result, err := floorDiv(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMod:
			b := stack.pop()
			a := stack.pop()
			result, err := modulo(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpLt:
			b := stack.pop()
			a := stack.pop()
			ord, err := cmp(a, b, "<")
			if err != nil {
				return Result{}, err
			}
			result := 0
			if ord == -1 {
				result = 1
			}
			stack.push(BInt(result))
		case OpGt:
			b := stack.pop()
			a := stack.pop()
			ord, err := cmp(a, b, ">")
			if err != nil {
				return Result{}, err
			}
			result := 0
			if ord == 1 {
				result = 1
			}
			stack.push(BInt(result))
		case OpLe:
			b := stack.pop()
			a := stack.pop()
			ord, err := cmp(a, b, "<=")
			if err != nil {
				return Result{}, err
			}
			result := 0
			if ord != 1 {
				result = 1
			}
			stack.push(BInt(result))
		case OpGe:
			b := stack.pop()
			a := stack.pop()
			ord, err := cmp(a, b, ">=")
			if err != nil {
				return Result{}, err
			}
			result := 0
			if ord != -1 {
				result = 1
			}
			stack.push(BInt(result))
		case OpEq:
			b := stack.pop()
			a := stack.pop()
			equal := eq(a, b)
			result := 0
			if equal {
				result = 1
			}
			stack.push(BInt(result))
		case OpNe:
			b := stack.pop()
			a := stack.pop()
			equal := eq(a, b)
			result := 0
			if !equal {
				result = 1
			}
			stack.push(BInt(result))
		case OpLoadAttr:
			// Base is loaded before field, so field pops first
			field := stack.pop()
			base := stack.pop()
			fieldStr := field.(BStr)
			baseObj, ok := base.(BObj)
			if !ok {
				return Result{}, fmt.Errorf("AttributeError: %s object has no attribute %s", base, field)
			}
			val, ok := baseObj[string(fieldStr)]
			if !ok {
				return Result{}, fmt.Errorf("AttributeError: %s object has no attribute %s", base, field)
			}
			stack.push(val)
		// ----------------Unary Operations------------------
		case OpUnaryPlus:
			a := stack.peek() // don't pop!
			switch a.(type) {
			case BInt, BFloat:
				// do nothing
			default:
				return Result{}, fmt.Errorf("TypeError: bad operand type for unary +: %s", a.Typename())
			}
		case OpUnaryMinus:
			switch a := stack.pop().(type) {
			case BInt:
				stack.push(BInt(-int64(a)))
			case BFloat:
				stack.push(BInt(-int64(a)))
			default:
				return Result{}, fmt.Errorf("TypeError: bad operand type for unary -: %s", a.Typename())
			}
		case OpUnaryNot:
			a := stack.pop()
			neg := !a.IsTruthy()
			result := 0
			if neg {
				result = 1
			}
			stack.push(BInt(result))
		// ----------------Conditional Operations------------------
		case OpBr:
			pc = int(code.IntVal)
			continue InstLoop
		case OpBrIf:
			a := stack.pop()
			if a.IsTruthy() {
				pc = int(code.IntVal)
				continue InstLoop
			} // else fallthrough
		case OpBrIfOrPop:
			a := stack.peek()
			if a.IsTruthy() {
				pc = int(code.IntVal)
				continue InstLoop
			} else {
				stack.pop()
			}
		case OpBrIfFalseOrPop:
			a := stack.peek()
			if !a.IsTruthy() {
				pc = int(code.IntVal)
				continue InstLoop
			} else {
				stack.pop()
			}
		// ----------------Function Operations------------------
		case OpCall:
			name := stack.pop()
			bFn, ok := name.(BFunc)
			if !ok {
				return Result{}, fmt.Errorf("Stack value %s is not a function", name)
			}
			// load params
			numParams := code.IntVal
			if numParams != bFn.NumArgs {
				return Result{}, fmt.Errorf("Fn %s passed wrong number of args, expected %d, got %d", bFn.Name, bFn.NumArgs, numParams)
			}
			params := make([]BVal, numParams)
			for i := 0; i < numParams; i++ {
				params[i] = stack.pop()
			}
			result, err := bFn.Fn(params)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		default:
			return Result{}, errors.New("not impl")
		}
		pc++
	}
	val := stack.pop()
	// vm.stack = stack
	return Result{Val: val}, nil
}

type Stack []BVal

// Make sure inlined
func (stack *Stack) pop() (result BVal) {
	n := len(*stack)
	result = (*stack)[n-1]
	*stack = (*stack)[:n-1]
	return result
}

// Make sure inlined
func (stack *Stack) peek() (result BVal) {
	n := len(*stack)
	result = (*stack)[n-1]
	return result
}

// Make sure inlined
func (stack *Stack) push(x BVal) {
	*stack = append(*stack, x)
}

type Result struct {
	Val BVal
}

// Configuring the VM
type Params struct {
	MaxInstructions int
	Debug           bool
}
