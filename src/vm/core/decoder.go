package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
)

// operandsMax defines the maximum number of operands, calculated as 2 raised to the power of 4.
// operandsMask is a bitmask derived from operandsMax, used for extracting operand-related values.
const (
// operandsMax  = 1 << 4
// operandsMask = operandsMax - 1
)

type DecoderData struct {
	offset   uint
	retrieve func(*Frame, uint) int
}

// Decoder represents an instruction decoder used to process bytecode in a virtual machine.
// It handles decoding operands and executing the associated instruction logic.
//
// --- ARCHITECTURAL NOTE ---
// The decoder uses a "reverse operand" logic. The list of operand widths is stored
// in reverse order. This allows the Decode function to read operands backwards from the
// instruction stream, which simplifies handling variable-sized instructions.
// Consequently, when an IOpExecutor calls decoder.Read(N), it's accessing the N-th operand
// from the *end* of the list. For example, decoder.Read(0) retrieves the LAST operand.
type Decoder struct {
	executor        IOpExecutor
	name            string
	execute         func(data *Decoder)
	fullWidth       int
	decodedOperands []int
	operands        []*DecoderData
	operandsMask    int
}

// NewDecoder creates a new Decoder instance with the specified execution function and operand widths.
func NewDecoder(executor IOpExecutor) (*Decoder, error) {
	operands := executor.Operands()
	operandsMask := computeBitmask(len(operands))

	sd := &Decoder{
		executor:        executor,
		execute:         executor.Execute,
		name:            executor.Name(),
		operandsMask:    operandsMask,
		fullWidth:       0,
		decodedOperands: make([]int, operandsMask+1),
	}

	idx := 0
	for i := len(operands) - 1; i >= 0; i-- {
		var retrieve func(*Frame, uint) int
		width := operands[i]
		switch width {
		case bytecode.Uint8Size:
			retrieve = sd.get8
		case bytecode.Uint16Size:
			retrieve = sd.get16
		case bytecode.Uint32Size:
			retrieve = sd.get32
		case bytecode.Uint64Size:
			retrieve = sd.get64
		default:
			return nil, fmt.Errorf("invalid operand width: %d", width)
		}
		sd.operands = append(sd.operands, &DecoderData{offset: uint(sd.fullWidth), retrieve: retrieve})
		sd.fullWidth += width
		idx++
	}
	if len(sd.operands) > operandsMask {
		return nil, fmt.Errorf("invalid operand mask: %d", operandsMask)
	}
	return sd, nil
}

// Name returns the name of the instruction.
func (d *Decoder) Name() string {
	return d.name
}

// Decode processes and decodes operands from the instruction pointer, updating decodedOperands and returning new ip.
func (d *Decoder) Decode(frame *Frame, ip int) int {
	if d.fullWidth == 0 {
		return ip
	}
	ip += d.fullWidth
	for idx, dd := range d.operands {
		d.decodedOperands[idx] = dd.retrieve(frame, uint(ip)-dd.offset)
	}
	return ip
}

// Execute runs the logic associated with the current instruction using the provided virtual machine instance.
func (d *Decoder) Execute() {
	d.execute(d)
}

// Read retrieves a decoded operand from the `decodedOperands` slice using a masked index derived from the input parameter.
//
// IMPORTANT: Due to the decoder's reverse logic, Read(0) accesses the LAST operand
// of the instruction, Read(1) the second to last, and so on
func (d *Decoder) Read(x int) int {
	return d.decodedOperands[x&d.operandsMask]
}

// get8 retrieves an 8-bit integer from the provided frame at the specified instruction pointer as a signed integer.
func (d *Decoder) get8(f *Frame, ip uint) int {
	return int(f.Get8(ip))
}

// get16 retrieves a 16-bit integer from the frame's instructions at the specified instruction pointer (ip) position.
func (d *Decoder) get16(f *Frame, ip uint) int {
	return int(f.Get16(ip))
}

// get32 retrieves a 32-bit signed integer from the provided frame at the specified instruction pointer position.
func (d *Decoder) get32(f *Frame, ip uint) int {
	return int(f.Get32(ip))
}

// get64 retrieves a 32-bit integer from the provided frame at the specified instruction pointer and converts it to int.
func (d *Decoder) get64(f *Frame, ip uint) int {
	return int(f.Get64(ip))
}

// computeBitmask calculates and returns a bitmask that can represent all possible operand combinations up to op.maxLen.
func computeBitmask(l int) int {
	bits := 0
	for (1 << bits) <= l {
		bits++
	}
	return (1 << bits) - 1
}
