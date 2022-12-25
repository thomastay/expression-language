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
	outJsonFilename := "./lexer_generated.json"
	outGoFilename := "./lexer_generated.go"
	jsonBlob, err := json.MarshalIndent(parser.GenLexerDefinition, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outJsonFilename, jsonBlob, 644)
	goFile, err := os.Create(outGoFilename)
	if err != nil {
		log.Fatalln(err)
	}
	cmd := exec.Command("participle", "gen", "lexer", "parser", outJsonFilename)
	cmd.Stdout = goFile
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
