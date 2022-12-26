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
	Typename() string
}

func (b BNull) isBVal()  {}
func (b BFloat) isBVal() {}
func (b BInt) isBVal()   {}
func (b BStr) isBVal()   {}
func (b BFunc) isBVal()  {}

type BNull struct{}
type BFloat float64
type BInt int64
type BStr string

type VMFunc func(args []BVal) (BVal, error)
type BFunc struct {
	Fn      VMFunc
	NumArgs int
	Name    string // for debugging
}

func (b BNull) String() string {
	return "null"
}
func (b BFloat) String() string {
	return fmt.Sprintf("%f", float64(b))
}
func (b BInt) String() string {
	return strconv.FormatInt(int64(b), 10)
}
func (b BStr) String() string {
	return string(b)
}
func (b BFunc) String() string {
	return fmt.Sprintf("Function %s taking in %d args", b.Name, b.NumArgs)
}

func (b BNull) IsTruthy() bool {
	return false
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

func (b BFunc) IsTruthy() bool {
	return true
}

func (b BNull) Typename() string {
	return "null"
}
func (b BFloat) Typename() string {
	return "float"
}
func (b BInt) Typename() string {
	return "int"
}
func (b BStr) Typename() string {
	return "string"
}
func (b BFunc) Typename() string {
	return "function"
}

func (b Bytecode) String() string {
	switch b.Inst {
	// No value
	case OpAdd, OpMul, OpDiv, OpAnd, OpMod, OpLt, OpGt, OpGe, OpLe, OpMinus, OpNe, OpEq:
		return b.Inst.String()
	default:
		if b.Val == nil {
			return fmt.Sprintf("%s %d", b.Inst, b.IntVal)
		}
		return fmt.Sprintf("%s %s", b.Inst, b.Val)
	}
}
