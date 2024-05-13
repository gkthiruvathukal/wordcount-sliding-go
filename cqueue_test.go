package main

import (
	"testing"
)

// newTestQueue creates a new instance of CircularQueue[int] with specified capacity.
func newTestQueue(capacity int) *CircularQueue[int] {
	q := NewCircularQueue[int](capacity)
	return q
}

func TestNewCircularQueue(t *testing.T) {
	q := newTestQueue(5)
	if q == nil {
		t.Error("NewCircularQueue() failed, got nil")
	}
	if len(q.queue) != 5 {
		t.Errorf("Expected queue length 5, got %d", len(q.queue))
	}
}

func TestEnqueue(t *testing.T) {
	q := newTestQueue(3)
	if err := q.Enqueue(1); err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}
	if err := q.Enqueue(2); err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}
	if err := q.Enqueue(3); err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}
	if err := q.Enqueue(4); err == nil {
		t.Error("Enqueue should fail when queue is full")
	}
}

func TestDequeue(t *testing.T) {
	q := newTestQueue(3)
	if _, err := q.Dequeue(); err == nil {
		t.Error("Dequeue should fail on empty queue")
	}
	q.Enqueue(1)
	q.Enqueue(2)
	if item, err := q.Dequeue(); err != nil || item != 1 {
		t.Errorf("Dequeue failed: expected 1, got %d, error: %v", item, err)
	}
	if _, err := q.Dequeue(); err != nil {
		t.Errorf("Dequeue failed on second element: %v", err)
	}
	if _, err := q.Dequeue(); err == nil {
		t.Error("Dequeue should fail now, queue should be empty")
	}
}

func TestPeek(t *testing.T) {
	q := newTestQueue(3)
	if _, err := q.Peek(); err == nil {
		t.Error("Peek should fail on empty queue")
	}
	q.Enqueue(1)
	if item, err := q.Peek(); err != nil || item != 1 {
		t.Errorf("Peek failed: expected 1, got %d, error: %v", item, err)
	}
}

func TestClear(t *testing.T) {
	q := newTestQueue(3)
	q.Enqueue(1)
	q.Clear()
	if !q.IsEmpty() {
		t.Error("Clear failed: queue is not empty")
	}
}

func TestIsFull(t *testing.T) {
	q := newTestQueue(2)
	if q.IsFull() {
		t.Error("IsFull failed: new queue should not be full")
	}
	q.Enqueue(1)
	q.Enqueue(2)
	if !q.IsFull() {
		t.Error("IsFull failed: queue should be full")
	}
}

func TestIsEmpty(t *testing.T) {
	q := newTestQueue(3)
	if !q.IsEmpty() {
		t.Error("IsEmpty failed: new queue should be empty")
	}
	q.Enqueue(1)
	q.Dequeue()
	if !q.IsEmpty() {
		t.Error("IsEmpty failed: queue should be empty after dequeue")
	}
}

func TestSize(t *testing.T) {
	q := newTestQueue(3)
	if expected, got := 0, q.Size(); expected != got {
		t.Errorf("Size failed: expected %d, got %d", expected, got)
	}
	q.Enqueue(1)
	q.Enqueue(2)
	if expected, got := 2, q.Size(); expected != got {
		t.Errorf("Size failed: expected %d, got %d", expected, got)
	}
}

func TestCapacity(t *testing.T) {
	q := newTestQueue(3)
	if expected, got := 3, q.Capacity(); expected != got {
		t.Errorf("Capacity failed: expected %d, got %d", expected, got)
	}
}
