package check

import (
	"testing"
)

func TestContains(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}

	var tests = []struct {
		word   string
		expect bool
	}{
		{"on", false},
		{"one", true},
		{"once", false},
		{"ones", false},
		{"two", true},
		{"three", true},
		{"four", false},
		{"zero", false},
	}

	// TODO: hate to have to copy the alphabet here.
	trie, _ := NewTrie(words, *NewAlphabet([]rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z', '-', '\'',}))

	for _, test := range tests {
		res := trie.Contains([]rune(test.word))
		if res != test.expect {
			t.Errorf("expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
	}
}
