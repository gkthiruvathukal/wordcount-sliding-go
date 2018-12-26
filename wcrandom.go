package slidingwordcount

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

	pretty := make([]string, len(wcSlice))
	pretty = pretty[:0]
	for _, wc := range wcSlice {
		pretty = append(pretty, fmt.Sprintf("%s: %d", wc.word, wc.count))
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
	fmt.Println(config)
	driver(&config)
}
