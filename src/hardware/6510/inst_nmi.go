package mos6510

// instOpNMI handles the Non-Maskable Interrupt (NMI) by setting the next instruction to instOpNMI1 if read is successful.
func instOpNMI(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpNMI1
}

// instOpNMI1 executes the first step of the Non-Maskable Interrupt (NMI) sequence and sets the next operation.
func instOpNMI1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpNMI2
}

// instOpNMI2 handles the second stage of the NMI interrupt by pushing the high byte of the return address onto the stack.
func instOpNMI2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpNMI3
}

// instOpNMI3 pushes the low byte of the return address onto the stack and updates the next CPU instruction to instOpNMI4.
func instOpNMI3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpNMI4
}

// instOpNMI4 handles the fourth step of the Non-Maskable Interrupt (NMI) sequence in the CPU's instruction set.
// It pushes the status register onto the stack, decrements the stack pointer, and disables further interrupts.
func instOpNMI4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instOpNMI5
}

// instOpNMI5 handles the non-maskable interrupt by reading the interrupt vector from address 0xfffa and updating the program counter.
// If reading fails, the function returns immediately without altering the CPU state.
// On success, it sets the next instruction handler to instOpNMI6.
func instOpNMI5(cpu *CPU) {
	//get irq vector from 0xfffa
	data, ok := cpu.read(0xfffa)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpNMI6
}

// instOpNMI6 handles the initialization of the Non-Maskable Interrupt (NMI) vector by reading the high byte from 0xfffb.
// It updates the program counter (PC) and sets the next CPU operation to instOpINI.
// The function halts execution if the memory read fails (e.g., invalid RDY state).
func instOpNMI6(cpu *CPU) {
	//get irq vector from 0xfffb
	data, ok := cpu.read(0xfffb)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}
