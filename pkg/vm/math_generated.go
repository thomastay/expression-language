// Code generated by vm/generate/main.go. DO NOT EDIT.
package vm

import (
	"fmt"
	"math"
	"strings"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

func add(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			result, ok := overflow.Add64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			result := float64(a) + float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("+", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			result := float64(a) + float64(b)
			return BFloat(result), nil
		case BFloat:
			result := float64(a) + float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("+", aVal, bVal)
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFloat:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BStr:
			result := BStr(a) + BStr(b)
			return BStr(result), nil
		case BObj:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("+", aVal, bVal)
		}
	case BObj:
		return nil, errTypeMismatch("+", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("+", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("+", aVal, bVal)
	case BArray:
		switch b := bVal.(type) {
		case BInt:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFloat:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BStr:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BArray:
			result := append([]BVal(a), []BVal(b)...)
			return BArray(result), nil
		}

	}
	panic(fmt.Sprintf("Unhandled operation between %s and %s: %s + %s", aVal.Typename(), bVal.Typename(), aVal, bVal))
}
func sub(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			result, ok := overflow.Sub64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			result := float64(a) - float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("-", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			result := float64(a) - float64(b)
			return BFloat(result), nil
		case BFloat:
			result := float64(a) - float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("-", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("-", aVal, bVal)
		}
	case BStr:
		return nil, errTypeMismatch("-", aVal, bVal)
	case BObj:
		return nil, errTypeMismatch("-", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("-", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("-", aVal, bVal)
	case BArray:
		return nil, errTypeMismatch("-", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) - %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func mul(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			result, ok := overflow.Mul64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			result := float64(a) * float64(b)
			return BFloat(result), nil
		case BStr:
			result := strings.Repeat(string(b), int(a))
			return BStr(result), nil
		case BObj:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BArray:
			result := repeatArr([]BVal(b), int(a))
			return BArray(result), nil
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			result := float64(a) * float64(b)
			return BFloat(result), nil
		case BFloat:
			result := float64(a) * float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("*", aVal, bVal)
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			result := strings.Repeat(string(a), int(b))
			return BStr(result), nil
		case BFloat:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BStr:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("*", aVal, bVal)
		}
	case BObj:
		return nil, errTypeMismatch("*", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("*", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("*", aVal, bVal)
	case BArray:
		switch b := bVal.(type) {
		case BInt:
			result := repeatArr([]BVal(a), int(b))
			return BArray(result), nil
		case BFloat:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BStr:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("*", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("*", aVal, bVal)
		}

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) * %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func div(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BFloat:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("/", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BFloat:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("/", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("/", aVal, bVal)
		}
	case BStr:
		return nil, errTypeMismatch("/", aVal, bVal)
	case BObj:
		return nil, errTypeMismatch("/", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("/", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("/", aVal, bVal)
	case BArray:
		return nil, errTypeMismatch("/", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) / %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func floorDiv(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			if b == 0 {
				return nil, errDivByZero
			}
			result := BInt(a) / BInt(b)
			return BInt(result), nil
		case BFloat:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("//", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BFloat:
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("//", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("//", aVal, bVal)
		}
	case BStr:
		return nil, errTypeMismatch("//", aVal, bVal)
	case BObj:
		return nil, errTypeMismatch("//", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("//", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("//", aVal, bVal)
	case BArray:
		return nil, errTypeMismatch("//", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) // %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func pow(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			result, ok := intPow(a, b)
			if !ok {
				return nil, errOverflow
			}
			return result, nil
		case BFloat:
			result := math.Pow(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("**", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			result := math.Pow(float64(a), float64(b))
			return BFloat(result), nil
		case BFloat:
			result := math.Pow(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("**", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("**", aVal, bVal)
		}
	case BStr:
		return nil, errTypeMismatch("**", aVal, bVal)
	case BObj:
		return nil, errTypeMismatch("**", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("**", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("**", aVal, bVal)
	case BArray:
		return nil, errTypeMismatch("**", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) ** %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func modulo(aVal, bVal BVal) (BVal, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			result := BInt(a) % BInt(b)
			return BInt(result), nil
		case BFloat:
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("%", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BFloat:
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BObj:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BFunc:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BNull:
			return nil, errTypeMismatch("%", aVal, bVal)
		case BArray:
			return nil, errTypeMismatch("%", aVal, bVal)
		}
	case BStr:
		return nil, errTypeMismatch("%", aVal, bVal)
	case BObj:
		return nil, errTypeMismatch("%", aVal, bVal)
	case BFunc:
		return nil, errTypeMismatch("%", aVal, bVal)
	case BNull:
		return nil, errTypeMismatch("%", aVal, bVal)
	case BArray:
		return nil, errTypeMismatch("%", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) %% %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}

// Returns -1 if a < b, 0 if a == b, 1 if a > b
// op is only used for debugging
func cmp(aVal, bVal BVal, op string) (int, error) {
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			aa, bb := BInt(a), BInt(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BFloat:
			aa, bb := BFloat(a), BFloat(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BStr:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BObj:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BFunc:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BNull:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BArray:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			aa, bb := BFloat(a), BFloat(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BFloat:
			aa, bb := BFloat(a), BFloat(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BStr:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BObj:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BFunc:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BNull:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BArray:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BFloat:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BStr:
			aa, bb := BStr(a), BStr(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BObj:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BFunc:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BNull:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		case BArray:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
	case BObj:
		return 0, errTypeMismatch("cmp", aVal, bVal)
	case BFunc:
		return 0, errTypeMismatch("cmp", aVal, bVal)
	case BNull:
		return 0, errTypeMismatch("cmp", aVal, bVal)
	case BArray:
		return 0, errTypeMismatch("cmp", aVal, bVal)

	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) %s %s(%s)", aVal, aVal.Typename(), op, bVal, bVal.Typename()))
}
