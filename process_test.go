package main

import (
	"github.com/tshprecher/gospell/check"
	"github.com/tshprecher/gospell/lang"
	"fmt"
	"io/ioutil"
	"testing"
)

var englishDict = check.NewTrie(lang.EnglishAlphabet)

func init() {
	for _, word := range lang.EnglishUsWords {
		englishDict.Add(word)
	}
}

type test struct {
	filename string
	expected []misspelling
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func runTests(t *testing.T, proc processor, dict check.Dict, tests []test) {
	for _, tst := range tests {
		src, _ := ioutil.ReadFile(tst.filename)
		result, _ := proc.run(tst.filename, src, dict)
		misp := result.misspellings
		fmt.Printf("%s: misp -> %v\n", tst.filename, misp)
		fmt.Printf("%s: expected misp -> %v\n", tst.filename, tst.expected)

		for m := 0; m < max(len(misp), len(tst.expected)); m++ {
			if m >= len(misp) {
				t.Errorf("%s:expected misspelling '%s' on line %d", tst.filename, tst.expected[m].word, tst.expected[m].line)
				continue
			}
			if m >= len(tst.expected) {
				t.Errorf("%s:unexpected misspelling '%s' on line %d", tst.filename, misp[m].word, misp[m].line)
				continue
			}
			if misp[m].line != tst.expected[m].line || misp[m].word != tst.expected[m].word {
				t.Errorf("%s:expected misspelling '%s' on line %d, received misspelling '%s' on line %d", tst.filename, tst.expected[m].word, tst.expected[m].line, misp[m].word, misp[m].line)
			}
		}
	}
}

func TestCStyleProcessor(t *testing.T) {
	tst := []test{
		test{
			"testdata/test_go.go",
			[]misspelling{
				misspelling{10, "helo", nil},
				misspelling{12, "wrold", nil},
			},
		},
		test{
			"testdata/test_c.c",
			[]misspelling{
				misspelling{7, "helo", nil},
				misspelling{9, "wrold", nil},
			},
		},
	}
	proc := cStyleCommentProcessor{checker}
	runTests(t, proc, englishDict, tst)
}
