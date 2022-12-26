package vm

import "errors"

// Cache some common errors
var errNotEnoughStackValues = errors.New("Not enough values on stack")
var errOverflow = errors.New("Overflow")
var errMismatchedTypes = errors.New("Mismatched Types")
var errDivByZero = errors.New("Divided by zero")
