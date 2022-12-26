// The instruction set for the VM
package bytecode

// Instructions based on https://craftinginterpreters.com/a-virtual-machine.html
// Note: these are imported wholesale into the VM, so always add the `Op` prefix in front
type Instruction int32

//go:generate stringer -type Instruction
const (
	OpConst Instruction = iota
	OpAdd
	OpMinus
	OpMul
	OpDiv
	OpAnd
	OpReturn
	// Unconditional branch
	OpBr
	// Conditional branch if top of stack is truthy. Also consume top of stack.
	OpBrIf
	// Conditional branch if top of stack is truthy. If so, doesn't consume, else it does.
	OpBrIfOrPop
	// Conditional branch if top of stack is falsey. If so, doesn't consume, else it does.
	OpBrIfFalseOrPop
)
