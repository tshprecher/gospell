package check

// A Trie represents the head node of a trie.
type Trie struct {
	alphabet Alphabet
	children   map[rune]*Trie
	terminates bool
}

// Add adds a word to the trie and returns true if
// the word already exists or is added successfully,
// false otherwise.
func (t *Trie) Add(word string) bool {
	letters := []rune(word)
	curTrie := t
	for _, l := range letters {
		if !t.alphabet.Contains(l) {
			return false
		}
		if curTrie.children[l] == nil {
			newChild := &Trie{alphabet: t.alphabet, children: make(map[rune]*Trie)}
			curTrie.children[l] = newChild
			curTrie = newChild
		} else {
			curTrie = curTrie.children[l]
		}
	}
	curTrie.terminates = true
	return true
}

func (t *Trie) Contains(word []rune) bool {
	for _, l := range word {
		if !t.alphabet.Contains(l) {
			return false
		}
		child := t.children[l]
		if child == nil {
			return false
		}
		t = child
	}

	if !t.terminates {
		return false
	}
	return true
}

func (t *Trie) Alphabet() Alphabet {
	return t.alphabet
}

// NewTrie create a Trie from a given Alphabet.
func NewTrie(alphabet Alphabet) *Trie {
	trie := &Trie{alphabet: alphabet, children: make(map[rune]*Trie)}
	return trie
}
