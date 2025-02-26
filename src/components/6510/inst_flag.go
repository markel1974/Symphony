package mos6510

// Flag

// instOpSEC sets the Carry flag (cFlag) to 1 and moves execution to the next instruction handler (instOpINI).
func instOpSEC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 1
	cpu.next = instOpINI
}

// instOpCLC clears the carry flag in the CPU and sets the next instruction to instOpINI. It halts if the current PC read fails.
func instOpCLC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 0
	cpu.next = instOpINI
}

// instOpSED sets the decimal mode flag (dFlag) to 1 and assigns the next instruction handler to instOpINI.
func instOpSED(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 1
	cpu.next = instOpINI
}

// instOpCLD clears the decimal mode flag (dFlag) and sets the next instruction handler to instOpINI.
func instOpCLD(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 0
	cpu.next = instOpINI
}

// instOpSEI sets the interrupt disable flag and updates CPU state to handle the next instruction cycle.
func instOpSEI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = instOpINI
}

// instOpCLI clears the interrupt disable flag and sets the next instruction to instOpINI if the current opcode is valid.
func instOpCLI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = instOpINI
}

// instOpCLV clears the overflow flag in the CPU and sets the next instruction handler to instOpINI.
// If the current PC address cannot be read, the operation is aborted.
func instOpCLV(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.vFlag = 0
	cpu.next = instOpINI
}

// instOpNOP is a no-operation function for the CPU that progresses the state to the next instruction without modifying registers.
func instOpNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpINI
}
