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
	"github.com/thomastay/expression_language/pkg/runtime"
)

var defaultMaxInstructions = 1000
var defaultMaxMemory = 1048576 // 2 ** 20

func New(params Params) VMState {
	if params.MaxInstructions == 0 {
		params.MaxInstructions = defaultMaxInstructions
	}
	if params.MaxMemory == 0 {
		params.MaxMemory = defaultMaxMemory
	}
	return VMState{params: params}
}

type VMEnv map[string]BVal

type VMState struct {
	params Params
}

// Convenience method if you just want to evaluate a string. Concatenates all compile errors into one
func (vm *VMState) EvalString(s string, env VMEnv) (Result, error) {
	expr, err := parser.ParseString(s)
	if err != nil {
		return Result{}, err
	}
	comp := compiler.Compile(expr, compiler.Params{})
	if len(comp.Errors) > 0 {
		var errString string
		for _, c := range comp.Errors {
			errString += c.Error()
		}
		return Result{}, errors.New(errString)
	}
	return vm.Eval(comp, env)
}

func (vm *VMState) Eval(compilation compiler.Compilation, env VMEnv) (Result, error) {
	executedInsts := 0
	memoryUsed := 0
	pc := 0
	variables := env
	if variables == nil {
		variables = make(VMEnv)
	}
	stack := make(Stack, 0, 64) // preallocate some space for items
	codes := compilation.Bytecode
	if vm.params.Debug {
		fmt.Println("Constant table:")
		fmt.Println("  ", compilation.Constants)
	}
InstLoop:
	for pc < codes.Len() && executedInsts < vm.params.MaxInstructions {
		executedInsts++
		inst := codes.Insts[pc]

		switch inst {
		case OpReturn:
			break InstLoop
		case OpConst:
			pos := codes.IntData[pc]
			// lookup from constant table
			stack.push(compilation.Constants[pos])
		case OpLoad:
			pos := codes.IntData[pc]
			identName := compilation.Constants[pos].(BStr)
			val, ok := variables[string(identName)]
			if !ok {
				return Result{}, fmt.Errorf("NameError: name %s is not defined", identName)
			}
			stack = append(stack, val)
		// ----------------Binary Operations------------------
		case OpAdd:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Add(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpAddImm:
			// Add immediately to a without popping from stack (fuse two inst into one)
			a := stack.pop()
			pos := codes.IntData[pc]
			// lookup from constant table
			b := compilation.Constants[pos]
			result, err := runtime.Add(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMinus:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Sub(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMinusImm:
			a := stack.pop()
			pos := codes.IntData[pc]
			// lookup from constant table
			b := compilation.Constants[pos]
			result, err := runtime.Sub(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMul:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Mul(a, b, vm.params.MaxMemory-memoryUsed)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMulImm:
			a := stack.pop()
			pos := codes.IntData[pc]
			// lookup from constant table
			b := compilation.Constants[pos]
			result, err := runtime.Mul(a, b, vm.params.MaxMemory-memoryUsed)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpDiv:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Div(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpDivImm:
			a := stack.pop()
			pos := codes.IntData[pc]
			// lookup from constant table
			b := compilation.Constants[pos]
			result, err := runtime.Div(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpFloorDiv:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.FloorDiv(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpFloorDivImm:
			a := stack.pop()
			pos := codes.IntData[pc]
			// lookup from constant table
			b := compilation.Constants[pos]
			result, err := runtime.FloorDiv(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpMod:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Modulo(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpPow:
			b := stack.pop()
			a := stack.pop()
			result, err := runtime.Pow(a, b)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		// ----------------Compare Operations------------------
		case OpLt:
			b := stack.pop()
			a := stack.pop()
			ord, err := runtime.Cmp(a, b, "<")
			if err != nil {
				return Result{}, err
			}
			result := false
			if ord == -1 {
				result = true
			}
			stack.push(BBool(result))
		case OpGt:
			b := stack.pop()
			a := stack.pop()
			ord, err := runtime.Cmp(a, b, ">")
			if err != nil {
				return Result{}, err
			}
			result := false
			if ord == 1 {
				result = true
			}
			stack.push(BBool(result))
		case OpLe:
			b := stack.pop()
			a := stack.pop()
			ord, err := runtime.Cmp(a, b, "<=")
			if err != nil {
				return Result{}, err
			}
			result := false
			if ord != 1 {
				result = true
			}
			stack.push(BBool(result))
		case OpGe:
			b := stack.pop()
			a := stack.pop()
			ord, err := runtime.Cmp(a, b, ">=")
			if err != nil {
				return Result{}, err
			}
			result := false
			if ord != -1 {
				result = true
			}
			stack.push(BBool(result))
		case OpEq:
			b := stack.pop()
			a := stack.pop()
			equal := runtime.Eq(a, b)
			result := false
			if equal {
				result = true
			}
			stack.push(BBool(result))
		case OpNe:
			b := stack.pop()
			a := stack.pop()
			equal := runtime.Eq(a, b)
			result := false
			if !equal {
				result = true
			}
			stack.push(BBool(result))
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
			case BBool:
				stack.pop()
				a = runtime.CastBoolToInt(a)
				stack.push(a)
			default:
				return Result{}, fmt.Errorf("TypeError: bad operand type for unary +: %s", a.Typename())
			}
		case OpUnaryMinus:
			a := stack.pop()
			result, err := runtime.Negate(a)
			if err != nil {
				return Result{}, err
			}
			stack.push(result)
		case OpUnaryNot:
			a := stack.pop()
			neg := !a.IsTruthy()
			result := false
			if neg {
				result = true
			}
			stack.push(BBool(result))
		// ----------------Conditional Operations------------------
		case OpBr:
			pc = codes.IntData[pc]
			continue InstLoop
		case OpBrIf:
			a := stack.pop()
			if a.IsTruthy() {
				pc = codes.IntData[pc]
				continue InstLoop
			} // else fallthrough
		case OpBrIfOrPop:
			a := stack.peek()
			if a.IsTruthy() {
				pc = codes.IntData[pc]
				continue InstLoop
			} else {
				stack.pop()
			}
		case OpBrIfFalseOrPop:
			a := stack.peek()
			if !a.IsTruthy() {
				pc = int(codes.IntData[pc])
				continue InstLoop
			} else {
				stack.pop()
			}
		// ----------------Function Operations------------------
		case OpCall:
			name := stack.pop()
			bFn, ok := name.(BFunc)
			if !ok {
				return Result{}, fmt.Errorf("InterpError: Stack value %s is not a function", name)
			}
			// load params
			numParams := codes.IntData[pc]
			if numParams != bFn.NumArgs {
				return Result{}, fmt.Errorf("RuntimeError: function %s passed wrong number of args, expected %d, got %d", bFn.Name, bFn.NumArgs, numParams)
			}
			params := make([]BVal, numParams)
			for i := 0; i < numParams; i++ {
				params[i] = stack.pop()
			}
			result, err := bFn.Fn(params)
			if err != nil {
				return Result{}, fmt.Errorf("RuntimeError: %w", err)
			}
			stack.push(result)
		// ----------------Array Operations------------------
		case OpNewArray:
			n := codes.IntData[pc]
			memoryUsed += n
			if memoryUsed > vm.params.MaxMemory {
				return Result{}, runtime.ErrOOM
			}
			vals := make([]BVal, n)
			for i := 0; i < n; i++ {
				vals[i] = stack.pop()
			}
			stack.push(BArray(vals))
		case OpLoadSubscript:
			b := stack.pop()
			a := stack.pop()
			arr, ok := a.(BArray)
			if !ok {
				return Result{}, fmt.Errorf("TypeError: %s object is not subscriptable", a.Typename())
			}
			b = runtime.CastBoolToInt(b)
			idx, ok := b.(BInt)
			if !ok {
				return Result{}, fmt.Errorf("List index must be an integer, found %s", b.Typename())
			}
			if idx < 0 || int(idx) >= len(arr) {
				return Result{}, fmt.Errorf("Array index %d out of bounds (len %d)", idx, len(arr))
			}
			stack.push(arr[idx])
		default:
			return Result{}, errors.New("Opcode not impl")
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
	// The maximum number of values that can be created in the VM
	MaxMemory int
	Debug     bool
}
