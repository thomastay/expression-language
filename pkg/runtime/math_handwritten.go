package runtime

import (
	"fmt"
	"math"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

//go:generate go run ./generate/main.go

// Returns true if a == b
// Unlike cmp, this function cannot return an error and must always compare values
func Eq(aVal BVal, bVal BVal) bool {
	aVal = CastBoolToInt(aVal)
	bVal = CastBoolToInt(bVal)
	// For want of a MATCH, the happiness was lost...
	switch a := aVal.(type) {
	case BInt:
		switch b := bVal.(type) {
		case BInt:
			// case 1: both int
			return a == b
		case BFloat:
			// Case 2a: one int, one float
			// Cast int to float and add
			aa, bb := float64(a), float64(b)
			return aa == bb
		default:
			return false
		}
	case BFloat:
		switch b := bVal.(type) {
		case BInt:
			aa, bb := float64(a), float64(b)
			return aa == bb
		case BFloat:
			aa, bb := float64(a), float64(b)
			return aa == bb
		default:
			return false
		}
	case BStr:
		b, ok := bVal.(BStr)
		if !ok {
			return false
		}
		aa, bb := string(a), string(b)
		return aa == bb
	case BNull:
		_, ok := bVal.(BNull)
		if !ok {
			return false
		}
		return true
	case BFunc:
		bfn, ok := bVal.(BFunc)
		if !ok {
			return false
		}
		// Go doesn't have a fast way of checking function pointers, so we cheat
		// and just check the names
		return a.Name == bfn.Name
	case BObj:
		bobj, ok := bVal.(BObj)
		if !ok {
			return false
		}
		if len(a) != len(bobj) {
			return false
		}
		for k, v := range a {
			// lookup
			v2, ok := bobj[k]
			if !ok {
				return false
			}
			if !Eq(v, v2) {
				return false
			}
		}
		return true
	case BArray:
		// compare element by element
		arr, ok := bVal.(BArray)
		if !ok {
			return false
		}
		if len(a) != len(arr) {
			return false
		}
		for i, val := range arr {
			if !Eq(val, a[i]) {
				return false
			}
		}
		return true
	default:
		panic("Not impl")
	}
}

// Note that 0^0 returns 1, mathematically this is undefined
func intPow(baseVal BInt, exp BInt) (BVal, bool) {
	// Exponentiation by squaring for positive integers
	// Taken from https://docs.rs/num-traits/latest/src/num_traits/pow.rs.html#189
	if exp < 0 {
		return BFloat(math.Pow(float64(baseVal), float64(exp))), true
	}
	if exp == 0 {
		return BInt(1), true
	}
	base := int64(baseVal) // for simplicity
	var ok bool

	// Fast path powers of two (most common exp)
	for exp&1 == 0 {
		base, ok = overflow.Mul64(base, base)
		if !ok {
			return nil, false
		}
		exp >>= 1
	}

	if exp == 1 {
		return BInt(base), true
	}

	// Based on the identity (start with y := 1)
	// y * x^n = | (yx) (x^2)^(n-1/2) if n is odd
	//           | y (x^2)^(n/2)      if n is even
	// In the code below, base is `y` and acc is `x`
	// Note that we assume that exp is now an odd number, since we shifted until odd above.
	// So we can immediately apply the first transformation.
	// base := yx (which equals base)
	// x = x ** 2
	acc := base
	for exp > 1 {
		exp >>= 1
		// acc **= 2
		acc, ok = overflow.Mul64(acc, acc)
		if !ok {
			return nil, false
		}
		if exp&1 == 1 {
			// base *= acc
			base, ok = overflow.Mul64(base, acc)
			if !ok {
				return nil, false
			}
		}
	}
	return BInt(base), true
}

func repeatArr(arr []BVal, n int) []BVal {
	result := make([]BVal, len(arr)*n)
	for i := 0; i < n; i++ {
		for j, val := range arr {
			result[i*len(arr)+j] = val
		}
	}
	return result
}

// Keep these two as leafs, easier for compiler to Inline
func BoolToInt(b BBool) BInt {
	if b {
		return 1
	}
	return 0
}

// Keep these two as leafs, easier for compiler to Inline
func CastBoolToInt(val BVal) BVal {
	v, ok := val.(BBool)
	if !ok {
		// fast path
		return val
	}
	if bool(v) {
		return BInt(1)
	}
	return BInt(0)
}

func OrdToBool(op string, ord int) bool {
	var result bool
	switch op {
	case "<":
		if ord == -1 {
			result = true
		}
	case ">":
		if ord == 1 {
			result = true
		}
	case "<=":
		if ord != 1 {
			result = true
		}
	case ">=":
		if ord != -1 {
			result = true
		}
	default:
		panic("No other cmp")
	}
	return result
}

func Negate(val BVal) (BVal, error) {
	switch a := val.(type) {
	case BInt:
		return BInt(-int64(a)), nil
	case BFloat:
		return BInt(-int64(a)), nil
	case BBool:
		return BInt(-BoolToInt(a)), nil
	default:
		return nil, fmt.Errorf("TypeError: bad operand type for unary -: %s", a.Typename())
	}
}
