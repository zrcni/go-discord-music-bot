package queue

import (
	"errors"
	"log"
)

// Queue struct
type Queue struct {
	items []interface{}
	len   int
}

// Add adds an item to the end of the queue
func (q *Queue) Add(item interface{}) {
	log.Print("Added an item to the queue")
	q.items = append(q.items, item)
	q.len = q.len + 1
}

// Remove removes an item from the queue
// https://github.com/golang/go/wiki/SliceTricks
// NOTE: If the type of the element in a slice is a pointer or a struct with pointer fields,
// which need to be garbage collected, the "normal" way of removing an element from a slice
// might have a potential memory leak problem: some elements with values are still referenced
// by the slice and thus can not be garbage collected.
// Normal way: q.items = append(q.items[:index], q.items[index+1]...)
func (q *Queue) Remove(item interface{}) error {
	if item == nil {
		return errors.New("provided pointer is nil")
	}

	for index, el := range q.items {
		if item == el {
			copy(q.items[index:], q.items[index+1:])
			q.items[len(q.items)-1] = nil
			q.items = q.items[:len(q.items)-1]
			q.len = q.len - 1
		}
	}

	return nil
}

// Shift removes the first item in the queue and returns it
func (q *Queue) Shift() interface{} {
	log.Printf("shifting... items: %+v", q.items)
	first, rest := q.items[0], q.items[1:]
	q.items = rest
	log.Printf("shifted... items: %+v", q.items)
	q.len = q.len - 1
	return first
}

// Pop removes the last items in the queue and returns it
func (q *Queue) Pop() interface{} {
	last, rest := q.items[len(q.items)-1], q.items[:len(q.items)-1]
	q.items = rest
	return last
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	return q.len
}
