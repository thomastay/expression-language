// Package vm implements the Virtual Machine that will interpret the Bytecode
// There are a few features that (will) be implemented for the VM
// 1. A VM is something that your program interfacts with, so you can do things like:
// vm.AddVariable("foo", IntType, 32) --> makes the foo variable an int in the vm which can be used by the expression language
// 1. VM state is serializable, so we can save and resume it from disk
// 1. VM state keeps track of the total number of instructions performed, so we can set a maximum # of instructions
package vm

import (
	"errors"

	"github.com/thomastay/expression_language/pkg/compiler"
	. "github.com/thomastay/expression_language/pkg/instructions"
)

func New(params Params) VMState {
	return VMState{params: params}
}

type VMState struct {
	stack     Stack
	variables map[string]int64 // TODO don't use int
	params    Params
}

func (vm *VMState) Eval(compilation compiler.Compilation) Result {
	pc := 0
	stack := vm.stack
	codes := compilation.Bytecode
InstLoop:
	for pc < len(codes) {
		code := codes[pc]
		switch code.Inst {
		case OpReturn:
			break InstLoop
		case OpConstant:
			stack = append(stack, code.Val)
		case OpAdd:
			n := len(stack)
			if n < 2 {
				return Result{Err: errors.New("Not enough values on stack")}
			}
			b := stack.pop()
			a := stack.pop()
			stack.push(a + b)
		case OpMinus:
			n := len(stack)
			if n < 2 {
				return Result{Err: errors.New("Not enough values on stack")}
			}
			b := stack.pop()
			a := stack.pop()
			stack.push(a - b)
		case OpMul:
			n := len(stack)
			if n < 2 {
				return Result{Err: errors.New("Not enough values on stack")}
			}
			b := stack.pop()
			a := stack.pop()
			stack.push(a * b)
		case OpDiv:
			n := len(stack)
			if n < 2 {
				return Result{Err: errors.New("Not enough values on stack")}
			}
			b := stack.pop()
			a := stack.pop()
			if b == 0 {
				return Result{Err: errors.New("Divide by Zero")}
			}
			stack.push(a / b) // TODO implement casting to float
		default:
			return Result{Err: errors.New("Not implemented")}
		}
		pc++
	}
	val := stack.pop()
	vm.stack = stack
	return Result{Val: val}
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
func (stack *Stack) push(x int64) {
	*stack = append(*stack, x)
}

type Result struct {
	Val int64 // TODO, what should an evaluator return?
	Err error
}

// Configuring the VM
type Params struct {
	MaxInstructions int
	Debug           bool
}
