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
	ignoreCase                                     bool
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

func driver(config *WordCountConfig) {
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

func parseCommandLine() WordCountConfig {
	config :=  WordCountConfig{lastNWords: 1000, showTop: 10, minWordLength: 5, everySteps: 1000, ignoreCase: false}
	flag.IntVar(&config.lastNWords, "last-n-words", config.lastNWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.showTop, "show-top", config.showTop, "show top n words")
	flag.IntVar(&config.minWordLength, "min-word-length", config.minWordLength, "minimum word length")
	flag.IntVar(&config.everySteps, "every-steps", config.everySteps, "minimum word length")
	flag.BoolVar(&config.ignoreCase, "ignore-case", config.ignoreCase, "treat all words as upper case")
	flag.Parse()
	return config
}

func main() {
	config := parseCommandLine()
	driver(&config)
}
