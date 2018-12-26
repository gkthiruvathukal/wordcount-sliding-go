package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type WordCountConfig struct {
	uniqueWords, lastNWords, wordsToGenerate, showTop int
}

func wordGenerator(wordBase string, uniqueWords int) func() string {
	return func() string {
		n := rand.Intn(uniqueWords)
		return wordBase + strconv.Itoa(n)
	}
}

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
	nextWord := wordGenerator("unique", config.uniqueWords)
	queue := new(CQueueString)
	queue.init(config.lastNWords)
	wc := make(map[string]int)

	for i := 0; i < config.wordsToGenerate; i++ {
		word := nextWord()
		//queue.show()
		if queue.isFull() {
			_, droppedWord := queue.remove()
			wc[droppedWord]--
			if wc[droppedWord] <= 0 {
				delete(wc, droppedWord)
			}
		}
		queue.add(word)
		wc[word]++
		showWordCounts(wc, config.showTop)
	}
}

func main() {
	config := WordCountConfig{}
	flag.IntVar(&config.uniqueWords, "unique", config.uniqueWords, "number of unique words")
	flag.IntVar(&config.lastNWords, "last_n_words", config.uniqueWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.wordsToGenerate, "generate", config.wordsToGenerate, "words to generate randomly")
	flag.IntVar(&config.showTop, "show_top", config.showTop, "show top n words")
	flag.Parse()
	driver(&config)
}
