package main

import (
	"bytes"
	"fmt"
	"github.com/tshprecher/gospell/check"
	"strings"
	"text/scanner"
)

type misspelling struct {
	line        int
	word        string
	suggestions []string
}

type fileResult struct {
	filename     string
	misspellings []misspelling
}

func (r fileResult) String() string {
	var buffer bytes.Buffer
	for _, misp := range r.misspellings {
		buffer.WriteString(fmt.Sprintf("%s:%d\t'%s'\t%v\n", r.filename, misp.line, misp.word, misp.suggestions))
	}
	return buffer.String()
}

type processor interface {
	run(filename string, src []byte, dict check.Dict) (*fileResult, error)
}

// TODO: benchmark object allocations
func sanitizeWord(word string, alphabet check.Alphabet) string {
	if len(word) == 0 {
		return word
	}
	letters := []rune(strings.ToLower(word))
	startIdx, endIdx := 0, len(letters)-1

	for startIdx < endIdx && !alphabet.Contains(letters[startIdx]) {
		startIdx++
	}
	for startIdx < endIdx && !alphabet.Contains(letters[endIdx]) {
		endIdx--
	}
	return string(letters[startIdx : endIdx+1])
}

type cStyleCommentProcessor struct {
	checker check.Checker
}

func handleComments(curLineNo int, comment string, dict check.Dict, res *fileResult) {
	var lines []string = strings.Split(comment, "\n")
	if len(lines) == 0 {
		lines = append(lines, comment)
	}
	for l := 0; l < len(lines); l++ {
		for _, word := range strings.Split(lines[l], " ") {
			sanitized := sanitizeWord(word, dict.Alphabet())
			if mis, sug := checker.IsMisspelled(sanitized, dict); mis {
				m := misspelling{curLineNo - len(lines) + 1 + l, sanitized, sug}
				res.misspellings = append(res.misspellings, m)
			}
		}
	}
}

func (p cStyleCommentProcessor) run(filename string, src []byte, dict check.Dict) (*fileResult, error) {
	scan := new(scanner.Scanner).Init(bytes.NewReader(src))
	scan.Mode = scanner.ScanComments
	res := &fileResult{filename, nil}

	for tok := scan.Scan(); tok != scanner.EOF; tok = scan.Scan() {
		if tok == scanner.Comment {
			handleComments(scan.Pos().Line, scan.TokenText(), dict, res)
		}
	}
	return res, nil
}
