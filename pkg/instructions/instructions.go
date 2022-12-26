// The instruction set for the VM
package instructions

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
	OpReturn
	OpPop
	// Unconditional branch
	OpBr
	// Conditional branch if top of stack is nonzero. Also consume top of stack.
	OpBrIf
)
