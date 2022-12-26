// Code generated by "stringer -type Instruction"; DO NOT EDIT.

package instructions

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpConst-0]
	_ = x[OpAdd-1]
	_ = x[OpMinus-2]
	_ = x[OpMul-3]
	_ = x[OpDiv-4]
	_ = x[OpReturn-5]
	_ = x[OpPop-6]
	_ = x[OpBr-7]
	_ = x[OpBrIf-8]
}

const _Instruction_name = "OpConstOpAddOpMinusOpMulOpDivOpReturnOpPopOpBrOpBrIf"

var _Instruction_index = [...]uint8{0, 7, 12, 19, 24, 29, 37, 42, 46, 52}

func (i Instruction) String() string {
	if i < 0 || i >= Instruction(len(_Instruction_index)-1) {
		return "Instruction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Instruction_name[_Instruction_index[i]:_Instruction_index[i+1]]
}
