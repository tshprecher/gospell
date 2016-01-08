package main

import (
	"bytes"
	"fmt"
)

type Misspelling struct {
	Word string
	Line int
}

type Result struct {
	Filename     string
	Misspellings []Misspelling
}

func (r Result) String() string {
	var buffer bytes.Buffer
	for m, misp := range r.Misspellings {
		// TODO: proper line number formatting?
		buffer.WriteString(fmt.Sprintf("%s:%d '%s'", r.Filename, misp.Line, misp.Word))
		if m != len(r.Misspellings)-1 {
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}
