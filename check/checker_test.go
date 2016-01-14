package check


import (
	"testing"
)

func TestStrictChecker(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}

	var tests = []struct {
		word        string
		misspelled  bool
		suggested   []string
	}{
		{"hone", true, nil},
		{"one", false, nil},
		{"two", false, nil},
	}

	sc := StrictChecker{}
	// TODO: hate to have to copy the Alphabet here to avoid import cycle: check -> lang -> check
	dict, _ := NewTrie(words, *NewAlphabet([]rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z', '-', '\'',}))

	for _, test := range tests {
		res, sug := sc.IsMisspelled(test.word, dict)
		if res != test.misspelled {
			t.Errorf("expected misspelled=%v, received %v for word '%v'.", test.misspelled, res, test.word)
		}
		if len(sug) != len(test.suggested) {
			t.Errorf("expected suggested words %v, received %v.", test.suggested, sug)
		} else {
			for s := range(sug) {
				if sug[s] != test.suggested[s] {
					t.Errorf("expected suggested='%s', received '%s'.", test.suggested[s], sug[s])
				}
			}
		}
	}
}

func TestDeltaChecker(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}

	var tests = []struct {
		word       string
		misspelled bool
		suggested  []string
	}{
		{"one", false, nil},
		{"oen", true, []string{"one"}},        // swapped chars           -> misspelled
		{"oenn", true, []string{"one"}},       // swapped + insertion     -> misspelled
		{"oene", true, []string{"one"}},       // insertion               -> misspelled
		{"oeene", false, nil},     // two insertions          -> not misspelled
		{"tree", true, []string{"three"}},       // deletion                -> misspelled
		{"ene", true, []string{"one"}},        // deletion + insertion    -> misspelled
		{"tee", false, nil},       // two deletions           -> misspelled
		{"twotwotwo", false, nil}, // too many insertions     -> not misspelled
	}

	dc := DeltaChecker{
		AllowedIns:   1,
		AllowedDel:   1,
		AllowedSwaps: 1}

	// TODO: again, hate to have to copy this
	dict, _ := NewTrie(words, *NewAlphabet([]rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z', '-', '\'',}))

	for _, test := range tests {
		res, sug := dc.IsMisspelled(test.word, dict)
		if res != test.misspelled {
			t.Errorf("expected %v, received %v for word '%v'.", test.misspelled, res, test.word)
		}
		if len(sug) != len(test.suggested) {
			t.Errorf("expected suggested words %v, received %v.", test.suggested, sug)
		} else {
			for s := range(sug) {
				if sug[s] != test.suggested[s] {
					t.Errorf("expected suggested='%s', received '%s'.", test.suggested[s], sug[s])
				}
			}
		}
	}
}
