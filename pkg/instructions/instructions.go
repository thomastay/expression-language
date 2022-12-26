// The instruction set for the VM
package instructions

type Instruction int32

// Instructions based on https://craftinginterpreters.com/a-virtual-machine.html
// Note: these are imported wholesale into the VM, so always add the `Op` prefix in front
const (
	OpReturn Instruction = iota
	OpConstant
	OpAdd
	OpMinus
)
