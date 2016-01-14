package check

import (
	"errors"
)

// TODO: make this package protected?
type Trie struct {
	// TODO: distinguish between root and subsequent
	//       nodes to avoid having to copy values?
	alphabet Alphabet
	children   map[rune]*Trie
	terminates bool
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

func NewTrie(words []string, alphabet Alphabet) (*Trie, error) {
	trie := &Trie{alphabet: alphabet, children: make(map[rune]*Trie)}
	for _, w := range words {
		letters := []rune(w)
		currentTrie := trie
		for _, l := range letters {
			if !alphabet.Contains(l) {
				return nil, errors.New("unicode char '%c' not found in alphabet.")
			}
			if currentTrie.children[l] == nil {
				newChild := &Trie{alphabet: alphabet, children: make(map[rune]*Trie)}
				currentTrie.children[l] = newChild
				currentTrie = newChild
			} else {
				currentTrie = currentTrie.children[l]
			}
		}
		currentTrie.terminates = true
	}
	return trie, nil
}
