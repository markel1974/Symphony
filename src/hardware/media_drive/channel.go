package media_drive

import (
	"github.com/markel1974/c64emu/src/common/fifo"
	"github.com/markel1974/c64emu/src/hardware/media_drive/adapters"
	"strings"
)

// Channel represents a communication medium capable of handling buffered data, using a specified adapter for storage.
// It supports data queuing operations, buffering, and customizable modes for handling files or streams.
type Channel struct {
	id       int
	adapter  adapters.IAdapter
	data     *fifo.CircularQueue
	buffer   []byte
	openMode uint8
	mode     adapters.FMode
	kind     adapters.FType
	name     string
	matcher  *Matcher
}

// NewChannel initializes and returns a new Channel instance with the given id and adapter.
// It sets up default values for the channel's data, buffer, mode, kind, and name fields.
func NewChannel(id int, adapter adapters.IAdapter) *Channel {
	c := &Channel{
		id:       id,
		adapter:  adapter,
		data:     fifo.NewCircularQueue(32),
		matcher:  NewMatcher(),
		buffer:   []byte{},
		openMode: 0,
		mode:     adapters.FModeUnknown,
		kind:     adapters.FTypeUnk,
		name:     "",
	}
	return c
}

// Reset clears the state of the Channel, including its data queue, buffer, openMode, mode, kind, and name.
func (c *Channel) Reset() {
	c.data.Initialize(0)
	c.buffer = []byte{}
	c.openMode = 0
	c.mode = adapters.FModeUnknown
	c.kind = adapters.FTypeUnk
	c.name = ""
}

// SetAdapter associates the given adapter object with the Channel instance for subsequent read/write operations.
func (c *Channel) SetAdapter(adapter adapters.IAdapter) {
	c.adapter = adapter
}

// Close writes buffer contents to file based on mode (write or append), resets channel state, and returns an error if any.
func (c *Channel) Close() error {
	var err error
	if c.mode == adapters.FModeWrite {
		err = c.adapter.WriteFile(c.name, c.buffer)
	} else if c.mode == adapters.FModeAppend {
		var data []byte
		if data, err = c.adapter.ReadFile(c.name); err == nil {
			data = append(data, c.buffer...)
			err = c.adapter.WriteFile(c.name, data)
		}
	}
	c.Reset()
	return err
}

// BufferAdd appends a single byte to the Channel's internal buffer.
func (c *Channel) BufferAdd(b uint8) {
	c.buffer = append(c.buffer, b)
}

// Buffer retrieves the current contents of the buffer associated with the Channel instance.
func (c *Channel) Buffer() []byte {
	return c.buffer
}

// SetError sets the error message for the channel. It converts the error to a byte slice and updates the channel's data.
func (c *Channel) SetError(err error) {
	if err == nil {
		c.dataSet([]byte("nil error"))
		return
	}
	c.dataSet([]byte(err.Error()))
}

// dataSet initializes the circular queue with the length of the input data and pushes each byte as an integer into the queue.
func (c *Channel) dataSet(data []byte) {
	c.data.Initialize(len(data))
	for _, k := range data {
		c.data.Push(int(k))
	}
}

// Read retrieves the next byte from the channel's data queue, returning the byte and a success status.
func (c *Channel) Read() (uint8, bool) {
	b, ok := c.data.Pop()
	return uint8(b), ok
}

// ReadIsEmpty checks if the channel's internal data queue is empty and returns true if no elements are present.
func (c *Channel) ReadIsEmpty() bool {
	l := c.data.IsEmpty()
	return l
}

// OpenModeSet sets the openMode field of the Channel to the provided value after resetting the Channel state.
func (c *Channel) OpenModeSet(d uint8) {
	c.openMode = d
}

// OpenModeGet retrieves and returns the current open mode of the channel as an unsigned 8-bit integer.
func (c *Channel) OpenModeGet() uint8 {
	return c.openMode
}

// OpenFile attempts to open a file by parsing its name and mode, handling read, write, and search operations accordingly.
// Returns an error if the file operation fails or if an unsupported type or mode is encountered.
func (c *Channel) OpenFile(realName string) error {
	c.name, c.mode, c.kind, _ = adapters.ParseFileName(realName)
	// Channel 0 is READ, channel 1 is WRITE
	if c.id == 0 || c.id == 1 {
		c.mode = adapters.FModeRead
		if c.id != 0 {
			c.mode = adapters.FModeWrite
		}
		if c.kind == adapters.FTypeDel {
			c.kind = adapters.FTypePrg
		}
	}
	if c.matcher.Contains(c.name) {
		if c.mode == adapters.FModeWrite || c.mode == adapters.FModeAppend {
			return adapters.Error(adapters.ErrSyntax33)
		}
		items, err := c.adapter.ReadDir()
		if err != nil {
			return adapters.Error(adapters.ErrFileNotFound)
		}
		found := false
		for _, item := range items {
			if !item.IsDir() {
				if found = c.matcher.Match(c.name, item.Name()); found {
					c.name = item.Name()
					break
				}
			}
		}
		if !found {
			return adapters.Error(adapters.ErrFileNotFound)
		}
	}
	if c.kind == adapters.FTypeRel {
		return adapters.Error(adapters.ErrUnimplemented)
	}
	if c.mode == adapters.FModeRead || c.mode == adapters.FModeM {
		data, err := c.adapter.ReadFile(c.name)
		if err != nil {
			return adapters.Error(adapters.ErrFileNotFound)
		}
		c.dataSet(data)
	}
	return nil
}

// OpenDirectory processes a directory pattern, retrieves entries via the adapter, and prepares the directory buffer.
func (c *Channel) OpenDirectory(pattern string) error {
	// Special treatment for "$0"
	if len(pattern) > 0 {
		if pattern[0] == '0' && len(pattern) == 1 {
			pattern = ""
		}
	}
	if p := strings.Index(pattern, ":"); p >= 0 {
		p++
		if len(pattern) < p {
			pattern = pattern[p:]
		}
	}
	entries, err := c.adapter.ReadDir()
	if err != nil {
		return err
	}
	buf := adapters.CreateDir(c.adapter.Name(), entries, pattern)
	c.dataSet(buf)
	return nil
}
