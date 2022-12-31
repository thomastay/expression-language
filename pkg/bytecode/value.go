package bytecode

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
func (b BBool) isBVal()  {}
func (b BStr) isBVal()   {}
func (b BFunc) isBVal()  {}
func (b BObj) isBVal()   {}
func (b BArray) isBVal() {}

type BNull struct{}
type BFloat float64
type BInt int64
type BBool bool
type BStr string
type BObj map[string]BVal
type BArray []BVal

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
func (b BBool) String() string {
	return strconv.FormatBool(bool(b))
}
func (b BStr) String() string {
	return "'" + string(b) + "'"
}
func (b BFunc) String() string {
	return fmt.Sprintf("Function %s taking in %d args", b.Name, b.NumArgs)
}
func (b BObj) String() string {
	// TODO nicer formatter which does newlines and maybe strings.Join?
	s := "{ "
	var arr []string
	for k, v := range map[string]BVal(b) {
		arr = append(arr, fmt.Sprintf("'%s': %s", k, v))
	}
	s += strings.Join(arr, ", ")
	s += "}"
	return s
}
func (b BArray) String() string {
	s := "["
	var arr []string
	for _, x := range []BVal(b) {
		arr = append(arr, x.String())
	}
	s += strings.Join(arr, ", ")
	s += "]"
	return s
}

// We just go off Python's Truthy:
// https://docs.python.org/3/library/stdtypes.html
// Here are most of the built-in objects considered false:
// - constants defined to be false: None and False.
// - zero of any numeric type: 0, 0.0, 0j, Decimal(0), Fraction(0, 1)
// - empty sequences and collections: '', (), [], {}, set(), range(0)

func (b BNull) IsTruthy() bool {
	return false
}
func (b BFloat) IsTruthy() bool {
	return float64(b) != 0 && !math.IsNaN(float64(b))
}
func (b BInt) IsTruthy() bool {
	return int64(b) != 0
}
func (b BBool) IsTruthy() bool {
	return bool(b)
}
func (b BStr) IsTruthy() bool {
	return len(string(b)) > 0
}
func (b BFunc) IsTruthy() bool {
	return true
}
func (b BObj) IsTruthy() bool {
	return len(b) > 0
}
func (b BArray) IsTruthy() bool {
	return len(b) > 0
}

func (b BNull) Typename() string {
	// NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT
	// NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT
	// NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT NULL IS NOT AN OBJECT
	// How can Javascript get it so wrong
	return "null"
}
func (b BFloat) Typename() string {
	return "float"
}
func (b BInt) Typename() string {
	return "int"
}
func (b BBool) Typename() string {
	return "bool"
}
func (b BStr) Typename() string {
	return "string"
}
func (b BFunc) Typename() string {
	return "function"
}
func (b BObj) Typename() string {
	return "object"
}
func (b BArray) Typename() string {
	return "array"
}

func (b Bytecode) String() string {
	switch b.Inst {
	// No value
	case OpAdd, OpMul, OpDiv, OpAnd, OpMod, OpLt, OpGt, OpGe, OpLe, OpMinus, OpNe, OpEq, OpUnaryMinus, OpUnaryPlus, OpUnaryNot, OpLoadAttr, OpPow, OpFloorDiv:
		return b.Inst.String()
	default:
		if b.Val == nil {
			return fmt.Sprintf("%s %d", b.Inst, b.IntVal)
		}
		return fmt.Sprintf("%s %s", b.Inst, b.Val)
	}
}
