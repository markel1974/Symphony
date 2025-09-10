package opcodes

// Instructions represents a sequence of bytecode instructions used in the virtual machine or interpreter.
type Instructions struct {
	data []byte
}

// NewInstructions creates a new Instructions instance initialized with the provided byte slice.
func NewInstructions(data []byte) *Instructions {
	return &Instructions{data: data}
}

// Copy creates a new Instructions instance and duplicates its internal data.
func (i *Instructions) Copy() *Instructions {
	out := NewInstructions(nil)
	out.data = append([]byte{}, i.data...)
	return out
}

// Data returns the byte slice representing the instruction set stored within the Instructions object.
func (i *Instructions) Data() []byte {
	return i.data
}

// Length returns the number of elements in the Instructions data slice.
func (i *Instructions) Length() int {
	return len(i.data)
}

// GetReverse reads a value of specified width in backward order starting from the given base position in the instruction data.
// Returns the extracted value as uint and an error if the width is invalid or if the base is out of bounds.
func (i *Instructions) GetReverse(width int, base uint) (uint, bool) {
	switch width {
	case 1:
		v, ok := i.Get8Reverse(base)
		return uint(v), ok
	case 2:
		v, ok := i.Get16Reverse(base)
		return uint(v), ok
	case 3:
		v, ok := i.Get32Reverse(base)
		return uint(v), ok
	case 4:
		v, ok := i.Get64Reverse(base)
		return uint(v), ok
	default:
		return 0, false
	}
}

// Get8Reverse retrieves an 8-bit unsigned integer from the specified index in the instruction data if within bounds.
func (i *Instructions) Get8Reverse(base uint) (uint8, bool) {
	if base >= uint(len(i.data)) {
		return 0, false
	}
	return i.data[base], true
}

// Get16Reverse retrieves a 16-bit value from the `data` byte slice in reverse order, starting at the given `base` index.
// Returns an error if the provided `base` or corresponding high index is out of bounds.
func (i *Instructions) Get16Reverse(base uint) (uint16, bool) {
	const backwardBytes = 1
	if base < backwardBytes || base >= uint(len(i.data)) {
		return 0, false
	}
	high := base - backwardBytes
	return uint16(i.data[base]) | uint16(i.data[high])<<8, true
}

// Get32Reverse reads a 32-bit unsigned integer in reverse order starting at the specified base index in the data slice.
// Returns the 32-bit value and an error if the base index is out of bounds.
func (i *Instructions) Get32Reverse(base uint) (uint32, bool) {
	const backwardBytes = 3
	if base < backwardBytes || base >= uint(len(i.data)) {
		return 0, false
	}
	byte1 := i.data[base-backwardBytes] // MSB
	byte2 := i.data[base-2]
	byte3 := i.data[base-1]
	byte4 := i.data[base]
	return uint32(byte4) | uint32(byte3)<<8 | uint32(byte2)<<16 | uint32(byte1)<<24, true
}

// Get64Reverse reads 8 bytes in reverse order starting from the given base index and combines them into a uint64 value.
// Returns an error if the base index is out of bounds or less than 7.
func (i *Instructions) Get64Reverse(base uint) (uint64, bool) {
	const backwardBytes = 7
	if base < backwardBytes || base >= uint(len(i.data)) {
		return 0, false
	}
	byte1 := i.data[base-backwardBytes] // MSB
	byte2 := i.data[base-6]
	byte3 := i.data[base-5]
	byte4 := i.data[base-4]
	byte5 := i.data[base-3]
	byte6 := i.data[base-2]
	byte7 := i.data[base-1]
	byte8 := i.data[base]
	return uint64(byte8) | uint64(byte7)<<8 | uint64(byte6)<<16 | uint64(byte5)<<24 |
		uint64(byte4)<<32 | uint64(byte3)<<40 | uint64(byte2)<<48 | uint64(byte1)<<56, true
}
