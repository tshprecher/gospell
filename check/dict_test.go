package check

import "testing"

func TestAlphabet(t *testing.T) {
	alphabet := NewAlphabet([]rune{'a', 'b', 'c'})

	if alphabet.Size() != 3 {
		t.Errorf("expected size -> 3, received size -> %d.", alphabet.Size())
	}

	var tests = []struct {
		rne         rune
		expectOk    bool
	}{
		{'a', true},
		{'b', true},
		{'c', true},
		{'d', false},
		{'0', false},
		{'A', false},
	}

	for _, test := range tests {
		ok := alphabet.Contains(test.rne)
		if ok != test.expectOk {
			t.Errorf("expected ok -> %v, received ok -> %v for rune '%c'", ok, test.expectOk, test.rne)
		}
	}
}
