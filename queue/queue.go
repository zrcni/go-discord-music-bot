package queue

import (
	"fmt"

	log "github.com/sirupsen/logrus"
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
	length := len(q.items)
	if length == q.maxLen {
		return fmt.Errorf("Queue is full (%v/%v)", length, q.maxLen)
	}
	q.items = append(q.items, item)
	log.Debugf("Added an item to the queue. (len: %v)", length)
	log.Debugf("queue: %+v", q.items)
	return nil
}

// Shift removes the first item in the queue and returns it
func (q *Queue) Shift() interface{} {
	first, rest := q.items[0], q.items[1:]
	q.items = rest
	log.Debugf("(Shift) Deleted an item from the queue. (len: %v)", len(q.items))
	return first
}

// Pop removes the last items in the queue and returns it
func (q *Queue) Pop() interface{} {
	last, rest := q.items[len(q.items)-1], q.items[:len(q.items)-1]
	q.items = rest
	log.Debugf("(Pop) Deleted an item from the queue. (len: %v)", len(q.items))
	return last
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	return len(q.items)
}

// Clear clears the array
func (q *Queue) Clear() {
	for i := range q.items {
		q.items[i] = nil
	}
	var emptySlice []interface{}
	q.items = emptySlice
	log.Debugf("Cleared the queue. (len: %v)", len(q.items))
}

// GetAt returns a item at index
func (q *Queue) GetAt(i int) (interface{}, error) {
	if len(q.items) < i+1 {
		return nil, fmt.Errorf("Can't get item at index %v. Queue length is %v", i, len(q.items))
	}
	return q.items[i], nil
}

// DeleteAt deletes item at index
func (q *Queue) DeleteAt(i int) bool {
	if len(q.items) < i+1 {
		return false
	}

	copy(q.items[i:], q.items[i+1:])
	q.items[len(q.items)-1] = nil
	q.items = q.items[:len(q.items)-1]

	// q.items = append(q.items[:i], q.items[i+1:]...)

	// log.Debugf("(DeleteAt) Deleted an item from the queue. (len: %v)", len(q.items))
	return true
}

// ReplaceAt replaces an item at index
func (q *Queue) ReplaceAt(i int, item interface{}) bool {
	if len(q.items) < i+1 {
		return false
	}

	log.Debugf("BEFORE: %+v", q.items[i])
	a := q.items[i]
	log.Debug("slightly after")
	b := &a
	log.Debug("a bit more after")
	*b = item
	q.items[i] = *b
	log.Debugf("AFTER: %+v", q.items[i])

	// iitm := itm
	// q.items[i] = iitm
	// log.Debugf("ITEM: %v", item)
	// log.Debugf("ATINDEX before: %v.", q.items[i])
	// for ind, e := range q.items {
	// 	if ind == i {
	// 		log.Debugf("ELEMENT: %v.", e)
	// 		addr := &e
	// 		log.Debugf("ADDR: %v.", addr)
	// 		*addr = &item
	// 		log.Debugf("ATINDEX after: %v.", q.items[i])
	// 	}
	// }

	// log.Debugf("BEFORE: %+v", q.items[i])
	// itm := q.items[i]
	// log.Debugf("slightly after")
	// itm = item
	// log.Debug("a bit more after")
	// q.items[i] = itm
	// log.Debugf("AFTER: %+v", q.items[i])

	// log.Debugf("REPLACEAT PREVIOUS ITEM: %+v", item)
	// q.items[i] = item
	// log.Debugf("REPLACEAT ITEM: %+v", q.items[i])

	// log.Debugf("(ReplaceAt) Replaced an item in the queue. (len: %v)", len(q.items)
	return true
}
