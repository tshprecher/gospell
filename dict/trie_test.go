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

	trie, _ := NewTrie(words)

	for _, test := range tests {
		res := trie.ContainsString(test.word)
		if res != test.expect {
			t.Errorf("ContainsString() expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
		res = trie.ContainsSlice([]rune(test.word))
		if res != test.expect {
			t.Errorf("ContainsSlice() expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
	}
}
