package fsdrive

type chunk struct {
	items []int
	first int
	last  int
	next  *chunk
}

func newChunk(chunkSize int) *chunk {
	return &chunk{
		items: make([]int, chunkSize),
		first: 0,
		last:  0,
		next:  nil,
	}
}

type Queue struct {
	head, tail *chunk // chunk head and tail
	count      int    // total amount of items in the queue
	max        int
}

// NewQueue creates a new and empty *fifo.Queue
func NewQueue(length int) (q *Queue) {
	initChunk := newChunk(length)
	q = &Queue{
		max:  length,
		head: initChunk,
		tail: initChunk,
	}
	return q
}

func (q *Queue) Len() (length int) {
	return q.count
}

func (q *Queue) AddMulti(data int, count int) bool {
	for x := 0; x < count; x++ {
		if !q.Add(data) {
			return false
		}
	}
	return true
}

func (q *Queue) Add(data int) bool {
	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= q.max {
		return false
		//q.tail.next = new(chunk)
		//q.tail = q.tail.next
	}
	// add item to the tail chunk at the last position
	q.tail.items[q.tail.last] = data
	q.tail.last++
	q.count++
	return true
}

func (q *Queue) Next() int {
	var item int
	if q.count == 0 {
		return 0
	}
	// FIXME: why would this check be required?
	if q.head.first >= q.head.last {
		return 0
	}
	item = q.head.items[q.head.first]
	q.head.first++
	q.count--
	if q.head.first >= q.head.last {
		// we're at the end of this chunk and we should do some maintenance
		// if there are no follow up chunks then reset the current one so it can be used again.
		if q.count == 0 {
			q.head.first = 0
			q.head.last = 0
			q.head.next = nil
		} else {
			// set queue's head chunk to the next chunk
			// old head will fall out of scope and be GC-ed
			q.head = q.head.next
		}
	}
	// return the retrieved item
	return item
}
