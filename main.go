package main

import (
	"fmt"
	"github.com/tshprecher/gospell/spellcheck"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func spellCheckFile(filename string) (*spellcheck.Result, error) {
	_, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)

	if err != nil {
		return nil, err
	}
	return nil, nil
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing file argument.")
	}
	if len(os.Args) > 2 {
		// TODO: allow multiple files
		log.Fatal("only one file argument allowed.")
	}

	res, err := spellCheckFile(os.Args[1])

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("spell check ran with result -> %v\n", res)
}
