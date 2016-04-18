package check

import (
	"strings"
)

// A Dict represents a dictionary of words.
type Dict interface {
	// Contains returns true iff the dictionary contains the word
	Contains(word []rune) bool
	// Alphabet returns the alphabet of the dictionary's language
	Alphabet() Alphabet
}

// An Alphabet represents the set of character allowed in a word.
type Alphabet string

// Size returns the number of characters.
func (ab Alphabet) Size() int {
	return len(string(ab))
}

// Contains returns true iff the character exists in the alphabet.
func (ab Alphabet) Contains(r rune) bool {
	return strings.ContainsRune(string(ab), r)
}

// Letter returns the character at a given index into the alphabet.
// ok == true iff a character at the index exists.
func (ab Alphabet) Letter(index int) (letter rune, ok bool) {
	if index < 0 || index >= ab.Size() {
		ok = false
		return
	}
	return []rune(string(ab))[index], true
}

// NewAlphabet creates an alphabet from the given slice of characters.
func NewAlphabet(letters string) *Alphabet {
	ab := Alphabet(letters)
	return &ab
}
