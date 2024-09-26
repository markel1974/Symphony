package fifo

type StaticFifo struct {
	head  *Element
	tail  *Element
	count int
}

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

func (q *StaticFifo) Len() (length int) {
	return q.count
}

func (q *StaticFifo) SetMulti(data int, count uint) bool {
	for x := uint(0); x < count; x++ {
		if !q.Set(data) {
			return false
		}
	}
	return true
}

func (q *StaticFifo) Set(data int) bool {
	if !q.tail.Push(data) {
		return false
	}
	q.count++
	return true
}

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
