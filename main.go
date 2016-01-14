// Package main runs this
// line two
// line three
package main

import (
	"flag"
	"fmt"
	"github.com/tshprecher/gospell/check"
	"github.com/tshprecher/gospell/lang"
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
	checker check.Checker
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func processDir(dir string, dict check.Dict) error {
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

func processFile(filename string, dict check.Dict) (*Result, error) {
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

func handleCommentGroup(cg *ast.CommentGroup, src []byte, dict check.Dict, res *Result) {
	if cg == nil {
		return
	}
	for _, com := range cg.List {
		line := stringFromPosition(src, com.Pos(), com.End())
		for _, word := range strings.Split(line, " ") {
			if word == "//" {
				continue
			}
			sanitized := dict.Alphabet().Sanitize(word)
			if len(sanitized) > *minLength {
				// TODO: can we do short circuit eval with fns returning mul values instead of nesting ifs?
				if m, _ := checker.IsMisspelled(sanitized, dict); m {
					// TODO: handle the suggestion(s)
					misp := Misspelling{sanitized, fileset.Position(com.Pos()).Line}
					res.Misspellings = append(res.Misspellings, misp)
				}
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

func main() {
	flag.Parse()
	checkPositiveArg(minLength, "l")
	checkPositiveArg(maxSwaps, "s")
	checkPositiveArg(maxIns, "i")
	checkPositiveArg(maxDel, "d")

	checker = &check.DeltaChecker{
		AllowedIns:   *maxIns,
		AllowedDel:   *maxDel,
		AllowedSwaps: *maxSwaps}

	for a := 0; a < flag.NArg(); a++ {
		filename := flag.Arg(a)
		fileInfo, err := os.Stat(filename)
		checkError(err)
		// TODO: define a Language struct encapsulating alphabet and words
		dict, err := check.NewTrie(lang.Words, lang.English)
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
