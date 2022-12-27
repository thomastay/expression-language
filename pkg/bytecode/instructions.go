// The instruction set for the VM
package bytecode

// Instructions based on https://craftinginterpreters.com/a-virtual-machine.html
// As well as https://docs.python.org/3/library/dis.html#python-bytecode-instructions
// Note: these are imported wholesale into the VM, so always add the `Op` prefix in front
type Instruction int32

//go:generate stringer -type Instruction
const (
	OpConst Instruction = iota
	// Load a variable
	OpLoad
	// Your usual run of the mill binary operators
	OpAdd
	OpMinus
	OpMul
	OpDiv
	OpAnd
	OpMod
	OpLt
	OpGt
	OpGe
	OpLe
	OpEq
	OpNe
	OpUnaryNot
	OpUnaryPlus
	OpUnaryMinus
	// Return from a stack frame. Currently there arent any, so this just halts the VM
	OpReturn
	// Call a function
	OpCall
	// A binary operator, loads base.field
	OpLoadAttr
	// Unconditional branch
	OpBr
	// Conditional branch if top of stack is truthy. Also consume top of stack.
	OpBrIf
	// Conditional branch if top of stack is truthy. If so, doesn't consume, else it does.
	OpBrIfOrPop
	// Conditional branch if top of stack is falsey. If so, doesn't consume, else it does.
	OpBrIfFalseOrPop
)
