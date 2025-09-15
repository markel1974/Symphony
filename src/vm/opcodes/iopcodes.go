package opcodes

// IOpcodes is an interface for managing and compiling operation codes in an instruction set.
// Id returns a string identifier for the opcode manager.
// Opcode fetches the Opcode definition for a given OpcodeId.
// Compile compiles an OpcodeId and its operands into a byte slice, returning an error if compilation fails.
// Mask returns an integer mask applied to determine opcode categories.
// Len returns the total number of opcodes managed.
type IOpcodes interface {
	Id() string

	Opcode(opcodeId OpcodeId) *Opcode

	Bytecode(opcodeId OpcodeId, operands ...int) ([]byte, error)

	Mask() int

	Len() int
}
