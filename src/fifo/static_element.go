package fifo

// Element represents a segment of memory for storing integers with indices for first and last entries.
type Element struct {
	items    []int
	firstIdx uint
	lastIdx  uint
}

// NewElement creates a new Element with a pre-allocated slice of integers of size chunkSize and initializes its indices.
func NewElement(chunkSize uint) *Element {
	return &Element{
		items:    make([]int, chunkSize),
		firstIdx: 0,
		lastIdx:  0,
	}
}

// Reset reinitializes the first and last indices of the Element to zero, effectively clearing its state.
func (e *Element) Reset() {
	e.firstIdx = 0
	e.lastIdx = 0
}

// Push inserts the provided data into the next available slot in the Element. Returns false if there is no space available.
func (e *Element) Push(data int) bool {
	if e.lastIdx >= uint(len(e.items)) {
		return false
	}
	e.items[e.lastIdx] = data
	e.lastIdx++
	return true
}

// Pop retrieves and removes the first element from the collection.
// Returns the element and true if successful, or 0 and false if the collection is empty.
func (e *Element) Pop() (int, bool) {
	if e.firstIdx >= e.lastIdx || e.firstIdx >= uint(len(e.items)) {
		return 0, false
	}
	item := e.items[e.firstIdx]
	e.firstIdx++
	return item, true
}
