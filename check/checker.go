package check

type Checker interface {
	// IsMisspelled returns false if the word is classified as misspelled.
	// If misspelled, suggested words may be nil or a slice with len > 0.
	IsMisspelled(word string, dict Dict) (res bool, suggested []string)
}

type StrictChecker struct{}

func (StrictChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if dict.Contains([]rune(word)) {
		return false, nil
	}
	return true, nil
}

type DeltaChecker struct {
	MinLength    int
	AllowedIns   int
	AllowedDel   int
	AllowedSwaps int
}

func (dc *DeltaChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if len(word) < dc.MinLength {
		return false, nil
	}
	wordSlice := []rune(word)
	if dict.Contains(wordSlice) {
		return false, nil
	}
	ab := dict.Alphabet()
	for _, r := range(wordSlice) {
		if !ab.Contains(r) {
			return false, nil
		}
	}

	buf := make([]rune, len(wordSlice)+dc.AllowedIns)
	seen := make(map[string]bool)
	copy(buf, wordSlice)
	return dc.isMisspelledDelta(seen, buf, dict, len(wordSlice), dc.AllowedIns, dc.AllowedDel, dc.AllowedSwaps)
}

// TODO: make this efficient and FAST!
// 1) do this iteratively instead of recursively
// 2) as small as slices may be, we can pass the beg and end indices to Contains
// 3) create a data structure that efficiently does the insertions and deletions,
//    allowing underlying Dict implementations like Tries to take advantage of their
//    representations
func (dc *DeltaChecker) isMisspelledDelta(seen map[string]bool, word []rune, dict Dict, len, ins, del, swaps int) (bool, []string) {
	wordStr := string(word)
	if seen[wordStr] {
		// avoid cycles
		return false, nil
	}
	if ins <= 0 && del <= 0 && swaps <= 0 {
		if dict.Contains(word[:len]) {
			return true, []string{string(word[:len])}
		} else {
			seen[wordStr] = true
			return false, nil
		}
	}

	for w := range word {
		// attempt swap
		if swaps > 0 && w < len-1 {
			word[w], word[w+1] = word[w+1], word[w]
			if dict.Contains(word[:len]) {
				return true, []string{string(word[:len])}
			}
			if misp, sug := dc.isMisspelledDelta(seen, word, dict, len, ins, del, swaps-1); misp {
				return true, sug
			}
			word[w], word[w+1] = word[w+1], word[w]
		}

		// attempt insertion
		if ins > 0 {
			for r := range dict.Alphabet().letters {
				for l := len; l > w; l-- {
					word[l] = word[l-1]
				}

				word[w] = r
				if dict.Contains(word[:len]) {
					return true, []string{string(word[:len])}
				}
				if misp, sug := dc.isMisspelledDelta(seen, word, dict, len+1, ins-1, del, swaps); misp {
					return true, sug
				}
				for l := w; l < len; l++ {
					word[l] = word[l+1]
				}
			}
		}

		// attempt deletion
		if del > 0 && len > 0 {
			deleted := word[w]
			for l := w + 1; l < len; l++ {
				word[l-1] = word[l]
			}
			if dict.Contains(word[:len]) {
				return true, []string{string(word[:len])}
			}
			if misp, sug := dc.isMisspelledDelta(seen, word, dict, len-1, ins, del-1, swaps); misp {
				return true, sug
			}
			for l := len - 1; l > w; l-- {
				word[l] = word[l-1]
			}
			word[w] = deleted
		}
	}
	seen[wordStr] = true
	return false, nil
}
