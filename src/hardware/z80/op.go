package z80

// OP1 -> inst.IT.operand1
// OP2 -> inst.IT.operand2)
// OFFSET -> inst.offset
// IMM -> inst.immediate

// computeADDFlags updates the CPU flags based on the result of an addition operation involving two operands.
// It sets or resets the Sign, Zero, Half Carry, Overflow/Parity, Carry, and undocumented Y/X flags accordingly.
func (cpu *Z80) computeADDFlags(result uint16, op1 uint16, op2 uint16) {
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, result&0xff == 0)
	// Half Carry per ADD (senza carry)
	cpu.SetFlagValue(FLAG_H, (op1&0x0f)+(op2&0x0f) > 0x0f)
	// Logica di overflow corretta per ADD
	isOverflow := ((op1 ^ result) & (op2 ^ result) & 0x80) != 0
	cpu.SetFlagValue(FLAG_P, isOverflow)
	cpu.ResetFlag(FLAG_N)
	cpu.SetFlagValue(FLAG_C, result > 0xff)
	// Flag non documentati
	cpu.SetFlagValue(FLAG_Y, result&0x20 != 0)
	cpu.SetFlagValue(FLAG_X, result&0x08 != 0)
}

// computeLdBlockFlags updates specific Z80 CPU flags based on a transferred byte during a block transfer operation.
// The method resets the H and N flags, updates the P flag based on the BC register, and sets undocumented flags X and Y.
// The X and Y flags are determined by bits 1 and 3 of the sum of the accumulator and the transferred byte.
func (cpu *Z80) computeLdBlockFlags(transferredByte uint8) {
	// I flag H (Half Carry) e N (Subtract) vengono sempre resettati.
	cpu.ResetFlag(FLAG_H)
	cpu.ResetFlag(FLAG_N)
	// Il flag P/V (Parity/Overflow) indica se il contatore BC non si è ancora azzerato.
	cpu.SetFlagValue(FLAG_P, cpu.GetBC() != 0)
	// I flag non documentati dipendono dalla somma di A + il dato appena trasferito.
	sum := transferredByte + cpu.GetA()
	cpu.SetFlagValue(FLAG_Y, sum&0x02 != 0) // bit 1
	cpu.SetFlagValue(FLAG_X, sum&0x08 != 0) // bit 3
}

// computeAdcFlags updates the CPU flags after an ADD with carry operation.
// It sets or clears flags S, Z, H, P/V (overflow), C, Y, and X based on the given operands and result.
func (cpu *Z80) computeAdcFlags(op1 uint16, op2 uint16, carry uint16, result uint16) {
	// flag settings
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, result&0xff == 0)
	cpu.SetFlagValue(FLAG_H, (op1&0x0f)+(op2&0x0f)+carry > 0x0f)
	// overflow logic
	isOverflow := ((op1 ^ result) & (op2 ^ result) & 0x80) != 0
	cpu.SetFlagValue(FLAG_P, isOverflow)
	// operation flags
	cpu.ResetFlag(FLAG_N)
	cpu.SetFlagValue(FLAG_C, result > 0xff)
	// Flags undocumented
	cpu.SetFlagValue(FLAG_Y, result&0x20 != 0) // bit 5
	cpu.SetFlagValue(FLAG_X, result&0x08 != 0) // bit 3
}

// instLD_I_N loads the immediate value into the memory address computed by adding a register's value and an offset.
func (cpu *Z80) instLD_I_N(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) + uint16(int16(inst.offset))
	cpu.WriteMem(addr, uint8(inst.immediate), cpu)
}

// instLD_MRR_N executes the LD (RR), N instruction, writing an immediate value to memory at the address stored in a word register.
// It calculates the memory address using the operand1 register and writes the provided immediate value to that address.
func (cpu *Z80) instLD_MRR_N(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) // Deve puntare a HL
	cpu.WriteMem(addr, uint8(inst.immediate), cpu)
}

// instLD_I_R writes the value from the specified register into memory at the address computed with offset and register.
func (cpu *Z80) instLD_I_R(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) + uint16(int16(inst.offset))
	val := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(addr, val, cpu)
}

// instLD_MRR_R handles the LD (address pointed by a register pair) <- (value from a register) Z80 CPU instruction.
// It writes the value from a specific CPU register to the memory address stored in a register pair (HL, BC, or DE).
func (cpu *Z80) instLD_MRR_R(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) // Deve puntare a HL, BC o DE
	val := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(addr, val, cpu)
}

// instLD_MNN_R writes a value from a CPU register to the memory at the address specified by the instruction's immediate value.
func (cpu *Z80) instLD_MNN_R(inst *Instruction) {
	result := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(uint16(inst.immediate), result, cpu)
}

// instLD_R_I executes the LD R, (I) operation, loading a value from memory at a calculated address into a register.
func (cpu *Z80) instLD_R_I(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand2)) + uint16(int16(inst.offset))
	val := cpu.ReadMem(addr, false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), val)
}

// instLD_R_MRR loads a byte from the memory address specified by register pair (HL, BC, or DE) into a target register.
func (cpu *Z80) instLD_R_MRR(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand2)) // Deve puntare a HL, BC o DE
	val := cpu.ReadMem(addr, false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), val)
}

// instLD_R_MNN loads a byte from memory at the address specified by the instruction's immediate value into a CPU register.
func (cpu *Z80) instLD_R_MNN(inst *Instruction) {
	result := cpu.ReadMem(uint16(inst.immediate), false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), result)
}

// instLD_R_N loads an immediate value into a target register specified by the instruction's first operand.
func (cpu *Z80) instLD_R_N(inst *Instruction) {
	cpu.SetByteReg(uint8(inst.IT.operand1), uint8(inst.immediate))
}

// instLD_R_R loads the value from a specified source register into a destination register in the CPU.
// If the source register is REG_I or REG_R, flag values are updated based on the resulting value and CPU state.
func (cpu *Z80) instLD_R_R(inst *Instruction) {
	result := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.SetByteReg(uint8(inst.IT.operand1), result)

	// Logica speciale per LD A, I e LD A, R
	if inst.IT.operand2 == REG_I || inst.IT.operand2 == REG_R {
		cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
		cpu.SetFlagValue(FLAG_Z, result == 0)
		cpu.SetFlagValue(FLAG_Y, result&0x20 != 0)
		cpu.ResetFlag(FLAG_H)
		cpu.SetFlagValue(FLAG_X, result&0x08 != 0)
		cpu.SetFlagValue(FLAG_P, cpu.iff2)
		cpu.ResetFlag(FLAG_N)
	}
}

// instLD_RR_MNN loads a 16-bit value from memory at the address specified by the instruction's immediate field.
// The loaded value is then stored into the specified 16-bit register pair.
func (cpu *Z80) instLD_RR_MNN(inst *Instruction) {
	result := uint16(cpu.ReadMem(uint16(inst.immediate), false, cpu))
	result |= uint16(cpu.ReadMem(uint16(inst.immediate+1), false, cpu)) << 8
	cpu.SetWordReg(uint8(inst.IT.operand1), result)
}

// instLD_RR_NN loads a 16-bit immediate value into a specified 16-bit register pair.
func (cpu *Z80) instLD_RR_NN(inst *Instruction) {
	cpu.SetWordReg(uint8(inst.IT.operand1), uint16(inst.immediate))
}

// instLD_RR_RR loads the value from the specified source register pair into the specified destination register pair.
func (cpu *Z80) instLD_RR_RR(inst *Instruction) {
	cpu.SetWordReg(uint8(inst.IT.operand1), cpu.GetWordReg(uint8(inst.IT.operand2)))
}

// instLD_MNN_RR writes the value of a specified register pair to memory at an immediate 16-bit address.
// The lower byte is written first, followed by the higher byte.
func (cpu *Z80) instLD_MNN_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.WriteMem(uint16(inst.immediate), uint8(result&0xff), cpu)
	cpu.WriteMem(uint16(inst.immediate+1), uint8(result>>8), cpu)
}

// instPOP_RR pops a 16-bit value from the stack into the specified register pair in the Z80 CPU.
func (cpu *Z80) instPOP_RR(inst *Instruction) {
	result := uint16(cpu.ReadMem(cpu.GetSP(), false, cpu))
	cpu.IncrementSP()
	result |= uint16(cpu.ReadMem(cpu.GetSP(), false, cpu)) << 8
	cpu.IncrementSP()
	cpu.SetWordReg(uint8(inst.IT.operand1), result)
}

// instPUSH_RR pushes the contents of a 16-bit register to the stack, decrementing the stack pointer appropriately.
func (cpu *Z80) instPUSH_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand1))
	cpu.DecrementSP()
	cpu.WriteMem(cpu.GetSP(), uint8(result>>8), cpu)
	cpu.DecrementSP()
	cpu.WriteMem(cpu.GetSP(), uint8(result&0xff), cpu)
}

// instCPI executes the CPI (Compare and Increment) instruction for the Z80 processor.
// It compares the value in the accumulator (A) with the value at the memory address pointed to by HL.
// The HL register pair is incremented, and the BC register pair is decremented as part of the operation.
// This method sets or clears condition flags (S, Z, H, P/V, N, X, Y) based on comparison results.
func (cpu *Z80) instCPI(_ *Instruction) {
	// 1. Esegui la comparazione
	op1 := cpu.GetA()
	op2 := cpu.ReadMem(cpu.GetHL(), false, cpu)
	result := int16(op1) - int16(op2)

	// 2. Aggiorna i registri principali
	cpu.SetHL(cpu.GetHL() + 1)
	cpu.SetBC(cpu.GetBC() - 1)

	// 3. Calcola i flag standard e speciali per CPI
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, byte(result) == 0)
	cpu.SetFlagValue(FLAG_H, (op1&0x0F) < (op2&0x0F))
	cpu.SetFlagValue(FLAG_P, cpu.GetBC() != 0) // P/V indica se BC è diventato non-zero
	cpu.SetFlag(FLAG_N)                        // N è sempre impostato a 1

	// 4. Calcola i flag non documentati X e Y
	// Il valore è basato sul risultato della sottrazione, corretto con il flag Half-Carry
	adjustedResult := byte(result)
	if cpu.FlagIsSet(FLAG_H) {
		adjustedResult--
	}
	cpu.SetFlagValue(FLAG_Y, adjustedResult&0x02 != 0) // Flag Y è il bit 1 del risultato corretto
	cpu.SetFlagValue(FLAG_X, adjustedResult&0x08 != 0) // Flag X è il bit 3 del risultato corretto
}

// instCPD executes the "Compare and Decrement" operation for the Z80 CPU, modifying HL, BC, and relevant flags.
func (cpu *Z80) instCPD(_ *Instruction) {
	// 1. Esegui la comparazione
	op1 := cpu.GetA()
	op2 := cpu.ReadMem(cpu.GetHL(), false, cpu)
	result := int16(op1) - int16(op2)

	// 2. Aggiorna i registri principali (unica differenza con CPI)
	cpu.SetHL(cpu.GetHL() - 1)
	cpu.SetBC(cpu.GetBC() - 1)

	// 3. Calcola i flag standard e speciali per CPD
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, byte(result) == 0)
	cpu.SetFlagValue(FLAG_H, (op1&0x0F) < (op2&0x0F))
	cpu.SetFlagValue(FLAG_P, cpu.GetBC() != 0) // P/V indica se BC è diventato non-zero
	cpu.SetFlag(FLAG_N)                        // N è sempre impostato a 1

	// 4. Calcola i flag non documentati X e Y
	adjustedResult := byte(result)
	if cpu.FlagIsSet(FLAG_H) {
		adjustedResult--
	}
	cpu.SetFlagValue(FLAG_Y, adjustedResult&0x02 != 0) // Flag Y è il bit 1 del risultato corretto
	cpu.SetFlagValue(FLAG_X, adjustedResult&0x08 != 0) // Flag X è il bit 3 del risultato corretto
}

// instCPDR executes the CPDR instruction, comparing A with memory at HL and updating flags and registers accordingly.
func (cpu *Z80) instCPDR(_ *Instruction) {
	// Logica quasi identica a CPIR, cambia solo l'operazione su HL
	// 1. Esegui la comparazione
	op1 := cpu.GetA()
	op2 := cpu.ReadMem(cpu.GetHL(), false, cpu)
	result := int16(op1) - int16(op2)

	// 2. Aggiorna i registri HL e BC
	cpu.SetHL(cpu.GetHL() - 1) // L'unica differenza
	cpu.SetBC(cpu.GetBC() - 1)

	// 3. Calcolo dei flag SPECIFICO per CPD/CPDR
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, byte(result) == 0)
	cpu.SetFlagValue(FLAG_H, (op1&0x0F) < (op2&0x0F))
	cpu.SetFlagValue(FLAG_P, cpu.GetBC() != 0) // P/V riflette BC
	cpu.SetFlag(FLAG_N)

	// Undocumented flags
	if cpu.FlagIsSet(FLAG_H) {
		result--
	}
	cpu.SetFlagValue(FLAG_Y, result&0x02 != 0)
	cpu.SetFlagValue(FLAG_X, result&0x08 != 0)

	// 4. Logica di ripetizione
	if cpu.GetBC() != 0 && !cpu.FlagIsSet(FLAG_Z) {
		cpu.SetPC(cpu.GetPC() - 2)
	}
}

// instCPIR executes the CPIR instruction, comparing the accumulator with the value at HL and repeating if BC != 0 and Z is unset.
func (cpu *Z80) instCPIR(_ *Instruction) {
	// 1. Esegui la comparazione
	op1 := cpu.GetA()
	op2 := cpu.ReadMem(cpu.GetHL(), false, cpu)
	result := int16(op1) - int16(op2)
	// 2. Aggiorna i registri HL e BC
	cpu.SetHL(cpu.GetHL() + 1)
	cpu.SetBC(cpu.GetBC() - 1)
	// 3. Calcolo dei flag SPECIFICO per CPI/CPIR
	cpu.SetFlagValue(FLAG_S, result&0x80 != 0)
	cpu.SetFlagValue(FLAG_Z, byte(result) == 0)
	cpu.SetFlagValue(FLAG_H, (op1&0x0F) < (op2&0x0F))
	// Ecco la logica chiave: P/V riflette lo stato di BC
	cpu.SetFlagValue(FLAG_P, cpu.GetBC() != 0)
	cpu.SetFlag(FLAG_N) // Flag di sottrazione sempre impostato
	// Undocumented flags
	if cpu.FlagIsSet(FLAG_H) {
		result--
	}
	cpu.SetFlagValue(FLAG_Y, result&0x02 != 0) // bit 1 del risultato corretto
	cpu.SetFlagValue(FLAG_X, result&0x08 != 0) // bit 3 del risultato corretto

	// 4. Logica di ripetizione (il tuo codice era già corretto)
	if cpu.GetBC() != 0 && !cpu.FlagIsSet(FLAG_Z) {
		cpu.SetPC(cpu.GetPC() - 2)
	}
}

// instEX_MRR_RR handles the EX operation between a memory location (specified by a register pair) and a register pair.
func (cpu *Z80) instEX_MRR_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand1))
	op1 := uint16(cpu.ReadMem(result, false, cpu))
	op1 |= uint16(cpu.ReadMem(result+1, false, cpu)) << 8
	op2 := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.WriteMem(result, uint8(op2&0xff), cpu)
	cpu.WriteMem(result+1, uint8(op2>>8), cpu)
	cpu.SetWordReg(uint8(inst.IT.operand2), op1)
}

// instEX_RR_RR exchanges the contents of two 16-bit registers specified in the instruction's operands.
func (cpu *Z80) instEX_RR_RR(inst *Instruction) {
	result1 := cpu.GetWordReg(uint8(inst.IT.operand1))
	result2 := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.SetWordReg(uint8(inst.IT.operand1), result2)
	cpu.SetWordReg(uint8(inst.IT.operand2), result1)
}

// instEXX exchanges the BC, DE, and HL register values with their alternate counterparts (BC', DE', HL').
func (cpu *Z80) instEXX(_ *Instruction) {
	bc := cpu.GetBC()
	de := cpu.GetDE()
	hl := cpu.GetHL()
	cpu.SetBC(cpu.GetWordReg(REG_BCP))
	cpu.SetDE(cpu.GetWordReg(REG_DEP))
	cpu.SetHL(cpu.GetWordReg(REG_HLP))
	cpu.SetWordReg(uint8(REG_BCP), bc)
	cpu.SetWordReg(uint8(REG_DEP), de)
	cpu.SetWordReg(uint8(REG_HLP), hl)
}

// instLDD executes the LDD instruction, transferring a byte of data from HL to DE and updating relevant registers and flags.
func (cpu *Z80) instLDD(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() - 1)
	cpu.SetDE(cpu.GetDE() - 1)
	cpu.SetBC(cpu.GetBC() - 1)

	cpu.computeLdBlockFlags(transferredByte)
}

// instLDI implements the LDI instruction of the Z80 CPU, transferring a byte of data from one memory location to another.
// The method increments the HL and DE registers, decrements the BC register, and updates flags accordingly.
func (cpu *Z80) instLDI(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() + 1)
	cpu.SetDE(cpu.GetDE() + 1)
	cpu.SetBC(cpu.GetBC() - 1)
	cpu.computeLdBlockFlags(transferredByte)
}

// instLDDR executes the LDDR instruction, transferring a byte from memory at HL to DE, decrementing HL, DE, and BC.
// It loops until BC reaches zero and adjusts the program counter for repetition. Updates flags based on the transfer.
func (cpu *Z80) instLDDR(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() - 1)
	cpu.SetDE(cpu.GetDE() - 1)
	cpu.SetBC(cpu.GetBC() - 1)
	if cpu.GetBC() != 0 {
		cpu.SetPC(cpu.GetPC() - 2)
	}
	cpu.computeLdBlockFlags(transferredByte)
}

// instLDIR executes the LDIR instruction on the Z80 CPU, transferring a byte from HL to DE and updating flags.
func (cpu *Z80) instLDIR(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() + 1)
	cpu.SetDE(cpu.GetDE() + 1)
	cpu.SetBC(cpu.GetBC() - 1)
	// 2. Logica di ripetizione
	if cpu.GetBC() != 0 {
		cpu.SetPC(cpu.GetPC() - 2)
	}
	cpu.computeLdBlockFlags(transferredByte)
}

// instADC_R_MRR adds the value from memory (referenced via a register and possible offset) to a register with carry.
// Reads the value from memory based on the register and offset specified in the instruction.
// Retrieves the carry flag and adjusts the result accordingly.
// Updates the target register with the result of the addition operation.
// Computes and sets the CPU flags based on the result of the operation.
func (cpu *Z80) instADC_R_MRR(inst *Instruction) {
	op2 := uint16(cpu.ReadMem(cpu.GetWordReg(uint8(inst.IT.operand2))+uint16(inst.offset), false, cpu))
	carry := uint16(0)
	if cpu.FlagIsSet(FLAG_C) {
		carry = 1
	}
	op1 := uint16(cpu.GetByteReg(uint8(inst.IT.operand1)))
	result := op1 + op2 + carry
	cpu.SetByteReg(uint8(inst.IT.operand1), uint8(result&0xff))
	cpu.computeAdcFlags(op1, op2, carry, result)
}

// instADC_R_I performs an add with carry operation using the accumulator, an immediate value, and the carry flag.
// Updates the accumulator with the result and sets flags based on the operation.
func (cpu *Z80) instADC_R_I(inst *Instruction) {
	op1 := uint16(cpu.GetA())
	op2 := uint16(inst.immediate)
	carry := uint16(0)
	if cpu.FlagIsSet(FLAG_C) {
		carry = 1
	}
	result := op1 + op2 + carry
	cpu.SetA(uint8(result & 0xff))
	cpu.computeAdcFlags(op1, op2, carry, result)
}

// instADC_R_N performs the ADC operation with an 8-bit register and an immediate value, including the carry flag.
// The result is stored back into the specified register, and flags are updated accordingly.
func (cpu *Z80) instADC_R_N(inst *Instruction) {
	op2 := uint16(inst.immediate)
	op1 := uint16(cpu.GetByteReg(uint8(inst.IT.operand1)))
	carry := uint16(0)
	if cpu.FlagIsSet(FLAG_C) {
		carry = 1
	}
	result := op1 + op2 + carry
	cpu.SetByteReg(uint8(inst.IT.operand1), uint8(result&0xff))
	cpu.computeAdcFlags(op1, op2, carry, result)
}

// instADC_R_R performs the ADC (Add with Carry) operation between two registers and updates flags accordingly.
// The destination register is always the accumulator (register A).
// It adds the source register value, destination register value, and carry flag, storing the result in register A.
// Flags affected include the carry, half-carry, zero, sign, overflow, and subtract flags based on the result.
func (cpu *Z80) instADC_R_R(inst *Instruction) {
	op2 := uint16(cpu.GetByteReg(uint8(inst.IT.operand2)))
	op1 := uint16(cpu.GetByteReg(uint8(inst.IT.operand1))) // Destinazione è sempre A
	carry := uint16(0)
	if cpu.FlagIsSet(FLAG_C) {
		carry = 1
	}
	result := op1 + op2 + carry
	cpu.SetByteReg(uint8(inst.IT.operand1), uint8(result&0xff))
	cpu.computeAdcFlags(op1, op2, carry, result)
}

// instADD_R_I performs an addition between the contents of the A register and an immediate value provided by the instruction.
// The result is stored back in the A register, and the operation flags are updated accordingly.
func (cpu *Z80) instADD_R_I(inst *Instruction) {
	op2 := uint16(inst.immediate)
	op1 := uint16(cpu.GetA())
	result := op1 + op2
	cpu.SetA(uint8(result & 0xff))
	cpu.computeADDFlags(op1, op2, result)
}

// instADD_R_MRR adds the value at a memory address to the A register and sets flags accordingly.
// The memory address is calculated using a register plus an optional offset.
func (cpu *Z80) instADD_R_MRR(inst *Instruction) {
	// Calcola l'indirizzo di memoria da cui leggere il secondo operando.
	// Gestisce sia (HL) con offset=0, sia (IX/IY+d).
	addr := cpu.GetWordReg(uint8(inst.IT.operand2)) + uint16(int16(inst.offset))
	op2 := uint16(cpu.ReadMem(addr, false, cpu))
	// Il primo operando (e destinazione) è SEMPRE il registro A.
	op1 := uint16(cpu.GetA())
	// Esegue la somma.
	result := op1 + op2
	// Il risultato viene SEMPRE salvato nel registro A.
	cpu.SetA(uint8(result & 0xff))
	cpu.computeADDFlags(op1, op2, result)
}

// -------------------
