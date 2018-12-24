package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

// newWord() shows how to build an iterator on string
// can make a stateful iterator by putting some variables before the return...

func newWord(wordBase string, uniqueWords int) func() string {
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
	var wcSlice []WC
	for word, count := range wc {
		wcSlice = append(wcSlice, WC{word, count})
	}
	sort.Slice(wcSlice, func(i, j int) bool {
		return wcSlice[i].count > wcSlice[j].count
	})
	fmt.Println("sorted words")
	if len(wcSlice) > showTop {
		wcSlice = wcSlice[:showTop]
	}
	for _, wc := range wcSlice {
		fmt.Printf("word: %s, count: %d\n", wc.word, wc.count)
	}

}

func driver(uniqueWords, lastNWords, wordsToGenerate, showTop int) {
	genWord := newWord("unique", uniqueWords)

	// Go has typed slices; equivalent of a FIFO
	var queue []string

	// Go has typed maps (must be allocated)
	var wc map[string]int
	wc = make(map[string]int)

	for i := 0; i < wordsToGenerate; i++ {
		word := genWord()
		if len(queue) >= lastNWords {
			queue = queue[1:]
			wc[word]--
			if wc[word] < 1 {
				delete(wc, word)
			}
		}
		queue = append(queue, word)
		wc[word]++
		showWordCounts(wc, showTop)
	}
}

func main() {
	uniqueWordsPtr := flag.Int("unique", 100, "number of unique words")
	lastNWordsPtr := flag.Int("last_n_words", 100, "last n words from current word (to count in word cloud)")
	wordsToGeneratePtr := flag.Int("generate", 1000000, "words to generate randomly")
	showTopPtr := flag.Int("show_top", 10, "show top n words")

	flag.Parse()
	fmt.Printf("unique: %d, last_n_words: %d, generate: %d, show_top: %d\n", *uniqueWordsPtr, *lastNWordsPtr, *wordsToGeneratePtr, *showTopPtr)
	driver(*uniqueWordsPtr, *lastNWordsPtr, *wordsToGeneratePtr, *showTopPtr)
}
