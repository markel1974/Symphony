package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/opcodes"
)

type OperandsDecoderData struct {
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
// Consequently, when an IOpExecutor calls decoder.Operand(N), it's accessing the N-th operand
// from the *end* of the list. For example, decoder.Operand(0) retrieves the LAST operand.
type Decoder struct {
	executor        IOpExecutor
	name            string
	execute         func(data *Decoder)
	operandsSize    uint
	decodedOperands []int
	operands        []*OperandsDecoderData
	operandsMask    int
}

// NewDecoder creates a new Decoder instance with the specified execution function and operand widths.
func NewDecoder(executor IOpExecutor) (*Decoder, error) {
	features := executor.Opcode().Operands()
	operandsBits := 0
	for (1 << operandsBits) <= len(features) {
		operandsBits++
	}
	operandsMask := (1 << operandsBits) - 1

	sd := &Decoder{
		executor:        executor,
		execute:         executor.Execute,
		name:            executor.Opcode().Name(),
		operandsMask:    operandsMask,
		operandsSize:    0,
		decodedOperands: make([]int, operandsMask+1),
	}

	idx := 0
	for i := len(features) - 1; i >= 0; i-- {
		var retrieve func(*Frame, uint) int
		width := features[i] & opcodes.SzMask
		switch width {
		case opcodes.SzUint8:
			retrieve = sd.get8Reverse
		case opcodes.SzUint16:
			retrieve = sd.get16Reverse
		case opcodes.SzUint32:
			retrieve = sd.get32Reverse
		case opcodes.SzUint64:
			retrieve = sd.get64Reverse
		default:
			return nil, fmt.Errorf("invalid operand width: %d", width)
		}
		sd.operands = append(sd.operands, &OperandsDecoderData{offset: uint(sd.operandsSize), retrieve: retrieve})
		sd.operandsSize += uint(width)
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

// OperandsSize returns the total size of the operands for the current instruction decoder instance.
func (d *Decoder) OperandsSize() uint {
	return d.operandsSize
}

// DecodeReverse processes and decodes operands from the instruction pointer, updating decodedOperands and returning new ip.
func (d *Decoder) DecodeReverse(frame *Frame, ip uint) {
	if d.operandsSize == 0 {
		return
	}
	for idx, dd := range d.operands {
		d.decodedOperands[idx] = dd.retrieve(frame, ip-dd.offset)
	}
}

// Execute runs the logic associated with the current instruction using the provided virtual machine instance.
func (d *Decoder) Execute() {
	d.execute(d)
}

// Operand retrieves a decoded operand from the `decodedOperands` slice using a masked index derived from the input parameter.
//
// IMPORTANT: Due to the decoder's reverse logic, Operand(0) accesses the LAST operand
// of the instruction, Operand(1) the second to last, and so on
func (d *Decoder) Operand(x int) int {
	return d.decodedOperands[x&d.operandsMask]
}

// get8Reverse retrieves an 8-bit integer from the frame's instructions, moving backward from the specified pointer position.
func (d *Decoder) get8Reverse(f *Frame, ip uint) int {
	return int(f.Get8Reverse(ip))
}

// get16Reverse retrieves a 16-bit signed integer from the frame's instructions using the specified instruction pointer.
func (d *Decoder) get16Reverse(f *Frame, ip uint) int {
	return int(f.Get16Reverse(ip))
}

// get32Reverse retrieves a 32-bit integer from the frame's instructions at a given position in reverse order.
func (d *Decoder) get32Reverse(f *Frame, ip uint) int {
	return int(f.Get32Reverse(ip))
}

// get64Reverse retrieves a 64-bit integer from the frame's instructions at the given instruction pointer offset.
func (d *Decoder) get64Reverse(f *Frame, ip uint) int {
	return int(f.Get64Reverse(ip))
}
