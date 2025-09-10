package bytecode

//const OpcodeWidth = 1

const (
	HeaderSizeBytes     = 1
	HeaderOpcodeIdBytes = 2
)

type OperandFeature int

const (
	SzUint8  OperandFeature = 1 << 0
	SzUint16 OperandFeature = 1 << 1
	SzUint32 OperandFeature = 1 << 2
	SzUint64 OperandFeature = 1 << 3
	SzMask   OperandFeature = (1 << 4) - 1

	IsRelocatable OperandFeature = 1 << 4
	//HintForGC      OperandFeature = 1 << 5
	//HintForJIT     OperandFeature = 1 << 6
	//IsRegisterHint OperandFeature = 1 << 7
)

const (
	Relocatable = SzUint16 | IsRelocatable
)

// Opcode represents the opcode of an opcode, including its identifier, its operand, and its name.
type Opcode struct {
	opcodeId      OpcodeId
	relocate      []int
	operands      []OperandFeature
	name          string
	operandsBytes int
	compiler      *Compiler
}

// NewOpcode creates a new Opcode instance, initializing its opcode, operands, and name fields.
func NewOpcode(opcodeId OpcodeId, operands []OperandFeature, name string) *Opcode {
	od := &Opcode{
		opcodeId:      opcodeId,
		operands:      operands,
		relocate:      []int{},
		name:          name,
		operandsBytes: 0,
	}
	for idx, of := range od.operands {
		size := of & SzMask
		od.operandsBytes += int(size)
		if of&IsRelocatable != 0 {
			od.relocate = append(od.relocate, idx)
		}
	}
	od.compiler = NewCompiler(od)
	return od
}

// OpcodeId returns the opcode associated with the Opcode instance.
func (od *Opcode) OpcodeId() OpcodeId {
	return od.opcodeId
}

// Name returns the name of the opcode as a string.
func (od *Opcode) Name() string {
	return od.name
}

// Operands returns the operand widths for the Opcode as a slice of integers.
func (od *Opcode) Operands() []OperandFeature {
	return od.operands
}

// OperandsLen returns the length of the operands for the Opcode instance.
func (od *Opcode) OperandsLen() int {
	return len(od.operands)
}

// OperandsBytes returns the total width of the operands for the Opcode instance.
func (od *Opcode) OperandsBytes() int {
	return od.operandsBytes
}

// IsRelocatable returns the relocatable value associated with the Opcode instance.
func (od *Opcode) IsRelocatable() bool {
	return len(od.relocate) > 0
}

// Relocate returns the slice of integers representing relocation information for the Opcode.
func (od *Opcode) Relocate() []int {
	return od.relocate
}

// Compile compiles the opcode into a sequence of bytes.
func (od *Opcode) Compile(operands []int) ([]byte, error) {
	if err := od.compiler.Compile(nil, operands); err != nil {
		return nil, err
	}
	return od.compiler.Instructions(), nil
}

// DecompileOperands extracts operands from the provided bytecode and returns them as a slice of integers or an error.
func (od *Opcode) DecompileOperands(headerBytes uint8, bytecode []byte) ([]int, error) {
	od.compiler.SetInstructions(bytecode)
	return od.compiler.DecompileOperands(headerBytes)
}
