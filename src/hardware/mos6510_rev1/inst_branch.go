package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/references"
)

// Jump - Branch

// InstOpJMP is the instruction handler for the JMP operation, updating the program counter and setting the next operation.
//
//go:nosplit
func InstOpJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstOpJMP1
}

// InstOpJMP1 sets the program counter to a new address using a high byte from memory combined with the current address register.
//
//go:nosplit
func InstOpJMP1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = InstOpINI
}

// InstOiJMP handles the JMP (Jump) instruction by reading an address from memory and updating the program counter.
//
//go:nosplit
func InstOiJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = InstOiJMP1
}

// InstOiJMP1 sets the program counter high byte using the data from memory and prepares the next instruction.
//
//go:nosplit
func InstOiJMP1(cpu *CPU) {
	data, ok := cpu.read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = InstOpINI
}

// InstOpJSR performs the Jump to Subroutine (JSR) operation by reading the target address and setting the next instruction.
//
//go:nosplit
func InstOpJSR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = InstOpJSR1
}

// InstOpJSR1 processes the first step of the JSR (Jump to Subroutine) instruction, verifying the stack before proceeding.
//
//go:nosplit
func InstOpJSR1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.next = InstOpJSR2
}

// InstOpJSR2 handles the second stage of the JSR instruction by pushing the high byte of the program counter onto the stack.
//
//go:nosplit
func InstOpJSR2(cpu *CPU) {
	cpu.busWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = InstOpJSR3
}

// InstOpJSR3 handles the third step of the JSR (Jump to Subroutine) instruction, writing the low byte of the program counter to the stack.
//
//go:nosplit
func InstOpJSR3(cpu *CPU) {
	cpu.busWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = InstOpJSR4
}

// InstOpJSR4 executes the JSR4 (Jump Subroutine) operation, updating the program counter and setting the next instruction.
//
//go:nosplit
func InstOpJSR4(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = InstOpINI
}

// InstOpRTS handles the RTS (Return from Subroutine) instruction, setting the next CPU operation to InstOpRTS1 if successful.
//
//go:nosplit
func InstOpRTS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpRTS1
}

// InstOpRTS1 increments the stack pointer and sets the next instruction handler to InstOpRTS2 if stack read is successful.
//
//go:nosplit
func InstOpRTS1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = InstOpRTS2
}

// InstOpRTS2 handles the RTS operation by fetching the return address from the stack and updating the program counter.
//
//go:nosplit
func InstOpRTS2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = InstOpRTS3
}

// InstOpRTS3 retrieves the high byte of the return address from the stack and updates the program counter.
// Proceeds to the next step in the RTS operation by setting the instruction handler to InstOpRTS4.
//
//go:nosplit
func InstOpRTS3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = InstOpRTS4
}

// InstOpRTS4 increments the program counter and sets the next instruction handler to InstOpINI if the current PC is readable.
//
//go:nosplit
func InstOpRTS4(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = InstOpINI
}

// InstOpRTI handles the RTI (Return from Interrupt) instruction by preparing the CPU for the next instruction in sequence.
//
//go:nosplit
func InstOpRTI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = InstOpRTI1
}

// InstOpRTI1 handles the first phase of the RTI instruction by reading a byte from the stack and advancing the stack pointer.
//
//go:nosplit
func InstOpRTI1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = InstOpRTI2
}

// InstOpRTI2 handles the second stage of the RTI (Return from Interrupt) operation, restoring CPU flags and updating SP.
//
//go:nosplit
func InstOpRTI2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = InstOpRTI3
}

// InstOpRTI3 restores the program counter from the stack and advances the stack pointer, preparing the next CPU state.
//
//go:nosplit
func InstOpRTI3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = InstOpRTI4
}

// InstOpRTI4 restores the most significant byte of the program counter from the stack and updates the next instruction.
//
//go:nosplit
func InstOpRTI4(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = InstOpINI
}

// InstOpBRK handles the BRK (Break) instruction by incrementing the program counter and setting the next execution step.
//
//go:nosplit
func InstOpBRK(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = InstOpBRK1
}

// InstOpBRK1 handles the first phase of the BRK (Break) operation by pushing the high byte of the program counter onto the stack.
// The stack pointer is decremented, and the next instruction phase is set to InstOpBRK2.
//
//go:nosplit
func InstOpBRK1(cpu *CPU) {
	cpu.busWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = InstOpBRK2
}

// InstOpBRK2 handles the second stage of the BRK instruction by storing the low byte of the program counter to the stack.
// It then decrements the stack pointer and prepares the CPU for the next instruction phase by setting `next` to InstOpBRK3.
//
//go:nosplit
func InstOpBRK2(cpu *CPU) {
	cpu.busWrite(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = InstOpBRK3
}

// InstOpBRK3 handles the BRK (break) instruction sequence by pushing flags, managing stack writes, and setting the next operation.
//
//go:nosplit
func InstOpBRK3(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.busWrite((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	if cpu.picHasNMI() {
		cpu.picClearNMI()     // Simulate an edge-triggered input
		cpu.next = InstOpNMI5 // Jump to NMI sequence
	} else {
		cpu.next = InstOpBRK4
	}
}

// InstOpBRK4 handles the BRK instruction by reading the interrupt vector from memory and updating the program counter.
//
//go:nosplit
func InstOpBRK4(cpu *CPU) {
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = InstOpBRK5
}

// InstOpBRK5 handles the BRK instruction for the CPU by updating the program counter and setting the next instruction handler.
//
//go:nosplit
func InstOpBRK5(cpu *CPU) {
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = InstOpINI
}

// InstOpBCS handles the Branch if Carry Set (BCS) instruction, branching if the carry flag is set and advancing the PC.
//
//go:nosplit
func InstOpBCS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBCC handles the branch if carry flag is clear, updating the program counter or setting the next operation accordingly.
//
//go:nosplit
func InstOpBCC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBEQ handles the BEQ (Branch if Equal) instruction, branching if the zero flag is set, or proceeding to the next operation.
//
//go:nosplit
func InstOpBEQ(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBNE handles the BNE (Branch if Not Equal) instruction, branching if the zero flag (zFlag) is not set.
//
//go:nosplit
func InstOpBNE(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBVS handles the BVS (Branch if Overflow Set) instruction, branching if the overflow flag is set.
//
//go:nosplit
func InstOpBVS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag == 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBVC handles the BVC (Branch if Overflow Clear) instruction, branching if the overflow flag (vFlag) is 0.
//
//go:nosplit
func InstOpBVC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag != 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBMI executes the BMI (Branch if Minus) instruction, branching if the negative flag (nFlag) is set.
//
//go:nosplit
func InstOpBMI(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBPL handles the BPL (Branch on Positive) operation by branching if the negative flag (N) is clear.
//
//go:nosplit
func InstOpBPL(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.next = InstOpINI
	} else {
		cpu.next = cpu.computeBranch(data)
	}
}

// InstOpBRAnp handles a branch operation without crossing a page boundary, updating the program counter and setting the next instruction.
//
//go:nosplit
func InstOpBRAnp(cpu *CPU) {
	// No page crossed
	cpu.opFlags |= references.OpFlagIntDelayed
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = InstOpINI
}

// InstOpBRAbp sets the program counter to the address register for backward branching after a page crossing.
//
//go:nosplit
func InstOpBRAbp(cpu *CPU) {
	// Page crossed (branch backwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = InstOpBRAbp1
}

// InstOpBRAbp1 handles branching logic when a specific overflow condition is met and sets the next instruction to execute.
//
//go:nosplit
func InstOpBRAbp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc + stackAddr); !ok {
		return
	}
	cpu.next = InstOpINI
}

// InstOpBRAfp handles the forward branch operation by updating the program counter to the address register (ar).
// If the branch crosses a page, it ensures the next instruction is set accordingly.
//
//go:nosplit
func InstOpBRAfp(cpu *CPU) {
	// Page crossed (branch forwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = InstOpBRAfp1
}

// InstOpBRAfp1 executes a branch operation if the required memory condition is met; sets the next instruction to InstOpINI.
//
//go:nosplit
func InstOpBRAfp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc - stackAddr); !ok {
		return
	}
	cpu.next = InstOpINI
}
