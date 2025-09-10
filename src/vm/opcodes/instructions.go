package opcodes

// Instructions represents a sequence of bytecode instructions used in the virtual machine or interpreter.
type Instructions struct {
	data []byte
}

// NewInstructions creates a new Instructions instance initialized with the provided byte slice.
func NewInstructions(data []byte) *Instructions {
	return &Instructions{data: data}
}

// Allocate reserves the specified number of bytes and initializes the data field with a new byte slice of length l.
func (i *Instructions) Allocate(l int) {
	i.data = make([]byte, l)
}

// Assign copies the provided data into the receiver's data field.
func (i *Instructions) Assign(data []byte) {
	//i.data = make([]byte, len(data))
	//copy(i.data, data)
	i.data = data
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

// Get retrieves a value of a specified width from the data slice at the provided offset and returns it as an integer.
// The width specifies the size of the operand (e.g., 8, 16, 32, or 64 bits). Returns an error if the offset is out of bounds.
func (i *Instructions) Get(width uint, offset uint) (int, bool) {
	if offset >= uint(len(i.data)) {
		return 0, false
	}
	switch width {
	case uint(SzUint8):
		return i.Get8(offset)
	case uint(SzUint16):
		return i.Get16(offset)
	case uint(SzUint32):
		return i.Get32(offset)
	case uint(SzUint64):
		return i.Get64(offset)
	default:
		return 0, false
	}
}

// Get8 retrieves an 8-bit unsigned integer from the data slice at the specified offset and returns it as an int.
// If the offset is out of bounds, an error is returned.
func (i *Instructions) Get8(offset uint) (int, bool) {
	if offset >= uint(len(i.data)) {
		return 0, false
	}
	return int(i.data[offset]), true
}

// Get16 retrieves a 16-bit integer from the given byte slice at the specified offset.
// Returns an error if the offset exceeds the data length or if two bytes can't be safely read.
func (i *Instructions) Get16(offset uint) (int, bool) {
	if offset+1 >= uint(len(i.data)) {
		return 0, false
	}
	val := uint16(i.data[offset])<<8 | uint16(i.data[offset+1])
	return int(val), true
}

// Get32 retrieves a 32-bit integer from the provided byte slice at the specified offset.
// Returns an error if there are not enough bytes available to read.
func (i *Instructions) Get32(offset uint) (int, bool) {
	if offset+3 >= uint(len(i.data)) {
		return 0, false
	}
	val := uint32(i.data[offset])<<24 |
		uint32(i.data[offset+1])<<16 |
		uint32(i.data[offset+2])<<8 |
		uint32(i.data[offset+3])
	return int(val), true
}

// Get64 extracts a 64-bit integer from the given byte slice starting at the specified offset.
// Returns an error if there are not enough bytes available in the slice.
func (i *Instructions) Get64(offset uint) (int, bool) {
	if offset+7 >= uint(len(i.data)) {
		return 0, false
	}
	val := uint64(i.data[offset])<<56 |
		uint64(i.data[offset+1])<<48 |
		uint64(i.data[offset+2])<<40 |
		uint64(i.data[offset+3])<<32 |
		uint64(i.data[offset+4])<<24 |
		uint64(i.data[offset+5])<<16 |
		uint64(i.data[offset+6])<<8 |
		uint64(i.data[offset+7])
	return int(val), true
}

// Set writes an operand value to the instructions buffer at a specified offset using a specified width in bytes.
func (i *Instructions) Set(operand uint, width uint, offset uint) bool {
	switch width {
	case uint(SzUint8):
		return i.Set8(operand, offset)
	case uint(SzUint16):
		return i.Set16(operand, offset)
	case uint(SzUint32):
		return i.Set32(operand, offset)
	case uint(SzUint64):
		return i.Set64(operand, offset)
	default:
		return false
	}
}

// Set8 sets a 1-byte operand value into the instructions slice at the specified offset, ensuring boundaries are respected.
func (i *Instructions) Set8(operand uint, offset uint) bool {
	const uint8Mask = (1 << (8 * SzUint8)) - 1
	if operand > uint8Mask {
		return false
	}
	if offset >= uint(len(i.data)) {
		return false
	}
	i.data[offset] = byte(operand)
	return true
}

// Set16 writes a 2-byte unsigned integer operand at the specified offset in the instructions slice in big-endian format.
func (i *Instructions) Set16(operand uint, offset uint) bool {
	const uint16Mask = (1 << (8 * SzUint16)) - 1
	if operand > uint16Mask {
		return false
	}
	if offset+1 >= uint(len(i.data)) {
		return false
	}
	n := uint16(operand)
	i.data[offset] = byte(n >> 8) // Most significant byte (Big Endian)
	i.data[offset+1] = byte(n)    // Least significant byte
	return true
}

// Set32 writes a 32-bit unsigned integer operand to the instructions at the specified offset in Big Endian order.
// Returns an error if the operand exceeds the 32-bit range or if the offset is out of bounds.
func (i *Instructions) Set32(operand uint, offset uint) bool {
	const uint32Mask = (1 << (8 * SzUint32)) - 1
	if operand > uint32Mask {
		return false
	}
	if offset+3 >= uint(len(i.data)) {
		return false
	}
	n := uint32(operand)
	i.data[offset] = byte(n >> 24)
	i.data[offset+1] = byte(n >> 16)
	i.data[offset+2] = byte(n >> 8)
	i.data[offset+3] = byte(n)
	return true
}

// Set64 writes an 8-byte (64-bit) unsigned integer operand at the specified offset in the instructions slice in big-endian order.
// Returns an error if the offset is out of bounds.
func (i *Instructions) Set64(operand uint, offset uint) bool {
	if offset+7 >= uint(len(i.data)) {
		return false
	}
	n := uint64(operand)
	i.data[offset] = byte(n >> 56)
	i.data[offset+1] = byte(n >> 48)
	i.data[offset+2] = byte(n >> 40)
	i.data[offset+3] = byte(n >> 32)
	i.data[offset+4] = byte(n >> 24)
	i.data[offset+5] = byte(n >> 16)
	i.data[offset+6] = byte(n >> 8)
	i.data[offset+7] = byte(n)
	return true
}

// Header parses a bytecode instruction, extracting the operation code and header metadata at the given instruction pointer.
// It returns the OpcodeId, header metadata as a uint8, and an error if extraction or decoding fails.
func (i *Instructions) Header(offset uint) (OpcodeId, uint8, bool) {
	headerBytes, ok := i.Get(HeaderSizeBytes, offset)
	if !ok {
		return OpUnknown, 0, false
	}
	opcodeId, ok := i.Get(HeaderOpcodeIdBytes, offset+HeaderSizeBytes)
	if !ok {
		return OpUnknown, 0, false
	}
	return OpcodeId(opcodeId), uint8(headerBytes), true
}
