// Package main runs this
// line two
// line three
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// TODO: put minLength in the delta checker
	minLength = flag.Int("l", 4, "filter out words less than 'l' characters")
	maxSwaps  = flag.Int("s", 1, "correct spelling up to 's' consecutive character swaps")
	maxIns    = flag.Int("i", 0, "correct spelling up to 'i' character insertions")
	maxDel    = flag.Int("d", 0, "correct spelling up to 'd' character deletions")

	fileset *token.FileSet = token.NewFileSet()
	checker Checker
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func processDir(dir string, dict Dict) error {
	var visit = func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TODO: generalize into a fn
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go") {
			res, err := processFile(path, dict)
			if err != nil {
				return err
			}
			fmt.Print(res)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return err
	}
	return nil
}

func processFile(filename string, dict Dict) (*Result, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	res := &Result{filename, nil}

	ast, err := parser.ParseFile(fileset, filename, src, parser.ParseComments)

	if err != nil {
		return nil, err
	}

	handleCommentGroup(ast.Doc, src, dict, res)
	for _, com := range ast.Comments {
		handleCommentGroup(com, src, dict, res)
	}
	return res, nil
}

func handleCommentGroup(cg *ast.CommentGroup, src []byte, dict Dict, res *Result) {
	if cg == nil {
		return
	}
	for _, com := range cg.List {
		line := stringFromPosition(src, com.Pos(), com.End())
		for _, word := range strings.Split(line, " ") {
			if word == "//" {
				continue
			}
			ab := dict.Alphabet()
			sanitized := ab.Sanitize(word)
			if len(sanitized) > *minLength && checker.IsMisspelled(sanitized, dict) {
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

func checkPositiveArg(value *int, arg string) {
	if *value < 0 {
		log.Fatalf("arg '%v' must be positive.", arg)
	}
}

// It's almost hard to beleive this works! :)
func main() {
	flag.Parse()
	checkPositiveArg(minLength, "l")
	checkPositiveArg(maxSwaps, "s")
	checkPositiveArg(maxIns, "i")
	checkPositiveArg(maxDel, "d")

	checker = &DeltaChecker{
		AllowedIns:   *maxIns,
		AllowedDel:   *maxDel,
		AllowedSwaps: *maxSwaps}

	for a := 0; a < flag.NArg(); a++ {
		filename := flag.Arg(a)
		fileInfo, err := os.Stat(filename)
		checkError(err)
		dict, err := NewTrie(words, English)
		checkError(err)

		if fileInfo.IsDir() {
			err = processDir(filename, dict)
			checkError(err)
		} else {
			res, err := processFile(filename, dict)
			checkError(err)
			fmt.Print(res)
		}
	}
}
