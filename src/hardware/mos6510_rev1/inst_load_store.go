package mos6510_rev1

// instApZER loads a zero-page address into the address register and sets the next instruction handler.
func instApZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = _opTable[cpu.op]
}

// instApZERx performs a zero-page addressing operation by reading a byte at the program counter and updating the address register.
// It increments the program counter and sets the next instruction to instApZERx1.
func instApZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERx1
}

// instApZERx1 performs a zero-page indexed addressing mode operation using the CPU's address register and X register.
// It updates the address register by adding the X register value, ensuring the result wraps around at 8-bit boundaries.
// The next instruction handler is set based on the current opcode. No operation occurs if the read fails.
func instApZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

// instApZERy loads a byte from memory at the program counter into the address register and sets the next instruction.
func instApZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERy1
}

// instApZERy1 adjusts the address register by adding the Y register value and wraps it within a byte boundary.
// Proceeds to the next operation if the address read was successful.
func instApZERy1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

// instApABS loads the next byte from memory into the address register and advances the program counter.
func instApABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABS1
}

// instApABS1 reads a byte from memory at the program counter, increments the PC, and updates the address register (AR).
// Then, it fetches the next instruction from the operation table based on the current opcode.
func instApABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

// instApABSx fetches a byte from the program counter, increments the PC, and stores the value in the address register.
// Sets the next instruction handler to instApABSx1.
func instApABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSx1
}

// instApABSx1 executes the first step of an absolute addressing mode with X offset, updating CPU registers and next instruction.
func instApABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSx2
	} else {
		cpu.next = instApABSx3
	}
}

// instApABSx2 retrieves data from the address specified by the address register and sets the next operation if successful.
func instApABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApABSx3 performs absolute addressing with an additional stack address adjustment and updates the next instruction.
// If the page is crossed, the function ensures proper handling by checking the address read operation for success.
func instApABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instApABSy loads a byte from memory at the current program counter into the address register and advances the program counter.
func instApABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSy1
}

// instApABSy1 reads a byte from memory, updates address registers, and determines the next instruction to execute.
func instApABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSy2
	} else {
		cpu.next = instApABSy3
	}
}

// instApABSy2 handles the execution flow for a specific CPU instruction without crossing a memory page.
// If the memory read operation fails, it terminates further execution for this step.
// Updates the CPU's next instruction handler based on the opcode.
func instApABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApABSy3 handles the execution of an operation, performs a page cross check, and sets the next instruction to execute.
func instApABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instApINDx performs the indirect indexed addressing mode operation, updating CPU state and setting the next instruction.
func instApINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDx1
}

// instApINDx1 executes the first stage of the Indexed Indirect (IND,X) addressing mode operation.
// It reads the memory at the address in cpu.ar2 and checks its availability. If unavailable, the CPU halts.
// Updates cpu.ar2 by adding the X register value, masking the result to fit within 8 bits.
// Sets the next instruction handler to instApINDx2.
func instApINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instApINDx2
}

// instApINDx2 reads a value from the address in ar2, sets it to ar if successful, and updates the next handler to instApINDx3.
func instApINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDx3
}

// instApINDx3 performs an indirect indexed addressing operation, updating the address register and setting the next instruction.
func instApINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

// instApINDy fetches a byte from memory at the program counter, increments the PC, and updates the AR2 register.
// It sets the next CPU instruction to instApINDy1.
// If memory reading fails, the function exits early.
func instApINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDy1
}

// instApINDy1 reads a value from the address in `ar2`, sets `ar` with the value, and transitions the CPU to `instApINDy2`.
func instApINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDy2
}

// instApINDy2 performs indirect indexed addressing by updating `ar` using `ar2` and `y`, and sets the next instruction.
// The function reads a byte from memory, updates `ar2`, adjusts `ar`, and determines the next handler based on conditions.
func instApINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApINDy3
	} else {
		cpu.next = instApINDy4
	}
}

// instApINDy3 handles indirect indexed addressing with Y-register offset without page crossing for the CPU.
// It reads from the address stored in the AR register, updates the instruction handler, and checks RDY state.
func instApINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApINDy4 performs an indirect indexed addressing mode operation with Y register and updates the CPU state.
// If a page boundary is crossed during execution, it ensures proper handling and advances to the next instruction.
func instApINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeABSx fetches an 8-bit immediate value, increments the program counter, and stores it in the address register (ar).
func instAeABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSx1
}

// instAeABSx1 handles the first step of the absolute-indexed addressing mode with the X register adjustment.
// It combines the X register with the address register and determines the next instruction to execute.
// Updates the address register (ar) and increments the program counter (pc).
func instAeABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSx2
	}
}

// instAeABSx2 handles the Absolute Indexed X addressing mode with page crossing, updating the address and next operation.
func instAeABSx2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeABSy initializes the address register with the value read from the program counter and sets the next instruction handler.
func instAeABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSy1
}

// instAeABSy1 fetches data, updates the address register, and determines the next instruction based on the address range.
func instAeABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSy2
	}
}

// instAeABSy2 adjusts the address register by adding a stack offset and sets the next instruction from the operation table.
// It ensures instruction execution respects memory boundary crossing conditions.
func instAeABSy2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeINDy executes the AE INDY instruction, updating the program counter and secondary address register, then sets the next operation.
func instAeINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instAeINDy1
}

// instAeINDy1 reads data from the memory address specified by `ar2` and stores it in `ar`. Sets the next instruction handler.
func instAeINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instAeINDy2
}

// instAeINDy2 executes the second phase of the AE indirect Y-indexed addressing mode.
// It reads a value from memory, combines it with the address register and Y register, and sets the next instruction.
func instAeINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeINDy3
	}
}

// instAeINDy3 handles the addressing mode operation for instructions that cross a page boundary and updates the CPU state.
func instAeINDy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instMpZER is a memory operation instruction that reads data from memory into the address register (ar) and sets the next operation.
func instMpZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpRMW
}

// instMpZERx executes a zero-page read operation, increments the program counter, and sets the next instruction.
func instMpZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpZERx1
}

// instMpZERx1 performs a memory read operation at the address in `cpu.ar` and modifies the address by adding `cpu.x`.
// If the memory read fails, the function returns immediately.
// The adjusted address is constrained to an 8-bit boundary, and control proceeds to the `instOpRMW` handler.
func instMpZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = instOpRMW
}

//func instMpZERy(cpu *CPU) {
//	data, ok := cpu.read(cpu.pc)
//	if !ok {
//		return
//	}
//	cpu.pc++
//	cpu.ar = uint16(data)
//	cpu.next = instMpZERy1
//}

//func instMpZERy1(cpu *CPU) {
//	data, ok := cpu.read(cpu.ar)
//	if !ok {
//		return
//	}
//	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
//	cpu.next = instOpRMW
//}

// instMpABS loads a byte from memory at the program counter into the address register and updates the next instruction.
func instMpABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABS1
}

// instMpABS1 represents an instruction handler that modifies the address register (AR) with data fetched from memory.
// It reads a byte from memory located at the program counter (PC), shifts it 8 bits left, and ORs it with the current AR.
// The program counter is incremented, and the next instruction handler is set to a read-modify-write (RMW) operation.
func instMpABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

// instMpABSx fetches the next byte from memory, increments the program counter, and sets it in the address register.
// It updates the CPU's next state to instMpABSx1.
func instMpABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSx1
}

// instMpABSx1 is a CPU instruction handler for addressing mode manipulation.
// It reads the next byte from memory and uses it as part of the absolute address computation.
// Updates the address register (ar) with the computed address and determines the next instruction handler.
func instMpABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSx2
	} else {
		cpu.next = instMpABSx3
	}
}

// instMpABSx2 performs an absolute memory read operation and checks for page crossing issues.
// It sets the next instruction to a read-modify-write operation if the memory read is successful.
func instMpABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpABSx3 performs an operation using absolute addressing with additional adjustments and transitions to the next instruction.
func instMpABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instMpABSy sets the address register (ar) based on the byte at the program counter, then sets the next instruction handler.
func instMpABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSy1
}

// instMpABSy1 performs memory page addressing mode with Y register offset, updating `ar` and determining the next instruction.
func instMpABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSy2
	} else {
		cpu.next = instMpABSy3
	}
}

// instMpABSy2 handles the zero-page no-cross memory access and assigns the next instruction to a read-modify-write operation.
// It reads from the address register; if unsuccessful, it stops. Otherwise, it sets the next handler to instOpRMW.
func instMpABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpABSy3 adjusts the address register and sets up the next operation, handling page crossing scenarios.
func instMpABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instMpINDx loads an operand from the program counter into the ar2 register and advances the program counter.
// Sets the next CPU instruction handler to instMpINDx1. Returns immediately if the read operation fails.
func instMpINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDx1
}

// instMpINDx1 performs indexed indirect addressing mode update. Adjusts ar2 with x, wraps it, and sets the next operation.
func instMpINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instMpINDx2
}

// instMpINDx2 reads data from the memory address in ar2, updates ar with this value if successful, and sets the next instruction.
func instMpINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDx3
}

// instMpINDx3 performs an indexed memory read, updates the address register, and sets the next instruction handler.
func instMpINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

// instMpINDy reads a byte from memory at the program counter, increments the program counter, and updates ar2 and next state.
func instMpINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDy1
}

// instMpINDy1 reads data using the CPU's ar2 register, updates the ar register, and sets the next operation to instMpINDy2.
func instMpINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDy2
}

// instMpINDy2 handles the indirect indexed addressing mode logic by adjusting the address register and setting the next instruction.
func instMpINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpINDy3
	} else {
		cpu.next = instMpINDy4
	}
}

// instMpINDy3 handles the memory instruction with indirect addressing, updating the CPU's next operation if successful.
func instMpINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpINDy4 handles indirect indexed addressing with Y register offset and handles potential page crossing.
// If a page boundary is crossed, the function adjusts the address register and sets the next operation to instOpRMW.
func instMpINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instOpRMW reads data from the address in the CPU's address register, stores it in the `rmw` buffer, and sets the next operation.
func instOpRMW(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.rmw = data
	cpu.next = instOpRMW1
}

// instOpRMW1 executes the second phase of a read-modify-write operation by writing the modified value and updating the next instruction.
func instOpRMW1(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// instOpLDA loads a byte from memory at the address in the address register (AR) into the accumulator (A).
// Updates the negative (N) and zero (Z) flags based on the loaded value.
// Sets the next instruction to `instOpINI` if reading from memory is successful; does nothing otherwise.
func instOpLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDA loads a byte from memory into the accumulator, updates the negative and zero flags, and sets the next instruction.
func instOiLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOpLDX loads a value from memory into the X register and updates the negative and zero flags.
func instOpLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDX loads a byte from memory into the X register, updating the negative and zero flags, and sets the next instruction.
func instOiLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOpLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDY loads a value into the Y register, updates the negative and zero flags, and sets the next instruction handler.
func instOiLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// Store

// instOpSTA stores the value in the accumulator (A register) into memory at the address stored in the address register (AR).
// It updates the next CPU instruction handler to instOpINI.
func instOpSTA(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = instOpINI
}

// instOpSTX stores the value of the X register into memory at the address specified by the AR register and sets the next instruction.
func instOpSTX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = instOpINI
}

// instOpSTY stores the value of the Y register into memory at the address in the address register and sets the next instruction.
func instOpSTY(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = instOpINI
}
