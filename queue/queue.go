package queue

// Queue struct
type Queue struct {
	items []interface{}
	len   int
}

// Add adds an item to the end of the queue
func (q *Queue) Add(item interface{}) {
	q.items = append(q.items, item)
	q.len = q.len + 1
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
