package vm

import (
	"errors"
	"fmt"

	"github.com/thomastay/expression_language/pkg/bytecode"
)

// Cache some common errors
// var errNotEnoughStackValues = errors.New("VMError: Not enough values on stack")
var errOverflow = errors.New("ArithmeticError: Overflow")
var errDivByZero = errors.New("ArithmeticError: Divided by zero")
var errOOM = errors.New("Out of Memory")

func errTypeMismatch(op string, v1 bytecode.BVal, v2 bytecode.BVal) error {
	return fmt.Errorf("TypeError: unsupported operand type(s) for %s: %s and %s", op, v1.Typename(), v2.Typename())
}
