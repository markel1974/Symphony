package media_drive

import "github.com/markel1974/c64emu/src/common/fifo"

// Channel is a data structure for managing buffered and sequential data processing with adjustable mode operations.
type Channel struct {
	data   *fifo.CircularQueue
	buffer []byte
	mode   uint8
}

// NewChannel initializes and returns a new instance of Channel with default properties and a circular queue of size 32.
func NewChannel() *Channel {
	c := &Channel{
		data:   fifo.NewCircularQueue(32),
		buffer: []byte{},
		mode:   0,
	}
	return c
}

// Reset clears the state of the Channel, reinitializing its data queue, emptying its buffer, and setting mode to 0.
func (c *Channel) Reset() {
	c.data.Initialize(0)
	c.buffer = []byte{}
	c.mode = 0
}

// Close resets the channel by clearing its buffer, setting mode to 0, and reinitializing the data queue with size 0.
func (c *Channel) Close() {
	c.data.Initialize(0)
	c.buffer = []byte{}
	c.mode = 0
}

// BufferAdd appends a single byte to the internal buffer of the channel.
func (c *Channel) BufferAdd(b uint8) {
	c.buffer = append(c.buffer, b)
}

// BufferGet retrieves the current content of the buffer associated with the channel as a slice of bytes.
func (c *Channel) BufferGet() []byte {
	return c.buffer
}

// DataSet initializes the channel's data queue and pushes each byte from the provided slice into the queue.
func (c *Channel) DataSet(data []byte) {
	c.data.Initialize(len(data))
	for _, k := range data {
		c.data.Push(int(k))
	}
}

// DataNext retrieves and removes the next byte from the channel's circular queue. Returns the byte and a status indicating success.
func (c *Channel) DataNext() (uint8, bool) {
	b, ok := c.data.Pop()
	return uint8(b), ok
}

// DataIsEmpty checks if the Channel's data queue is empty and returns true if no data is present, otherwise false.
func (c *Channel) DataIsEmpty() bool {
	l := c.data.IsEmpty()
	return l
}

// ModeSet sets the mode of the Channel to the specified uint8 value provided as an argument.
func (c *Channel) ModeSet(d uint8) {
	c.mode = d
}

// ModeGet returns the current mode of the Channel as an 8-bit unsigned integer.
func (c *Channel) ModeGet() uint8 {
	return c.mode
}
