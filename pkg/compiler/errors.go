package compiler

import "fmt"

func errUnaryType(op, typename string) error {
	return fmt.Errorf("TypeError: bad operand type for unary %s: %s", op, typename)
}
