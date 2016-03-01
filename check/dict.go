package check

// A Dict represents a dictionary of words.
type Dict interface {
	// Contains returns true iff the dictionary contains the word
	Contains(word []rune) bool
	// Alphabet returns the alphabet of the dictionary's language
	Alphabet() Alphabet
}

// An Alphabet represents the set of character allowed in a word.
type Alphabet struct {
	letters []rune
}

// Size returns the number of characters.
func (ab Alphabet) Size() int {
	return len(ab.letters)
}

// Contains returns true iff the character exists in the alphabet.
func (ab Alphabet) Contains(r rune) bool {
	for l := range(ab.letters) {
		if ab.letters[l] == r {
			return true
		}
	}
	return false
}

// Letter returns the character at a given index into the alphabet.
// ok == true iff the character exists.
func (ab Alphabet) Letter(index int) (letter rune, ok bool) {
	if index < 0 || index >= ab.Size() {
		ok = false
		return
	}
	return ab.letters[index], true
}

// NewAlphabet creates an alphabet from the given slice of characters.
func NewAlphabet(letters []rune) *Alphabet {
	return &Alphabet{letters}
}
