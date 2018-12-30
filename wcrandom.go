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

func driver(config *WordCountConfig) {
	regex := regexp.MustCompile(`\w+`)
	scanner := bufio.NewScanner(os.Stdin)
	queue := new(CQueueString)
	queue.init(config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0

	for scanner.Scan() {
		text := scanner.Text()
		matches := regex.FindAllString(text, -1)
		for _, word := range matches {
			if queue.isFull() {
				_, droppedWord := queue.remove()
				wc[droppedWord]--
				if wc[droppedWord] <= 0 {
					delete(wc, droppedWord)
				}
			}
			wordPosition++
			queue.add(word)
			// the minimum word test applies to counting only, not to last N
			storeWord := word
			if config.ignoreCase {
				storeWord = strings.ToUpper(word)
			}
			if len(word) >= config.minWordLength {
				wc[storeWord]++
			}
			if wordPosition%config.everySteps == 0 {
				fmt.Printf("%d: ", wordPosition)
				showWordCounts(wc, config.showTop)
			}
		}
	}
}

func main() {
	config := WordCountConfig{1000, 10, 5, 1000, false}

	flag.IntVar(&config.lastNWords, "last-n-words", config.lastNWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.showTop, "show-top", config.showTop, "show top n words")
	flag.IntVar(&config.minWordLength, "min-word-length", config.minWordLength, "minimum word length")
	flag.IntVar(&config.everySteps, "every-steps", config.everySteps, "minimum word length")
	flag.BoolVar(&config.ignoreCase, "ignore-case", config.ignoreCase, "treat all words as upper case")
	flag.Parse()
	driver(&config)
}
