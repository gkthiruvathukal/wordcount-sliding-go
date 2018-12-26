package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// WordCountConfig The configuration options available.
type WordCountConfig struct {
	uniqueWords, lastNWords, wordsToGenerate, showTop int
}

// newWord() shows how to build an iterator on string
// can make a stateful iterator by putting some variables before the return...

func wordGenerator(wordBase string, uniqueWords int) func() string {
	return func() string {
		n := rand.Intn(uniqueWords)
		return wordBase + strconv.Itoa(n)
	}
}

func showWordCounts(wc map[string]int, showTop int) {
	type WC struct {
		word  string
		count int
	}

	wcSlice := make([]WC, len(wc))
	wcSlice = wcSlice[:0]
	for word, count := range wc {
		wcSlice = append(wcSlice, WC{word, count})
	}
	sort.Slice(wcSlice, func(i, j int) bool {
		return wcSlice[i].count > wcSlice[j].count
	})
	if len(wcSlice) > showTop {
		wcSlice = wcSlice[:showTop]
	}

	//fmt.Printf("wc len %d, show top %d, slice len %d, slice cap %d", len(wc), showTop, len(wcSlice), cap(wcSlice))
	// printing words...

	pretty := make([]string, len(wcSlice))
	pretty = pretty[:0]
	for _, wc := range wcSlice {
		pretty = append(pretty, fmt.Sprintf("%s: %d", wc.word, wc.count))
	}
	fmt.Printf("words { %s }\n", strings.Join(pretty, ", "))
}

// CQueueString is a circular queue
type CQueueString struct {
	queue                 []string
	storePos, retrievePos int
	count                 int
}

func (cq *CQueueString) init(size int) {
	cq.queue = make([]string, size)
	cq.storePos = 0
	cq.retrievePos = 0
	cq.count = 0
}
func (cq *CQueueString) add(s string) int {
	if cq.count > len(cq.queue) {
		return -1
	}
	cq.queue[cq.storePos] = s
	cq.storePos = (cq.storePos + 1) % len(cq.queue)
	cq.count++
	return cq.count
}

func (cq *CQueueString) remove() (int, string) {
	if cq.count <= 0 {
		return -1, ""
	}
	item := cq.queue[cq.retrievePos]
	cq.retrievePos = (cq.retrievePos + 1) % len(cq.queue)
	cq.count--
	return cq.count, item
}

func (cq *CQueueString) isFull(s string) bool {
	return len(cq.queue) == cq.count
}

func driver(config *WordCountConfig) {
	nextWord := wordGenerator("unique", config.uniqueWords)
	queue := make([]string, config.lastNWords)
	wc := make(map[string]int)

	queue = queue[:0]
	for i := 0; i < config.wordsToGenerate; i++ {
		word := nextWord()
		//fmt.Printf("queue length %d\n", len(queue))
		if len(queue) >= config.lastNWords {
			droppedWord := queue[0]
			queue = queue[1:]
			fmt.Printf("queue len = %d, cap = %d", len(queue), cap(queue))
			wc[droppedWord]--
			if wc[droppedWord] <= 0 {
				delete(wc, droppedWord)
			}
		}
		queue = append(queue, word)
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
	fmt.Println(config)
	driver(&config)
}
