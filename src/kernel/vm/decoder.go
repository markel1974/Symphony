package vm

// operandsMax defines the maximum number of operands, calculated as 2 raised to the power of 4.
// operandsMask is a bitmask derived from operandsMax, used for extracting operand-related values.
const (
	operandsMax  = 1 << 4
	operandsMask = operandsMax - 1
)

// Decoder represents an instruction decoder used to process bytecode in a virtual machine.
// It handles decoding operands and executing the associated instruction logic.
type Decoder struct {
	execute         func(vm *VM, data *Decoder)
	operands        []func(*Frame, int) (int, int)
	fullWidth       int
	decodedOperands []int
}

// NewDecoder creates a new Decoder instance with the specified execution function and operand widths.
func NewDecoder(execute func(vm *VM, data *Decoder), operands []int) *Decoder {
	sd := &Decoder{
		execute:         execute,
		operands:        make([]func(*Frame, int) (int, int), len(operands)),
		fullWidth:       0,
		decodedOperands: make([]int, operandsMax),
	}
	idx := 0
	for i := len(operands) - 1; i >= 0; i-- {
		width := operands[i]
		switch width {
		case 1:
			sd.operands[idx] = func(frame *Frame, ip int) (int, int) { return int(frame.Get8(ip)), 1 }
		case 2:
			sd.operands[idx] = func(frame *Frame, ip int) (int, int) { return int(frame.Get16(ip)), 2 }
		}
		sd.fullWidth += width
		idx++
	}
	return sd
}

// Decode processes and decodes operands from the instruction pointer, updating decodedOperands and returning new ip.
func (d *Decoder) Decode(frame *Frame, ip int) int {
	if d.fullWidth == 0 {
		return ip
	}
	ip += d.fullWidth
	readOffset := ip
	for idx, fn := range d.operands {
		val, width := fn(frame, readOffset)
		d.decodedOperands[idx] = val
		readOffset -= width
	}
	return ip
}

// Execute runs the logic associated with the current instruction using the provided virtual machine instance.
func (d *Decoder) Execute(v *VM) {
	d.execute(v, d)
}

// Read retrieves a decoded operand from the `decodedOperands` slice using a masked index derived from the input parameter.
func (d *Decoder) Read(x int) int {
	return d.decodedOperands[x&operandsMask]
}
