package objects

import "fmt"

// Instructions represents a collection of bytecode instructions stored as a byte slice.
type Instructions struct {
	data []byte
}

// NewInstructions creates a new instance of Instructions with the provided byte slice data.
func _newInstructions(data []byte) *Instructions {
	return &Instructions{data: data}
}

// Copy creates and returns a new Instructions instance with a duplicated copy of the original data slice.
func (i *Instructions) Copy() *Instructions {
	out := _newInstructions(nil)
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

// Get8 retrieves a single byte from the instructions at the specified index.
// Returns an error if the index is out of bounds.
func (i *Instructions) Get8(index int) (uint8, error) {
	if index < 0 || index >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction index: %d", index)
	}
	return i.data[index], nil
}

// Get16 retrieves a 16-bit unsigned integer from two byte positions in the byte slice, given their low and high indices.
// Returns an error if the provided indices are out of bounds.
func (i *Instructions) Get16(low int, high int) (uint16, error) {
	if low < 0 || low >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction low index: %d", low)
	}
	if high < 0 || high >= len(i.data) {
		return 0, fmt.Errorf("invalid instruction high index: %d", high)
	}
	return uint16(i.data[low]) | uint16(i.data[high])<<8, nil
}
