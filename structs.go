package main

import (
	"go/token"
)

type Misspelling struct {
	Word string
	Positions []token.Pos
}

type Result struct {
	Filename string
	Misspellings []Misspelling
}
