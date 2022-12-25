package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

func main() {
	outFileName := os.Args[1]
	file, err := json.MarshalIndent(GenLexerDefinition, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outFileName, file, 644)
}

var operators = [...]string{
	// Sort in descending length order
	"and",
	"not",
	"or",
	"\\+=",
	"\\-=",
	"\\*=",
	"\\/=",
	":",
	"\\-",
	"\\+",
	"\\*",
	"\\/",
	"\\(",
	"\\?",
	"\\.",
	"\\[",
}

var endExpr = [...]string{
	",",
	"\\)",
	"\\]",
}

var operatorString = fmt.Sprintf("%s", strings.Join(operators[:], "|"))
var endExprString = fmt.Sprintf("%s", strings.Join(endExpr[:], "|"))

var GenLexerDefinition = lexer.MustStateful(lexer.Rules{
	"Root": {
		lexer.Include("Expr"),
	},
	"Expr": {
		{"DoubleString", `"`, lexer.Push("DoubleString")},
		{`whitespace`, `\s+`, nil},
		{`Op`, operatorString, nil},
		{`EndExpr`, endExprString, nil},
		{"Ident", `\w+`, nil},
		{"Int", `[0-9][1-9]*`, nil},
		// {"ExprEnd", `}`, lexer.Pop()},
	},
	"DoubleString": {
		// TODO string escapes
		// {"ExprEscaped", `\\.`, nil},
		{"DoubleStringEnd", `"`, lexer.Pop()},
		// {"Expr", `\${`, lexer.Push("Expr")},
		// {"ExprChar", `[^$"\\]+`, nil},
	},
	// "ExprReference": {
	// 	{"ExprDot", `\.`, nil},
	// 	{"Ident", `\w+`, nil},
	// 	lexer.Return(),
	// },
})
