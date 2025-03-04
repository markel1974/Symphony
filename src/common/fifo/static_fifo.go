package fifo

// StaticFifo represents a fixed-size first-in-first-out (FIFO) buffer implemented with linked Elements.
// It supports efficient insertion and removal of integers while maintaining the order of elements.
// The structure uses linked Element nodes, and each Element stores a preallocated slice of integers.
// The count field tracks the total number of elements present in the FIFO.
// The head field points to the first Element for retrieval, and the tail points to the last Element for insertion.
type StaticFifo struct {
	head  *Element
	tail  *Element
	count int
}

// NewStaticFifo creates a new StaticFifo structure with specified length, initializing its head and tail segments.
func NewStaticFifo(length uint) (q *StaticFifo) {
	if length < 1 {
		panic("length must be greater than 0")
	}
	root := NewElement(length)
	q = &StaticFifo{
		head: root,
		tail: root,
	}
	return q
}

// Len returns the current number of elements stored in the StaticFifo.
func (q *StaticFifo) Len() (length int) {
	return q.count
}

// SetMulti attempts to add a value to the FIFO multiple times based on the count and returns false if any addition fails.
func (q *StaticFifo) SetMulti(data int, count uint) bool {
	for x := uint(0); x < count; x++ {
		if !q.Set(data) {
			return false
		}
	}
	return true
}

// Set adds the provided data to the FIFO queue by pushing it to the tail. Returns true if successful, otherwise false.
func (q *StaticFifo) Set(data int) bool {
	if !q.tail.Push(data) {
		return false
	}
	q.count++
	return true
}

// Next retrieves and removes the next element from the StaticFifo. Returns the element and true, or 0 and false if empty.
func (q *StaticFifo) Next() (int, bool) {
	if q.count == 0 {
		return 0, false
	}
	q.count--
	item, ok := q.head.Pop()
	if !ok {
		return 0, false
	}
	if q.count == 0 && (q.head.firstIdx >= q.head.lastIdx) {
		q.head.Reset()
	}
	return item, true
}
