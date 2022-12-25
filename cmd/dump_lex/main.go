package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/thomastay/expression_language/pkg/parser"
)

func main() {
	// THIS FUNCTION IS MEANT TO BE RUN BY GO GENERATE
	outFileName := "./lexer_generated.json"
	jsonBlob, err := json.MarshalIndent(parser.GenLexerDefinition, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outFileName, jsonBlob, 644)
	cmd := exec.Command("participle", "gen", "lexer", "parser", outFileName)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
