package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/thomastay/kaitai-expr-lang/pkg/parser"
)

func main() {
	expr, err := parser.ParseString(strings.Join(os.Args[1:], " "))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(expr.String())
}
