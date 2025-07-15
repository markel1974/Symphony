package mos6510_rev1

// Arithmetic

// InstOpADC performs the ADC (Add with Carry) operation by reading data from memory and calling the ADC handler.
//
//go:nosplit
func InstOpADC(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.doADC(data)
	cpu.next = InstOpINI
}

// InstOiADC executes the ADC instruction at the current program counter, updating the accumulator and advancing the PC.
// Invokes the next instruction handler after performing the operation.
//
//go:nosplit
func InstOiADC(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doADC(data)
	cpu.next = InstOpINI
}

// InstOpSBC executes the SBC (Subtract with Carry) operation using data from the address register and updates the next instruction.
//
//go:nosplit
func InstOpSBC(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.doSBC(data)
	cpu.next = InstOpINI
}

// InstOiSBC handles the SBC (Subtract with Carry) instruction. It reads data at the program counter, increments the PC,
// performs the SBC operation using the read data, and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOiSBC(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = InstOpINI
}

// Increment, decrement

// InstOpINX increments the X register, updates the negative and zero flags, and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpINX(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = InstOpINI
}

// InstOpDEX decrements the X register, setting the negative and zero flags, and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpDEX(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = InstOpINI
}

// InstOpINY increments the Y register, updates the negative and zero flags, and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpINY(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = InstOpINI
}

// InstOpDEY decrements the Y register, updating the negative and zero flags based on the resulting value.
//
//go:nosplit
func InstOpDEY(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = InstOpINI
}

// InstOpINC increments a value, updates the N and Z flags, writes the result to memory, and sets the next operation.
//
//go:nosplit
func InstOpINC(cpu *CPU) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.bus.Write(cpu.ar, v)
	cpu.next = InstOpINI
}

// InstOpDEC performs a decrement operation on the value in the RMW register, updating CPU flags and writing the result.
//
//go:nosplit
func InstOpDEC(cpu *CPU) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.bus.Write(cpu.ar, v)
	cpu.next = InstOpINI
}

// InstOpAND performs a bitwise AND operation between the accumulator and memory at the address in the address register.
// Updates the negative and zero flags based on the result. Sets the next instruction handler.
//
//go:nosplit
func InstOpAND(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOiAND performs the AND operation between the accumulator and fetched data. Updates N and Z flags accordingly.
//
//go:nosplit
func InstOiAND(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpORA performs the ORA (logical OR with accumulator) operation, updates the negative/zero flags, and sets the next instruction.
//
//go:nosplit
func InstOpORA(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOiOPA is an instruction that performs a bitwise OR operation between the accumulator and a fetched operand.
// Updates negative and zero flags based on the result. Sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOiOPA(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpEOR performs the EOR (Exclusive OR) operation on the accumulator with a value from memory and updates CPU flags.
//
//go:nosplit
func InstOpEOR(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOiEOR executes the Exclusive OR (EOR) operation on the accumulator with a fetched memory value and updates flags.
//
//go:nosplit
func InstOiEOR(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpCMP performs a comparison between the accumulator and a memory value, updating CPU flags based on the result.
//
//go:nosplit
func InstOpCMP(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	// Set Carry flag if result is >= 0 (no borrow occurred).
	// The result of (A - M) is stored in cpu.ar.
	// If cpu.ar < 0x100, it means the high byte is 0, so no borrow.
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// InstOiCMP performs a comparison between the CPU's accumulator and the operand, updating CPU flags accordingly.
//
//go:nosplit
func InstOiCMP(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// InstOpCPX executes the CPX (Compare with X Register) instruction, updating the CPU's flags based on the result.
//
//go:nosplit
func InstOpCPX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// InstOiCPX executes the OiCPX instruction by reading a value from memory and performing a subtraction with register X.
// It updates the CPU's flags (negative, zero, carry) and the address register (AR) as per the operation's result.
// Finally, it sets the next instruction handler to InstOpINI.
//
//go:nosplit
func InstOiCPX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// InstOpCPY handles the CPY (Compare Y Register) instruction, updating flags based on comparison with memory data.
//
//go:nosplit
func InstOpCPY(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// InstOiCPY handles the CPY (Compare Y Register with Memory) instruction in immediate addressing mode.
// It updates the negative, zero, and carry flags based on the comparison result.
// The next instruction to execute is set to InstOpINI.
//
//go:nosplit
func InstOiCPY(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = InstOpINI
}

// Bit-test

// InstOpBIT performs the BIT instruction, updating the CPU flags based on a bitwise AND between the accumulator and memory data.
//
//go:nosplit
func InstOpBIT(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = InstOpINI
}

// InstOpASL performs the ASL (Arithmetic Shift Left) operation on the CPU's RMW buffer and updates relevant CPU flags.
// It shifts the RMW buffer left by one bit, sets the carry flag to the original top bit, and updates zero/negative flags.
// The result is written back to memory at the address pointed to by the address register (ar).
// Finally, sets the next CPU instruction to InstOpINI for subsequent execution.
//
//go:nosplit
func InstOpASL(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.bus.Write(cpu.ar, v)
	cpu.next = InstOpINI
}

// InstOaASL executes the ASL (Arithmetic Shift Left) operation on the accumulator, updating flags and the next instruction.
//
//go:nosplit
func InstOaASL(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpLSR performs the Logical Shift Right (LSR) operation on the CPU's RMW register, updating flags and writing the result.
//
//go:nosplit
func InstOpLSR(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.bus.Write(cpu.ar, v)
	cpu.next = InstOpINI
}

// InstOaLSR performs a logical shift right (LSR) on the A register, updating the carry, zero, and negative flags accordingly.
//
//go:nosplit
func InstOaLSR(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpROL performs the ROL (Rotate Left) operation on the `rmw` register, updating CPU flags and memory state.
//
//go:nosplit
func InstOpROL(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw << 1) | 0x1
	} else {
		t = cpu.rmw << 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.bus.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x80
	cpu.next = InstOpINI
}

// InstOaROL performs a rotate left operation on the accumulator and updates the CPU flags accordingly.
//
//go:nosplit
func InstOaROL(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x80
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a << 1) | 0x1
	} else {
		cpu.a = cpu.a << 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = InstOpINI
}

// InstOpROR performs a Rotate Right operation on the CPU's RMW register, updates flags, and writes the result to memory.
//
//go:nosplit
func InstOpROR(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw >> 1) | 0x80
	} else {
		t = cpu.rmw >> 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.bus.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x1
	cpu.next = InstOpINI
}

// InstOaROR performs the ROR (Rotate Right) operation on the accumulator,
// updating the negative, zero, and carry flags and setting the next instruction.
//
//go:nosplit
func InstOaROR(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x1
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a >> 1) | 0x80
	} else {
		cpu.a = cpu.a >> 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = InstOpINI
}
