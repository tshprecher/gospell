package dict

import (
	"errors"
)

// TODO: make this package protected?
type Trie struct {
	// TODO: make this an arbitrary alphabet
	// only handle lowercase english letters
	children   [26]*Trie
	terminates bool
}

func (t *Trie) ContainsSlice(word []rune) bool {
	for _, l := range word {
		idx, ok := toChildIndex(l)
		if !ok {
			return false
		}
		child := t.children[idx]
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

func (t *Trie) ContainsString(word string) bool {
	return t.ContainsSlice([]rune(word))
}

func toChildIndex(r rune) (index int, ok bool) {
	if r < 'a' || r > 'z' {
		return 0, false
	}
	return int(r - 'a'), true
}

func NewTrie(words []string) (*Trie, error) {
	trie := new(Trie)
	for _, w := range words {
		runes := []rune(w)
		currentTrie := trie
		for _, r := range runes {
			idx, ok := toChildIndex(r)
			if !ok {
				return nil, errors.New("only characters from 'a' to 'z' allowed in trie")
			}
			if currentTrie.children[idx] == nil {
				nTrie := new(Trie)
				currentTrie.children[idx] = nTrie
				currentTrie = nTrie
			} else {
				currentTrie = currentTrie.children[idx]
			}
		}
		currentTrie.terminates = true
	}

	return trie, nil
}
