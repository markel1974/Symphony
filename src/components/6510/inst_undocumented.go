package mos6510

import (
	"github.com/markel1974/c64emu/src/common/conversion"
	"log"
	"os"
)

// Undocumented functions

// NOP

// instOiNOP increments the program counter and sets the next instruction to instOpINI if the current read is successful.
func instOiNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpINI
}

// instOaNOP is a no-operation instruction handler that validates memory accessibility and sets the next instruction to instOpINI.
func instOaNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpINI
}

// Load A/X

// instOpLAX handles the LAX instruction, which loads data into both the A and X registers and updates the N and Z flags.
func instOpLAX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// Store A/X

// instOpSAX executes the SAX (Store AND X) operation, storing the logical AND of registers A and X into memory.
// It writes the result to the address register (ar) and sets the next instruction to instOpINI.
func instOpSAX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = instOpINI
}

// ASL/ORA

// instOpSLO executes the SLO instruction, which shifts left the RMW value, updates flags, and performs OR with A register.
func instOpSLO(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// ROL/AND

// instOpRLA performs a ROL (Rotate Left) operation on memory with accumulator AND logic and updates CPU flags.
func instOpRLA(cpu *CPU) {
	tmp := cpu.rmw & 0x80
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw << 1) | 0x1
	} else {
		cpu.rmw = cpu.rmw << 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a &= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// LSR/EOR

// instOpSRE performs the SRE instruction: shifts memory right, XORs with accumulator, and updates flags and memory.
func instOpSRE(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	cpu.rmw >>= 1
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// ROR/ADC

// instOpRRA executes the RRA (Rotate Right and Add) operation, updating carry, performing ADC, and setting the next instruction.
func instOpRRA(cpu *CPU) {
	tmp := cpu.rmw & 0x1
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw >> 1) | 0x80
	} else {
		cpu.rmw = cpu.rmw >> 1
	}
	cpu.cFlag = tmp
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doADC(cpu.rmw)
	cpu.next = instOpINI
}

// DEC/CMP

// instOpDCP performs a decrement and compare operation, modifying flags and memory based on the CPU state.
func instOpDCP(cpu *CPU) {
	cpu.rmw--
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpISB executes the ISB instruction, writing data, updating the RMW buffer, performing SBC, and setting the next op.
func instOpISB(cpu *CPU) {
	cpu.rmw++
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = instOpINI
}

// instOiANC performs a logical AND operation between the accumulator and a fetched value, updating CPU flags accordingly.
func instOiANC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.next = instOpINI
}

// instOiASR performs an AND operation with a fetched byte and shifts the accumulator right by one bit, updating CPU flags.
func instOiASR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiARR performs a bitwise AND operation between the accumulator and a memory value, followed by a right shift.
// Updates accumulator and various flags (N, Z, C, V) based on the CPU state and operation result.
func instOiARR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	data &= cpu.a
	if cpu.cFlag != 0 {
		cpu.a = (data >> 1) | 0x80
	} else {
		cpu.a = data >> 1
	}
	if cpu.dFlag == 0 {
		cpu.nFlag = cpu.a
		cpu.zFlag = cpu.a
		cpu.cFlag = cpu.a & 0x40
		cpu.vFlag = (cpu.a & 0x40) ^ ((cpu.a & 0x20) << 1)
	} else {
		if cpu.cFlag != 0 {
			cpu.nFlag = 0x80
		} else {
			cpu.nFlag = 0
		}
		cpu.zFlag = cpu.a
		cpu.vFlag = (data ^ cpu.a) & 0x40
		if ((data & 0xf) + (data & 0x1)) > 5 {
			cpu.a = (cpu.a & 0xf0) | ((cpu.a + 6) & 0xf)
		}
		k := uint16((data)+(uint8(data)&0x10)) & 0x1f0
		cpu.cFlag = uint8(k)
		if k > 0x50 {
			cpu.a += 0x60
		}
	}
	cpu.next = instOpINI
}

// instOiANE performs a bitwise operation on the accumulator with a constant and memory data, updates flags, and sets the next instruction.
func instOiANE(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiLXA performs a bitwise OR operation between the accumulator and 0xee, then ANDs the result with a fetched byte.
// Updates the X and A registers, and sets the negative and zero flags based on the new accumulator value.
// Advances the program counter and assigns the next instruction handler.
func instOiLXA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = (cpu.a | 0xee) & data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiSBX executes the OiSBX instruction, performing a logical AND of A and X, subtraction with memory, and updates flags.
func instOiSBX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = (uint16(cpu.x) & uint16(cpu.a)) - uint16(data)
	cpu.x = uint8(cpu.ar)
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpLAS performs a logical AND operation between the stack pointer and memory, updating the X and A registers and flags.
func instOpLAS(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.sp = data & cpu.sp
	cpu.x = cpu.sp
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpSHS performs the SHS (Store Stack Pointer and Memory) operation, combining the A and X registers to set SP and memory.
func instOpSHS(cpu *CPU) {
	cpu.sp = cpu.a & cpu.x
	d := uint8((cpu.ar2 + 1) & uint16(cpu.sp))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

// instOpSHY computes a data value using the Y register and writes it to a memory address determined by CPU state.
// It sets the next instruction to instOpINI.
func instOpSHY(cpu *CPU) {
	d := uint8(uint16(cpu.y) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

// instOpSHX performs the SHX instruction, calculating a value from the X register and AR2, then writing it to memory.
func instOpSHX(cpu *CPU) {
	d := uint8(uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

// instOpSHA stores the logical AND of registers A, X, and (ar2 + 1) into memory at address specified by ar.
// It then sets the next CPU instruction to `instOpINI`.
func instOpSHA(cpu *CPU) {
	d := uint8(uint16(cpu.a) & uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.banks.Write(cpu.ar, d)
	cpu.next = instOpINI
}

// instOpJAM logs an illegal opcode error with CPU context, resets the CPU, and exits the application.
func instOpJAM(cpu *CPU) {
	log.Printf("[%s] illegal opcode %02x at %04x.", cpu.id, cpu.op, cpu.pc-1)
	//TODO EVENT
	cpu.Reset()
	os.Exit(1)
}
