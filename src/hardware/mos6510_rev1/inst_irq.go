package mos6510_rev1

// InstOpIRQ handles the IRQ (Interrupt Request) operation by checking the read status and transitioning to the next state.
//
//go:nosplit
func InstOpIRQ(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpIRQ1
}

// InstOpIRQ1 handles the initial stage of processing an IRQ (Interrupt Request) and sets up the next instruction state.
//
//go:nosplit
func InstOpIRQ1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpIRQ2
}

// InstOpIRQ2 pushes the high byte of the program counter onto the stack and prepares the next interrupt handler step.
//
//go:nosplit
func InstOpIRQ2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.bankWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = InstOpIRQ3
}

// InstOpIRQ3 handles the intermediate step of the IRQ handler by pushing the low byte of the return address onto the stack.
// Updates the stack pointer and sets the next operation to `InstOpIRQ4`.
//
//go:nosplit
func InstOpIRQ3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.bankWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = InstOpIRQ4
}

// InstOpIRQ4 handles the IRQ interrupt by pushing the status register onto the stack and updating the CPU state.
//
//go:nosplit
func InstOpIRQ4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.bankWrite((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = InstOpIRQ5
}

// InstOpIRQ5 handles the IRQ interrupt vector by fetching the address from memory at 0xfffe, updating the program counter,
// and setting the next instruction to InstOpIRQ6.
//
//go:nosplit
func InstOpIRQ5(cpu *CPU) {
	//get irq vector from 0xfffe
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = InstOpIRQ6
}

// InstOpIRQ6 handles the IRQ interrupt by reading the vector at address 0xFFFF, updating the program counter, and setting the next instruction.
//
//go:nosplit
func InstOpIRQ6(cpu *CPU) {
	//get irq vector from 0xffff
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = InstOpINI
}
