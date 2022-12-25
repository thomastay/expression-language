package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/thomastay/expression_language/pkg/parser"
)

func main() {
	outFileName := os.Args[1]
	file, err := json.MarshalIndent(parser.GenLexerDefinition, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outFileName, file, 644)
}
