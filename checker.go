package main

import (
//	"fmt"
)

// TODO: add ability to return optional suggestions
type Checker interface {
	// IsMisspelled returns false if the word is classified as misspelled.
	IsMisspelled(word string, dict Dict) bool
}

type StrictChecker struct{}

func (StrictChecker) IsMisspelled(word string, dict Dict) bool {
	if dict.Contains([]rune(word)) {
		return false
	}
	return true
}

type DeltaChecker struct {
	AllowedIns   int
	AllowedDel   int
	AllowedSwaps int
}

func (dc *DeltaChecker) IsMisspelled(word string, dict Dict) bool {
	wordSlice := []rune(word)
	if dict.Contains(wordSlice) {
		return false
	}
	ab := dict.Alphabet()
	for _, r := range(wordSlice) {
		if _, ok := ab.Index(r); !ok {
			return false
		}
	}

	buf := make([]rune, len(wordSlice)+dc.AllowedIns)
	copy(buf, wordSlice)
	return dc.isMisspelledDelta(buf, dict, len(wordSlice), dc.AllowedIns, dc.AllowedDel, dc.AllowedSwaps)
}

// TODO: make this efficient and FAST!
// 1) do this iteratively instead of recursively
// 2) as small as slices may be, we can pass the beg and end indices to Contains
// 3) create a data structure that efficiently does the insertions and deletions,
//    allowing underlying Dict implementations like Tries to take advantage of their
//    representation
func (dc *DeltaChecker) isMisspelledDelta(word []rune, dict Dict, len, ins, del, swaps int) bool {
	if ins <= 0 && del <= 0 && swaps <= 0 {
		return dict.Contains(word[:len])
	}

	for w := range word {
		// attempt swap
		if swaps > 0 && w < len-1 {
			word[w], word[w+1] = word[w+1], word[w]
			if dict.Contains(word[:len]) || dc.isMisspelledDelta(word, dict, len, ins, del, swaps-1) {
//				fmt.Printf("suggested swap: '%v'\n", string(word[:len]))
				return true
			}
			word[w], word[w+1] = word[w+1], word[w]
		}

		// attempt insertion
		if ins > 0 {
			for _, r := range dict.Alphabet().letters {
				for l := len; l > w; l-- {
					word[l] = word[l-1]
				}

				word[w] = r
				if dict.Contains(word[:len]) || dc.isMisspelledDelta(word, dict, len+1, ins-1, del, swaps) {
//					fmt.Printf("suggested insertion: '%v'\n", string(word[:len]))
					return true
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

			if dict.Contains(word[:len]) || dc.isMisspelledDelta(word, dict, len-1, ins, del-1, swaps) {
//				fmt.Printf("suggested deletion: '%v'\n", string(word[:len]))
				return true
			}

			for l := len - 1; l > w; l-- {
				word[l] = word[l-1]
			}
			word[w] = deleted
		}
	}
	return false
}
