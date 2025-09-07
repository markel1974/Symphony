package bytecode

import "fmt"

// Compiler is responsible for constructing binary instructions based on the provided opcodes and operands.
type Compiler struct {
	opcode       *Opcode
	instructions []byte
}

// NewCompiler initializes a new Compiler with the given Opcode. It prepares the Compiler for assembling bytecode.
func NewCompiler(opcode *Opcode) *Compiler {
	return &Compiler{
		opcode: opcode,
	}
}

// Instructions returns the compiled bytecode as a slice of bytes from the compiler instance.
func (c *Compiler) Instructions() []byte {
	return c.instructions
}

// Compile converts a list of operands into bytecode based on the opcode and writes the result to the instructions buffer.
// An error is returned if the number or size of operands does not match the expectations of the opcode.
func (c *Compiler) Compile(operands []int) error {
	features := c.opcode.Operands()
	if len(operands) != len(features) {
		return fmt.Errorf("wrong number of operands for %s: want %d, got %d", c.opcode.Name(), len(features), len(operands))
	}
	totalLen := c.opcode.FullWidth()
	c.instructions = make([]byte, totalLen)

	offset := uint(0)
	if err := c.set(uint(c.opcode.OpcodeId()), uint(OpcodeWidth), offset); err != nil {
		return fmt.Errorf("failed to set opcode: %w", err)
	}
	offset += uint(OpcodeWidth)
	for idx, operand := range operands {
		size := uint(features[idx] & SzMask)
		if err := c.set(uint(operand), size, offset); err != nil {
			return fmt.Errorf("failed to set operand %d: %w", idx, err)
		}
		offset += size
	}
	return nil
}

// SetInstructions sets the compiler's instructions to a new byte slice. If input is empty, it sets instructions to an empty slice.
func (c *Compiler) SetInstructions(v []byte) {
	if len(v) == 0 {
		c.instructions = []byte{}
		return
	}
	c.instructions = make([]byte, len(v))
	copy(c.instructions, v)
}

// Decompile parses the instructions in the compiler and returns a slice of integers representing the opcode and operands.
func (c *Compiler) Decompile() ([]int, error) {
	var out []int
	offset := uint(0)
	v, err := c.get(uint(OpcodeWidth), offset)
	if err != nil {
		return nil, fmt.Errorf("failed to set opcode: %w", err)
	}
	out = append(out, v)
	offset += uint(OpcodeWidth)
	features := c.opcode.Operands()
	for _, feature := range features {
		size := uint(feature & SzMask)
		v, err = c.get(size, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to set operand: %w", err)
		}
		out = append(out, v)
		offset += size
	}
	return out, nil
}

// get retrieves an operand of the specified width from the instructions at the given offset. Returns an error on failure.
func (c *Compiler) get(width uint, offset uint) (int, error) {
	if offset >= uint(len(c.instructions)) {
		return 0, fmt.Errorf("offset %d is out of bounds for instruction length %d", offset, len(c.instructions))
	}
	switch width {
	case uint(SzUint8):
		return c.get8(offset)
	case uint(SzUint16):
		return c.get16(offset)
	case uint(SzUint32):
		return c.get32(offset)
	case uint(SzUint64):
		return c.get64(offset)
	default:
		return 0, fmt.Errorf("unsupported operand width: %d", width)
	}
}

// set writes an operand value to the instructions buffer at a specified offset using a specified width in bytes.
func (c *Compiler) set(operand uint, width uint, offset uint) error {
	switch width {
	case uint(SzUint8):
		if err := c.set8(operand, offset); err != nil {
			return err
		}
	case uint(SzUint16):
		if err := c.set16(operand, offset); err != nil {
			return err
		}
	case uint(SzUint32):
		if err := c.set32(operand, offset); err != nil {
			return err
		}
	case uint(SzUint64):
		if err := c.set64(operand, offset); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operand width %d", width)
	}
	return nil
}

// set8 sets a 1-byte operand value into the instructions slice at the specified offset, ensuring boundaries are respected.
func (c *Compiler) set8(operand uint, offset uint) error {
	const uint8Mask = (1 << (8 * SzUint8)) - 1
	if operand > uint8Mask {
		return fmt.Errorf("operand value %d out of 1-byte range", operand)
	}
	if offset >= uint(len(c.instructions)) {
		return fmt.Errorf("offset %d out of range", offset)
	}
	c.instructions[offset] = byte(operand)
	return nil
}

// set16 writes a 2-byte unsigned integer operand at the specified offset in the instructions slice in big-endian format.
func (c *Compiler) set16(operand uint, offset uint) error {
	const uint16Mask = (1 << (8 * SzUint16)) - 1
	if operand > uint16Mask {
		return fmt.Errorf("operand value %d out of 2-byte range", operand)
	}
	if offset+1 >= uint(len(c.instructions)) {
		return fmt.Errorf("offset %d out of range", offset)
	}
	n := uint16(operand)
	c.instructions[offset] = byte(n >> 8) // Most significant byte (Big Endian)
	c.instructions[offset+1] = byte(n)    // Least significant byte
	return nil
}

// set32 writes a 32-bit unsigned integer operand to the instructions at the specified offset in Big Endian order.
// Returns an error if the operand exceeds the 32-bit range or if the offset is out of bounds.
func (c *Compiler) set32(operand uint, offset uint) error {
	const uint32Mask = (1 << (8 * SzUint32)) - 1
	if operand > uint32Mask {
		return fmt.Errorf("operand value %d out of 4-byte range", operand)
	}
	if offset+3 >= uint(len(c.instructions)) {
		return fmt.Errorf("offset %d out of range", offset)
	}
	n := uint32(operand)
	c.instructions[offset] = byte(n >> 24)
	c.instructions[offset+1] = byte(n >> 16)
	c.instructions[offset+2] = byte(n >> 8)
	c.instructions[offset+3] = byte(n)
	return nil
}

// set64 writes an 8-byte (64-bit) unsigned integer operand at the specified offset in the instructions slice in big-endian order.
// Returns an error if the offset is out of bounds.
func (c *Compiler) set64(operand uint, offset uint) error {
	if offset+7 >= uint(len(c.instructions)) {
		return fmt.Errorf("offset %d out of range", offset)
	}
	n := uint64(operand)
	c.instructions[offset] = byte(n >> 56)
	c.instructions[offset+1] = byte(n >> 48)
	c.instructions[offset+2] = byte(n >> 40)
	c.instructions[offset+3] = byte(n >> 32)
	c.instructions[offset+4] = byte(n >> 24)
	c.instructions[offset+5] = byte(n >> 16)
	c.instructions[offset+6] = byte(n >> 8)
	c.instructions[offset+7] = byte(n)
	return nil
}

// get8 retrieves the 8-bit integer value at the specified offset from the instructions slice.
// Returns an error if the offset is out of bounds.
func (c *Compiler) get8(offset uint) (int, error) {
	if offset >= uint(len(c.instructions)) {
		return 0, fmt.Errorf("offset %d is out of bounds for instruction length %d", offset, len(c.instructions))
	}
	return int(c.instructions[offset]), nil
}

// get16 retrieves a 16-bit integer (Big Endian) from the instructions slice at the specified offset.
// Returns an error if the offset exceeds the bounds of the instructions slice.
func (c *Compiler) get16(offset uint) (int, error) {
	if offset+1 >= uint(len(c.instructions)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 2 bytes for 16-bit operand at offset %d", offset)
	}
	val := int(c.instructions[offset])<<8 | int(c.instructions[offset+1])
	return val, nil
}

// get32 retrieves a 32-bit integer from the instructions starting at the specified offset. It returns an error if out of bounds.
func (c *Compiler) get32(offset uint) (int, error) {
	if offset+3 >= uint(len(c.instructions)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 4 bytes for 32-bit operand at offset %d", offset)
	}
	val := int(c.instructions[offset])<<24 |
		int(c.instructions[offset+1])<<16 |
		int(c.instructions[offset+2])<<8 |
		int(c.instructions[offset+3])
	return val, nil
}

// get64 retrieves a 64-bit integer from the instruction byte array starting at the specified offset.
// Returns an error if the offset is out of bounds or insufficient bytes are available.
func (c *Compiler) get64(offset uint) (int, error) {
	if offset+7 >= uint(len(c.instructions)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 8 bytes for 64-bit operand at offset %d", offset)
	}
	val := uint64(c.instructions[offset])<<56 |
		uint64(c.instructions[offset+1])<<48 |
		uint64(c.instructions[offset+2])<<40 |
		uint64(c.instructions[offset+3])<<32 |
		uint64(c.instructions[offset+4])<<24 |
		uint64(c.instructions[offset+5])<<16 |
		uint64(c.instructions[offset+6])<<8 |
		uint64(c.instructions[offset+7])
	return int(val), nil
}
