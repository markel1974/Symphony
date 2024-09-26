package fifo

type Element struct {
	items    []int
	firstIdx uint
	lastIdx  uint
}

func NewElement(chunkSize uint) *Element {
	return &Element{
		items:    make([]int, chunkSize),
		firstIdx: 0,
		lastIdx:  0,
	}
}

func (e *Element) Reset() {
	e.firstIdx = 0
	e.lastIdx = 0
}

func (e *Element) Push(data int) bool {
	if e.lastIdx >= uint(len(e.items)) {
		return false
	}
	e.items[e.lastIdx] = data
	e.lastIdx++
	return true
}

func (e *Element) Pop() (int, bool) {
	if e.firstIdx >= e.lastIdx || e.firstIdx >= uint(len(e.items)) {
		return 0, false
	}
	item := e.items[e.firstIdx]
	e.firstIdx++
	return item, true
}
