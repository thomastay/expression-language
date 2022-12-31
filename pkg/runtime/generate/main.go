// Based on https://github.com/antonmedv/expr/blob/master/vm/runtime/helpers/main.go

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

var allowedCases = map[Case]CaseResult{
	// Add
	{"+", "BInt", "BInt"}: {
		s: `result, ok := overflow.Add64(int64(a), int64(b))
			if !ok { return nil, errOverflow }`},
	{"+", "BInt", "BFloat"}:   defaultFloat("+"),
	{"+", "BFloat", "BInt"}:   defaultFloat("+"),
	{"+", "BFloat", "BFloat"}: defaultFloat("+"),
	{"+", "BStr", "BStr"}:     defaultOp("+", "BStr"),
	{"+", "BArray", "BArray"}: {
		s:  "result := append([]BVal(a), []BVal(b)...)",
		tp: "BArray",
	},
	// Sub
	{"-", "BInt", "BInt"}: {
		s: `result, ok := overflow.Sub64(int64(a), int64(b))
			if !ok { return nil, errOverflow }`},
	{"-", "BInt", "BFloat"}:   defaultFloat("-"),
	{"-", "BFloat", "BInt"}:   defaultFloat("-"),
	{"-", "BFloat", "BFloat"}: defaultFloat("-"),
	// Mul
	{"*", "BInt", "BInt"}: {
		s: `result, ok := overflow.Mul64(int64(a), int64(b))
			if !ok { return nil, errOverflow }`},
	{"*", "BInt", "BFloat"}:   defaultFloat("*"),
	{"*", "BFloat", "BInt"}:   defaultFloat("*"),
	{"*", "BFloat", "BFloat"}: defaultFloat("*"),
	{"*", "BInt", "BStr"}:     mulIntStr("a", "b"),
	{"*", "BStr", "BInt"}:     mulIntStr("b", "a"),
	{"*", "BInt", "BArray"}:   mulIntArr("a", "b"),
	{"*", "BArray", "BInt"}:   mulIntArr("b", "a"),
	// Power
	{"**", "BInt", "BInt"}: {
		s: `result, ok := intPow(a, b)
			if !ok { return nil, errOverflow }`},
	{"**", "BInt", "BFloat"}:   defaultExp,
	{"**", "BFloat", "BInt"}:   defaultExp,
	{"**", "BFloat", "BFloat"}: defaultExp,
	// div
	{"/", "BInt", "BInt"}:     defaultFloat("/"),
	{"/", "BInt", "BFloat"}:   defaultFloat("/"),
	{"/", "BFloat", "BInt"}:   defaultFloat("/"),
	{"/", "BFloat", "BFloat"}: defaultFloat("/"),
	// Integer division
	{"//", "BInt", "BInt"}:     defaultOp("/", "BInt"),
	{"//", "BInt", "BFloat"}:   defaultFloat("/"),
	{"//", "BFloat", "BInt"}:   defaultFloat("/"),
	{"//", "BFloat", "BFloat"}: defaultFloat("/"),
	// mod
	{"%", "BInt", "BInt"}:     defaultOp("%%", "BInt"),
	{"%", "BInt", "BFloat"}:   {s: `result := math.Mod(float64(a), float64(b))`, tp: "BFloat"},
	{"%", "BFloat", "BInt"}:   {s: `result := math.Mod(float64(a), float64(b))`, tp: "BFloat"},
	{"%", "BFloat", "BFloat"}: {s: `result := math.Mod(float64(a), float64(b))`, tp: "BFloat"},
	// cmp
	{"cmp", "BInt", "BInt"}:     defaultCmp("BInt"),
	{"cmp", "BInt", "BFloat"}:   defaultCmp("BFloat"),
	{"cmp", "BFloat", "BInt"}:   defaultCmp("BFloat"),
	{"cmp", "BFloat", "BFloat"}: defaultCmp("BFloat"),
	{"cmp", "BStr", "BStr"}:     defaultCmp("BStr"),
	// eq is hardcoded since it's different from the rest
}

func cases(op string) string {
	var out string
	echoMain := func(s string, xs ...interface{}) {
		out += fmt.Sprintf(s, xs...) + "\n"
	}
	for _, aType := range types {
		echoMain("case %s:", aType)
		// write to bOut, if turns out that it's all blank, don't use bOut and use default.
		bOut := ""
		hasAnyCase := false
		{
			echo := func(s string, xs ...interface{}) {
				bOut += fmt.Sprintf(s, xs...) + "\n"
			}
			echo("switch b := bVal.(type) {")
			for _, bType := range types {
				echo("case %s:", bType)
				result, ok := allowedCases[Case{op, aType, bType}]
				if !ok {
					if op == "cmp" {
						echo("return 0, errTypeMismatch(\"%s\", aVal, bVal)", op)
					} else {
						echo("return nil, errTypeMismatch(\"%s\", aVal, bVal)", op)
					}
				} else {
					hasAnyCase = true
					if op == "%" || op == "/" || op == "//" {
						// div by zero
						echo(`if b == 0 { return nil, errDivByZero }`)
					}
					echo(result.s)
					if result.tp == "" {
						result.tp = aType
					}
					if op == "cmp" {
						// echo(result.s) // return is inside
					} else if op == "**" && aType == "BInt" && bType == "BInt" {
						// Special case this, since intPow can return a float or an int
						echo("return result, nil")
					} else {
						echo("return %s(result), nil", result.tp)
					}
				}
			}
		}
		if hasAnyCase {
			out += bOut
			echoMain("}")
		} else {
			if op == "cmp" {
				echoMain("return 0, errTypeMismatch(\"%s\", aVal, bVal)", op)
			} else {
				echoMain("return nil, errTypeMismatch(\"%s\", aVal, bVal)", op)
			}
		}
	}
	return out
}

func main() {
	var b bytes.Buffer
	err := template.Must(
		template.New("helpers").
			Funcs(template.FuncMap{
				"cases": func(op string) string { return cases(op) },
			}).
			Parse(helpers),
	).Execute(&b, types)
	if err != nil {
		panic(err)
	}

	outGoFilename := "./math_generated.go"
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	// os.WriteFile(outGoFilename, b.Bytes(), 644)
	os.WriteFile(outGoFilename, formatted, 644)
}

var types = []string{
	"BInt",
	"BFloat",
	"BStr",
	"BObj",
	"BFunc",
	"BNull",
	"BArray",
}

type Case struct {
	op string
	t1 string
	t2 string
}

type CaseResult struct {
	s string
	// If empty string, it's the same as t1
	tp string
}

func defaultOp(op, tp string) CaseResult {
	return CaseResult{
		s:  fmt.Sprintf("result := %s(a) %s %s(b)", tp, op, tp),
		tp: tp,
	}
}
func defaultFloat(op string) CaseResult {
	return CaseResult{
		s:  fmt.Sprintf("result := float64(a) %s float64(b)", op),
		tp: "BFloat",
	}
}
func boolToInt(boolName, otherName, op, tp string) CaseResult {
	return CaseResult{
		s: fmt.Sprintf(`integerValueOfBool := 0
	if %s {
		integerValueOfBool = 1
	}
	result := %s(integerValueOfBool) %s %s(%s)`, boolName, tp, op, tp, otherName),
		tp: tp,
	}
}

var mulIntStrTemplateStr = `var result string
		if int({{.IntName}}) <= 0 || len({{.StrName}}) == 0 {
			result = ""
		} else {
			memoryUsed, ok := overflow.Mul(len({{.StrName}}), int({{.IntName}}))
			if !ok || memoryUsed >= memoryLimit {
				return nil, errOOM
			}
			result = strings.Repeat(string({{.StrName}}), int({{.IntName}}))
		}`
var mulIntArrTemplateStr = `var result []BVal
		if int({{.IntName}}) > 0 && len({{.StrName}}) > 0 {
			memoryUsed, ok := overflow.Mul(len({{.StrName}}), int({{.IntName}}))
			if !ok || memoryUsed >= memoryLimit {
				return nil, errOOM
			}
			result = repeatArr([]BVal({{.StrName}}), int({{.IntName}}))
		}`

type MulTemplData struct {
	StrName string
	IntName string
}

func mulIntStr(intName, strName string) CaseResult {
	data := MulTemplData{strName, intName}
	var b bytes.Buffer
	tmpl := template.Must(template.New("MulIntStr").Parse(mulIntStrTemplateStr))
	err := tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return CaseResult{
		s:  b.String(),
		tp: "BStr",
	}
}
func mulIntArr(intName, arrName string) CaseResult {
	data := MulTemplData{arrName, intName}
	var b bytes.Buffer
	tmpl := template.Must(template.New("MulIntStr").Parse(mulIntArrTemplateStr))
	err := tmpl.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return CaseResult{
		s:  b.String(),
		tp: "BArray",
	}
}

var defaultExp = CaseResult{
	s:  "result := math.Pow(float64(a), float64(b))",
	tp: "BFloat",
}

func defaultCmp(tp string) CaseResult {
	return CaseResult{
		s: fmt.Sprintf(`aa, bb := %s(a), %s(b)
			if aa < bb {
				return -1, nil
			} else if aa == bb {
				return 0, nil
			}
			return 1, nil`, tp, tp),
		tp: tp,
	}
}

const helpers = `// Code generated by runtime/generate/main.go. DO NOT EDIT.
package runtime
import (
	"fmt"
	"math"
	"strings"

	"github.com/johncgriffin/overflow"
	. "github.com/thomastay/expression_language/pkg/bytecode"
)

func Add(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "+" }}
	}
	panic(fmt.Sprintf("Unhandled operation between %s and %s: %s + %s", aVal.Typename(), bVal.Typename(), aVal, bVal))
}
func Sub(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "-" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) - %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func Mul(aVal, bVal BVal, memoryLimit int) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "*" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) * %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func Div(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "/" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) / %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func FloorDiv(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "//" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) // %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func Pow(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "**" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) ** %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}
func Modulo(aVal, bVal BVal) (BVal, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "%" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) %% %s(%s)", aVal, aVal.Typename(), bVal, bVal.Typename()))
}

// Returns -1 if a < b, 0 if a == b, 1 if a > b
// op is only used for debugging
func Cmp(aVal, bVal BVal, op string) (int, error) {
	aVal = castBoolToInt(aVal)
	bVal = castBoolToInt(bVal)
	switch a := aVal.(type) {
	{{ cases "cmp" }}
	}
	panic(fmt.Sprintf("Unhandled operation: %s(%s) %s %s(%s)", aVal, aVal.Typename(), op, bVal, bVal.Typename()))
}
`
