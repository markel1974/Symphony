package opcodes

import (
	"fmt"
)

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
func (c *Compiler) Compile(meta []uint8, operands []int) error {
	operandsFeatures := c.opcode.Operands()
	if len(operands) != len(operandsFeatures) {
		return fmt.Errorf("wrong number of operands for %s: want %d, got %d", c.opcode.Name(), len(operandsFeatures), len(operands))
	}
	headerBytes := HeaderSizeBytes + HeaderOpcodeIdBytes + len(meta)
	totalBytes := headerBytes + c.opcode.OperandsBytes()

	c.instructions = make([]byte, totalBytes)
	//Header Start
	offset := uint(0)
	if err := c.set(uint(headerBytes), HeaderSizeBytes, offset); err != nil {
		return fmt.Errorf("failed to set opcode: %w", err)
	}
	offset += HeaderSizeBytes
	if err := c.set(uint(c.opcode.OpcodeId()), HeaderOpcodeIdBytes, offset); err != nil {
		return fmt.Errorf("failed to set opcode: %w", err)
	}
	offset += HeaderOpcodeIdBytes
	if len(meta) > 0 {
		for _, v := range meta {
			if err := c.set(uint(v), 1, offset); err != nil {
				return fmt.Errorf("failed to set metadata: %w", err)
			}
		}
		offset += uint(len(meta))
	}
	//Header End
	for idx, operand := range operands {
		size := uint(operandsFeatures[idx] & SzMask)
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

func (c *Compiler) DecompileHeader() (OpcodeId, uint8, error) {
	return DecompileHeader(0, c.instructions)
}

func (c *Compiler) DecompileOperands(headerBytes uint8) ([]int, error) {
	var out []int
	offset := uint(headerBytes)
	features := c.opcode.Operands()
	for _, feature := range features {
		size := uint(feature & SzMask)
		v, err := Get(size, offset, c.instructions)
		if err != nil {
			return nil, fmt.Errorf("failed to set operand: %w", err)
		}
		out = append(out, v)
		offset += size
	}
	return out, nil
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

// DecompileHeader parses a bytecode instruction, extracting the operation code and header metadata at the given instruction pointer.
// It returns the OpcodeId, header metadata as a uint8, and an error if extraction or decoding fails.
func DecompileHeader(offset uint, data []byte) (OpcodeId, uint8, error) {
	headerBytes, err := Get(HeaderSizeBytes, offset, data)
	if err != nil {
		return OpUnknown, 0, err
	}
	opcodeId, err := Get(HeaderOpcodeIdBytes, offset+HeaderSizeBytes, data)
	if err != nil {
		return OpUnknown, 0, err
	}
	return OpcodeId(opcodeId), uint8(headerBytes), nil
}

// Get retrieves a value of a specified width from the data slice at the provided offset and returns it as an integer.
// The width specifies the size of the operand (e.g., 8, 16, 32, or 64 bits). Returns an error if the offset is out of bounds.
func Get(width uint, offset uint, data []uint8) (int, error) {
	if offset >= uint(len(data)) {
		return 0, fmt.Errorf("offset %d is out of bounds for instruction length %d", offset, len(data))
	}
	switch width {
	case uint(SzUint8):
		return get8(offset, data)
	case uint(SzUint16):
		return get16(offset, data)
	case uint(SzUint32):
		return get32(offset, data)
	case uint(SzUint64):
		return get64(offset, data)
	default:
		return 0, fmt.Errorf("unsupported operand width: %d", width)
	}
}

// get8 retrieves an 8-bit unsigned integer from the data slice at the specified offset and returns it as an int.
// If the offset is out of bounds, an error is returned.
func get8(offset uint, data []uint8) (int, error) {
	if offset >= uint(len(data)) {
		return 0, fmt.Errorf("offset %d is out of bounds for instruction length %d", offset, len(data))
	}
	return int(data[offset]), nil
}

// get16 retrieves a 16-bit integer from the given byte slice at the specified offset.
// Returns an error if the offset exceeds the data length or if two bytes can't be safely read.
func get16(offset uint, data []uint8) (int, error) {
	if offset+1 >= uint(len(data)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 2 bytes for 16-bit operand at offset %d", offset)
	}
	val := uint16(data[offset])<<8 | uint16(data[offset+1])
	return int(val), nil
}

// get32 retrieves a 32-bit integer from the provided byte slice at the specified offset.
// Returns an error if there are not enough bytes available to read.
func get32(offset uint, data []uint8) (int, error) {
	if offset+3 >= uint(len(data)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 4 bytes for 32-bit operand at offset %d", offset)
	}
	val := uint32(data[offset])<<24 |
		uint32(data[offset+1])<<16 |
		uint32(data[offset+2])<<8 |
		uint32(data[offset+3])
	return int(val), nil
}

// get64 extracts a 64-bit integer from the given byte slice starting at the specified offset.
// Returns an error if there are not enough bytes available in the slice.
func get64(offset uint, data []uint8) (int, error) {
	if offset+7 >= uint(len(data)) {
		return 0, fmt.Errorf("unexpected end of bytecode, expected 8 bytes for 64-bit operand at offset %d", offset)
	}
	val := uint64(data[offset])<<56 |
		uint64(data[offset+1])<<48 |
		uint64(data[offset+2])<<40 |
		uint64(data[offset+3])<<32 |
		uint64(data[offset+4])<<24 |
		uint64(data[offset+5])<<16 |
		uint64(data[offset+6])<<8 |
		uint64(data[offset+7])
	return int(val), nil
}
