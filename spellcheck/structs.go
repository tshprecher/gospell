// TODO: is the naming of the package and file here correct?
package spellcheck

import (
	"github.com/tshprecher/gospell/dict"
	"go/token"
)

type Misspelling struct {
	Word string
	Positions []token.Pos
}

type Result struct {
	Filename string
	Misspellings []Misspelling
}

// TODO: add ability to return optional suggestions
type Checker interface {
	Check(word string, dict dict.Dict) bool
}

type StrictChecker struct{}

func (s *StrictChecker) Check(word string, dict dict.Dict) bool {
	if dict.Contains([]rune(word)) {
		return true
	}
	return false
}
