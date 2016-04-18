package check

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// A Checker classifies misspelled words.
type Checker interface {
	// IsMisspelled returns false if the word is classified as misspelled.
	// If misspelled, suggested words may be nil or a slice with len > 0.
	IsMisspelled(word string, dict Dict) (res bool, suggested []string)
}

func And(checkers ...Checker) Checker {
	return intersectChecker{checkers}
}

func Or(checkers ...Checker) Checker {
	return unionChecker{checkers}
}

// A unionChecker takes a slice of checkers and classifies a word as
// misspelled if it's classified as such by at least one of the checkers.
type unionChecker struct{
	checkers []Checker
}

func (c unionChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	var suggestions = make(map[string]bool)
	var result bool
	for _, c := range(c.checkers) {
		if ok, sug := c.IsMisspelled(word, dict); ok {
			for s := range sug {
				suggestions[sug[s]] = true
			}
			result = true
		}
	}
	if result {
		sug := make([]string, 0, len(suggestions))
		for k := range suggestions {
			sug = append(sug, k)
		}
		return result, sug
	} else {
		return result, nil
	}
}

// An intersectChecker takes a slice of checkers and classifies a word as
// misspelled if it's classified as such by all of the checkers.
type intersectChecker struct {
	checkers []Checker
}

func (c intersectChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	var suggestions []string
	for _, c := range(c.checkers) {
		if ok, sug := c.IsMisspelled(word, dict); !ok {
			return false, nil
		} else {
			suggestions = sug
		}
	}
	return true, suggestions
}

// A StrictChecker classifies a word as misspelled if it does not exist in
// the dictionary provided.
type StrictChecker struct{}

func (StrictChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if dict.Contains([]rune(word)) {
		return false, nil
	}
	return true, nil
}

// A MinLengthChecker classifies a word as misspelled iff the length of the word
// is greater than a given value.
type MinLengthChecker uint8

func (c MinLengthChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if len(word) >= int(c) {
		return true, nil
	}
	return false, nil
}

// A DeltaChecker classifies a word as misspelled if its length is >= MinLength
// and is within a given number of character deletions, inserts, and consecutive
// swaps.
type DeltaChecker struct {
	AllowedIns   uint
	AllowedDel   uint
	AllowedSwaps uint
	AllowedMods uint
}

func (c DeltaChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	wordSlice := []rune(word)
	if dict.Contains(wordSlice) {
		return false, nil
	}

	for _, r := range(wordSlice) {
		if !dict.Alphabet().Contains(r) {
			return false, nil
		}
	}

	buf := make([]rune, uint(len(wordSlice))+c.AllowedIns)
	copy(buf, wordSlice)
	return c.isMisspelledDelta(buf, dict, uint(len(wordSlice)), c.AllowedIns, c.AllowedDel, c.AllowedSwaps, c.AllowedMods)
}

func (c *DeltaChecker) isMisspelledDelta(word []rune, dict Dict, length, ins, del, swaps, mods uint) (bool, []string) {
	if ins+del+swaps+mods == 0 {
		return !dict.Contains(word[:length]), nil
	}

	stack := make([]interface{}, ins+del+swaps+mods)
	depth := 0

	// 4 types of nodes: insertion, deletion, swap, modification
	type insertion struct {
		locationIndex int
		letterIndex int
	}
	type deletion struct {
		locationIndex int
		letter rune
	}
	type swap struct {
		index int
	}
	type mod struct{
		replaced rune
		locationIndex int
		letterIndex int
	}

	// push the nodes onto the stack
	for i := uint(0); i < ins; i++ {
		stack[depth] = &insertion{-1, -1}
		depth++
	}
	for d := uint(0); d < del; d++ {
		stack[depth] = &deletion{-1, -1}
		depth++
	}
	for s := uint(0); s < swaps; s++ {
		stack[depth] = &swap{-1}
		depth++
	}
	for m := uint(0); m < mods; m++ {
		stack[depth] = &mod{-1, -1, -1}
		depth++
	}

	// iterate through all paths with some paths visited multiple times based
	// on the input parameters.
	var curDepth = 0
	for curDepth >= 0 {
		var currentWord []rune
		// get the current type and go to the next state
		switch op := stack[curDepth].(type) {
		case *swap:
			if op.index == -1 {
				op.index = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.index > int(length)-2 {
				if op.index-1 >= 0 {
					word[op.index-1], word[op.index] = word[op.index], word[op.index-1]
				}
				*op = swap{-1}
				curDepth--;
				continue;
			} else {
				if op.index > 0 {
					word[op.index-1], word[op.index] = word[op.index], word[op.index-1]
				}
				word[op.index], word[op.index+1] = word[op.index+1], word[op.index]
				op.index++
				curDepth = min(curDepth+1, depth-1)
			}
		case *mod:
			if op.locationIndex == -1 {
				op.locationIndex = 0
				op.letterIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.locationIndex == int(length) {
				*op = mod{-1, -1, -1}
				curDepth--;
				continue;
			} else if op.letterIndex >= dict.Alphabet().Size() {
				word[op.locationIndex] = op.replaced
				op.locationIndex++
				op.letterIndex = 0
				op.replaced = -1
				continue;
			} else {
				if op.replaced == -1 {
					op.replaced = word[op.locationIndex]
				}
				modLetter, _ := dict.Alphabet().Letter(op.letterIndex)
				if modLetter == op.replaced {
					op.letterIndex++
					continue;
				} else {
					word[op.locationIndex] = modLetter
					curDepth = min(curDepth+1, depth-1)
					op.letterIndex++
				}
			}

		case *insertion:
			if op.locationIndex == -1 {
				op.locationIndex = 0
				op.letterIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.locationIndex == int(length) {
				*op = insertion{-1, -1}
				curDepth--;
				continue;
			} else if op.letterIndex >= dict.Alphabet().Size() {
				for i := op.locationIndex+1; i < int(length); i++ {
					word[i-1] = word[i]
				}
				length--
				op.locationIndex++
				op.letterIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else {
				if op.letterIndex == 0 {
					for i := int(length); i > op.locationIndex; i-- {
						word[i] = word[i-1]
					}
					length++
				}
				word[op.locationIndex], _ = dict.Alphabet().Letter(op.letterIndex)
				op.letterIndex++
				curDepth = min(curDepth+1, depth-1)
			}
		case *deletion:
			if op.locationIndex == -1 {
				if length == 0 {
					curDepth--
				} else {
					op.locationIndex = 0
					curDepth = min(curDepth+1, depth-1)
				}
				continue;
			} else if op.locationIndex > int(length) {
				word[length] = op.letter
				*op = deletion{-1, -1}
				length++
				curDepth--;
				continue;
			} else {
				if op.locationIndex == 0 {
					op.letter = word[op.locationIndex]
					for i := 0; i < int(length)-1; i++ {
						word[i] = word[i+1]
					}
					length--
				} else {
					newLetter := word[op.locationIndex-1]
					word[op.locationIndex-1] = op.letter
					op.letter = newLetter
				}
				op.locationIndex++
				curDepth = min(curDepth+1, depth-1)
			}
		}

		currentWord = word[:length]
		if dict.Contains(currentWord) {
			return true, []string{string(currentWord)}
		}
	}
	return false, nil
}
