package cq

import (
	"errors"
	"log"
)

type CircularQueue[T comparable] struct {
	queue                 []T
	storePos, retrievePos int
	count                 int
	default_zero          T
}

func NewCircularQueue[T comparable](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		queue: make([]T, size),
	}
}

func (cq *CircularQueue[T]) Show() {
	log.Printf("storePos = %d, retrievePos = %d, queue = %v\n", cq.storePos, cq.retrievePos, cq.queue)
}

func (cq *CircularQueue[T]) Enqueue(s T) error {
	if cq.IsFull() {
		return errors.New("queue is full")
	}
	cq.queue[cq.storePos] = s
	cq.storePos = (cq.storePos + 1) % len(cq.queue)
	cq.count++
	return nil
}

func (cq *CircularQueue[T]) Dequeue() (T, error) {
	if cq.IsEmpty() {
		return cq.default_zero, errors.New("queue is empty")
	}
	item := cq.queue[cq.retrievePos]
	cq.queue[cq.retrievePos] = cq.default_zero
	cq.retrievePos = (cq.retrievePos + 1) % len(cq.queue)
	cq.count--
	return item, nil
}

func (cq *CircularQueue[T]) Peek() (T, error) {
	if cq.IsEmpty() {
		return cq.default_zero, errors.New("queue is empty")
	}
	return cq.queue[cq.retrievePos], nil
}

func (cq *CircularQueue[T]) Clear() {
	cq.storePos = 0
	cq.retrievePos = 0
	cq.count = 0

	for i := range cq.queue {
		cq.queue[i] = cq.default_zero
	}
}

func (cq *CircularQueue[T]) IsFull() bool {
	return cq.count == len(cq.queue)
}

func (cq *CircularQueue[T]) IsEmpty() bool {
	return cq.count == 0
}

func (cq *CircularQueue[T]) Size() int {
	return cq.count
}

func (cq *CircularQueue[T]) Capacity() int {
	return len(cq.queue)
}

func (cq *CircularQueue[T]) Length() int {
	return len(cq.queue)
}
