package opcodes

import (
	"fmt"
)

// Compiler is responsible for constructing binary instructions based on the provided opcodes and operands.
type Compiler struct {
	opcode       *Opcode
	instructions *Instructions
}

// NewCompiler initializes a new Compiler with the given Opcode. It prepares the Compiler for assembling bytecode.
func NewCompiler(opcode *Opcode) *Compiler {
	return &Compiler{
		opcode:       opcode,
		instructions: NewInstructions([]byte{}),
	}
}

// Instructions returns the compiled bytecode as a slice of bytes from the compiler instance.
func (c *Compiler) Instructions() []byte {
	return c.instructions.Code()
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

	c.instructions.Allocate(totalBytes)
	//Header Start
	offset := uint(0)
	if ok := c.instructions.Set(uint(headerBytes), HeaderSizeBytes, offset); !ok {
		return fmt.Errorf("failed to set opcode %d at offset %d", c.opcode.OpcodeId(), offset)
	}
	offset += HeaderSizeBytes
	if ok := c.instructions.Set(uint(c.opcode.OpcodeId()), HeaderOpcodeIdBytes, offset); !ok {
		return fmt.Errorf("failed to set opcode %d at offset %d", c.opcode.OpcodeId(), offset)
	}
	offset += HeaderOpcodeIdBytes
	if len(meta) > 0 {
		for _, v := range meta {
			if ok := c.instructions.Set(uint(v), 1, offset); !ok {
				return fmt.Errorf("failed to set metadata %d at offset %d", v, offset)
			}
		}
		offset += uint(len(meta))
	}
	//Header End
	for idx, operand := range operands {
		size := uint(operandsFeatures[idx] & SzMask)
		if ok := c.instructions.Set(uint(operand), size, offset); !ok {
			return fmt.Errorf("failed to set metadata %d at offset %d", operand, idx)
		}
		offset += size
	}
	return nil
}

// SetInstructions sets the compiler's instructions to a new byte slice. If input is empty, it sets instructions to an empty slice.
func (c *Compiler) SetInstructions(v []byte) {
	c.instructions.Assign(v)
}

// DecompileOperands extracts operands from the compiler's instruction buffer using the opcode metadata and headerBytes.
func (c *Compiler) DecompileOperands(headerBytes uint8) ([]int, error) {
	var out []int
	offset := uint(headerBytes)
	features := c.opcode.Operands()
	for _, feature := range features {
		size := uint(feature & SzMask)
		v, ok := c.instructions.Get(size, offset)
		if !ok {
			return nil, fmt.Errorf("failed to get operand size %d at offset %d", size, offset)
		}
		out = append(out, v)
		offset += size
	}
	return out, nil
}
