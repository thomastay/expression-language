package vm

import (
	"math"
	"strings"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

func add(aVal BVal, bVal BVal) (BVal, error) {
	// For lack of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both ints
			result, ok := overflow.Add64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			result := float64(a) + float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			result := float64(a) + float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			result := float64(a) + float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			return nil, errMismatchedTypes
		case BFloat:
			return nil, errMismatchedTypes
		case BStr:
			s := a + b
			return BStr(s), nil
		default:
			panic("Unreachable")
		}
	default:
		panic("Unreachable")
	}
}

func sub(aVal BVal, bVal BVal) (BVal, error) {
	// For lack of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both ints
			result, ok := overflow.Sub64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			result := float64(a) - float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			result := float64(a) - float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			result := float64(a) - float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BStr:
		return nil, errMismatchedTypes
	default:
		panic("Unreachable")
	}
}

func mul(aVal BVal, bVal BVal) (BVal, error) {
	// For lack of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both ints
			result, ok := overflow.Mul64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			result := float64(a) * float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			s := strings.Repeat(string(b), int(a))
			return BStr(s), nil
		default:
			panic("Unreachable")
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			result := float64(a) * float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			result := float64(a) * float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			s := strings.Repeat(string(a), int(b))
			return BStr(s), nil
		case BFloat:
			return nil, errMismatchedTypes
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	default:
		panic("Unreachable")
	}
}

func div(aVal BVal, bVal BVal) (BVal, error) {
	// For lack of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			if b == 0 {
				return nil, errDivByZero
			}
			// case 1: both ints
			result, ok := overflow.Div64(int64(a), int64(b))
			if !ok {
				return nil, errOverflow
			}
			return BInt(result), nil
		case BFloat:
			if b == 0 {
				return nil, errDivByZero
			}
			// Case 2a: one int, one float
			// Cast int to float and add
			result := float64(a) / float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			if b == 0 {
				return nil, errDivByZero
			}
			result := float64(a) / float64(b)
			// Note: find another lib to check for overflow
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BStr:
		return nil, errMismatchedTypes
	default:
		panic("Unreachable")
	}
}

func modulo(aVal BVal, bVal BVal) (BVal, error) {
	// For lack of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both ints
			return BInt(a % b), nil
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			result := math.Mod(float64(a), float64(b))
			return BFloat(result), nil
		case BStr:
			return nil, errMismatchedTypes
		default:
			panic("Unreachable")
		}
	case BStr:
		return nil, errMismatchedTypes
	default:
		panic("Unreachable")
	}
}
