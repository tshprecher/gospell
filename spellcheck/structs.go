// TODO: is the naming of the package and file here correct?
package spellcheck

import (
	"github.com/tshprecher/gospell/dict"
)

// TODO: should I be strict and make these struct immutable outside the package by adding methods?
type Position struct {
	Line int
	Column int
}

type Misspelling struct {
	Word string
	Suggestions []string
	Locations []Position
}

type Result struct {
	Filename string
	Misspellings []Misspelling
}

type SpellChecker interface {
	Check(word string, dict dict.Dict) (ok bool, suggestions *[]string)
}

type StrictSpellChecker struct{}

func (s *StrictSpellChecker) Check(word string, dict dict.Dict) (bool, *[]string) {
	if dict.Contains([]rune(word)) {
		return true, nil
	}
	return false, nil
}
