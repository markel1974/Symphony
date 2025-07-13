package mos6510_rev1

// InstApZER loads a zero-page address into the address register and sets the next instruction handler.
//
//go:nosplit
func InstApZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = cpu.opTable[cpu.op]
}

// InstApZERx performs a zero-page addressing operation by reading a byte at the program counter and updating the address register.
// It increments the program counter and sets the next instruction to InstApZERx1.
//
//go:nosplit
func InstApZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstApZERx1
}

// InstApZERx1 performs a zero-page indexed addressing mode operation using the CPU's address register and X register.
// It updates the address register by adding the X register value, ensuring the result wraps around at 8-bit boundaries.
// The next instruction handler is set based on the current opcode. No operation occurs if the read fails.
//
//go:nosplit
func InstApZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = cpu.opTable[cpu.op]
}

// InstApZERy loads a byte from memory at the program counter into the address register and sets the next instruction.
//
//go:nosplit
func InstApZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstApZERy1
}

// InstApZERy1 adjusts the address register by adding the Y register value and wraps it within a byte boundary.
// Proceeds to the next operation if the address read was successful.
//
//go:nosplit
func InstApZERy1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABS loads the next byte from memory into the address register and advances the program counter.
//
//go:nosplit
func InstApABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstApABS1
}

// InstApABS1 reads a byte from memory at the program counter, increments the PC, and updates the address register (AR).
// Then, it fetches the next instruction from the operation table based on the current opcode.
//
//go:nosplit
func InstApABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSx fetches a byte from the program counter, increments the PC, and stores the value in the address register.
// Sets the next instruction handler to InstApABSx1.
//
//go:nosplit
func InstApABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstApABSx1
}

// InstApABSx1 executes the first step of an absolute addressing mode with X offset, updating CPU registers and next instruction.
//
//go:nosplit
func InstApABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = InstApABSx2
	} else {
		cpu.next = InstApABSx3
	}
}

// InstApABSx2 retrieves data from the address specified by the address register and sets the next operation if successful.
//
//go:nosplit
func InstApABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSx3 performs absolute addressing with an additional stack address adjustment and updates the next instruction.
// If the page is crossed, the function ensures proper handling by checking the address read operation for success.
//
//go:nosplit
func InstApABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSy loads a byte from memory at the current program counter into the address register and advances the program counter.
//
//go:nosplit
func InstApABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstApABSy1
}

// InstApABSy1 reads a byte from memory, updates address registers, and determines the next instruction to execute.
//
//go:nosplit
func InstApABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = InstApABSy2
	} else {
		cpu.next = InstApABSy3
	}
}

// InstApABSy2 handles the execution flow for a specific CPU instruction without crossing a memory page.
// If the memory read operation fails, it terminates further execution for this step.
// Updates the CPU's next instruction handler based on the opcode.
//
//go:nosplit
func InstApABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSy3 handles the execution of an operation, performs a page cross check, and sets the next instruction to execute.
//
//go:nosplit
func InstApABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstApINDx performs the indirect indexed addressing mode operation, updating CPU state and setting the next instruction.
//
//go:nosplit
func InstApINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = InstApINDx1
}

// InstApINDx1 executes the first stage of the Indexed Indirect (IND,X) addressing mode operation.
// It reads the memory at the address in cpu.ar2 and checks its availability. If unavailable, the CPU halts.
// Updates cpu.ar2 by adding the X register value, masking the result to fit within 8 bits.
// Sets the next instruction handler to InstApINDx2.
//
//go:nosplit
func InstApINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = InstApINDx2
}

// InstApINDx2 reads a value from the address in ar2, sets it to ar if successful, and updates the next handler to InstApINDx3.
//
//go:nosplit
func InstApINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = InstApINDx3
}

// InstApINDx3 performs an indirect indexed addressing operation, updating the address register and setting the next instruction.
//
//go:nosplit
func InstApINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = cpu.opTable[cpu.op]
}

// InstApINDy fetches a byte from memory at the program counter, increments the PC, and updates the AR2 register.
// It sets the next CPU instruction to InstApINDy1.
// If memory reading fails, the function exits early.
//
//go:nosplit
func InstApINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = InstApINDy1
}

// InstApINDy1 reads a value from the address in `ar2`, sets `ar` with the value, and transitions the CPU to `InstApINDy2`.
//
//go:nosplit
func InstApINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = InstApINDy2
}

// InstApINDy2 performs indirect indexed addressing by updating `ar` using `ar2` and `y`, and sets the next instruction.
// The function reads a byte from memory, updates `ar2`, adjusts `ar`, and determines the next handler based on conditions.
//
//go:nosplit
func InstApINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = InstApINDy3
	} else {
		cpu.next = InstApINDy4
	}
}

// InstApINDy3 handles indirect indexed addressing with Y-register offset without page crossing for the CPU.
// It reads from the address stored in the AR register, updates the instruction handler, and checks RDY state.
//
//go:nosplit
func InstApINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApINDy4 performs an indirect indexed addressing mode operation with Y register and updates the CPU state.
// If a page boundary is crossed during execution, it ensures proper handling and advances to the next instruction.
//
//go:nosplit
func InstApINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeABSx fetches an 8-bit immediate value, increments the program counter, and stores it in the address register (ar).
//
//go:nosplit
func InstAeABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstAeABSx1
}

// InstAeABSx1 handles the first step of the absolute-indexed addressing mode with the X register adjustment.
// It combines the X register with the address register and determines the next instruction to execute.
// Updates the address register (ar) and increments the program counter (pc).
//
//go:nosplit
func InstAeABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = InstAeABSx2
	}
}

// InstAeABSx2 handles the Absolute Indexed X addressing mode with page crossing, updating the address and next operation.
//
//go:nosplit
func InstAeABSx2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeABSy initializes the address register with the value read from the program counter and sets the next instruction handler.
//
//go:nosplit
func InstAeABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstAeABSy1
}

// InstAeABSy1 fetches data, updates the address register, and determines the next instruction based on the address range.
//
//go:nosplit
func InstAeABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = InstAeABSy2
	}
}

// InstAeABSy2 adjusts the address register by adding a stack offset and sets the next instruction from the operation table.
// It ensures instruction execution respects memory boundary crossing conditions.
//
//go:nosplit
func InstAeABSy2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeINDy executes the AE INDY instruction, updating the program counter and secondary address register, then sets the next operation.
//
//go:nosplit
func InstAeINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = InstAeINDy1
}

// InstAeINDy1 reads data from the memory address specified by `ar2` and stores it in `ar`. Sets the next instruction handler.
//
//go:nosplit
func InstAeINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = InstAeINDy2
}

// InstAeINDy2 executes the second phase of the AE indirect Y-indexed addressing mode.
// It reads a value from memory, combines it with the address register and Y register, and sets the next instruction.
//
//go:nosplit
func InstAeINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = InstAeINDy3
	}
}

// InstAeINDy3 handles the addressing mode operation for instructions that cross a page boundary and updates the CPU state.
//
//go:nosplit
func InstAeINDy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstMpZER is a memory operation instruction that reads data from memory into the address register (ar) and sets the next operation.
//
//go:nosplit
func InstMpZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstOpRMW
}

// InstMpZERx executes a zero-page read operation, increments the program counter, and sets the next instruction.
//
//go:nosplit
func InstMpZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpZERx1
}

// InstMpZERx1 performs a memory read operation at the address in `cpu.ar` and modifies the address by adding `cpu.x`.
// If the memory read fails, the function returns immediately.
// The adjusted address is constrained to an 8-bit boundary, and control proceeds to the `InstOpRMW` handler.
//
//go:nosplit
func InstMpZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = InstOpRMW
}

// InstMpABS loads a byte from memory at the program counter into the address register and updates the next instruction.
//
//go:nosplit
func InstMpABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpABS1
}

// InstMpABS1 represents an instruction handler that modifies the address register (AR) with data fetched from memory.
// It reads a byte from memory located at the program counter (PC), shifts it 8 bits left, and ORs it with the current AR.
// The program counter is incremented, and the next instruction handler is set to a read-modify-write (RMW) operation.
//
//go:nosplit
func InstMpABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = InstOpRMW
}

// InstMpABSx fetches the next byte from memory, increments the program counter, and sets it in the address register.
// It updates the CPU's next state to InstMpABSx1.
//
//go:nosplit
func InstMpABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpABSx1
}

// InstMpABSx1 is a CPU instruction handler for addressing mode manipulation.
// It reads the next byte from memory and uses it as part of the absolute address computation.
// Updates the address register (ar) with the computed address and determines the next instruction handler.
//
//go:nosplit
func InstMpABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = InstMpABSx2
	} else {
		cpu.next = InstMpABSx3
	}
}

// InstMpABSx2 performs an absolute memory read operation and checks for page crossing issues.
// It sets the next instruction to a read-modify-write operation if the memory read is successful.
//
//go:nosplit
func InstMpABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = InstOpRMW
}

// InstMpABSx3 performs an operation using absolute addressing with additional adjustments and transitions to the next instruction.
//
//go:nosplit
func InstMpABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = InstOpRMW
}

// InstMpABSy sets the address register (ar) based on the byte at the program counter, then sets the next instruction handler.
//
//go:nosplit
func InstMpABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpABSy1
}

// InstMpABSy1 performs memory page addressing mode with Y register offset, updating `ar` and determining the next instruction.
//
//go:nosplit
func InstMpABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = InstMpABSy2
	} else {
		cpu.next = InstMpABSy3
	}
}

// InstMpABSy2 handles the zero-page no-cross memory access and assigns the next instruction to a read-modify-write operation.
// It reads from the address register; if unsuccessful, it stops. Otherwise, it sets the next handler to InstOpRMW.
//
//go:nosplit
func InstMpABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = InstOpRMW
}

// InstMpABSy3 adjusts the address register and sets up the next operation, handling page crossing scenarios.
//
//go:nosplit
func InstMpABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = InstOpRMW
}

// InstMpINDx loads an operand from the program counter into the ar2 register and advances the program counter.
// Sets the next CPU instruction handler to InstMpINDx1. Returns immediately if the read operation fails.
//
//go:nosplit
func InstMpINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = InstMpINDx1
}

// InstMpINDx1 performs indexed indirect addressing mode update. Adjusts ar2 with x, wraps it, and sets the next operation.
//
//go:nosplit
func InstMpINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = InstMpINDx2
}

// InstMpINDx2 reads data from the memory address in ar2, updates ar with this value if successful, and sets the next instruction.
//
//go:nosplit
func InstMpINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = InstMpINDx3
}

// InstMpINDx3 performs an indexed memory read, updates the address register, and sets the next instruction handler.
//
//go:nosplit
func InstMpINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = InstOpRMW
}

// InstMpINDy reads a byte from memory at the program counter, increments the program counter, and updates ar2 and next state.
//
//go:nosplit
func InstMpINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = InstMpINDy1
}

// InstMpINDy1 reads data using the CPU's ar2 register, updates the ar register, and sets the next operation to InstMpINDy2.
//
//go:nosplit
func InstMpINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = InstMpINDy2
}

// InstMpINDy2 handles the indirect indexed addressing mode logic by adjusting the address register and setting the next instruction.
//
//go:nosplit
func InstMpINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = InstMpINDy3
	} else {
		cpu.next = InstMpINDy4
	}
}

// InstMpINDy3 handles the memory instruction with indirect addressing, updating the CPU's next operation if successful.
//
//go:nosplit
func InstMpINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = InstOpRMW
}

// InstMpINDy4 handles indirect indexed addressing with Y register offset and handles potential page crossing.
// If a page boundary is crossed, the function adjusts the address register and sets the next operation to InstOpRMW.
//
//go:nosplit
func InstMpINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = InstOpRMW
}

// InstOpRMW reads data from the address in the CPU's address register, stores it in the `rmw` buffer, and sets the next operation.
//
//go:nosplit
func InstOpRMW(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.rmw = data
	cpu.next = InstOpRMW1
}

// InstOpRMW1 executes the second phase of a read-modify-write operation by writing the modified value and updating the next instruction.
//
//go:nosplit
func InstOpRMW1(cpu *CPU) {
	cpu.bankWrite(cpu.ar, cpu.rmw)
	cpu.next = cpu.opTable[cpu.op]
}

// InstOpLDA loads a byte from memory at the address in the address register (AR) into the accumulator (A).
// Updates the negative (N) and zero (Z) flags based on the loaded value.
// Sets the next instruction to `InstOpINI` if reading from memory is successful; does nothing otherwise.
//
//go:nosplit
func InstOpLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

// InstOiLDA loads a byte from memory into the accumulator, updates the negative and zero flags, and sets the next instruction.
//
//go:nosplit
func InstOiLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

// InstOpLDX loads a value from memory into the X register and updates the negative and zero flags.
//
//go:nosplit
func InstOpLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

// InstOiLDX loads a byte from memory into the X register, updating the negative and zero flags, and sets the next instruction.
//
//go:nosplit
func InstOiLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

//go:nosplit
func InstOpLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

// InstOiLDY loads a value into the Y register, updates the negative and zero flags, and sets the next instruction handler.
//
//go:nosplit
func InstOiLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = InstOpINI
}

// Store

// InstOpSTA stores the value in the accumulator (A register) into memory at the address stored in the address register (AR).
// It updates the next CPU instruction handler to InstOpINI.
//
//go:nosplit
func InstOpSTA(cpu *CPU) {
	cpu.bankWrite(cpu.ar, cpu.a)
	cpu.next = InstOpINI
}

// InstOpSTX stores the value of the X register into memory at the address specified by the AR register and sets the next instruction.
//
//go:nosplit
func InstOpSTX(cpu *CPU) {
	cpu.bankWrite(cpu.ar, cpu.x)
	cpu.next = InstOpINI
}

// InstOpSTY stores the value of the Y register into memory at the address in the address register and sets the next instruction.
//
//go:nosplit
func InstOpSTY(cpu *CPU) {
	cpu.bankWrite(cpu.ar, cpu.y)
	cpu.next = InstOpINI
}

/*
//go:nosplit
func InstMpZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpZERy1
}

//go:nosplit
func InstMpZERy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = InstOpRMW
}
*/
