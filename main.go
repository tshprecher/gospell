// Package main runs this
// line two
// line three
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/tshprecher/gospell/check"
	"github.com/tshprecher/gospell/lang"
	"io/ioutil"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// tunable variables
	minLength = flag.Int("vL", 4, "filter out words less than 'l' characters")
	maxSwaps  = flag.Int("vS", 1, "correct spelling up to 's' consecutive character swaps")
	maxIns    = flag.Int("vI", 0, "correct spelling up to 'i' character insertions")
	maxDel    = flag.Int("vD", 0, "correct spelling up to 'd' character deletions")

	// flags for other languages (comments only)
	langC     = flag.Bool("c", false, "process C files")
	langCpp   = flag.Bool("cpp", false, "process C++ files")
	langJava  = flag.Bool("java", false, "process Java files")
	langScala  = flag.Bool("scala", false, "process Scala files")

	fileSuffix string = ".go" // TODO: adapt for multiple suffixes, e.g. ".cpp" and ".cc"
	fileset *token.FileSet = token.NewFileSet() // TODO: put this in process.go?
	checker check.Checker
	proc  processor
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func isInputFile(finfo os.FileInfo) bool {
	return !finfo.IsDir() &&
		!strings.HasPrefix(finfo.Name(), ".") &&
		strings.HasSuffix(finfo.Name(), fileSuffix)
}

func processFile(filename string, dict check.Dict) (*fileResult, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
	 	return nil, err
	}
	return proc.process(filename, src, dict)
}

func processDir(dir string, dict check.Dict) error {
	var fileFound = false
	var visit = func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isInputFile(f) {
			fileFound = true
//			fmt.Printf("debug: processing file %s\n", path)
			res, err := processFile(path, dict)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Print(res)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)
	if err != nil {
		return err
	}
	if !fileFound {
		fmt.Println("no input files found.")
	}
	return nil
}

func checkPositiveArg(value *int, arg string) {
	if *value < 0 {
		log.Fatalf("arg '%v' must be positive.", arg)
	}
}

// handleFlags parses the command line flags and sets the processor
// on success, returns an error otherwise.
func handleFlags() error {
	log.SetFlags(0) // do not prefix message with timestamp
	flag.Parse()
	checkPositiveArg(minLength, "l")
	checkPositiveArg(maxSwaps, "s")
	checkPositiveArg(maxIns, "i")
	checkPositiveArg(maxDel, "d")


	proc = goProcessor{}
	fileSuffix = ".go"

	var altLang uint8 = 0
	var langs = []struct{
		isSet bool
		suffix string
	}{
		{*langC, ".c"},
		{*langCpp, ".cc"},
		{*langJava, ".java"},
		{*langScala, ".scala"},
	}

	for _, l := range(langs) {
		if l.isSet {
			altLang++
			fileSuffix = l.suffix
		}
	}

	if altLang > 1 {
		return errors.New("cannot specify multiple languages.")
	}
	if altLang > 0 {
		proc = cStyleCommentProcessor{}
	}
	return nil
}


func main() {
	checkError(handleFlags())

	checker = &check.DeltaChecker{
		MinLength:    *minLength,
		AllowedIns:   *maxIns,
		AllowedDel:   *maxDel,
		AllowedSwaps: *maxSwaps}

	for a := 0; a < flag.NArg(); a++ {
		filename := flag.Arg(a)
		fileInfo, err := os.Stat(filename)
		checkError(err)
		// TODO: define a Language struct encapsulating alphabet and words
		dict, err := check.NewTrie(lang.Words, lang.English)
		checkError(err)

		if fileInfo.IsDir() {
			err = processDir(filename, dict)
			checkError(err)
		} else {
			res, err := processFile(filename, dict)
			checkError(err)
			fmt.Print(res)
		}
	}
}
