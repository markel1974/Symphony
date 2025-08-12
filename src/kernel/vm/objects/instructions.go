package objects

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
func (i *Instructions) Get(index int) byte {
	return i.data[index]
}
