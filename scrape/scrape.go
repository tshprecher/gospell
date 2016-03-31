package scrape

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	mwDefinitionEndpointFormat = "http://www.merriam-webster.com/dictionary/%s"
	mwPopularityEndpointFormat = "http://stats.merriam-webster.com/pop-score-redesign.php?word=%s&t=%d&id=popularity-score"
)

var (
	definitionNotFound = errors.New("definition not found")
	bufStdout          = bufio.NewWriter(os.Stdout)
	outMutex           sync.Mutex
)

func getWordFromBody(body string) (string, error) {
	lines := strings.Split(string(body), "\n")
	for _, l := range lines {
		if strings.Contains(l, "<title>") {
			word := strings.Split(l, "<title>")[1]
			word = strings.Split(word, "|")[0]
			word = strings.TrimSpace(strings.ToLower(word))
			return word, nil
		}
	}
	return "", definitionNotFound
}

func getPopularityFromBody(body string) (score int, label string, err error) {
	lines := strings.Split(string(body), "\n")
	for _, l := range lines {
		if strings.Contains(l, "label:") {
			label = strings.TrimSuffix(strings.SplitAfter(l, ":")[1], ",")
		}
		if strings.Contains(l, "score:") {
			score, err = strconv.Atoi(strings.TrimSuffix(strings.SplitAfter(l, " ")[1], ","))
			if err != nil {
				return
			}
		}
	}
	return
}

func lookupWord(word string) (statusCode int, mappedWord string, score int, label string, err error) {
	escWord := url.QueryEscape(word)
	defEp := fmt.Sprintf(mwDefinitionEndpointFormat, escWord)

	// get the definition
	defResp, er := http.Get(defEp)
	if er != nil {
		err = er
		return
	}
	defer defResp.Body.Close()

	statusCode = defResp.StatusCode
	if statusCode != 200 {
		return
	}

	defBody, _ := ioutil.ReadAll(defResp.Body)
	mappedWord, er = getWordFromBody(string(defBody))

	// word exists, so get its popularity
	popEp := fmt.Sprintf(mwPopularityEndpointFormat, mappedWord, time.Now().Nanosecond()/1000)
	popResp, er := http.Get(popEp)
	if er != nil {
		err = er
		return
	}
	defer popResp.Body.Close()

	popBody, _ := ioutil.ReadAll(popResp.Body)
	score, label, err = getPopularityFromBody(string(popBody))
	return
}

func exitError(message string) {
	fmt.Fprintf(os.Stderr, message)
	os.Exit(1)
}

func lookupLoop(outputPrefix string, words <-chan string) {
	for word := range words {
		var response string
		st, wrd, scr, lbl, er := lookupWord(word)
		if er != nil || st != 200 {
			response = fmt.Sprintf("%s\t%d\t%s\t%s\terror: %v\n", outputPrefix, st, word, wrd, er)
		} else {
			response = fmt.Sprintf("%s\t%d\t%s\t%s\t%d\t%s\n", outputPrefix, st, word, wrd, scr, lbl)
		}
		outMutex.Lock()
		bufStdout.WriteString(response)
		outMutex.Unlock()
	}
}

func Main() {
	runtime.GOMAXPROCS(30)
	inputWords := make(chan string, 1000)

	var wg sync.WaitGroup

	// multiple lookup goroutines
	for nl := 1; nl <= 25; nl++ {
		wg.Add(1)
		prefixString := fmt.Sprintf("loop%d", nl)
		go func() {
			defer wg.Done()
			lookupLoop(prefixString, inputWords)
		}()
	}

	// main goroutine pushes all the words to the input channel
	byts, _ := ioutil.ReadAll(os.Stdin)
	buffer := bytes.NewBuffer(byts)
	word, err := buffer.ReadString('\n')
	for err != io.EOF {
		word = strings.Trim(word, " \n\t")
		if err != nil {
			exitError(err.Error())
		}
		inputWords <- word
		word, err = buffer.ReadString('\n')
	}
	close(inputWords)
	wg.Wait()
	bufStdout.Flush()
}
