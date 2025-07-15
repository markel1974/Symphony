package mos6510_rev1

// Transfer

// InstOpTAX performs the TAX instruction, transferring the value of the accumulator (A) to the X register.
// Updates the negative (nFlag) and zero (zFlag) flags based on the value of A. Sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpTAX(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpTXA transfers the X register to the A register, updating the negative and zero flags based on the value of X.
//
//go:nosplit
func InstOpTXA(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = InstOpINI
}

// InstOpTAY transfers the value of the accumulator (A) to the Y register and updates the negative and zero flags.
//
//go:nosplit
func InstOpTAY(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpTYA transfers the value of the Y register to the A register and updates the negative and zero flags accordingly.
//
//go:nosplit
func InstOpTYA(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = InstOpINI
}

// InstOpTSX loads the stack pointer into the X register, updates the negative and zero flags, and sets the next instruction.
//
//go:nosplit
func InstOpTSX(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = InstOpINI
}

// InstOpTXS transfers the value from the X register to the stack pointer and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpTXS(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.sp = cpu.x
	cpu.next = InstOpINI
}
