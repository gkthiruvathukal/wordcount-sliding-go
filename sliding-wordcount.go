package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// WordCountConfig is used for CLI arguments
type WordCountConfig struct {
	lastNWords, showTop, minWordLength, everySteps int
	ignoreCase, idiomatic                          bool
}

// WC is used for sorting at presentation layer (top N words in word cloud)
type WC struct {
	word  string
	count int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func showWordCounts(wc map[string]int, showTop int) {

	// wcData is a slice of (word, count) pairs (as struct WC)
	wcData := make([]WC, len(wc))
	i := 0
	for word, count := range wc {
		wcData[i] = WC{word, count}
		i++
	}

	// sort.Slice() sorts by count field, descending
	sort.Slice(wcData, func(i, j int) bool {
		return wcData[i].count > wcData[j].count
	})

	// create slice of strings for presentation as JSON like output
	numberToPrint := min(len(wcData), showTop)
	pretty := make([]string, numberToPrint)

	for i := 0; i < numberToPrint; i++ {
		pretty[i] = fmt.Sprintf("%s: %d", wcData[i].word, wcData[i].count)
	}
	fmt.Printf("words { %s }\n", strings.Join(pretty, ", "))
}

func wordUp(wordCloud map[string]int, addWord string) {
	wordCloud[addWord]++
}

func wordDown(wordCloud map[string]int, dropWord string) {
	wordCloud[dropWord]--
	if wordCloud[dropWord] <= 0 {
		delete(wordCloud, dropWord)
	}
}

// The imperative pipline is a more traditional way of writing the solution
// It is monolithic, but the code is still fairly intuitive.

func imperativePipeline(config *WordCountConfig) {
	regex := regexp.MustCompile(`\p{L}+`)
	scanner := bufio.NewScanner(os.Stdin)
	queue := NewCircularQueue[string](config.lastNWords)
	//#queue.init(config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0
	for scanner.Scan() {
		text := scanner.Text()
		matches := regex.FindAllString(text, -1)
		for _, word := range matches {
			// ignore words below the minimum length altogether
			if len([]rune(word)) >= config.minWordLength {
				if config.ignoreCase {
					word = strings.ToLower(word)
				}
				wordUp(wc, word)
				if queue.IsFull() {
					droppedWord, _ := queue.Dequeue()
					wordDown(wc, droppedWord)
				}
				wordPosition++
				queue.Enqueue(word)
				if wordPosition%config.everySteps == 0 {
					fmt.Printf("%d: ", wordPosition)
					showWordCounts(wc, config.showTop)
				}
			}
		}
	}
}

// The "idiomatic" pipeline uses go routines and channels to suggest a more functional style, similar to Scala and others.
// Note that Go does not support most of FP and nevertheless provides a delightfully composable approach.

func goIdiomaticPipeline(config *WordCountConfig) {
	words := generateWords(config)
	filteredWords := filterBasedOnCommandLine(config, words)
	slidingAnalysis(config, filteredWords)
}

func generateWords(config *WordCountConfig) <-chan string {
	regex := regexp.MustCompile(`\p{L}+`)
	scanner := bufio.NewScanner(os.Stdin)

	out := make(chan string)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			matches := regex.FindAllString(text, -1)
			for _, word := range matches {
				out <- word
			}
		}
		close(out)
	}()
	return out
}

func filterBasedOnCommandLine(config *WordCountConfig, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for word := range in {
			newWord := word
			if len([]rune(word)) < config.minWordLength {
				continue
			}
			if config.ignoreCase {
				newWord = strings.ToLower(newWord)
			}
			out <- newWord
		}
		close(out)
	}()
	return out
}

// last stage of pipeline

func slidingAnalysis(config *WordCountConfig, in <-chan string) {
	queue := NewCircularQueue[string](config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0
	for word := range in {
		wordUp(wc, word)
		if queue.IsFull() {
			droppedWord, _ := queue.Dequeue()
			wordDown(wc, droppedWord)
		}
		wordPosition++
		queue.Enqueue(word)
		if wordPosition%config.everySteps == 0 {
			fmt.Printf("%d: ", wordPosition)
			showWordCounts(wc, config.showTop)
		}
	}
}

// We will have one main that can select the versionw with or without go-routines.
func parseCommandLine() *WordCountConfig {
	config := WordCountConfig{lastNWords: 1000, showTop: 10, minWordLength: 5, everySteps: 1000, ignoreCase: false, idiomatic: false}
	flag.IntVar(&config.lastNWords, "last-n-words", config.lastNWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.showTop, "show-top", config.showTop, "show top n words")
	flag.IntVar(&config.minWordLength, "min-word-length", config.minWordLength, "minimum word length")
	flag.IntVar(&config.everySteps, "every-steps", config.everySteps, "minimum word length")
	flag.BoolVar(&config.ignoreCase, "ignore-case", config.ignoreCase, "treat all words as upper case")
	flag.BoolVar(&config.idiomatic, "idiomatic", config.ignoreCase, "use Go routines to emulate a functional-style pipeline")
	flag.Parse()
	return &config
}

func driver(config *WordCountConfig) {
	if config.idiomatic {
		fmt.Println("Running driver with idiomatic")
		goIdiomaticPipeline(config)
	} else {
		fmt.Println("Running driver with imperative")
		imperativePipeline(config)
	}
}

func main() {
	config := parseCommandLine()
	driver(config)
}
