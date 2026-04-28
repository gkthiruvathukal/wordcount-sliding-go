package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"cq"
)

type WordCountConfig struct {
	lastNWords, showTop, minWordLength, everySteps int
	ignoreCase, goRoutines                         bool
	channelSize                                    int
	circularQueue, channelQueue                    bool
}

type WC struct {
	word  string
	count int
}

func showWordCounts(wc map[string]int, showTop int) {
	wcData := make([]WC, len(wc))
	i := 0
	for word, count := range wc {
		wcData[i] = WC{word, count}
		i++
	}

	sort.Slice(wcData, func(i, j int) bool {
		return wcData[i].count > wcData[j].count
	})

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

func imperativePipelineCQ(config *WordCountConfig) {
	regex := regexp.MustCompile(`\p{L}+`)
	scanner := bufio.NewScanner(os.Stdin)
	queue := cq.NewCircularQueue[string](config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0
	for scanner.Scan() {
		text := scanner.Text()
		matches := regex.FindAllString(text, -1)
		for _, word := range matches {
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

func slidingAnalysisCQ(config *WordCountConfig, in <-chan string) {
	queue := cq.NewCircularQueue[string](config.lastNWords)
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

func imperativePipelineChannel(config *WordCountConfig) {
	regex := regexp.MustCompile(`\p{L}+`)
	scanner := bufio.NewScanner(os.Stdin)
	window := make(chan string, config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0
	for scanner.Scan() {
		text := scanner.Text()
		matches := regex.FindAllString(text, -1)
		for _, word := range matches {
			if len([]rune(word)) >= config.minWordLength {
				if config.ignoreCase {
					word = strings.ToLower(word)
				}
				wordUp(wc, word)
				if len(window) == cap(window) {
					droppedWord := <-window
					wordDown(wc, droppedWord)
				}
				wordPosition++
				window <- word
				if wordPosition%config.everySteps == 0 {
					fmt.Printf("%d: ", wordPosition)
					showWordCounts(wc, config.showTop)
				}
			}
		}
	}
}

func slidingAnalysisChannel(config *WordCountConfig, in <-chan string) {
	window := make(chan string, config.lastNWords)
	wc := make(map[string]int)
	wordPosition := 0
	for word := range in {
		wordUp(wc, word)
		if len(window) == cap(window) {
			droppedWord := <-window
			wordDown(wc, droppedWord)
		}
		wordPosition++
		window <- word
		if wordPosition%config.everySteps == 0 {
			fmt.Printf("%d: ", wordPosition)
			showWordCounts(wc, config.showTop)
		}
	}
}

func generateWords(config *WordCountConfig) <-chan string {
	regex := regexp.MustCompile(`\p{L}+`)
	scanner := bufio.NewScanner(os.Stdin)

	out := make(chan string, config.channelSize)
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
	out := make(chan string, config.channelSize)
	go func() {
		for word := range in {
			if len([]rune(word)) < config.minWordLength {
				continue
			}
			if config.ignoreCase {
				word = strings.ToLower(word)
			}
			out <- word
		}
		close(out)
	}()
	return out
}

func driver(config *WordCountConfig) {
	switch {
	case !config.goRoutines && config.circularQueue:
		fmt.Println("Imperative + circular queue")
		imperativePipelineCQ(config)
	case !config.goRoutines && config.channelQueue:
		fmt.Println("Imperative + channel window")
		imperativePipelineChannel(config)
	case config.goRoutines && config.circularQueue:
		fmt.Println("Goroutine pipeline + circular queue")
		words := generateWords(config)
		filtered := filterBasedOnCommandLine(config, words)
		slidingAnalysisCQ(config, filtered)
	case config.goRoutines && config.channelQueue:
		fmt.Println("Goroutine pipeline + channel window")
		words := generateWords(config)
		filtered := filterBasedOnCommandLine(config, words)
		slidingAnalysisChannel(config, filtered)
	}
}

func parseCommandLine() *WordCountConfig {
	config := WordCountConfig{
		lastNWords:    1000,
		showTop:       10,
		minWordLength: 5,
		everySteps:    1000,
		ignoreCase:    false,
		goRoutines:    false,
		channelSize:   10,
		circularQueue: false,
		channelQueue:  false,
	}
	flag.IntVar(&config.lastNWords, "last-n-words", config.lastNWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.showTop, "show-top", config.showTop, "show top n words")
	flag.IntVar(&config.minWordLength, "min-word-length", config.minWordLength, "minimum word length")
	flag.IntVar(&config.everySteps, "every-steps", config.everySteps, "print word cloud every N words")
	flag.BoolVar(&config.ignoreCase, "ignore-case", config.ignoreCase, "treat all words as lower case")
	flag.BoolVar(&config.goRoutines, "go-routines", config.goRoutines, "use goroutine/channel pipeline")
	flag.IntVar(&config.channelSize, "channel-size", config.channelSize, "pipeline channel buffer size (goroutine mode)")
	flag.BoolVar(&config.circularQueue, "circular-queue", config.circularQueue, "use circular queue for sliding window (default)")
	flag.BoolVar(&config.channelQueue, "channel", config.channelQueue, "use buffered channel for sliding window")
	flag.Parse()

	if config.circularQueue && config.channelQueue {
		log.Fatal("error: -circular-queue and -channel are mutually exclusive")
	}
	if !config.circularQueue && !config.channelQueue {
		config.circularQueue = true
	}
	return &config
}

func main() {
	config := parseCommandLine()
	driver(config)
}
