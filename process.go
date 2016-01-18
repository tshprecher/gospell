package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tshprecher/gospell/check"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"text/scanner"
)

// TODO: make everything here package private.
type misspelling struct {
	Word string
	Line int
}

type fileResult struct {
	filename     string
	misspellings []misspelling
}

func (r fileResult) String() string {
	var buffer bytes.Buffer
	for _, misp := range r.misspellings {
		buffer.WriteString(fmt.Sprintf("%s:%d\t'%s'\n", r.filename, misp.Line, misp.Word))
	}
	return buffer.String()
}

// TODO: is processor the right name for this?
type processor interface {
	process(filename string, src []byte, dict check.Dict) (*fileResult, error)
}

func sanitizeWord(word string, alphabet check.Alphabet) string {
	if len(word) == 0 {
		return word
	}

	letters := []rune(strings.ToLower(word))
	startIdx, endIdx := 0, len(letters)-1

	if !alphabet.Contains(letters[startIdx]) {
		startIdx++
	}
	if !alphabet.Contains(letters[endIdx]) {
		endIdx--
	}

	if endIdx + 1 <= startIdx {
		return ""
	}
	return string(letters[startIdx:endIdx+1])
}

type goProcessor struct {}

func (p goProcessor) process(filename string, src []byte, dict check.Dict) (*fileResult, error) {
	res := &fileResult{filename, nil}
	ast, err := parser.ParseFile(fileset, filename, src, parser.ParseComments)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s: could not parse go file.", filename))
	}

	handleCommentGroup(ast.Doc, src, dict, res)
	for _, com := range ast.Comments {
		handleCommentGroup(com, src, dict, res)
	}
	return res, nil
}

func stringFromPosition(src []byte, start, end token.Pos) string {
	// TODO: can we do this without creating so many objects?
	beginOffset, endOffset := fileset.Position(start).Offset, fileset.Position(end).Offset
	return string(src[beginOffset:endOffset])
}

// TODO: encapsulate without a go processor?
func handleCommentGroup(cg *ast.CommentGroup, src []byte, dict check.Dict, res *fileResult) {
	if cg == nil {
		return
	}
	for _, com := range cg.List {
		line := stringFromPosition(src, com.Pos(), com.End())
		for _, word := range strings.Split(line, " ") {
			sanitized := sanitizeWord(word, dict.Alphabet())
			if m, _ := checker.IsMisspelled(sanitized, dict); m {
				// TODO: handle the suggestion(s)
				misp := misspelling{sanitized, fileset.Position(com.Pos()).Line}
				res.misspellings = append(res.misspellings, misp)
			}
		}
	}
}

type cStyleCommentProcessor struct{}

func (p cStyleCommentProcessor) process(filename string, src []byte, dict check.Dict) (*fileResult, error) {
	scan := new(scanner.Scanner).Init(bytes.NewReader(src))
	scan.Mode = scanner.ScanComments // only scan the comments
	res := &fileResult{filename, nil}

	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		if tok == scanner.Comment {
			for _, word := range strings.Split(scan.TokenText(), " ") {
				sanitized := sanitizeWord(word, dict.Alphabet())
				if m, _ := checker.IsMisspelled(sanitized, dict); m {
					// TODO: handle the suggestion(s)
					misp := misspelling{sanitized, scan.Line}
					res.misspellings = append(res.misspellings, misp)
				}
			}
		}
	}
	return res, nil
}
