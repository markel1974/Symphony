package oto_render

// minSize defines the minimum allowed size for initializing the CircularQueue's capacity to ensure proper functionality.
const minimumSize = 32

// CircularQueue is a fixed-size queue implemented in a circular buffer format to efficiently manage FIFO operations.
type CircularQueue2 struct {
	data    []*[]float32
	start   int
	end     int
	counter int
}

// NewCircularQueue2 creates and initializes a new CircularQueue with the specified capacity.
func NewCircularQueue2(capacity int, entrySize int) *CircularQueue2 {
	if capacity < minimumSize {
		capacity = minimumSize
	}
	cq := &CircularQueue2{
		start:   0,
		end:     0,
		counter: 0,
	}
	cq.data = make([]*[]float32, capacity)
	for x := range cq.data {
		v := make([]float32, entrySize)
		cq.data[x] = &v
	}
	return cq
}

// Push attempts to add the given element to the circular queue and returns true on success or false if the queue is full.
func (r *CircularQueue2) Push(elem *[]float32) bool {
	nextEnd := (r.end + 1) % len(r.data)
	if nextEnd == r.start {
		//max reached
		return false
	}
	r.counter++
	copy(*r.data[r.end], *elem)
	r.end = nextEnd
	return true
}

// Pop removes and returns the oldest element in the circular queue and a boolean indicating success.
// Returns (0, false) if the queue is empty.
func (r *CircularQueue2) Pop() (*[]float32, bool) {
	if r.start == r.end {
		return nil, false
	}
	ptr := r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.counter--
	return ptr, true
}

func (r *CircularQueue2) Counter() int {
	return r.counter
}
