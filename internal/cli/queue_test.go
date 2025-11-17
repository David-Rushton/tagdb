package cli

import (
	"slices"
	"testing"
)

func Test_queue_EnqueueDequeueRoundTrips(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	q := newQueue(expected)

	for _, exp := range expected {
		if q.len() == 0 {
			t.Fatalf("expected queue to have items")
		}

		actual := q.dequeue()
		if actual != exp {
			t.Fatalf("dequeue() = %v, expected %v", actual, exp)
		}
	}

	if q.len() != 0 {
		t.Fatalf("unexpected items in queue after dequeues")
	}
}

func Test_queue_DequeuePanicsWhenEmpty(t *testing.T) {
	q := newQueue([]int{})
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("dequeue() did not panic on empty queue")
		}
	}()
	q.dequeue()
}

func Test_queue_PeekReturnsNext(t *testing.T) {
	expected := []string{"a", "b", "c"}
	q := newQueue(expected)

	for i := range expected {
		q.dequeue()

		if q.len() > 0 {
			actual := q.peek()
			if actual != expected[i+1] {
				t.Fatalf("peek() = %v, expected %v", actual, expected[i+1])
			}
		}
	}
}

func Test_queue_PeekPanicsWhenEmpty(t *testing.T) {
	q := newQueue([]int{})
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("dequeue() did not panic on empty queue")
		}
	}()
	q.peek()
}

func Test_queue_ToSlice(t *testing.T) {
	expected := []int{10, 20, 30, 40}
	q := newQueue(expected)

	for i := range expected {
		actualSlice := q.toSlice()

		if !slices.Equal(expected[i:], actualSlice) {
			t.Fatalf("toSlice() = %v, expected %v", actualSlice, expected[i:])
		}

		q.dequeue()
	}

	actualSlice := q.toSlice()
	if !slices.Equal(actualSlice, []int{}) {
		t.Fatalf("toSlice() = %v, expected empty slice", actualSlice)
	}
}

func Test_queue_String(t *testing.T) {
	expected := []string{"x", "y", "z"}
	q := newQueue(expected)

	expectedStr := "[x y z]"
	actualStr := q.String()
	if actualStr != expectedStr {
		t.Fatalf("String() = %v, expected %v", actualStr, expectedStr)
	}
	q.dequeue()

	expectedStr = "[y z]"
	actualStr = q.String()
	if actualStr != expectedStr {
		t.Fatalf("String() = %v, expected %v", actualStr, expectedStr)
	}
	q.dequeue()

	expectedStr = "[z]"
	actualStr = q.String()
	if actualStr != expectedStr {
		t.Fatalf("String() = %v, expected %v", actualStr, expectedStr)
	}
	q.dequeue()

	expectedStr = "[]"
	actualStr = q.String()
	if actualStr != expectedStr {
		t.Fatalf("String() = %v, expected %v", actualStr, expectedStr)
	}
}
