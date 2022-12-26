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
)
