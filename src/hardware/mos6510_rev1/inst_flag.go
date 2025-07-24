package mos6510_rev1

// Flag

// InstOpSEC sets the Carry flag (cFlag) to 1 and moves execution to the next instruction handler (InstOpINI).
//
//go:nosplit
func (er *Executor) InstOpSEC(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 1
	cpu.next = er.InstOpINI
}

// InstOpCLC clears the carry flag in the CPU and sets the next instruction to InstOpINI. It halts if the current PC read fails.
func (er *Executor) InstOpCLC(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 0
	cpu.next = er.InstOpINI
}

// InstOpSED sets the decimal mode flag (dFlag) to 1 and assigns the next instruction handler to InstOpINI.
func (er *Executor) InstOpSED(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 1
	cpu.next = er.InstOpINI
}

// InstOpCLD clears the decimal mode flag (dFlag) and sets the next instruction handler to InstOpINI.
func (er *Executor) InstOpCLD(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 0
	cpu.next = er.InstOpINI
}

// InstOpSEI sets the interrupt disable flag and updates CPU state to handle the next instruction cycle.
func (er *Executor) InstOpSEI(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = er.InstOpINI
}

// InstOpCLI clears the interrupt disable flag and sets the next instruction to InstOpINI if the current opcode is valid.
func (er *Executor) InstOpCLI(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= OpFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = er.InstOpINI
}

// InstOpCLV clears the overflow flag in the CPU and sets the next instruction handler to InstOpINI.
// If the current PC address cannot be read, the operation is aborted.
func (er *Executor) InstOpCLV(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.vFlag = 0
	cpu.next = er.InstOpINI
}

// InstOpNOP is a no-operation function for the CPU that progresses the state to the next instruction without modifying registers.
func (er *Executor) InstOpNOP(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.next = er.InstOpINI
}
