package main

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

func NewAlphabet(letters []rune) *Alphabet {
	lettersCopy := make([]rune, len(letters))
	copy(lettersCopy, letters)
	return &Alphabet{lettersCopy}
}
