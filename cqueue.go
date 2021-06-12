package main

import "fmt"

// CQueueString is a circular queue

type CQueueString struct {
	queue                 []string
	storePos, retrievePos int
	count                 int
}

func (cq *CQueueString) show() {
	fmt.Printf("storePos = %d, retrievePos = %d, queue = ", cq.storePos, cq.retrievePos)
	fmt.Println(cq.queue)
}

func (cq *CQueueString) init(size int) {
	cq.queue = make([]string, size)
	cq.storePos = 0
	cq.retrievePos = 0
	cq.count = 0
}
func (cq *CQueueString) add(s string) int {
	if cq.isFull() {
		return -1
	} else {
		cq.queue[cq.storePos] = s
		cq.storePos = (cq.storePos + 1) % len(cq.queue)
		cq.count++
		return cq.count
	}
}

func (cq *CQueueString) remove() (int, string) {
	if cq.isEmpty() {
		return -1, ""
	} else {
		item := cq.queue[cq.retrievePos]
		cq.retrievePos = (cq.retrievePos + 1) % len(cq.queue)
		cq.count--
		return cq.count, item
	}
}

func (cq *CQueueString) isFull() bool {
	return cq.count == len(cq.queue)
}

func (cq *CQueueString) isEmpty() bool {
	return cq.count == 0
}

func (cq *CQueueString) size() int {
	return cq.count
}

func (cq *CQueueString) length() int {
	return len(cq.queue)
}
