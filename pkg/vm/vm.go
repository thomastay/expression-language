// Package vm implements the Virtual Machine that will interpret the Bytecode
// There are a few features that (will) be implemented for the VM
// 1. A VM is something that your program interfacts with, so you can do things like:
// vm.AddVariable("foo", IntType, 32) --> makes the foo variable an int in the vm which can be used by the expression language
// 1. VM state is serializable, so we can save and resume it from disk
// 1. VM state keeps track of the total number of instructions performed, so we can set a maximum # of instructions
package vm

import (
	"errors"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
	"github.com/thomastay/expression_language/pkg/compiler"
	"github.com/thomastay/expression_language/pkg/parser"
)

var defaultMaxInstructions = 1000

func New(params Params) VMState {
	if params.MaxInstructions == 0 {
		params.MaxInstructions = defaultMaxInstructions
	}
	return VMState{params: params}
}

type VMState struct {
	stack     Stack
	variables map[string]int64 // TODO don't use int
	params    Params
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
	stack := vm.stack
	codes := compilation.Bytecode
InstLoop:
	for pc < len(codes) && executedInsts < vm.params.MaxInstructions {
		executedInsts++
		code := codes[pc]
		switch code.Inst {
		case OpReturn:
			break InstLoop
		case OpConst:
			stack = append(stack, code.Val)
		// ----------------Binary Operations------------------
		case OpAdd:
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Add64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpMinus:
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Sub64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpMul:
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Mul64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpDiv:
			b := stack.pop()
			a := stack.pop()
			if b == 0 {
				return Result{}, errors.New("Divide by zero")
			}
			result, ok := overflow.Div64(a, b) // TODO cast to float
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		// ----------------Conditional Operations------------------
		case OpBr:
			pc = int(code.Val)
			continue InstLoop
		case OpBrIf:
			a := stack.pop()
			if a != 0 {
				pc = int(code.Val)
				continue InstLoop
			} // else fallthrough
		case OpBrIfOrPop:
			a := stack.peek()
			if a != 0 {
				pc = int(code.Val)
				continue InstLoop
			} else {
				stack.pop()
			}
		case OpBrIfFalseOrPop:
			a := stack.peek()
			if a == 0 {
				pc = int(code.Val)
				continue InstLoop
			} else {
				stack.pop()
			}
		default:
			return Result{}, errors.New("not impl")
		}
		pc++
	}
	val := stack.pop()
	vm.stack = stack
	return Result{Val: val}, nil
}

type Stack []int64

// Make sure inlined
func (stack *Stack) pop() (result int64) {
	n := len(*stack)
	result = (*stack)[n-1]
	*stack = (*stack)[:n-1]
	return result
}

// Make sure inlined
func (stack *Stack) peek() (result int64) {
	n := len(*stack)
	result = (*stack)[n-1]
	return result
}

// Make sure inlined
func (stack *Stack) push(x int64) {
	*stack = append(*stack, x)
}

type Result struct {
	Val int64 // TODO, what should an evaluator return?
}

// Configuring the VM
type Params struct {
	MaxInstructions int
	Debug           bool
}
