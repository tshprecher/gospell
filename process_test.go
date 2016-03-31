package main

import (
	"fmt"
	"github.com/tshprecher/gospell/check"
	"github.com/tshprecher/gospell/lang"
	"io/ioutil"
	"testing"
)

var englishDict = check.NewTrie(lang.EnglishAlphabet)

func init() {
	for _, word := range lang.EnglishUsWords {
		englishDict.Add(word.Word)
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
				t.Errorf("%s:expected misspelling '%s' on line %d, received misspelling '%s' on line %d",
					tst.filename, tst.expected[m].word, tst.expected[m].line, misp[m].word, misp[m].line)
			}
		}
	}
}

func TestCStyleProcessor(t *testing.T) {
	tst := []test{
		test{
			"testdata/test_go.go",
			[]misspelling{
				misspelling{10, "helllo", []string{"hello"}},
				misspelling{12, "wrold", []string{"world", "wold"}},
			},
		},
		test{
			"testdata/test_c.c",
			[]misspelling{
				misspelling{7, "helllo", []string{"hello"}},
				misspelling{9, "wrold", []string{"world", "wold"}},
			},
		},
	}
	proc := cStyleCommentProcessor{checker}
	runTests(t, proc, englishDict, tst)
}
