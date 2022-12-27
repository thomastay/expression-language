// The instruction set for the VM
package bytecode

// Instructions based on https://craftinginterpreters.com/a-virtual-machine.html
// As well as https://docs.python.org/3/library/dis.html#python-bytecode-instructions
// Note: these are imported wholesale into the VM, so always add the `Op` prefix in front
type Instruction int32

//go:generate stringer -type Instruction
const (
	OpConst Instruction = iota
	OpLoad
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
	OpReturn
	OpCall
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
