package check

import "testing"

func TestAlphabet(t *testing.T) {
	alphabet := NewAlphabet([]rune{'a', 'b', 'c'})

	if alphabet.Size() != 3 {
		t.Errorf("expected size -> 3, received size -> %d.", alphabet.Size())
	}

	var containsTests = []struct {
		rne         rune
		expectedOk    bool
	}{
		{'a', true},
		{'b', true},
		{'c', true},
		{'d', false},
		{'0', false},
		{'A', false},
	}

	for _, test := range containsTests {
		ok := alphabet.Contains(test.rne)
		if ok != test.expectedOk {
			t.Errorf("expected ok -> %v, received ok -> %v for rune '%c'", ok, test.expectedOk, test.rne)
		}
	}

	var letterTests = []struct{
		index int
		expectRune rune
		expectedOk bool
	}{
		{-1, 0, false},
		{0, 'a', true},
		{1, 'b', true},
		{2, 'c', true},
		{3, 0, false},
	}

	for _, test := range letterTests {
		letter, ok := alphabet.Letter(test.index)
		if letter != test.expectRune || ok != test.expectedOk {
			t.Errorf("expected letter -> '%c', ok -> %v, received letter -> '%c', ok -> %v",
				test.expectRune, test.expectedOk, letter, ok)
		}
	}
}
