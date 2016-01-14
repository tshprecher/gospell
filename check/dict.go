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

// sanitizeWord returns a lowercase copy of the word with the last character
// removed if it's not included in the alphabet.
// TODO: this should probably be attached to a language construct
//   instead of an alphabet
func (ab Alphabet) Sanitize(word string) string {
	// TODO: this is crap code!
	// be smarter with allocations
	if len(word) == 0 {
		return word
	}
	if !ab.Contains(rune(word[len(word)-1])) {
		return strings.ToLower(word[0:len(word)-1])
	}
	return strings.ToLower(word)
}

func NewAlphabet(letters []rune) *Alphabet {
	lettersMap := make(map[rune]bool)
	for l := range letters {
		lettersMap[letters[l]] = true
	}
	return &Alphabet{lettersMap}
}
