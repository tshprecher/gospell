// Package main runs this
// line two
// line three
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	fileset *token.FileSet = token.NewFileSet()
	checker Checker        = new(StrictChecker)
)

func processFile(filename string, dict Dict) (res *Result, err error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	res = &Result{filename, nil}

	ast, err := parser.ParseFile(fileset, filename, src, parser.ParseComments)

	// fmt.Printf("doc -> %v\n", ast.Doc)
	// fmt.Printf("package -> %v\n", ast.Name)
	// fmt.Printf("decls -> %v\n", ast.Decls)
	// fmt.Printf("scope -> %v\n", ast.Scope)
	// fmt.Printf("scope.outer -> %v\n", ast.Scope.Outer)
	// fmt.Printf("imports -> %v\n", ast.Imports)
	// fmt.Printf("unresolved -> %v\n", ast.Unresolved)
	// fmt.Printf("comments -> %v\n", ast.Comments)

	handleCommentGroup(ast.Doc, src, dict, res)
	for _, com := range(ast.Comments) {
		handleCommentGroup(com, src, dict, res)
	}
	return
}

func handleCommentGroup(cg *ast.CommentGroup, src []byte, dict Dict, res *Result) {
	if cg == nil {
		return
	}
	for _, com := range cg.List {
		line := stringFromPosition(src, com.Pos(), com.End())
//		fmt.Println("comment line ", line)
		for _, word := range(strings.Split(line, " ")) {
			if word == "//" {
				continue
			}
			ab := dict.Alphabet()
			sanitized := ab.Sanitize(word)
			if len(sanitized) > 4 && checker.IsMisspelled(sanitized, dict) {
				misp := Misspelling{sanitized, fileset.Position(com.Pos()).Line}
				res.Misspellings = append(res.Misspellings, misp)
			}
		}
	}
}

func stringFromPosition(src []byte, start, end token.Pos) string {
	// TODO: can we do this without creating so many objects?
	beginOffset, endOffset := fileset.Position(start).Offset, fileset.Position(end).Offset
	return string(src[beginOffset:endOffset])
}

// It's almost hard to beleive this works! :)
func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing file argument.")
	}
	if len(os.Args) > 2 {
		// TODO: allow multiple files
		log.Fatal("only one file argument allowed.")
	}

	dict, err := NewTrie(words, English)

	if err != nil {
		log.Fatalf("%v", err)
	}

	checker = &DeltaChecker{1, 1, 1}
	res, err := processFile(os.Args[1], dict)

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println(res)
}
