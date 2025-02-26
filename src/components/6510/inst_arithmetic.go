package mos6510

import "github.com/markel1974/c64emu/src/conversion"

// Arithmetic

// instOpADC performs the ADC (Add with Carry) operation by reading data from memory and calling the ADC handler.
func instOpADC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doADC(data)
	cpu.next = instOpINI
}

// instOiADC executes the ADC instruction at the current program counter, updating the accumulator and advancing the PC.
// Invokes the next instruction handler after performing the operation.
func instOiADC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doADC(data)
	cpu.next = instOpINI
}

// instOpSBC executes the SBC (Subtract with Carry) operation using data from the address register and updates the next instruction.
func instOpSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doSBC(data)
	cpu.next = instOpINI
}

// instOiSBC handles the SBC (Subtract with Carry) instruction. It reads data at the program counter, increments the PC,
// performs the SBC operation using the read data, and sets the next instruction to instOpINI.
func instOiSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = instOpINI
}

// Increment, decrement

// instOpINX increments the X register, updates the negative and zero flags, and sets the next instruction to instOpINI.
func instOpINX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpDEX decrements the X register, setting the negative and zero flags, and sets the next instruction to instOpINI.
func instOpDEX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpINY increments the Y register, updates the negative and zero flags, and sets the next instruction to instOpINI.
func instOpINY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpDEY decrements the Y register, updating the negative and zero flags based on the resulting value.
func instOpDEY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpINC increments a value, updates the N and Z flags, writes the result to memory, and sets the next operation.
func instOpINC(cpu *CPU) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOpDEC performs a decrement operation on the value in the RMW register, updating CPU flags and writing the result.
func instOpDEC(cpu *CPU) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOpAND performs a bitwise AND operation between the accumulator and memory at the address in the address register.
// Updates the negative and zero flags based on the result. Sets the next instruction handler.
func instOpAND(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiAND performs the AND operation between the accumulator and fetched data. Updates N and Z flags accordingly.
func instOiAND(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpORA performs the ORA (logical OR with accumulator) operation, updates the negative/zero flags, and sets the next instruction.
func instOpORA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiOPA is an instruction that performs a bitwise OR operation between the accumulator and a fetched operand.
// Updates negative and zero flags based on the result. Sets the next instruction to instOpINI.
func instOiOPA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpEOR performs the EOR (Exclusive OR) operation on the accumulator with a value from memory and updates CPU flags.
func instOpEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiEOR executes the Exclusive OR (EOR) operation on the accumulator with a fetched memory value and updates flags.
func instOiEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpCMP performs a comparison between the accumulator and a memory value, updating CPU flags based on the result.
func instOpCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCMP performs a comparison between the CPU's accumulator and the operand, updating CPU flags accordingly.
func instOiCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpCPX executes the CPX (Compare with X Register) instruction, updating the CPU's flags based on the result.
func instOpCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCPX executes the OiCPX instruction by reading a value from memory and performing a subtraction with register X.
// It updates the CPU's flags (negative, zero, carry) and the address register (AR) as per the operation's result.
// Finally, it sets the next instruction handler to instOpINI.
func instOiCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpCPY handles the CPY (Compare Y Register) instruction, updating flags based on comparison with memory data.
func instOpCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCPY handles the CPY (Compare Y Register with Memory) instruction in immediate addressing mode.
// It updates the negative, zero, and carry flags based on the comparison result.
// The next instruction to execute is set to instOpINI.
func instOiCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// Bit-test

// instOpBIT performs the BIT instruction, updating the CPU flags based on a bitwise AND between the accumulator and memory data.
func instOpBIT(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = instOpINI
}

// instOpASL performs the ASL (Arithmetic Shift Left) operation on the CPU's RMW buffer and updates relevant CPU flags.
// It shifts the RMW buffer left by one bit, sets the carry flag to the original top bit, and updates zero/negative flags.
// The result is written back to memory at the address pointed to by the address register (ar).
// Finally, sets the next CPU instruction to instOpINI for subsequent execution.
func instOpASL(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOaASL executes the ASL (Arithmetic Shift Left) operation on the accumulator, updating flags and the next instruction.
func instOaASL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpLSR performs the Logical Shift Right (LSR) operation on the CPU's RMW register, updating flags and writing the result.
func instOpLSR(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOaLSR performs a logical shift right (LSR) on the A register, updating the carry, zero, and negative flags accordingly.
func instOaLSR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpROL performs the ROL (Rotate Left) operation on the `rmw` register, updating CPU flags and memory state.
func instOpROL(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw << 1) | 0x1
	} else {
		t = cpu.rmw << 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x80
	cpu.next = instOpINI
}

// instOaROL performs a rotate left operation on the accumulator and updates the CPU flags accordingly.
func instOaROL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
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
	cpu.next = instOpINI
}

// instOpROR performs a Rotate Right operation on the CPU's RMW register, updates flags, and writes the result to memory.
func instOpROR(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw >> 1) | 0x80
	} else {
		t = cpu.rmw >> 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x1
	cpu.next = instOpINI
}

// instOaROR performs the ROR (Rotate Right) operation on the accumulator,
// updating the negative, zero, and carry flags and setting the next instruction.
func instOaROR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
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
	cpu.next = instOpINI
}
