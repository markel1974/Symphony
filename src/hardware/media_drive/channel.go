package media_drive

import "github.com/markel1974/c64emu/src/common/fifo"

type Channel struct {
	data   *fifo.CircularQueue
	buffer []byte
	mode   uint8
}

func NewChannel() *Channel {
	c := &Channel{
		data:   fifo.NewCircularQueue(32),
		buffer: []byte{},
		mode:   0,
	}
	return c
}

func (c *Channel) Reset() {
	c.data.Initialize(0)
	c.buffer = []byte{}
	c.mode = 0
}

func (c *Channel) Close() {
	c.data.Initialize(0)
	c.buffer = []byte{}
	c.mode = 0
}

func (c *Channel) BufferAdd(b uint8) {
	c.buffer = append(c.buffer, b)
}

func (c *Channel) BufferGet() []byte {
	return c.buffer
}

func (c *Channel) DataSet(data []byte) {
	c.data.Initialize(len(data))
	for _, k := range data {
		c.data.Push(int(k))
	}
}

func (c *Channel) DataNext() (uint8, bool) {
	b, ok := c.data.Pop()
	return uint8(b), ok
}

func (c *Channel) DataIsEmpty() bool {
	l := c.data.IsEmpty()
	return l
}

//func (c *Channel) DataLen() int {
//	l := c.data.Len()
//	return l
//}

func (c *Channel) ModeSet(d uint8) {
	c.mode = d
}

func (c *Channel) ModeGet() uint8 {
	return c.mode
}
