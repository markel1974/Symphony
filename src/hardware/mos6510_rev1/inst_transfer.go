package mos6510_rev1

// Transfer

// instOpTAX performs the TAX instruction, transferring the value of the accumulator (A) to the X register.
// Updates the negative (nFlag) and zero (zFlag) flags based on the value of A. Sets the next instruction to instOpINI.
func instOpTAX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpTXA transfers the X register to the A register, updating the negative and zero flags based on the value of X.
func instOpTXA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpTAY transfers the value of the accumulator (A) to the Y register and updates the negative and zero flags.
func instOpTAY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpTYA transfers the value of the Y register to the A register and updates the negative and zero flags accordingly.
func instOpTYA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpTSX loads the stack pointer into the X register, updates the negative and zero flags, and sets the next instruction.
func instOpTSX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = instOpINI
}

// instOpTXS transfers the value from the X register to the stack pointer and sets the next instruction to instOpINI.
func instOpTXS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.sp = cpu.x
	cpu.next = instOpINI
}
