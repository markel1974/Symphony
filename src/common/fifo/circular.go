package fifo

import "fmt"

// minSize defines the minimum allowed size for initializing the CircularQueue's capacity to ensure proper functionality.
const minSize = 32

// CircularQueue is a fixed-size queue implemented in a circular buffer format to efficiently manage FIFO operations.
type CircularQueue struct {
	data    []int
	isFull  bool
	isEmpty bool
	start   int
	end     int
}

// NewCircularQueue creates and initializes a new CircularQueue with the specified capacity.
func NewCircularQueue(capacity int) *CircularQueue {
	cq := &CircularQueue{}
	cq.Initialize(capacity)
	return cq
}

// Initialize sets up the CircularQueue with a specified capacity or a minimum size if the capacity is too small.
func (r *CircularQueue) Initialize(capacity int) {
	if capacity < minSize {
		capacity = minSize
	}
	r.data = make([]int, capacity)
	r.Reset()
}

// Reset clears the CircularQueue, resetting state to empty, with start and end indices set to 0.
func (r *CircularQueue) Reset() {
	r.isFull = false
	r.isEmpty = true
	r.start = 0
	r.end = 0
}

// IsEmpty checks if the CircularQueue is currently empty and returns true if it contains no elements.
func (r *CircularQueue) IsEmpty() bool {
	return r.isEmpty
}

// IsFull checks if the circular queue is full and returns true if it is, otherwise false.
func (r *CircularQueue) IsFull() bool {
	return r.isFull
}

// Capacity returns the total number of elements the CircularQueue can hold.
func (r *CircularQueue) Capacity() int {
	return len(r.data)
}

// Push attempts to add the given element to the circular queue and returns true on success or false if the queue is full.
func (r *CircularQueue) Push(elem int) bool {
	if r.isFull {
		return false
	}
	r.data[r.end] = elem
	r.end = (r.end + 1) % len(r.data)
	r.isFull = r.end == r.start
	r.isEmpty = false
	return true
}

// Pop removes and returns the oldest element in the circular queue and a boolean indicating success.
// Returns (0, false) if the queue is empty.
func (r *CircularQueue) Pop() (int, bool) {
	if r.isEmpty {
		return 0, false
	}
	res := r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.isFull = false
	r.isEmpty = r.start == r.end
	return res, true
}

// String returns a string representation of the CircularQueue, showing its state, size, boundaries, and stored data.
func (r *CircularQueue) String() string {
	return fmt.Sprintf(
		"[CQueue full:%v empty:%v size:%d start:%d end:%d data:%v]",
		r.isFull,
		r.isEmpty,
		len(r.data),
		r.start,
		r.end,
		r.data)
}
