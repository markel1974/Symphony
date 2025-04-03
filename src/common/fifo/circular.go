package fifo

import "fmt"

// CircularQueue is a circular queue implementation backed by a slice of integers.
// It efficiently handles enqueue and dequeue operations with fixed capacity.
// When full, new additions will not overwrite existing elements unless explicitly managed.
// The queue tracks its start and end indices, looping over the slice as necessary.
// The isFull field indicates whether the queue is at full capacity.
type CircularQueue struct {
	data   []int
	isFull bool
	start  int
	end    int
}

// NewCircularQueue creates a new RingQueue instance with the specified capacity.
func NewCircularQueue(capacity int) *CircularQueue {
	return &CircularQueue{
		data:   make([]int, capacity),
		isFull: false,
		start:  0,
		end:    0,
	}
}

// Push adds an element to the end of the ring queue, returning an error if the queue is full.
func (r *CircularQueue) Push(elem int) bool {
	if r.isFull {
		return false
	}
	r.data[r.end] = elem
	r.end = (r.end + 1) % len(r.data)
	r.isFull = r.end == r.start
	return true
}

func (r *CircularQueue) IsEmpty() bool {
	return !r.isFull && r.start == r.end
}

// Pop removes and returns the front element of the queue. Returns an error if the queue is empty.
func (r *CircularQueue) Pop() (int, bool) {
	if !r.isFull && r.start == r.end {
		return 0, false
	}
	res := r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.isFull = false
	return res, true
}

// String returns a string representation of the RingQueue, including its status, size, indices, and data contents.
func (r *CircularQueue) String() string {
	return fmt.Sprintf(
		"[RRQ full:%v size:%d start:%d end:%d data:%v]",
		r.isFull,
		len(r.data),
		r.start,
		r.end,
		r.data)
}
