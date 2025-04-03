package fs_drive

import "github.com/markel1974/c64emu/src/common/fifo"

type Channel struct {
	data   *fifo.StaticFifo
	buffer []byte
	state  uint16
}

func NewChannel() *Channel {
	c := &Channel{
		data:   fifo.NewStaticFifo(32),
		buffer: []byte{},
		state:  0,
	}
	return c
}

func (c *Channel) Reset() {
	c.data = fifo.NewStaticFifo(32)
	c.buffer = []byte{}
	c.state = 0
}

func (c *Channel) Close() {
	c.data = fifo.NewStaticFifo(32)
	c.buffer = []byte{}
	c.state = 0
}

func (c *Channel) BufferAdd(b uint8) {
	c.buffer = append(c.buffer, b)
}

func (c *Channel) BufferGet() []byte {
	return c.buffer
}

func (c *Channel) DataSet(data []byte) {
	c.data = fifo.NewStaticFifo(uint(len(data)))
	for _, k := range data {
		c.data.Set(int(k))
	}
}

func (c *Channel) DataNext() (uint8, bool) {
	b, ok := c.data.Next()
	return uint8(b), ok
}

func (c *Channel) DataLen() int {
	l := c.data.Len()
	return l
}
