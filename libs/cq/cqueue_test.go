package cq

// George trying to figure out how to have basic unit tests & property tests in Go.

import (
	"testing"
	"testing/quick"
)

// Property Tests

func TestEnqueueDequeueProperty(t *testing.T) {
	property := func(a []int) bool {
		q := NewCircularQueue[int](len(a) + 1) // Ensure there's enough space
		for _, item := range a {
			if err := q.Enqueue(item); err != nil {
				return false // Should never fail to enqueue in this setup
			}
		}

		for _, expected := range a {
			if item, err := q.Dequeue(); err != nil || item != expected {
				return false
			}
		}
		return true
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error("Failed Enqueue/Dequeue property check:", err)
	}
}

// Test the consistency of the Peek method.

func TestPeekProperty(t *testing.T) {
	property := func(a []int) bool {
		if len(a) == 0 {
			return true // Nothing to test if no elements
		}
		q := NewCircularQueue[int](len(a))
		for _, item := range a {
			if err := q.Enqueue(item); err != nil {
				return false // Stop if full (though setup should prevent this)
			}
		}
		expected, _ := q.Peek() // Ignore error for non-empty queue
		item, _ := q.Dequeue()  // Ignore error for non-empty queue
		return expected == item
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error("Failed Peek property check:", err)
	}
}

// Test the size property of the queue.
func TestSizeProperty(t *testing.T) {
	property := func(a []int) bool {
		q := NewCircularQueue[int](len(a) * 2) // Make sure we don't run out of space
		count := 0
		for _, item := range a {
			if q.Enqueue(item) == nil {
				count++
			}
		}
		return q.Size() == count
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error("Failed Size property check:", err)
	}
}

// Test that the IsFull and IsEmpty methods report correctly.
func TestFullEmptyProperty(t *testing.T) {
	property := func(a []int) bool {
		capacity := len(a)
		if capacity == 0 {
			return true // Trivially true for empty input
		}
		q := NewCircularQueue[int](capacity)
		for i := 0; i < capacity; i++ {
			if q.Enqueue(a[i]) != nil {
				return false // Should not fail to enqueue up to capacity
			}
		}
		if !q.IsFull() {
			return false
		}
		for i := 0; i < capacity; i++ {
			if _, err := q.Dequeue(); err != nil {
				return false // Should not fail to dequeue when items are present
			}
		}
		return q.IsEmpty()
	}

	if err := quick.Check(property, nil); err != nil {
		t.Error("Failed Full/Empty property check:", err)
	}
}

// These are unit tests (not property tests)

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
