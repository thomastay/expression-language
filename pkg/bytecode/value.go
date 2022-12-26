package bytecode

import (
	"fmt"
	"strconv"
)

type Bytecode struct {
	Inst   Instruction
	IntVal int // Fast path jump indices since it's used much more than actual constants
	Val    BVal
}
type BVal interface {
	fmt.Stringer
	isBVal()
}

func (b BFloat) isBVal() {}
func (b BInt) isBVal()   {}
func (b BStr) isBVal()   {}

type BFloat float64
type BInt int64
type BStr string

func (b BFloat) String() string {
	return fmt.Sprintf("%f", float64(b))
}
func (b BInt) String() string {
	return strconv.FormatInt(int64(b), 10)
}
func (b BStr) String() string {
	return string(b)
}

func (b Bytecode) String() string {
	switch b.Inst {
	// No value
	case OpAdd:
		fallthrough
	case OpMul:
		fallthrough
	case OpDiv:
		fallthrough
	case OpAnd:
		fallthrough
	case OpMinus:
		return fmt.Sprintf("%s", b.Inst)
	default:
		return fmt.Sprintf("%s %s", b.Inst, b.Val)
	}
}
