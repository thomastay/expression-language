package vm

import (
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

func isTruthy(x BVal) bool {
	switch a := x.(type) {
	case BInt:
		return int64(a) != 0
	case BFloat:
		return float64(a) != 0
	case BStr:
		return len(string(a)) > 0
	}
	return false
}
