package fifo

type Queue struct {
	head  *Element
	tail  *Element
	count int
	max   int
}

func NewQueue(length int) (q *Queue) {
	root := NewElement(length)
	q = &Queue{
		max:  length,
		head: root,
		tail: root,
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
	if q.tail.last >= q.max {
		return false
	}
	q.tail.items[q.tail.last] = data
	q.tail.last++
	q.count++
	return true
}

func (q *Queue) AddElement(data int) {
	if q.tail.last >= q.max {
		q.tail.next = NewElement(q.max)
		q.tail = q.tail.next
	}
	q.tail.items[q.tail.last] = data
	q.tail.last++
	q.count++
}

func (q *Queue) Next() (int, bool) {
	var item int
	if q.count == 0 {
		return 0, false
	}
	if q.head.first >= q.head.last {
		return 0, false
	}
	item = q.head.items[q.head.first]
	q.head.first++
	q.count--
	if q.head.first >= q.head.last {
		if q.count == 0 {
			q.head.first = 0
			q.head.last = 0
			q.head.next = nil
		} else {
			q.head = q.head.next
		}
	}
	return item, true
}
