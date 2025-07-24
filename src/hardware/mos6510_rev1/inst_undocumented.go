package mos6510_rev1

import (
	"fmt"
	"log"
	"os"
)

// Undocumented functions

// NOP

// InstOiNOP increments the program counter and sets the next instruction to InstOpINI if the current read is successful.
//
//go:nosplit
func (er *Executor) InstOiNOP(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = er.InstOpINI
}

// InstOaNOP is a no-operation instruction handler that validates memory accessibility and sets the next instruction to InstOpINI.
//
//go:nosplit
func (er *Executor) InstOaNOP(cpu *CPU) {
	if _, ok := cpu.bus.Read(cpu.ar); !ok {
		return
	}
	cpu.next = er.InstOpINI
}

// Load A/X

// InstOpLAX handles the LAX instruction, which loads data into both the A and X registers and updates the N and Z flags.
//
//go:nosplit
func (er *Executor) InstOpLAX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// Store A/X

// InstOpSAX executes the SAX (Store AND X) operation, storing the logical AND of registers A and X into memory.
// It writes the result to the address register (ar) and sets the next instruction to InstOpINI.
//
//go:nosplit
func (er *Executor) InstOpSAX(cpu *CPU) {
	cpu.bus.Write(cpu.ar, cpu.a&cpu.x)
	cpu.next = er.InstOpINI
}

// ASL/ORA

// InstOpSLO executes the SLO instruction, which shifts left the RMW value, updates flags, and performs OR with A register.
//
//go:nosplit
func (er *Executor) InstOpSLO(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	cpu.rmw <<= 1
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.a |= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// ROL/AND

// InstOpRLA performs a ROL (Rotate Left) operation on memory with accumulator AND logic and updates CPU flags.
//
//go:nosplit
func (er *Executor) InstOpRLA(cpu *CPU) {
	tmp := cpu.rmw & 0x80
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw << 1) | 0x1
	} else {
		cpu.rmw = cpu.rmw << 1
	}
	cpu.cFlag = tmp
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.a &= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// LSR/EOR

// InstOpSRE performs the SRE instruction: shifts memory right, XORs with accumulator, and updates flags and memory.
//
//go:nosplit
func (er *Executor) InstOpSRE(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	cpu.rmw >>= 1
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.a ^= cpu.rmw
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// ROR/ADC

// InstOpRRA executes the RRA (Rotate Right and add) operation, updating carry, performing ADC, and setting the next instruction.
//
//go:nosplit
func (er *Executor) InstOpRRA(cpu *CPU) {
	tmp := cpu.rmw & 0x1
	if cpu.cFlag != 0 {
		cpu.rmw = (cpu.rmw >> 1) | 0x80
	} else {
		cpu.rmw = cpu.rmw >> 1
	}
	cpu.cFlag = tmp
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.doADC(cpu.rmw)
	cpu.next = er.InstOpINI
}

// DEC/CMP

// InstOpDCP performs a decrement and compare operation, modifying flags and memory based on the CPU state.
//
//go:nosplit
func (er *Executor) InstOpDCP(cpu *CPU) {
	cpu.rmw--
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.ar = uint16(cpu.a) - uint16(cpu.rmw)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	//cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = er.InstOpINI
}

// InstOpISB executes the ISB instruction, writing data, updating the RMW buffer, performing SBC, and setting the next op.
//
//go:nosplit
func (er *Executor) InstOpISB(cpu *CPU) {
	cpu.rmw++
	cpu.bus.Write(cpu.ar, cpu.rmw)
	cpu.doSBC(cpu.rmw)
	cpu.next = er.InstOpINI
}

// InstOiANC performs a logical AND operation between the accumulator and a fetched value, updating CPU flags accordingly.
//
//go:nosplit
func (er *Executor) InstOiANC(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = cpu.nFlag & 0x80
	cpu.next = er.InstOpINI
}

// InstOiASR performs an AND operation with a fetched byte and shifts the accumulator right by one bit, updating CPU flags.
//
//go:nosplit
func (er *Executor) InstOiASR(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// InstOiARR performs a bitwise AND operation between the accumulator and a memory value, followed by a right shift.
// Updates accumulator and various flags (N, Z, C, V) based on the CPU state and operation result.
//
//go:nosplit
func (er *Executor) InstOiARR(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
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
	cpu.next = er.InstOpINI
}

// InstOiANE performs a bitwise operation on the accumulator with a constant and memory data, updates flags, and sets the next instruction.
//
//go:nosplit
func (er *Executor) InstOiANE(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = (cpu.a | 0xee) & cpu.x & data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// InstOiLXA performs a bitwise OR operation between the accumulator and 0xee, then ANDs the result with a fetched byte.
// Updates the X and A registers, and sets the negative and zero flags based on the new accumulator value.
// Advances the program counter and assigns the next instruction handler.
//
//go:nosplit
func (er *Executor) InstOiLXA(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = (cpu.a | 0xee) & data
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// InstOiSBX executes the OiSBX instruction, performing a logical AND of A and X, subtraction with memory, and updates flags.
//
//go:nosplit
func (er *Executor) InstOiSBX(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = (uint16(cpu.x) & uint16(cpu.a)) - uint16(data)
	cpu.x = uint8(cpu.ar)
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	//cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	if cpu.ar < stackAddr {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.next = er.InstOpINI
}

// InstOpLAS performs a logical AND operation between the stack pointer and memory, updating the X and A registers and flags.
//
//go:nosplit
func (er *Executor) InstOpLAS(cpu *CPU) {
	data, ok := cpu.bus.Read(cpu.ar)
	if !ok {
		return
	}
	cpu.sp = data & cpu.sp
	cpu.x = cpu.sp
	cpu.a = cpu.x
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = er.InstOpINI
}

// InstOpSHS performs the SHS (Store Stack Pointer and Memory) operation, combining the A and X registers to set SP and memory.
//
//go:nosplit
func (er *Executor) InstOpSHS(cpu *CPU) {
	cpu.sp = cpu.a & cpu.x
	d := uint8((cpu.ar2 + 1) & uint16(cpu.sp))
	cpu.bus.Write(cpu.ar, d)
	cpu.next = er.InstOpINI
}

// InstOpSHY computes a data value using the Y register and writes it to a memory address determined by CPU state.
// It sets the next instruction to InstOpINI.
//
//go:nosplit
func (er *Executor) InstOpSHY(cpu *CPU) {
	d := uint8(uint16(cpu.y) & (cpu.ar2 + 1))
	cpu.bus.Write(cpu.ar, d)
	cpu.next = er.InstOpINI
}

// InstOpSHX performs the SHX instruction, calculating a value from the X register and AR2, then writing it to memory.
//
//go:nosplit
func (er *Executor) InstOpSHX(cpu *CPU) {
	d := uint8(uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.bus.Write(cpu.ar, d)
	cpu.next = er.InstOpINI
}

// InstOpSHA stores the logical AND of registers A, X, and (ar2 + 1) into memory at address specified by ar.
// It then sets the next CPU instruction to `InstOpINI`.
//
//go:nosplit
func (er *Executor) InstOpSHA(cpu *CPU) {
	d := uint8(uint16(cpu.a) & uint16(cpu.x) & (cpu.ar2 + 1))
	cpu.bus.Write(cpu.ar, d)
	cpu.next = er.InstOpINI
}

// InstOpJAM logs an illegal opcode error with CPU context, resets the CPU, and exits the application.
//
//go:nosplit
func (er *Executor) InstOpJAM(cpu *CPU) {
	err := fmt.Errorf("[%s] unknown opcode %02x at %04x", cpu.GetId(), cpu.op, cpu.pc-1)
	log.Println(err.Error())
	os.Exit(1)
	//TODO EVENT
	//cpu.Reset()
}
