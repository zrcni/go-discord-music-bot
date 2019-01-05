package queue

import (
	"fmt"
)

// Queue struct
type Queue struct {
	items  []interface{}
	len    int
	maxLen int
}

// New creates a new queue wih specific max length
func New(maxLen int) Queue {
	if maxLen == 0 {
		maxLen = 20
	}

	return Queue{
		maxLen: maxLen,
	}
}

// Add adds an item to the end of the queue
func (q *Queue) Add(item interface{}) error {
	if q.len == q.maxLen {
		return fmt.Errorf("Queue is full (%v/%v)", q.len, q.maxLen)
	}
	q.items = append(q.items, item)
	q.len = q.len + 1

	return nil
}

// Shift removes the first item in the queue and returns it
func (q *Queue) Shift() interface{} {
	first, rest := q.items[0], q.items[1:]
	q.items = rest
	q.len = q.len - 1
	return first
}

// Pop removes the last items in the queue and returns it
func (q *Queue) Pop() interface{} {
	last, rest := q.items[len(q.items)-1], q.items[:len(q.items)-1]
	q.items = rest
	q.len = q.len - 1
	return last
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	return q.len
}

// Clear clears the array
func (q *Queue) Clear() {
	for i := range q.items {
		q.items[i] = nil
	}
	var emptySlice []interface{}
	q.items = emptySlice
	q.len = 0
}
