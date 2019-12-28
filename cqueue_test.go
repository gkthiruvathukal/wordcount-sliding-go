package main

import (
    "testing"
    "strconv"
)

func TestNew(t *testing.T) {
    q := new(CQueueString)
    expected := 5
    q.init(expected)
    if !q.isEmpty() {
       t.Errorf("CQueueString not empty (size is %d; length is %d).", q.size(), q.length())
    }
    if q.length() != expected {
       t.Errorf("CQueueString not %d (length is %d).", expected, q.length())
    }
}


func TestFill(t *testing.T) {
    q := new(CQueueString)
    expected := 5
    q.init(expected)
    for i :=0; i < expected-1; i++ {
        q.add( strconv.Itoa(i))
        if q.isFull() {
           t.Errorf("CQueueString should not be full yet (has %d elements, length %d).", q.size(), q.length())
        }
    }
    q.add(strconv.Itoa(expected))
    if !q.isFull() {
       t.Errorf("CQueueString did not reach expected capacity %d (size is %d).", expected, q.size())
    }
}

