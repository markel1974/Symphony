package mos6510_rev1

// InstApZER loads a zero-page address into the address register and sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstApZER(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
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
func (er *Executor) InstApZERx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstApZERx1
}

// InstApZERx1 performs a zero-page indexed addressing mode operation using the CPU's address register and X register.
// It updates the address register by adding the X register value, ensuring the result wraps around at 8-bit boundaries.
// The next instruction handler is set based on the current opcode. No operation occurs if the read fails.
//
//go:nosplit
func (er *Executor) InstApZERx1(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = cpu.opTable[cpu.op]
}

// InstApZERy loads a byte from memory at the program counter into the address register and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstApZERy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstApZERy1
}

// InstApZERy1 adjusts the address register by adding the Y register value and wraps it within a byte boundary.
// Proceeds to the next operation if the address read was successful.
//
//go:nosplit
func (er *Executor) InstApZERy1(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABS loads the next byte from memory into the address register and advances the program counter.
//
//go:nosplit
func (er *Executor) InstApABS(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstApABS1
}

// InstApABS1 reads a byte from memory at the program counter, increments the PC, and updates the address register (AR).
// Then, it fetches the next instruction from the operation table based on the current opcode.
//
//go:nosplit
func (er *Executor) InstApABS1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
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
func (er *Executor) InstApABSx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstApABSx1
}

// InstApABSx1 executes the first step of an absolute addressing mode with X offset, updating CPU registers and next instruction.
//
//go:nosplit
func (er *Executor) InstApABSx1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = er.InstApABSx2
	} else {
		cpu.next = er.InstApABSx3
	}
}

// InstApABSx2 retrieves data from the address specified by the address register and sets the next operation if successful.
//
//go:nosplit
func (er *Executor) InstApABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSx3 performs absolute addressing with an additional stack address adjustment and updates the next instruction.
// If the page is crossed, the function ensures proper handling by checking the address read operation for success.
//
//go:nosplit
func (er *Executor) InstApABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSy loads a byte from memory at the current program counter into the address register and advances the program counter.
//
//go:nosplit
func (er *Executor) InstApABSy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstApABSy1
}

// InstApABSy1 reads a byte from memory, updates address registers, and determines the next instruction to execute.
//
//go:nosplit
func (er *Executor) InstApABSy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = er.InstApABSy2
	} else {
		cpu.next = er.InstApABSy3
	}
}

// InstApABSy2 handles the execution flow for a specific CPU instruction without crossing a memory page.
// If the memory read operation fails, it terminates further execution for this step.
// Updates the CPU's next instruction handler based on the opcode.
//
//go:nosplit
func (er *Executor) InstApABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApABSy3 handles the execution of an operation, performs a page cross check, and sets the next instruction to execute.
//
//go:nosplit
func (er *Executor) InstApABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstApINDx performs the indirect indexed addressing mode operation, updating CPU state and setting the next instruction.
//
//go:nosplit
func (er *Executor) InstApINDx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = er.InstApINDx1
}

// InstApINDx1 executes the first stage of the Indexed Indirect (IND,X) addressing mode operation.
// It reads the memory at the address in cpu.ar2 and checks its availability. If unavailable, the CPU halts.
// Updates cpu.ar2 by adding the X register value, masking the result to fit within 8 bits.
// Sets the next instruction handler to InstApINDx2.
//
//go:nosplit
func (er *Executor) InstApINDx1(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = er.InstApINDx2
}

// InstApINDx2 reads a value from the address in ar2, sets it to ar if successful, and updates the next handler to InstApINDx3.
//
//go:nosplit
func (er *Executor) InstApINDx2(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = er.InstApINDx3
}

// InstApINDx3 performs an indirect indexed addressing operation, updating the address register and setting the next instruction.
//
//go:nosplit
func (er *Executor) InstApINDx3(cpu *CPU) {
	data, ok := cpu.bus.Read((cpu.ar2 + 1) & 0xff)
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
func (er *Executor) InstApINDy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = er.InstApINDy1
}

// InstApINDy1 reads a value from the address in `ar2`, sets `ar` with the value, and transitions the CPU to `InstApINDy2`.
//
//go:nosplit
func (er *Executor) InstApINDy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = er.InstApINDy2
}

// InstApINDy2 performs indirect indexed addressing by updating `ar` using `ar2` and `y`, and sets the next instruction.
// The function reads a byte from memory, updates `ar2`, adjusts `ar`, and determines the next handler based on conditions.
//
//go:nosplit
func (er *Executor) InstApINDy2(cpu *CPU) {
	data, ok := cpu.bus.Read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = er.InstApINDy3
	} else {
		cpu.next = er.InstApINDy4
	}
}

// InstApINDy3 handles indirect indexed addressing with Y-register offset without page crossing for the CPU.
// It reads from the address stored in the AR register, updates the instruction handler, and checks RDY state.
//
//go:nosplit
func (er *Executor) InstApINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = cpu.opTable[cpu.op]
}

// InstApINDy4 performs an indirect indexed addressing mode operation with Y register and updates the CPU state.
// If a page boundary is crossed during execution, it ensures proper handling and advances to the next instruction.
//
//go:nosplit
func (er *Executor) InstApINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeABSx fetches an 8-bit immediate value, increments the program counter, and stores it in the address register (ar).
//
//go:nosplit
func (er *Executor) InstAeABSx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstAeABSx1
}

// InstAeABSx1 handles the first step of the absolute-indexed addressing mode with the X register adjustment.
// It combines the X register with the address register and determines the next instruction to execute.
// Updates the address register (ar) and increments the program counter (pc).
//
//go:nosplit
func (er *Executor) InstAeABSx1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = er.InstAeABSx2
	}
}

// InstAeABSx2 handles the Absolute Indexed X addressing mode with page crossing, updating the address and next operation.
//
//go:nosplit
func (er *Executor) InstAeABSx2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeABSy initializes the address register with the value read from the program counter and sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstAeABSy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstAeABSy1
}

// InstAeABSy1 fetches data, updates the address register, and determines the next instruction based on the address range.
//
//go:nosplit
func (er *Executor) InstAeABSy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = er.InstAeABSy2
	}
}

// InstAeABSy2 adjusts the address register by adding a stack offset and sets the next instruction from the operation table.
// It ensures instruction execution respects memory boundary crossing conditions.
//
//go:nosplit
func (er *Executor) InstAeABSy2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstAeINDy executes the AE INDY instruction, updating the program counter and secondary address register, then sets the next operation.
//
//go:nosplit
func (er *Executor) InstAeINDy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = er.InstAeINDy1
}

// InstAeINDy1 reads data from the memory address specified by `ar2` and stores it in `ar`. Sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstAeINDy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = er.InstAeINDy2
}

// InstAeINDy2 executes the second phase of the AE indirect Y-indexed addressing mode.
// It reads a value from memory, combines it with the address register and Y register, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstAeINDy2(cpu *CPU) {
	data, ok := cpu.bus.Read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = cpu.opTable[cpu.op]
	} else {
		cpu.next = er.InstAeINDy3
	}
}

// InstAeINDy3 handles the addressing mode operation for instructions that cross a page boundary and updates the CPU state.
//
//go:nosplit
func (er *Executor) InstAeINDy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = cpu.opTable[cpu.op]
}

// InstMpZER is a memory operation instruction that reads data from memory into the address register (ar) and sets the next operation.
//
//go:nosplit
func (er *Executor) InstMpZER(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstOpRMW
}

// InstMpZERx executes a zero-page read operation, increments the program counter, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstMpZERx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstMpZERx1
}

// InstMpZERx1 performs a memory read operation at the address in `cpu.ar` and modifies the address by adding `cpu.x`.
// If the memory read fails, the function returns immediately.
// The adjusted address is constrained to an 8-bit boundary, and control proceeds to the `InstOpRMW` handler.
//
//go:nosplit
func (er *Executor) InstMpZERx1(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = er.InstOpRMW
}

// InstMpABS loads a byte from memory at the program counter into the address register and updates the next instruction.
//
//go:nosplit
func (er *Executor) InstMpABS(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstMpABS1
}

// InstMpABS1 represents an instruction handler that modifies the address register (AR) with data fetched from memory.
// It reads a byte from memory located at the program counter (PC), shifts it 8 bits left, and ORs it with the current AR.
// The program counter is incremented, and the next instruction handler is set to a read-modify-write (RMW) operation.
//
//go:nosplit
func (er *Executor) InstMpABS1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = er.InstOpRMW
}

// InstMpABSx fetches the next byte from memory, increments the program counter, and sets it in the address register.
// It updates the CPU's next state to InstMpABSx1.
//
//go:nosplit
func (er *Executor) InstMpABSx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstMpABSx1
}

// InstMpABSx1 is a CPU instruction handler for addressing mode manipulation.
// It reads the next byte from memory and uses it as part of the absolute address computation.
// Updates the address register (ar) with the computed address and determines the next instruction handler.
//
//go:nosplit
func (er *Executor) InstMpABSx1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = er.InstMpABSx2
	} else {
		cpu.next = er.InstMpABSx3
	}
}

// InstMpABSx2 performs an absolute memory read operation and checks for page crossing issues.
// It sets the next instruction to a read-modify-write operation if the memory read is successful.
//
//go:nosplit
func (er *Executor) InstMpABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = er.InstOpRMW
}

// InstMpABSx3 performs an operation using absolute addressing with additional adjustments and transitions to the next instruction.
//
//go:nosplit
func (er *Executor) InstMpABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = er.InstOpRMW
}

// InstMpABSy sets the address register (ar) based on the byte at the program counter, then sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstMpABSy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = er.InstMpABSy1
}

// InstMpABSy1 performs memory page addressing mode with Y register offset, updating `ar` and determining the next instruction.
//
//go:nosplit
func (er *Executor) InstMpABSy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = er.InstMpABSy2
	} else {
		cpu.next = er.InstMpABSy3
	}
}

// InstMpABSy2 handles the zero-page no-cross memory access and assigns the next instruction to a read-modify-write operation.
// It reads from the address register; if unsuccessful, it stops. Otherwise, it sets the next handler to InstOpRMW.
//
//go:nosplit
func (er *Executor) InstMpABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = er.InstOpRMW
}

// InstMpABSy3 adjusts the address register and sets up the next operation, handling page crossing scenarios.
//
//go:nosplit
func (er *Executor) InstMpABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = er.InstOpRMW
}

// InstMpINDx loads an operand from the program counter into the ar2 register and advances the program counter.
// Sets the next CPU instruction handler to InstMpINDx1. Returns immediately if the read operation fails.
//
//go:nosplit
func (er *Executor) InstMpINDx(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = er.InstMpINDx1
}

// InstMpINDx1 performs indexed indirect addressing mode update. Adjusts ar2 with x, wraps it, and sets the next operation.
//
//go:nosplit
func (er *Executor) InstMpINDx1(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = er.InstMpINDx2
}

// InstMpINDx2 reads data from the memory address in ar2, updates ar with this value if successful, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstMpINDx2(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = er.InstMpINDx3
}

// InstMpINDx3 performs an indexed memory read, updates the address register, and sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstMpINDx3(cpu *CPU) {
	data, ok := cpu.bus.Read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = er.InstOpRMW
}

// InstMpINDy reads a byte from memory at the program counter, increments the program counter, and updates ar2 and next state.
//
//go:nosplit
func (er *Executor) InstMpINDy(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = er.InstMpINDy1
}

// InstMpINDy1 reads data using the CPU's ar2 register, updates the ar register, and sets the next operation to InstMpINDy2.
//
//go:nosplit
func (er *Executor) InstMpINDy1(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = er.InstMpINDy2
}

// InstMpINDy2 handles the indirect indexed addressing mode logic by adjusting the address register and setting the next instruction.
//
//go:nosplit
func (er *Executor) InstMpINDy2(cpu *CPU) {
	data, ok := cpu.bus.Read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = er.InstMpINDy3
	} else {
		cpu.next = er.InstMpINDy4
	}
}

// InstMpINDy3 handles the memory instruction with indirect addressing, updating the CPU's next operation if successful.
//
//go:nosplit
func (er *Executor) InstMpINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = er.InstOpRMW
}

// InstMpINDy4 handles indirect indexed addressing with Y register offset and handles potential page crossing.
// If a page boundary is crossed, the function adjusts the address register and sets the next operation to InstOpRMW.
//
//go:nosplit
func (er *Executor) InstMpINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = er.InstOpRMW
}

// InstOpRMW reads data from the address in the CPU's address register, stores it in the `rmw` buffer, and sets the next operation.
//
//go:nosplit
func (er *Executor) InstOpRMW(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.rmw = data
	cpu.next = er.InstOpRMW1
}

// InstOpRMW1 executes the second phase of a read-modify-write operation by writing the modified value and updating the next instruction.
//
//go:nosplit
func (er *Executor) InstOpRMW1(cpu *CPU) {
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.next = cpu.opTable[cpu.op]
}

// InstOpLDA loads a byte from memory at the address in the address register (AR) into the accumulator (A).
// Updates the negative (N) and zero (Z) flags based on the loaded value.
// Sets the next instruction to `InstOpINI` if reading from memory is successful; does nothing otherwise.
//
//go:nosplit
func (er *Executor) InstOpLDA(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

// InstOiLDA loads a byte from memory into the accumulator, updates the negative and zero flags, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstOiLDA(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

// InstOpLDX loads a value from memory into the X register and updates the negative and zero flags.
//
//go:nosplit
func (er *Executor) InstOpLDX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

// InstOiLDX loads a byte from memory into the X register, updating the negative and zero flags, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstOiLDX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

//go:nosplit
func (er *Executor) InstOpLDY(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

// InstOiLDY loads a value into the Y register, updates the negative and zero flags, and sets the next instruction handler.
//
//go:nosplit
func (er *Executor) InstOiLDY(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = er.InstOpINI
}

// Store

// InstOpSTA stores the value in the accumulator (A register) into memory at the address stored in the address register (AR).
// It updates the next CPU instruction handler to InstOpINI.
//
//go:nosplit
func (er *Executor) InstOpSTA(cpu *CPU) {
	cpu.bus.Write(cpu.ar, cpu.a)
	cpu.next = er.InstOpINI
}

// InstOpSTX stores the value of the X register into memory at the address specified by the AR register and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstOpSTX(cpu *CPU) {
	cpu.bus.Write(cpu.ar, cpu.x)
	cpu.next = er.InstOpINI
}

// InstOpSTY stores the value of the Y register into memory at the address in the address register and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstOpSTY(cpu *CPU) {
	cpu.bus.Write(cpu.ar, cpu.y)
	cpu.next = er.InstOpINI
}

/*
//go:nosplit
func (er *Executor) InstMpZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstMpZERy1
}

//go:nosplit
func (er *Executor) InstMpZERy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = InstOpRMW
}
*/
