package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/references"
)

// Stack

// instOpPHA handles the PHA (Push Accumulator) operation, verifying CPU state and setting the next instruction step.
func instOpPHA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPHA1
}

// instOpPHA1 pushes the accumulator onto the stack, updates the stack pointer, and sets the next operation to instOpINI.
func instOpPHA1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, cpu.a)
	cpu.sp--
	cpu.next = instOpINI
}

// instOpPLA prepares the CPU for the PLA (Pull Accumulator) instruction by setting the next state for execution.
func instOpPLA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPLA1
}

// instOpPLA1 handles the PLA (Pull Accumulator) step 1 by reading the stack and preparing for the next operation.
func instOpPLA1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpPLA2
}

// instOpPLA2 is an instruction handler that pulls a byte from the stack into the accumulator and updates flags.
func instOpPLA2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpPHP initiates the PHP (Push Processor Status) instruction by preparing to save the processor flags onto the stack.
func instOpPHP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPHP1
}

// instOpPHP1 pushes the CPU status flags onto the stack, decrements the stack pointer, and sets the next instruction to instOpINI.
func instOpPHP1(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.next = instOpINI
}

// instOpPLP processes the PLP (Pull Processor Status) instruction by advancing the CPU state to the next handler.
func instOpPLP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpPLP1
}

// instOpPLP1 reads a byte from the stack. If unsuccessful, exits; otherwise increments the stack pointer and sets instOpPLP2.
func instOpPLP1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpPLP2
}

// instOpPLP2 retrieves a byte from the stack, updates CPU flags, and manages IRQ enable/disable transitions.
func instOpPLP2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
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
	cpu.next = instOpINI
}
