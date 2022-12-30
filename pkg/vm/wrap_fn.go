package vm

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/thomastay/expression_language/pkg/bytecode"
)

var ErrInvalidNumParams = errors.New("Invalid number of arguments provided")

var bValType = reflect.TypeOf((*BVal)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func WrapFn(ff any) VMFuncWithArgs {
	fn := reflect.ValueOf(ff)
	fnType := fn.Type()
	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("VM WrapFn must receive a function, got %s", fnType.Kind()))
	}
	if fnType.NumOut() > 2 {
		panic("Functions should only have 0, 1, or 2 return values")
	}

	numIn := fnType.NumIn()
	for i := 0; i < numIn; i++ {
		t := fnType.In(i)
		if !t.Implements(bValType) {
			panic("Function param does not implement BVal")
		}
	}
	numOut := fnType.NumOut()
	if numOut > 0 {
		t := fnType.Out(0)
		if !t.Implements(bValType) {
			panic("Function return value 1 does not implement BVal")
		}
	}
	if numOut > 1 {
		t := fnType.Out(1)
		if !t.Implements(errorType) {
			panic("Function return value 2 does not implement Error")
		}
	}
	f := func(args []BVal) (BVal, error) {
		reflectArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
		reflectVals := fn.Call(reflectArgs)
		switch len(reflectVals) {
		case 0:
			return nil, nil
		case 1:
			return reflectVals[0].Interface().(BVal), nil
		case 2:
			return reflectVals[0].Interface().(BVal), reflectVals[1].Interface().(error)
		default:
			panic("Should not reach this point, wrong number of returns")
		}
	}
	return VMFuncWithArgs{Fn: f, NumArgs: fnType.NumIn()}
}
