package dict

import "testing"

func TestAlphabet(t *testing.T) {
	alphabet := NewAlphabet([]rune{'a', 'b', 'c'})

	if alphabet.Size() != 3 {
		t.Errorf("expected size -> 3, received size -> %d.", alphabet.Size())
	}

	var tests = []struct {
		rne         rune
		expectIndex int
		expectOk    bool
	}{
		{'a', 0, true},
		{'b', 1, true},
		{'c', 2, true},
		{'d', 0, false},
		{'0', 0, false},
		{'A', 0, false},
	}

	for _, test := range tests {
		idx, ok := alphabet.Index(test.rne)
		if ok != test.expectOk {
			t.Errorf("expected ok -> %v, received ok -> %v for rune '%c'", ok, test.expectOk, test.rne)
		}
		if idx != test.expectIndex {
			t.Errorf("expected index -> %d, received index -> %d for rune '%c'", idx, test.expectIndex, test.rne)
		}
	}
}
