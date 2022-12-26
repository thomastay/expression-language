package vm

import (
	"math"
	"strings"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

func add(aVal BVal, bVal BVal) (BVal, error) {
	// For want of a MATCH, the happiness was lost...
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
		default:
			return nil, errTypeMismatch("+", aVal, bVal)
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
		default:
			return nil, errTypeMismatch("+", aVal, bVal)
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BFloat:
			return nil, errTypeMismatch("+", aVal, bVal)
		case BStr:
			s := a + b
			return BStr(s), nil
		default:
			return nil, errTypeMismatch("+", aVal, bVal)
		}
	default:
		return nil, errTypeMismatch("+", aVal, bVal)
	}
}

func sub(aVal BVal, bVal BVal) (BVal, error) {
	// For want of a MATCH, the happiness was lost...
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
		default:
			return nil, errTypeMismatch("-", aVal, bVal)
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
		default:
			return nil, errTypeMismatch("-", aVal, bVal)
		}
	default:
		return nil, errTypeMismatch("-", aVal, bVal)
	}
}

func mul(aVal BVal, bVal BVal) (BVal, error) {
	// For want of a MATCH, the happiness was lost...
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
			return nil, errTypeMismatch("*", aVal, bVal)
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
		default:
			return nil, errTypeMismatch("*", aVal, bVal)
		}
	case BStr:
		switch b := bVal.(type) {
		case BInt:
			s := strings.Repeat(string(a), int(b))
			return BStr(s), nil
		case BFloat:
			return nil, errTypeMismatch("*", aVal, bVal)
		default:
			return nil, errTypeMismatch("*", aVal, bVal)
		}
	default:
		return nil, errTypeMismatch("*", aVal, bVal)
	}
}

func div(aVal BVal, bVal BVal) (BVal, error) {
	// For want of a MATCH, the happiness was lost...
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
		default:
			return nil, errTypeMismatch("/", aVal, bVal)
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
		default:
			return nil, errTypeMismatch("/", aVal, bVal)
		}
	default:
		return nil, errTypeMismatch("/", aVal, bVal)
	}
}

func modulo(aVal BVal, bVal BVal) (BVal, error) {
	// For want of a MATCH, the happiness was lost...
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
		default:
			return nil, errTypeMismatch("%", aVal, bVal)
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
		default:
			return nil, errTypeMismatch("%", aVal, bVal)
		}
	default:
		return nil, errTypeMismatch("%", aVal, bVal)
	}
}

// Returns -1 if a < b, 0 if a == b, 1 if a > b
func cmp(aVal BVal, bVal BVal) (int, error) {
	// For want of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both int
			if a < b {
				return -1, nil
			} else if a == b {
				return 0, nil
			}
			return 1, nil
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			aa, bb := float64(a), float64(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		default:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			// case 2b: one each
			aa, bb := float64(a), float64(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		case BFloat:
			// Case 3: both floats
			// Cast int to float and add
			aa, bb := float64(a), float64(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil
		default:
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
	case BStr:
		b, ok := bVal.(BStr)
		if !ok {
			return 0, errTypeMismatch("cmp", aVal, bVal)
		}
		aa, bb := string(a), string(b)
		if aa < bb {
			return -1, nil
		} else if aa == bb {
			return 0, nil
		}
		return 1, nil
	default:
		return 0, errTypeMismatch("cmp", aVal, bVal)
	}
}
