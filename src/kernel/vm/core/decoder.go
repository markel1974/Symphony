package core

// operandsMax defines the maximum number of operands, calculated as 2 raised to the power of 4.
// operandsMask is a bitmask derived from operandsMax, used for extracting operand-related values.
const (
	operandsMax  = 1 << 4
	operandsMask = operandsMax - 1
)

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
	execute         func(vm *VM, data *Decoder)
	fullWidth       int
	operands        []int
	decodedOperands []int
	//fnOperands      []func(*Frame, int) (int, int)
}

// NewDecoder creates a new Decoder instance with the specified execution function and operand widths.
func NewDecoder(executor IOpExecutor) *Decoder {
	operands := executor.Operands()
	sd := &Decoder{
		executor:        executor,
		execute:         executor.Execute,
		name:            executor.Name(),
		operands:        make([]int, len(operands)),
		fullWidth:       0,
		decodedOperands: make([]int, operandsMax),
		//fnOperands:      make([]func(*Frame, int) (int, int), len(operands)),
	}
	idx := 0
	for i := len(operands) - 1; i >= 0; i-- {
		sd.operands[idx] = operands[i]
		width := operands[i]
		//switch width {
		//case 1:
		//	sd.fnOperands[idx] = func(frame *Frame, ip int) (int, int) { return int(frame.Get8(ip)), 1 }
		//ase 2:
		//	sd.fnOperands[idx] = func(frame *Frame, ip int) (int, int) { return int(frame.Get16(ip)), 2 }
		//}
		sd.fullWidth += width
		idx++
	}
	return sd
}

// Name returns the name of the instruction.
func (d *Decoder) Name() string {
	return d.name
}

func (d *Decoder) Decode(frame *Frame, ip int) int {
	if d.fullWidth == 0 {
		return ip
	}
	ip += d.fullWidth
	readOffset := ip
	for idx, width := range d.operands {
		switch width {
		case 1:
			d.decodedOperands[idx] = int(frame.Get8(readOffset))
		case 2:
			d.decodedOperands[idx] = int(frame.Get8(readOffset))
		}
		readOffset -= width
	}
	return ip
}

// Execute runs the logic associated with the current instruction using the provided virtual machine instance.
func (d *Decoder) Execute(v *VM) {
	d.execute(v, d)
}

// Read retrieves a decoded operand from the `decodedOperands` slice using a masked index derived from the input parameter.
//
// IMPORTANT: Due to the decoder's reverse logic, Read(0) accesses the LAST operand
// of the instruction, Read(1) the second to last, and so on
func (d *Decoder) Read(x int) int {
	return d.decodedOperands[x&operandsMask]
}

/*
// Decode processes and decodes operands from the instruction pointer, updating decodedOperands and returning new ip.
func (d *Decoder) Decode2(frame *Frame, ip int) int {
	if d.fullWidth == 0 {
		return ip
	}
	ip += d.fullWidth
	readOffset := ip
	for idx, fn := range d.fnOperands {
		val, width := fn(frame, readOffset)
		d.decodedOperands[idx] = val
		readOffset -= width
	}
	return ip
}
*/
