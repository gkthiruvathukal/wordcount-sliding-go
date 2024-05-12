package main

import (
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	expected := 5
	q := NewCircularQueue[string](expected)
	if !q.IsEmpty() {
		t.Errorf("CircularQueue[string] not empty (size is %d; length is %d).", q.Size(), q.Length())
	}
	if q.Length() != expected {
		t.Errorf("CircularQueue[string] not %d (length is %d).", expected, q.Length())
	}
}

func TestFill(t *testing.T) {
	expected := 5
	q := NewCircularQueue[string](expected)
	for i := 0; i < expected-1; i++ {
		q.Enqueue(strconv.Itoa(i))
		if q.IsFull() {
			t.Errorf("CircularQueue[string] should not be full yet (has %d elements, length %d).", q.Size(), q.Length())
		}
	}
	q.Enqueue(strconv.Itoa(expected))
	if !q.IsFull() {
		t.Errorf("CircularQueue[string] did not reach expected capacity %d (size is %d).", expected, q.Size())
	}
}
