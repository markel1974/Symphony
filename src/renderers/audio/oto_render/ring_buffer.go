package oto_render

import (
	"sync"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// minSize defines the minimum allowed size for initializing the CircularQueue's capacity to ensure proper functionality.
const minSize = 32

// CircularQueue is a fixed-size queue implemented in a circular buffer format to efficiently manage FIFO operations.
type CircularQueue struct {
	data        []*[]float32
	start       int
	end         int
	ready       bool
	readyTarget int
	lock        sync.RWMutex
	counter     int
}

// NewCircularQueue creates and initializes a new CircularQueue with the specified capacity.
func NewCircularQueue(capacity int, entrySize int, readyTarget int) *CircularQueue {
	if capacity < minSize {
		capacity = minSize
	}
	cq := &CircularQueue{
		start:       0,
		end:         0,
		ready:       false,
		counter:     0,
		readyTarget: readyTarget,
	}
	cq.data = make([]*[]float32, capacity)
	for x := range cq.data {
		v := make([]float32, entrySize)
		cq.data[x] = &v
	}
	return cq
}

// Push attempts to add the given element to the circular queue and returns true on success or false if the queue is full.
func (r *CircularQueue) Push(elem *[]float32) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	nextEnd := (r.end + 1) % len(r.data)
	if nextEnd == r.start {
		//max reached
		return false
	}
	r.counter++
	if r.counter >= r.readyTarget {
		r.ready = true
	}
	copy(*r.data[r.end], *elem)
	r.end = nextEnd
	return true
}

// Pop removes and returns the oldest element in the circular queue and a boolean indicating success.
// Returns (0, false) if the queue is empty.
func (r *CircularQueue) Pop() (*[]float32, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.ready {
		//fmt.Println("[CircularQueue] NOT READY")
		return nil, false
	}
	if r.start == r.end {
		r.ready = false
		return nil, false
	}
	ptr := r.data[r.start]
	r.start = (r.start + 1) % len(r.data)
	r.counter--
	return ptr, true
}

// FillRatio computes and returns the fraction of the CircularQueue's capacity that is currently filled.
func (r *CircularQueue) FillRatio() float64 {
	r.lock.RLock() // Usiamo RLock per una lettura non bloccante
	defer r.lock.RUnlock()
	capacity := len(r.data)
	if capacity == 0 {
		return 0
	}
	return float64(r.counter) / float64(capacity)
}

/*
// Pop removes and returns the oldest element in the circular queue and a boolean indicating success.
// Returns (0, false) if the queue is empty.
func (r *CircularQueue) Pop() (*[]float32, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.ready {
		fmt.Println("[CircularQueue] NOT READY")
		return nil, false
	}
	nextStart := (r.start + 1) % len(r.data)
	if nextStart == r.end {
		r.ready = false
		return nil, false
	}
	ptr := r.data[r.start]
	r.start = nextStart
	r.counter--
	return ptr, true
}


*/
