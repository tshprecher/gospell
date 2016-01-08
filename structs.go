package main

import (
	"bytes"
	"fmt"
)

type Misspelling struct {
	Word string
	Line []int
}

type Result struct {
	Filename     string
	Misspellings []Misspelling
}

func (r Result) String() string {
	var buffer bytes.Buffer
	for _, misp := range r.Misspellings {
		for _, line := range misp.Line {
			// TODO: proper line number formatting?
			buffer.WriteString(fmt.Sprintf("%s:%d '%s'\n", r.Filename, line, misp.Word))
		}
	}
	return buffer.String()
}
