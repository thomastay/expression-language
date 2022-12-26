package vm

import "errors"

// Cache some common errors
var errNotEnoughStackValues = errors.New("Not enough values on stack")
var errOverflow = errors.New("Overflow")
