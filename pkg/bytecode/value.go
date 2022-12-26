package bytecode

import (
	"fmt"
	"math"
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
	IsTruthy() bool
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

func (b BFloat) IsTruthy() bool {
	return int64(b) != 0
}
func (b BInt) IsTruthy() bool {
	return float64(b) != 0 && !math.IsNaN(float64(b))
}
func (b BStr) IsTruthy() bool {
	return len(string(b)) > 0
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
	case OpMod:
		fallthrough
	case OpLt:
		fallthrough
	case OpGt:
		fallthrough
	case OpGe:
		fallthrough
	case OpLe:
		fallthrough
	case OpMinus:
		return b.Inst.String()
	default:
		if b.Val == nil {
			return fmt.Sprintf("%s %d", b.Inst, b.IntVal)
		}
		return fmt.Sprintf("%s %s", b.Inst, b.Val)
	}
}
