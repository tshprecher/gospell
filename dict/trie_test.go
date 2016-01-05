package dict

import "testing"

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

	trie, _ := NewTrie(words, English)

	for _, test := range tests {
		res := trie.Contains([]rune(test.word))
		if res != test.expect {
			t.Errorf("expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
	}
}
