This package is a library implementing a very simple **expression language** written in Go for interpreting one-liners.

```
> 5 * 10 == 50 ? 'hello' : 'world'
'hello'
```

The expression language is just like what you see above. Just a one liner language which can evaluate numbers, strings, objects, and arrays. In particular, there are no:

1. No loops
1. No user defined variables
1. No user defined functions

# Use cases

Why would you ever want to use something so simple? One use case is as a DSL for configuration files. You can let users enter in a string in your JSON file, and interpret it using the interpreter provided.

In this case, the limitations of the language are a strength. Everything has to fit into one line, so it works well with JSON. You can only jump forwards in this language, so users' one liners can't run indefinitely. This language is intentionally not Turing complete.

```json
{
  "radius": 2,
  "circumference": "3.14158 * diameter * 2"
}
```

## Where do the variables come from?

As you can see above, you can `evaluate` variables, but users just can't define them **in the language**. So, they have to be provided by the host environment. In the above example, the host environment has to provide `diameter` to the VM, by parsing the JSON file. Here's how the above snippet would be implemented in main.go

```go
type Circle struct {
	diameter int,
	circumference string
}
func main() {
	var circle Circle
	err := json.Marshal(&circle, fileBytes)
	// handle err
	m := vm.New()
	vm.AddInt("diameter", circle.diameter)
	vmResult, err := vm.EvalString(circle.circumference)
	computedCircumference := vmResult.Val
	fmt.Println(computedCircumference)
	// prints 12.566320
}
```

## What else is in the language?

See `vm_test.go`
