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

// A StrictChecker classifies a word as misspelled if it does not exist in
// the dictionary provided.
type StrictChecker struct{}

func (StrictChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if dict.Contains([]rune(word)) {
		return false, nil
	}
	return true, nil
}

// A UnionChecker takes a slice of checkers and classifies a word as
// misspelled if it's classified as such by at least one of the checkers.
type UnionChecker struct{
	Checkers []Checker
}

func (c UnionChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	var suggestions = make(map[string]bool)
	var result bool
	for _, c := range(c.Checkers) {
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

// A DeltaChecker classifies a word as misspelled if its length is >= MinLength
// and is within a given number of character deletions, inserts, and consecutive
// swaps.
type DeltaChecker struct {
	MinLength    int
	AllowedIns   int
	AllowedDel   int
	AllowedSwaps int
}

func (dc DeltaChecker) IsMisspelled(word string, dict Dict) (bool, []string) {
	if len(word) < dc.MinLength {
		return false, nil
	}
	wordSlice := []rune(word)
	if dict.Contains(wordSlice) {
		return false, nil
	}
	if dc.AllowedIns+dc.AllowedDel+dc.AllowedSwaps == 0 {
		return !dict.Contains(wordSlice), nil
	}
	for _, r := range(wordSlice) {
		if !dict.Alphabet().Contains(r) {
			return false, nil
		}
	}

	buf := make([]rune, len(wordSlice)+dc.AllowedIns)
	copy(buf, wordSlice)
	return dc.isMisspelledDelta(buf, dict, len(wordSlice), dc.AllowedIns, dc.AllowedDel, dc.AllowedSwaps)
}

func (dc *DeltaChecker) isMisspelledDelta(word []rune, dict Dict, len, ins, del, swaps int) (bool, []string) {
	return dc.isMisspelledDeltaIter(word, dict, len, ins, del, swaps)
}

func (dc *DeltaChecker) isMisspelledDeltaIter(word []rune, dict Dict, length, ins, del, swaps int) (bool, []string) {
	// fmt.Printf("checking word '%s'\n", string(word))
	stack := make([]interface{}, ins+del+swaps)
	depth := 0

	// 3 types of nodes: insertion, deletion, swap
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

	// push the nodes onto the stack
	for i := 0; i < ins; i++ {
		stack[depth] = &insertion{-1, -1}
		depth++
	}
	for d := 0; d < del; d++ {
		stack[depth] = &deletion{-1, -1}
		depth++
	}
	for s := 0; s < swaps; s++ {
		stack[depth] = &swap{-1}
		depth++
	}

	// iterate through all paths with some paths visited multiple times based
	// on the input parameters.
	var loopRan = false
	var curDepth = 0
	for curDepth >= 0 || !loopRan {
		// fmt.Printf("current depth: %d\n", curDepth)
		loopRan = true
		var currentWord []rune
		// get the current type and go to the next state
		switch op := stack[curDepth].(type) {
		case *swap:
			if op.index == -1 {
				op.index = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.index > length-2 {
				if op.index-1 >= 0 {
					word[op.index-1], word[op.index] = word[op.index], word[op.index-1]
				}
				*op = swap{-1}
				curDepth--;
				continue;
			} else {
				// fmt.Printf("current swap index: %d\n", op.index)
				if op.index > 0 {
					word[op.index-1], word[op.index] = word[op.index], word[op.index-1]
				}
				word[op.index], word[op.index+1] = word[op.index+1], word[op.index]
				// fmt.Printf("current swap length: %d\n", length)
				// fmt.Printf("current swap: %v\n", *op)
				op.index++
				curDepth = min(curDepth+1, depth-1)
			}
		case *insertion:
			if op.locationIndex == -1 {
				op.locationIndex = 0
				op.letterIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.locationIndex == length {
				*op = insertion{-1, -1}
				curDepth--;
				continue;
			} else if op.letterIndex >= dict.Alphabet().Size() {
				for i := op.locationIndex+1; i < length; i++ {
					word[i-1] = word[i]
				}
				length--
				op.locationIndex++
				op.letterIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else {
				// fmt.Printf("current insertion length: %d\n", length)
				// fmt.Printf("current insertion index: %d\n", op.letterIndex)
				if op.letterIndex == 0 {
					for i := length; i > op.locationIndex; i-- {
						// fmt.Printf("current insertion shift index: %d\n", i)
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
				op.locationIndex = 0
				curDepth = min(curDepth+1, depth-1)
				continue;
			} else if op.locationIndex > length {
				word[length] = op.letter
				*op = deletion{-1, -1}
				length++
				curDepth--;
				continue;
			} else {
				// fmt.Printf("current deletion length: %d\n", length)
				// fmt.Printf("current deletion index: %d\n", op.locationIndex)
				if op.locationIndex == 0 {
					op.letter = word[op.locationIndex]
					for i := 0; i < length-1; i++ {
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
		// fmt.Printf("current word: %s\n", string(currentWord))
		if dict.Contains(currentWord) {
			return true, []string{string(currentWord)}
		}
	}
	return false, nil
}
