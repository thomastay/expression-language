package vm

import (
	"math"

	. "github.com/thomastay/expression_language/pkg/bytecode"
)

//go:generate go run ./generate/main.go

// Returns true if a == b
// Unlike cmp, this function cannot return an error and must always compare values
func eq(aVal BVal, bVal BVal) bool {
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
	default:
		panic("Not impl")
	}
}

// Note that 0^0 returns 1, mathematically this is undefined
func intPow(a BInt, b BInt) BVal {
	// Exponentiation by squaring for positive integers
	// based on https://docs.rs/num-traits/latest/src/num_traits/pow.rs.html#189
	if b < 0 {
		return BFloat(math.Pow(float64(a), float64(b)))
	}
	if b == 0 {
		return BInt(1)
	}

	base := a
	for b > 1 {
		b >>= 1
		base = base * base
		if b&1 == 1 {
			// odd
			a *= base
		}
	}
	return a
}
