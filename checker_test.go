package main

import "testing"

func TestStrictChecker(t *testing.T) {
	words := []string{
		"one",
		"two",
		"three",
	}

	var tests = []struct {
		word   string
		expect bool
	}{
		{"hone", true},
		{"one", false},
		{"two", false},
	}

	sc := StrictChecker{}
	dict, _ := NewTrie(words, English)

	for _, test := range tests {
		res := sc.IsMisspelled(test.word, dict)
		if res != test.expect {
			t.Errorf("expected %v, received %v for word '%v'", test.expect, res, test.word)
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
		word   string
		expect bool
	}{
		{"one", false},
		{"oen", true},        // swapped chars           -> misspelled
		{"oenn", true},       // swapped + insertion     -> misspelled
		{"oene", true},       // insertion               -> misspelled
		{"oeene", false},     // two insertions          -> not misspelled
		{"tree", true},       // deletion                -> misspelled
		{"ene", true},        // deletion + insertion    -> misspelled
		{"tee", false},       // two deletions           -> misspelled
		{"twotwotwo", false}, // too many insertions     -> not misspelled
	}

	dc := DeltaChecker{
		AllowedIns:   1,
		AllowedDel:   1,
		AllowedSwaps: 1}

	dict, _ := NewTrie(words, English)

	for _, test := range tests {
		res := dc.IsMisspelled(test.word, dict)
		if res != test.expect {
			t.Errorf("expected %v, received %v for word '%v'", test.expect, res, test.word)
		}
	}
}
