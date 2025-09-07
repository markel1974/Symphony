package bytecode

import "fmt"

// Compiler is responsible for compiling OpCodes into a sequence of byte instructions.
type Compiler struct {
	instructions []byte
}

// NewCompiler initializes and returns a new instance of Compiler for generating and managing bytecode instructions.
func NewCompiler() *Compiler {
	return &Compiler{}
}

// Instructions returns the current set of compiled bytecode instructions from the Compiler instance.
func (c *Compiler) Instructions() []byte {
	return c.instructions
}

// Compile generates an instruction byte sequence based on the opcode and provided operands, ensuring correct widths and offsets.
// Returns an error if the number or values of operands are invalid or if encoding fails.
func (c *Compiler) Compile(opcode *Opcode, operands []int) error {
	operandsWidth := opcode.Operands()
	if len(operands) != len(operandsWidth) {
		return fmt.Errorf("wrong number of operands for %s: want %d, got %d", opcode.Name(), len(operandsWidth), len(operands))
	}
	totalLen := OpcodeWidth
	totalLen += opcode.Offset()
	c.instructions = make([]byte, totalLen)

	offset := uint(0)
	if err := c.set(uint(opcode.OpcodeId()), OpcodeWidth, offset); err != nil {
		return fmt.Errorf("failed to set opcode: %w", err)
	}
	offset += OpcodeWidth
	for idx, operand := range operands {
		width := uint(operandsWidth[idx])
		if err := c.set(uint(operand), width, offset); err != nil {
			return fmt.Errorf("failed to set operand %d: %w", idx, err)
		}
		offset += width
	}
	return nil
}

// set writes the given operand into the instruction byte slice at the specified offset with the defined width.
// Returns an error if the width is invalid or if the operand or offset values are out of bounds.
func (c *Compiler) set(operand uint, width uint, offset uint) error {
	switch width {
	case Uint8Size:
		if err := c.set8(operand, offset); err != nil {
			return err
		}
	case Uint16Size:
		if err := c.set16(operand, offset); err != nil {
			return err
		}
	case Uint32Size:
		if err := c.set32(operand, offset); err != nil {
			return err
		}
	case Uint64Size:
		if err := c.set64(operand, offset); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid operand width %d", width)
	}
	return nil
}

// set8 sets a 1-byte operand value at the specified offset in the instruction slice of the Compiler.
// The operand must fit within 1 byte, and the offset must be within the bounds of the instructions slice.
// Returns an error if the operand is out of range or the offset exceeds the slice boundary.
func (c *Compiler) set8(operand uint, offset uint) error {
	if operand > uint8Mask {
		return fmt.Errorf("operand value %d out of 1-byte range", operand)
	}
	if offset >= uint(len(c.instructions)) {
		return fmt.Errorf("offset %d out of range", offset)
	}
	c.instructions[offset] = byte(operand)
	return nil
}

// set16 writes a 16-bit operand at a specific offset in the instructions slice in Big Endian format.
// Returns an error if the operand exceeds the 16-bit limit or if the offset is out of range.
func (c *Compiler) set16(operand uint, offset uint) error {
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

// set32 encodes a 32-bit operand into the instructions slice at the specified offset in Big Endian order.
// Returns an error if the operand value exceeds 32-bit range or if the offset is out of bounds.
func (c *Compiler) set32(operand uint, offset uint) error {
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

// set64 encodes a 64-bit unsigned integer into the instructions at the specified offset in big-endian order.
// Returns an error if the offset plus the operand size exceeds the length of the instructions slice.
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
