package listqueue

import (
	"container/list"
	"errors"
	"fmt"
)

// Queue struct
type Queue struct {
	list      *list.List
	maxLength int
}

// New creates a new queue wih specific max length
func New(maxLength int) Queue {
	if maxLength == 0 {
		maxLength = 20
	}

	return Queue{
		list:      list.New(),
		maxLength: maxLength,
	}
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	return q.list.Len()
}

func (q *Queue) Process() (interface{}, error) {
	if q.Length() == 0 {
		return nil, errors.New("can't process queue, it's empty")
	}
	element := q.list.Front()
	item := q.list.Remove(element)
	return item, nil
}

func (q *Queue) Clear() {
	q.list.Init()
}

// Add adds an item to the queue
func (q *Queue) Add(item interface{}) error {
	if q.Length() >= q.maxLength {
		return fmt.Errorf("the queue is full (%v/%v)", q.Length(), q.maxLength)
	}

	q.list.PushBack(item)
	return nil
}

func (q *Queue) Remove(item interface{}) bool {
	element := q.find(item)
	if element == nil {
		return false
	}
	itm := q.list.Remove(element)
	if itm == nil {
		return false
	}
	return true
}

func (q *Queue) find(item interface{}) *list.Element {
	element := q.list.Front()
	for element != nil {
		if element.Value == item {
			return element
		}
		element = element.Next()
	}

	return nil
}

func (q *Queue) GetFirstItem() (interface{}, bool) {
	element := q.list.Front()
	if element == nil {
		return nil, false
	}
	if element.Value == nil {
		elem := q.list.Front()
		q.list.Remove(elem)
		elem = q.list.Front()
		return elem.Value, true
	}

	return element.Value, true
}

func (q *Queue) ReplaceFirstItem(item interface{}) bool {
	element := q.list.Front()
	q.list.Remove(element)
	el := q.list.PushFront(item)
	return el != nil
}

func (q *Queue) RemoveFirstItem() bool {
	element := q.list.Front()
	if element == nil {
		return false
	}

	item := q.list.Remove(element)
	return item != nil
}

func (q *Queue) GetList() *list.List {
	return q.list
}
