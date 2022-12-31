This package is a library implementing a very simple **expression language** written in Go for interpreting one-liners.

```
> 5 * 10 == 50 ? 'hello' : 'world'
'hello'
```

_Note that this library is very new and is mostly educational_. For a more refined implementation, see AntonMedv's Expr library here: https://github.com/antonmedv/expr

The expression language is just like what you see above. Just a one liner language which can evaluate numbers, strings, objects, and arrays. In particular, there are no:

1. No loops
1. No user defined variables
1. No user defined functions

# Use cases

Why would you ever want to use something so simple? One use case is as a DSL for configuration files. You can let users enter in a string in your JSON file, and interpret it using the interpreter provided.

In this case, the limitations of the language are a strength. Everything has to fit into one line, so it works well with JSON. You can only jump forwards in this language, so users' one liners can't run indefinitely. This language is intentionally not Turing complete.

```json
{
  "radius": 3,
  "circumference": "3.14158 * radius * 2",
  "area": "3.14158 * radius * radius"
}
```

## Where do the variables come from?

As you can see above, you can _evaluate_ variables, but users just can't define them **in the language**. So, they have to be provided by the host environment. In the above example, the host environment has to provide `radius` to the VM, by parsing the JSON file. Here's how the above snippet would be implemented in main.go

```go
type Circle struct {
	radius int64
	circumference string
	area string
}
func main() {
	var circle Circle
	err := json.Marshal(&circle, fileBytes)
	m := vm.New()
	// Host environment provides the variable for use
	// BInt is the type that the interpreter uses
	env := vm.VMEnv{
		"radius": bytecode.BInt(circle.radius)
	}

	vmResult, err := m.EvalString(circle.circumference, env)
	circumference := vmResult.Val
	fmt.Println("circumference", circumference)

	vmResult, err := m.EvalString(circle.area, env)
	area := vmResult.Val
	fmt.Println("area", area)
}
// prints:
// circumference 18.849480
// area 28.274220
```

## What else is in the language?

See `vm_test.go`

# Speed

Roughly ~300-500x slower than native compiled code

```
goos: windows
goarch: amd64
pkg: github.com/thomastay/expression_language/pkg/vm
cpu: Intel(R) Core(TM) i5-7200U CPU @ 2.50GHz
BenchmarkCollatz-4                  7858            163786 ns/op          132617 B/op        321 allocs/op
BenchmarkCollatzRegular-4        3935415               331.4 ns/op             0 B/op          0 allocs/op
```
