package parser

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

//go:generate go run ../../cmd/dump_lex/main.go

// This file is not used in the parser build but is used by cmd/dump_lex/main.go to generate the lexer.json file

var operators = [...]string{
	// Sort in descending length order
	"and",
	"not",
	"or",
	`\+=`,
	`\-=`,
	`\*=`,
	`\/=`,
	":",
	`\-`,
	`\+`,
	`\*`,
	`\/`,
	`\(`,
	`\?`,
	`\.`,
	`\[`,
}

var endExpr = [...]string{
	",",
	`\)`,
	`\]`,
}

var operatorString = fmt.Sprintf("%s", strings.Join(operators[:], "|"))
var endExprString = fmt.Sprintf("%s", strings.Join(endExpr[:], "|"))

var GenLexerDefinition = lexer.MustStateful(lexer.Rules{
	"Root": {
		lexer.Include("Expr"),
	},
	"Expr": {
		{"DoubleString", `"`, lexer.Push("DoubleString")},
		{"SingleString", `'[^\']*'`, nil},
		{`whitespace`, `\s+`, nil},
		{`Bool`, "true|false", nil},
		{`Op`, operatorString, nil},
		{`EndExpr`, endExprString, nil},
		{"Ident", `[a-zA-Z]\w*`, nil},
		{"Float", `\d*\.\d+(e\d+)?`, nil},
		{"Int", `[1-9][\d_]*`, nil},
		{"HexInt", `0x[0-9a-fA-F_]+`, nil},
		{"OctInt", `0o[0-7_]+`, nil},
		{"BinInt", `0b[01_]+`, nil},
		// {"ExprEnd", `}`, lexer.Pop()},
	},
	"DoubleString": {
		// TODO string escapes
		// {"ExprEscaped", `\\.`, nil},
		{"DoubleStringEnd", `"`, lexer.Pop()},
		// {"Expr", `\${`, lexer.Push("Expr")},
		// {"ExprChar", `[^$`\]+`, nil},
	},
	// "ExprReference": {
	// 	{"ExprDot", `\.`, nil},
	// 	{"Ident", `\w+`, nil},
	// 	lexer.Return(),
	// },
})
