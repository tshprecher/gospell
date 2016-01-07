// Package main runs this
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
)

var (
	fileset *token.FileSet = token.NewFileSet()
	checker = new(StrictChecker)
)

func processFile(filename string) (res *Result, err error) {
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

	handleCommentGroup(ast.Doc, src, res)
	return
}

func handleCommentGroup(cg *ast.CommentGroup, src []byte, res *Result) {
	if cg == nil {
		return
	}
	for  _, com := range(cg.List) {
		fmt.Printf("checking CommentGroup -> %v\n", stringFromPosition(src, com.Pos(), com.End()))
	}
}

func stringFromPosition(src []byte, start, end token.Pos) string {
	// TODO: can we do this without creating so many objects?
	beginOffset, endOffset := fileset.Position(start).Offset, fileset.Position(end).Offset
	return string(src[beginOffset:endOffset])
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("missing file argument.")
	}
	if len(os.Args) > 2 {
		// TODO: allow multiple files
		log.Fatal("only one file argument allowed.")
	}

	res, err := processFile(os.Args[1])

	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("spell check ran with result -> %v\n", res)
}
