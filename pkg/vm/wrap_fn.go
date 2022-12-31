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

// Wrap a function and its name
// A function must have the following return types. Any number of parameters accepted, but no variadics for now.
//
//	func (a bytecode.BVal)
//	func (a bytecode.BVal) bytecode.BVal
//	func (a bytecode.BVal) (bytecode.BVal, error)
func WrapFn(name string, ff any) BFunc {
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
	var err error
	f := func(args []BVal) (BVal, error) {
		reflectArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
		reflectVals := fn.Call(reflectArgs)
		switch len(reflectVals) {
		case 0:
			return BNull{}, nil
		case 1:
			return toBVal(reflectVals[0]), nil
		case 2:
			err, _ = reflectVals[1].Interface().(error)
			return toBVal(reflectVals[0]), err
		default:
			panic("Should not reach this point, wrong number of returns")
		}
	}
	return BFunc{Fn: f, NumArgs: fnType.NumIn(), Name: name}
}

func toBVal(val reflect.Value) BVal {
	if val.IsNil() {
		return BNull{}
	}
	return val.Interface().(BVal)
}

// Clones an env
func CloneEnv(env VMEnv) VMEnv {
	result := make(VMEnv, len(env))
	for k, v := range env {
		result[k] = v
	}
	return result
}
