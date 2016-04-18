package check

import (
	"testing"
)

func TestTrie(t *testing.T) {
	words := []string{
		"1",
		"12",
		"123",
		"1234",
	}

	var tests = []struct {
		word   string
		expect bool
	}{
		{"", false},
		{"1", true},
		{"12", true},
		{"123", true},
		{"1234", true},
		{"234", false},
		{"23", false},
		{"12345", false},
	}

	trie := NewTrie(*NewAlphabet("1234"))

	for w := range(words) {
		if ok := trie.Add(words[w]); !ok {
			t.Errorf("expected true, received false.")
		}
	}

	if ok := trie.Add("not_ok"); ok {
		t.Errorf("expected false, received true.")
	}

	for _, test := range tests {
		res := trie.Contains([]rune(test.word))
		if res != test.expect {
			t.Errorf("expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
	}
}
