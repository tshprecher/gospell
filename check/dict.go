package check

import (
	"strings"
)

// TODO: need to add something more here to help with efficiencies
// based on different implementations, like a trie, etc?
type Dict interface {
	Contains(word []rune) bool
	Alphabet() Alphabet
}

type Alphabet struct {
	letters map[rune]bool
}

func (ab Alphabet) Size() int {
	return len(ab.letters)
}

func (ab Alphabet) Contains(r rune) bool {
	_, ok := ab.letters[r]
	return ok
}

func NewAlphabet(letters []rune) *Alphabet {
	lettersMap := make(map[rune]bool)
	for l := range letters {
		lettersMap[letters[l]] = true
	}
	return &Alphabet{lettersMap}
}
