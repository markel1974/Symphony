package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/references"
)

// Stack

// InstOpPHA handles the PHA (Push Accumulator) operation, verifying CPU state and setting the next instruction step.
//
//go:nosplit
func InstOpPHA(cpu *CPU) {
	if _, ok := cpu.busRead(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpPHA1
}

// InstOpPHA1 pushes the accumulator onto the stack, updates the stack pointer, and sets the next operation to InstOpINI.
//
//go:nosplit
func InstOpPHA1(cpu *CPU) {
	cpu.busWrite(uint16(cpu.sp)|stackAddr, cpu.a)
	cpu.sp--
	cpu.next = InstOpINI
}

// InstOpPLA prepares the CPU for the PLA (Pull Accumulator) instruction by setting the next state for execution.
//
//go:nosplit
func InstOpPLA(cpu *CPU) {
	if _, ok := cpu.busRead(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpPLA1
}

// InstOpPLA1 handles the PLA (Pull Accumulator) step 1 by reading the stack and preparing for the next operation.
//
//go:nosplit
func InstOpPLA1(cpu *CPU) {
	if _, ok := cpu.busRead(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = InstOpPLA2
}

// InstOpPLA2 is an instruction handler that pulls a byte from the stack into the accumulator and updates flags.
//
//go:nosplit
func InstOpPLA2(cpu *CPU) {
	data, ok := cpu.busRead(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = InstOpINI
}

// InstOpPHP initiates the PHP (Push Processor Status) instruction by preparing to save the processor flags onto the stack.
//
//go:nosplit
func InstOpPHP(cpu *CPU) {
	if _, ok := cpu.busRead(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpPHP1
}

// InstOpPHP1 pushes the CPU status flags onto the stack, decrements the stack pointer, and sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpPHP1(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.busWrite((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.next = InstOpINI
}

// InstOpPLP processes the PLP (Pull Processor Status) instruction by advancing the CPU state to the next handler.
//
//go:nosplit
func InstOpPLP(cpu *CPU) {
	if _, ok := cpu.busRead(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpPLP1
}

// InstOpPLP1 reads a byte from the stack. If unsuccessful, exits; otherwise increments the stack pointer and sets InstOpPLP2.
//
//go:nosplit
func InstOpPLP1(cpu *CPU) {
	if _, ok := cpu.busRead(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = InstOpPLP2
}

// InstOpPLP2 retrieves a byte from the stack, updates CPU flags, and manages IRQ enable/disable transitions.
//
//go:nosplit
func InstOpPLP2(cpu *CPU) {
	data, ok := cpu.busRead(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	iFlagPrev := cpu.iFlag
	cpu.popFlags(data)
	if iFlagPrev == 0 && cpu.iFlag != 0 {
		cpu.opFlags |= references.OpFlagIrqDisabled
	} else if iFlagPrev != 0 && cpu.iFlag == 0 {
		cpu.opFlags |= references.OpFlagIrqEnabled
	}
	cpu.next = InstOpINI
}
