package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/tshprecher/gospell/check"
	"github.com/tshprecher/gospell/lang"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultMinLength = 5
	defaultMaxSwaps  = 1
	defaultMaxIns    = 0
	defaultMaxDel    = 0
)

var (
	minLength = flag.Int("ml", defaultMinLength, "filter out words less than 'ml' characters")
	maxSwaps  = flag.Int("ms", defaultMaxSwaps, "correct spelling up to 'ms' consecutive character swaps")
	maxIns    = flag.Int("mi", defaultMaxIns, "correct spelling up to 'mi' character insertions")
	maxDel    = flag.Int("md", defaultMaxDel, "correct spelling up to 'md' character deletions")

	// flags for other languages (comments only)
	langC     = flag.Bool("c", false, "process C files")
	langCpp   = flag.Bool("cpp", false, "process C++ files")
	langJava  = flag.Bool("java", false, "process Java files")
	langScala = flag.Bool("scala", false, "process Scala files")

	fileRegexp *regexp.Regexp
	checker    check.Checker = check.UnionChecker{
		[]check.Checker{
			check.DeltaChecker{5, 0, 0, 1},
			check.DeltaChecker{5, 0, 1, 0},
			check.DeltaChecker{5, 1, 0, 0},
		},
	}
	proc processor = cStyleCommentProcessor{checker}
)

func checkError(err error) {
	if err != nil {
		log.Fatal("error: " + err.Error())
	}
}

func isInputFile(finfo os.FileInfo) bool {
	return !finfo.IsDir() &&
		!strings.HasPrefix(finfo.Name(), ".") &&
		fileRegexp.MatchString(finfo.Name())
}

func processFile(filename string, dict check.Dict) (*fileResult, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return proc.run(filename, src, dict)
}

func processDir(dir string, dict check.Dict) error {
	var fileFound = false
	var visit = func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isInputFile(f) {
			fileFound = true
			res, err := processFile(path, dict)
			if err != nil {
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
		log.Fatalf("error: arg '%v' must be positive.", arg)
	}
}

// handleFlags parses the command line flags and sets the processor
// on success, returns an error otherwise.
func handleFlags() error {
	log.SetFlags(0) // do not prefix message with timestamp
	flag.Parse()

	checkPositiveArg(minLength, "ml")
	checkPositiveArg(maxSwaps, "ms")
	checkPositiveArg(maxIns, "mi")
	checkPositiveArg(maxDel, "md")

	if *minLength != defaultMinLength ||
		*maxSwaps != defaultMaxSwaps ||
		*maxIns != defaultMaxIns ||
		*maxDel != defaultMaxDel {
		// if any of the default value are set, override
		// the default checker
		checker = check.DeltaChecker{
			*minLength,
			*maxIns,
			*maxDel,
			*maxSwaps}
	}

	r, _ := regexp.Compile(".*\\.go$")
	fileRegexp = r

	var altLang uint8 = 0
	var langs = []struct {
		isSet   bool
		pattern string
	}{
		{*langC, ".*\\.c$"},
		{*langCpp, "(.*\\.cc$)|(.*\\.cpp$)"},
		{*langJava, ".*\\.java$"},
		{*langScala, ".*\\.scala$"},
	}

	for _, l := range langs {
		if l.isSet {
			altLang++
			r, _ := regexp.Compile(l.pattern)
			fileRegexp = r
		}
	}

	if altLang > 1 {
		return errors.New("cannot specify multiple languages.")
	}
	if altLang > 0 {
		proc = cStyleCommentProcessor{checker}
	}
	return nil
}

func main() {
	checkError(handleFlags())

	if flag.NArg() == 0 {
		log.Fatal("error: no files found.")
	}

	for a := 0; a < flag.NArg(); a++ {
		filename := flag.Arg(a)
		fileInfo, err := os.Stat(filename)
		checkError(err)
		dict := check.NewTrie(lang.EnglishAlphabet)
		for _, word := range lang.EnglishUsWords {
			dict.Add(word.Word)
		}
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
