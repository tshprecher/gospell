package main

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
	letters []rune
}

func (ab *Alphabet) Size() int {
	return len(ab.letters)
}

// TODO: benchmark map vs for loop. both constant time, but which
// one is generally faster when considering memory, cpu cache?
func (ab *Alphabet) Index(r rune) (index int, ok bool) {
	for i, _ := range ab.letters {
		if ab.letters[i] == r {
			return i, true
		}
	}
	return 0, false
}

// sanitizeWord returns a lowercase copy of the word with the last character
// removed if it's not included in the alphabet.
// TODO: this should probably be attached to a language construct
//   instead of an alphabet
func (ab *Alphabet) Sanitize(word string) string {
	// TODO: this is crap code!
	// be smarter with allocations
	if len(word) == 0 {
		return word
	}
	if _, ok := ab.Index(rune(word[len(word)-1])); !ok {
		return strings.ToLower(word[0:len(word)-1])
	}
	return strings.ToLower(word)
}

func NewAlphabet(letters []rune) *Alphabet {
	lettersCopy := make([]rune, len(letters))
	copy(lettersCopy, letters)
	return &Alphabet{lettersCopy}
}
