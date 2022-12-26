package vm

import (
	"errors"

	. "github.com/thomastay/expression_language/pkg/bytecode"
)

var ErrInvalidNumParams = errors.New("Invalid number of arguments provided")

func Wrap0(fn func() BVal) VMFunc {
	return func(args []BVal) (BVal, error) {
		if len(args) != 0 {
			return nil, ErrInvalidNumParams
		}
		fn()
		return BNull{}, nil
	}
}

func Wrap1(fn func(p1 BVal) BVal) VMFunc {
	return func(args []BVal) (BVal, error) {
		if len(args) != 1 {
			return nil, ErrInvalidNumParams
		}
		return fn(args[0]), nil
	}
}

func Wrap2(fn func(p1, p2 BVal) BVal) VMFunc {
	return func(args []BVal) (BVal, error) {
		if len(args) != 2 {
			return nil, ErrInvalidNumParams
		}
		return fn(args[0], args[1]), nil
	}
}

func Wrap3(fn func(p1, p2, p3 BVal) BVal) VMFunc {
	return func(args []BVal) (BVal, error) {
		if len(args) != 3 {
			return nil, ErrInvalidNumParams
		}
		return fn(args[0], args[1], args[2]), nil
	}
}
