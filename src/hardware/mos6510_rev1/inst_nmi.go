package mos6510_rev1

// InstOpNMI handles the Non-Maskable Interrupt (NMI) by setting the next instruction to InstOpNMI1 if read is successful.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI(cpu *CPU) {
	//internal operation
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.next = er.InstOpNMI1
}

// InstOpNMI1 executes the first step of the Non-Maskable Interrupt (NMI) sequence and sets the next operation.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.next = er.InstOpNMI2
}

// InstOpNMI2 handles the second stage of the NMI interrupt by pushing the high byte of the return address onto the stack.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.bus.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = er.InstOpNMI3
}

// InstOpNMI3 pushes the low byte of the return address onto the stack and updates the next CPU instruction to InstOpNMI4.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.bus.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = er.InstOpNMI4
}

// InstOpNMI4 handles the fourth step of the Non-Maskable Interrupt (NMI) sequence in the CPU's instruction set.
// It pushes the status register onto the stack, decrements the stack pointer, and disables further interrupts.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.bus.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = er.InstOpNMI5
}

// InstOpNMI5 handles the non-maskable interrupt by reading the interrupt vector from address 0xfffa and updating the program counter.
// If reading fails, the function returns immediately without altering the CPU state.
// On success, it sets the next instruction handler to InstOpNMI6.
//
//go:nosplit
func (er *ControlUnit) InstOpNMI5(cpu *CPU) {
	//get irq vector from 0xfffa
	data, ok := cpu.bus.Read(0xfffa)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = er.InstOpNMI6
}

// InstOpNMI6 handles the initialization of the Non-Maskable Interrupt (NMI) vector by reading the high byte from 0xfffb.
// It updates the program counter (PC) and sets the next CPU operation to InstOpINI.
// The function halts execution if the memory read fails (e.g., invalid RDY state).
//
//go:nosplit
func (er *ControlUnit) InstOpNMI6(cpu *CPU) {
	//get irq vector from 0xfffb
	data, ok := cpu.bus.Read(0xfffb)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = er.InstOpINI
}
