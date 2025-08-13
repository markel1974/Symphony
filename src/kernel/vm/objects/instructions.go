package objects

import "fmt"

// Instructions represents a collection of bytecode instructions stored as a byte slice.
type Instructions struct {
	data []byte
}

// NewInstructions creates a new instance of Instructions with the provided byte slice data.
func NewInstructions(data []byte) *Instructions {
	return &Instructions{data: data}
}

// Copy creates and returns a new Instructions instance with a duplicated copy of the original data slice.
func (i *Instructions) Copy() *Instructions {
	out := NewInstructions(nil)
	out.data = append([]byte{}, i.data...)
	return out
}

// Data returns the internal byte slice representing the instructions.
func (i *Instructions) Data() []byte {
	return i.data
}

// Length returns the number of elements in the Instructions' data slice.
func (i *Instructions) Length() int {
	return len(i.data)
}

// Get retrieves the byte at the specified index in the Instructions data slice.
func (i *Instructions) Get(index int) (int, error) {
	if index < 0 || index >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction index: %d", index)
	}
	return int(i.data[index]), nil
}

// Pos returns the position of the instruction at the specified indices.
func (i *Instructions) Pos(x int, y int) (int, error) {
	a, err := i.Get(x)
	if err != nil {
		return 0, err
	}
	b, err := i.Get(y)
	if err != nil {
		return 0, err
	}
	return a | b<<8, nil
}
