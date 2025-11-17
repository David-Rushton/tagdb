package cli

import "container/list"

// A simple queue implementation.
type queue[T any] struct {
	items *list.List
}

func newQueue[T any](s []T) *queue[T] {
	result := &queue[T]{
		items: list.New(),
	}

	for _, item := range s {
		result.enqueue(item)
	}

	return result
}

// Adds an item to the end of the queue.
func (q *queue[T]) enqueue(item T) {
	q.items.PushBack(item)
}

// Removes the item at the front of the queue and returns it.
// Panics if the queue is empty.
func (q *queue[T]) dequeue() T {
	if q.isEmpty() {
		panic("cannot dequeue from an empty queue")
	}

	result := q.items.Front().Value
	q.items.Remove(q.items.Front())

	return result.(T)
}

// Returns the item at the front of the queue without removing it.
// Panics if the queue is empty.
func (q *queue[T]) peek() T {
	if q.isEmpty() {
		panic("cannot peek from an empty queue")
	}

	return q.items.Front().Value.(T)
}

// Returns true if the queue is empty.
func (q *queue[T]) isEmpty() bool {
	return q.len() == 0
}

// Returns the number of items in the queue.
func (q *queue[T]) len() int {
	return q.items.Len()
}

func (q *queue[T]) toSlice() []T {
	result := make([]T, 0, q.len())
	for e := q.items.Front(); e != nil; e = e.Next() {
		result = append(result, e.Value.(T))
	}
	return result
}

// Returns a string representation of the queue.
func (q queue[T]) String() string {
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
