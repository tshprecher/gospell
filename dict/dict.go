package dict

type Dict interface {
	ContainsSlice(word []rune) bool
	ContainsString(word string) bool
}

type Alphabet struct {
	letters []rune
}

func (alph *Alphabet) Size() int {
	return len(alph.letters)
}

func (alph *Alphabet) Index(r rune) (index int, ok bool) {
	for i, _ := range alph.letters {
		if alph.letters[i] == r {
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
