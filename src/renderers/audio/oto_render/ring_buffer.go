package oto_render

import "math"

// minimumSize defines the smallest allowable size required for the related operation or setting.
const minimumSize = 2

// CircularQueue is a circular buffer implementation for storing and managing float32 array pointers in a fixed-size queue.
// It supports operations to push elements to the end, pop elements from the start, and track the current count of elements.
// The internal data structure wraps around when reaching the maximum capacity, overwriting behavior depending on usage logic.
type CircularQueue struct {
	data    []*[]float32
	start   uint8
	end     uint8
	counter int
}

// NewCircularQueue creates and initializes a new CircularQueue with a specified entry size for each element.
func NewCircularQueue(entrySize int) *CircularQueue {
	cq := &CircularQueue{
		start:   0,
		end:     0,
		counter: 0,
	}
	cq.data = make([]*[]float32, math.MaxUint8+1)
	for x := range cq.data {
		v := make([]float32, entrySize)
		cq.data[x] = &v
	}
	return cq
}

// Push attempts to add the given element to the circular queue.
// Returns true if the operation is successful, otherwise false if the queue is full.
func (r *CircularQueue) Push(elem *[]float32) bool {
	nextEnd := r.end + 1
	if nextEnd == r.start {
		//max reached
		return false
	}
	r.counter++
	copy(*r.data[r.end], *elem)
	r.end = nextEnd
	return true
}

// Pop removes and returns the chunk at the front of the queue and a success flag. It returns false if the queue is empty.
func (r *CircularQueue) Pop() (*[]float32, bool) {
	if r.start == r.end {
		return nil, false
	}
	ptr := r.data[r.start]
	r.start = r.start + 1
	r.counter--
	return ptr, true
}

// Counter returns the current count of elements in the CircularQueue.
func (r *CircularQueue) Counter() int {
	return r.counter
}
