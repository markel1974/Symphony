package bytecode

// CompileInstruction returns a bytecode for an opcode and the operands.
func CompileInstruction(opcode Opcode, operands ...int) []byte {
	numOperands := OpcodeToOperands(opcode)
	totalLen := 1
	for _, w := range numOperands {
		totalLen += w
	}
	instruction := make([]byte, totalLen)
	instruction[0] = opcode
	offset := 1
	for i, o := range operands {
		width := numOperands[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			n := uint16(o)
			instruction[offset] = byte(n >> 8)
			instruction[offset+1] = byte(n)
		}
		offset += width
	}
	return instruction
}
