package z80

// OP1 -> inst.IT.operand1
// OP2 -> inst.IT.operand2)
// OFFSET -> inst.offset
// IMM -> inst.immediate

// computeADDFlags sets CPU flags based on the result of an ADD operation for 8-bit operands without carry.
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

// computeLdBlockFlags updates Z80 CPU flags based on a transferred byte during LD block instructions.
// Resets FLAG_H and FLAG_N, sets FLAG_P if BC != 0, and updates undocumented flags FLAG_X and FLAG_Y.
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

// computeAdcFlags calculates and sets the appropriate CPU flags for the ADC instruction based on the operation result.
// It processes operands, carry, and result to update flags like S, Z, H, P/V, N, C, Y, and X.
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

// instLD_I_N loads an immediate value into memory at an address calculated from a word register and an offset.
func (cpu *Z80) instLD_I_N(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) + uint16(int16(inst.offset))
	cpu.WriteMem(addr, uint8(inst.immediate), cpu)
}

// instLD_MRR_N loads an immediate 8-bit value into the memory location addressed by a 16-bit register pair.
func (cpu *Z80) instLD_MRR_N(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) // Deve puntare a HL
	cpu.WriteMem(addr, uint8(inst.immediate), cpu)
}

// instLD_I_R writes the value from a CPU register to memory at an address derived from an instruction's operands and offset.
func (cpu *Z80) instLD_I_R(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) + uint16(int16(inst.offset))
	val := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(addr, val, cpu)
}

// instLD_MRR_R handles loading a value from a register into memory at an address specified by a word register (HL, BC, or DE).
func (cpu *Z80) instLD_MRR_R(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand1)) // Deve puntare a HL, BC o DE
	val := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(addr, val, cpu)
}

// instLD_MNN_R writes the value of a register to a memory location defined by an immediate value in the instruction.
func (cpu *Z80) instLD_MNN_R(inst *Instruction) {
	result := cpu.GetByteReg(uint8(inst.IT.operand2))
	cpu.WriteMem(uint16(inst.immediate), result, cpu)
}

// instLD_R_I loads a value from memory into a CPU register based on a calculated address and instruction operands.
func (cpu *Z80) instLD_R_I(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand2)) + uint16(int16(inst.offset))
	val := cpu.ReadMem(addr, false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), val)
}

// instLD_R_MRR loads a value from the memory address pointed to by a register pair (e.g., HL, BC, DE) into a single register.
func (cpu *Z80) instLD_R_MRR(inst *Instruction) {
	addr := cpu.GetWordReg(uint8(inst.IT.operand2)) // Deve puntare a HL, BC o DE
	val := cpu.ReadMem(addr, false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), val)
}

// instLD_R_MNN loads a byte from memory at the specified address (immediate value) into a specified register.
func (cpu *Z80) instLD_R_MNN(inst *Instruction) {
	result := cpu.ReadMem(uint16(inst.immediate), false, cpu)
	cpu.SetByteReg(uint8(inst.IT.operand1), result)
}

// instLD_R_N loads an immediate value into the specified 8-bit register as defined by the instruction template.
func (cpu *Z80) instLD_R_N(inst *Instruction) {
	cpu.SetByteReg(uint8(inst.IT.operand1), uint8(inst.immediate))
}

// instLD_R_R performs the LD r, r' instruction, copying data from one register to another in the Z80 CPU.
// If the source register is I or R, it also updates flags (S, Z, H, P, N, Y, and X) based on the result.
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

// instLD_RR_MNN loads a 16-bit value from memory at the given address into a register pair.
// The address is specified using the immediate field of the instruction.
// The method reads two consecutive bytes from memory, combines them into a 16-bit value, and stores it in the register.
// The register pair is determined by the operand1 field of the instruction's template.
func (cpu *Z80) instLD_RR_MNN(inst *Instruction) {
	result := uint16(cpu.ReadMem(uint16(inst.immediate), false, cpu))
	result |= uint16(cpu.ReadMem(uint16(inst.immediate+1), false, cpu)) << 8
	cpu.SetWordReg(uint8(inst.IT.operand1), result)
}

// instLD_RR_NN loads a 16-bit immediate value into a specified register pair based on the instruction template operand.
func (cpu *Z80) instLD_RR_NN(inst *Instruction) {
	cpu.SetWordReg(uint8(inst.IT.operand1), uint16(inst.immediate))
}

// instLD_RR_RR loads the value from a 16-bit source register to a 16-bit destination register as specified in the instruction.
func (cpu *Z80) instLD_RR_RR(inst *Instruction) {
	cpu.SetWordReg(uint8(inst.IT.operand1), cpu.GetWordReg(uint8(inst.IT.operand2)))
}

// instLD_MNN_RR stores a 16-bit value from a register into two consecutive memory locations specified by an immediate.
func (cpu *Z80) instLD_MNN_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.WriteMem(uint16(inst.immediate), uint8(result&0xff), cpu)
	cpu.WriteMem(uint16(inst.immediate+1), uint8(result>>8), cpu)
}

// instPOP_RR executes the POP instruction, retrieving a 16-bit value from the stack into a register pair.
func (cpu *Z80) instPOP_RR(inst *Instruction) {
	result := uint16(cpu.ReadMem(cpu.GetSP(), false, cpu))
	cpu.IncrementSP()
	result |= uint16(cpu.ReadMem(cpu.GetSP(), false, cpu)) << 8
	cpu.IncrementSP()
	cpu.SetWordReg(uint8(inst.IT.operand1), result)
}

// instPUSH_RR stores the content of a 16-bit register pair onto the stack in high-byte-first order.
func (cpu *Z80) instPUSH_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand1))
	cpu.DecrementSP()
	cpu.WriteMem(cpu.GetSP(), uint8(result>>8), cpu)
	cpu.DecrementSP()
	cpu.WriteMem(cpu.GetSP(), uint8(result&0xff), cpu)
}

// instCPI esegue la comparazione A-(HL), incrementa HL, decrementa BC e imposta i flag.
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

// instCPD esegue la comparazione A-(HL), decrementa HL, decrementa BC e imposta i flag.
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

// instCPDR performs the CPDR instruction, comparing the accumulator with the memory value at HL and updating flags.
// It decrements HL and BC, adjusts specific CPD flags, and repeats until BC is zero or zero flag is set.
// Unique to this method, HL is decremented after each iteration instead of incremented.
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

// instCPIR executes the CPIR instruction, which compares A with the byte at HL, updates HL and BC, and repeats if needed.
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

func (cpu *Z80) instEX_MRR_RR(inst *Instruction) {
	result := cpu.GetWordReg(uint8(inst.IT.operand1))
	op1 := uint16(cpu.ReadMem(result, false, cpu))
	op1 |= uint16(cpu.ReadMem(result+1, false, cpu)) << 8
	op2 := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.WriteMem(result, uint8(op2&0xff), cpu)
	cpu.WriteMem(result+1, uint8(op2>>8), cpu)
	cpu.SetWordReg(uint8(inst.IT.operand2), op1)
}

func (cpu *Z80) instEX_RR_RR(inst *Instruction) {
	result1 := cpu.GetWordReg(uint8(inst.IT.operand1))
	result2 := cpu.GetWordReg(uint8(inst.IT.operand2))
	cpu.SetWordReg(uint8(inst.IT.operand1), result2)
	cpu.SetWordReg(uint8(inst.IT.operand2), result1)
}

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

func (cpu *Z80) instLDD(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() - 1)
	cpu.SetDE(cpu.GetDE() - 1)
	cpu.SetBC(cpu.GetBC() - 1)

	cpu.computeLdBlockFlags(transferredByte)
}

func (cpu *Z80) instLDI(_ *Instruction) {
	transferredByte := cpu.ReadMem(cpu.GetHL(), false, cpu)
	cpu.WriteMem(cpu.GetDE(), transferredByte, cpu)
	cpu.SetHL(cpu.GetHL() + 1)
	cpu.SetDE(cpu.GetDE() + 1)
	cpu.SetBC(cpu.GetBC() - 1)
	cpu.computeLdBlockFlags(transferredByte)
}

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

// Z80: LDIR
// T-States: 21T (se BC != 0 dopo il decremento) / 16T (se BC = 0)
// M-Cycles (per iterazione): La sequenza di base è Lettura da (HL) e Scrittura in (DE).
//
//	Se BC != 0, vengono aggiunti cicli macchina per ricaricare il PC e ripetere.
//
// Desc: Copia un byte da (HL) a (DE), incrementa HL, incrementa DE e decrementa BC.
//
//	Se BC non è zero, ripete l'intera istruzione (decrementando il PC di 2).
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

// Z80: ADC A, (HL)
// T-States: 7T
// M-Cycles: M1(4T) [Fetch] + M2(3T) [Lettura Memoria]
// ---
// Z80: ADC A, (IX+d)  /  ADC A, (IY+d)
// T-States: 19T
// M-Cycles: M1(4T)+M2(4T) [Fetch] + M3(3T) [Lettura offset d] + M4(5T) [Calcolo addr] + M5(3T) [Lettura Memoria]
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

// ADC_R_I (ADC R, n)
// T-States: 7T
// M-Cycles: M1(4T) [Fetch opcode] + M2(3T) [Lettura immediato n]
// Desc: Somma il valore immediato 'n' e il flag Carry al registro A.
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

// instADC_R_N (ADC r, n)
// T-States: 7T
// M-Cycles: M1(4T) [Fetch] + M2(3T) [Lettura immediato n]
// Desc: Somma il valore immediato 'n' e il flag Carry al registro A.
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

// Z80: ADC A, r   (dove r = B, C, D, E, H, L, A)
// T-States: 4T
// M-Cycles: M1(4T) [Fetch]
// Desc: Somma il valore del registro 'r' e il flag Carry al registro A.
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

// Z80: ADD A, (HL)
// T-States: 7T
// M-Cycles: M1(4T) + M2(3T)
// ---
// Z80: ADD A, (IX+d) / ADD A, (IY+d)
// T-States: 19T
// M-Cycles: M1(4T)+M2(4T)+M3(3T)+M4(5T)+M5(3T)
func (cpu *Z80) instADD_R_I(inst *Instruction) {
	op2 := uint16(inst.immediate)
	op1 := uint16(cpu.GetA())
	result := op1 + op2
	cpu.SetA(uint8(result & 0xff))
	cpu.computeADDFlags(op1, op2, result)
}

// Z80: ADD A, (HL)
// T-States: 7T
// M-Cycles: M1(4T) [Fetch] + M2(3T) [Lettura Memoria]
// ---
// Z80: ADD A, (IX+d) / ADD A, (IY+d)
// T-States: 19T
// M-Cycles: M1(4T)+M2(4T) [Fetch] + M3(3T) [Lettura offset] + M4(5T) [Calcolo addr] + M5(3T) [Lettura Memoria]
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
