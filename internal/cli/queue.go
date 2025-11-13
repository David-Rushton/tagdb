package cli

import "container/list"

// A simple queue implementation.
type queue struct {
	items *list.List
}

func newQueue[T any](s []T) *queue {
	result := &queue{
		items: list.New(),
	}

	for _, item := range s {
		result.enqueue(item)
	}

	return result
}

// Adds an item to the end of the queue.
func (q *queue) enqueue(item interface{}) {
	q.items.PushBack(item)
}

// Removes the item at the front of the queue and returns it.
// Panics if the queue is empty.
func (q *queue) dequeue() interface{} {
	if q.isEmpty() {
		panic("cannot dequeue from an empty queue")
	}

	result := q.items.Front().Value
	q.items.Remove(q.items.Front())

	return result
}

// Returns the item at the front of the queue without removing it.
// Panics if the queue is empty.
func (q *queue) peek() interface{} {
	if q.isEmpty() {
		panic("cannot peek from an empty queue")
	}

	return q.items.Front().Value
}

// Returns true if the queue is empty.
func (q *queue) isEmpty() bool {
	return q.len() == 0
}

// Returns the number of items in the queue.
func (q *queue) len() int {
	return q.items.Len()
}

// Returns a string representation of the queue.
func (q queue) String() string {
	result := "["
	for e := q.items.Front(); e != nil; e = e.Next() {
		if len(result) > 1 {
			result += " "
		}

		result += e.Value.(string)
	}
	result += "]"

	return result
}
