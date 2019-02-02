package chanqueue

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// the number of items processed concurrently
const concurrentProcesses int = 1

// Queue struct
type Queue struct {
	Items     chan interface{}
	QuitChan  chan bool
	processFn func(interface{})
	maxLength int
}

// New creates a new queue wih specific max length
func New(maxLength int) Queue {
	if maxLength == 0 {
		maxLength = 20
	}

	return Queue{
		Items:     make(chan interface{}, concurrentProcesses),
		QuitChan:  make(chan bool),
		maxLength: maxLength,
	}
}

func (q *Queue) ProcessItem(processFn func(interface{})) {
	q.processFn = processFn
}

func (q *Queue) Start() {
	go func() {
		for {
			select {
			// Get item from Items channel
			case item := <-q.Items:
				log.Infof("Item: %+v", item)
				// process item
				if q.processFn != nil {
					q.processFn(item)
				} else {
					log.Error("the queue doesnt have a processFn!")
				}

			case <-q.QuitChan:
				log.Infof("worker stopping\n")
				return
			}
		}
	}()
}

func (q *Queue) Stop() {
	// goroutine, because the queue will quit only after the items have been processed?
	// TODO: quit cancels items in queue
	go func() {
		q.QuitChan <- true
	}()
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	return len(q.Items)
}

// Add adds an item to the queue
func (q *Queue) Add(item interface{}) error {
	if len(q.Items) >= q.maxLength {
		return fmt.Errorf("the queue is full (%v/%v)", len(q.Items), q.maxLength)
	}

	q.Items <- item

	return nil
}
