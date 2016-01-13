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
	for _, misp := range r.Misspellings {
		buffer.WriteString(fmt.Sprintf("%s:%d '%s'\n", r.Filename, misp.Line, misp.Word))
	}
	return buffer.String()
}
