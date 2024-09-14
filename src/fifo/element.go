package fifo

type Element struct {
	items []int
	first int
	last  int
	next  *Element
}

func NewElement(chunkSize int) *Element {
	return &Element{
		items: make([]int, chunkSize),
		first: 0,
		last:  0,
		next:  nil,
	}
}
