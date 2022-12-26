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

func (vm *VMState) Eval(compilation compiler.Compilation) (Result, error) {
	pc := 0
	stack := vm.stack
	codes := compilation.Bytecode
InstLoop:
	for pc < len(codes) {
		code := codes[pc]
		switch code.Inst {
		case OpReturn:
			break InstLoop
		case OpConst:
			stack = append(stack, code.Val)
		case OpAdd:
			n := len(stack)
			if n < 2 {
				return Result{}, errNotEnoughStackValues
			}
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Add64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpMinus:
			n := len(stack)
			if n < 2 {
				return Result{}, errNotEnoughStackValues
			}
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Sub64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpMul:
			n := len(stack)
			if n < 2 {
				return Result{}, errNotEnoughStackValues
			}
			b := stack.pop()
			a := stack.pop()
			result, ok := overflow.Mul64(a, b)
			if !ok {
				return Result{}, errOverflow
			}
			stack.push(result)
		case OpDiv:
			n := len(stack)
			if n < 2 {
				return Result{}, errNotEnoughStackValues
			}
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
