package mos6510

import (
	"github.com/markel1974/c64emu/src/references"
)

// Jump - Branch

// instOpJMP is the instruction handler for the JMP operation, updating the program counter and setting the next operation.
func instOpJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpJMP1
}

// instOpJMP1 sets the program counter to a new address using a high byte from memory combined with the current address register.
func instOpJMP1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc = (uint16(data) << 8) | cpu.ar
	cpu.next = instOpINI
}

// instOiJMP handles the JMP (Jump) instruction by reading an address from memory and updating the program counter.
func instOiJMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOiJMP1
}

// instOiJMP1 sets the program counter high byte using the data from memory and prepares the next instruction.
func instOiJMP1(cpu *CPU) {
	data, ok := cpu.read(((cpu.ar + 1) & 0xff) | (cpu.ar & 0xff00))
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

// instOpJSR performs the Jump to Subroutine (JSR) operation by reading the target address and setting the next instruction.
func instOpJSR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpJSR1
}

// instOpJSR1 processes the first step of the JSR (Jump to Subroutine) instruction, verifying the stack before proceeding.
func instOpJSR1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.next = instOpJSR2
}

// instOpJSR2 handles the second stage of the JSR instruction by pushing the high byte of the program counter onto the stack.
func instOpJSR2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpJSR3
}

// instOpJSR3 handles the third step of the JSR (Jump to Subroutine) instruction, writing the low byte of the program counter to the stack.
func instOpJSR3(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpJSR4
}

// instOpJSR4 executes the JSR4 (Jump Subroutine) operation, updating the program counter and setting the next instruction.
func instOpJSR4(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.pc = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpINI
}

// instOpRTS handles the RTS (Return from Subroutine) instruction, setting the next CPU operation to instOpRTS1 if successful.
func instOpRTS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpRTS1
}

// instOpRTS1 increments the stack pointer and sets the next instruction handler to instOpRTS2 if stack read is successful.
func instOpRTS1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpRTS2
}

// instOpRTS2 handles the RTS operation by fetching the return address from the stack and updating the program counter.
func instOpRTS2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = instOpRTS3
}

// instOpRTS3 retrieves the high byte of the return address from the stack and updates the program counter.
// Proceeds to the next step in the RTS operation by setting the instruction handler to instOpRTS4.
func instOpRTS3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpRTS4
}

// instOpRTS4 increments the program counter and sets the next instruction handler to instOpINI if the current PC is readable.
func instOpRTS4(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpINI
}

// instOpRTI handles the RTI (Return from Interrupt) instruction by preparing the CPU for the next instruction in sequence.
func instOpRTI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpRTI1
}

// instOpRTI1 handles the first phase of the RTI instruction by reading a byte from the stack and advancing the stack pointer.
func instOpRTI1(cpu *CPU) {
	if _, ok := cpu.read(uint16(cpu.sp) | stackAddr); !ok {
		return
	}
	cpu.sp++
	cpu.next = instOpRTI2
}

// instOpRTI2 handles the second stage of the RTI (Return from Interrupt) operation, restoring CPU flags and updating SP.
func instOpRTI2(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.popFlags(data)
	cpu.sp++
	cpu.next = instOpRTI3
}

// instOpRTI3 restores the program counter from the stack and advances the stack pointer, preparing the next CPU state.
func instOpRTI3(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.sp++
	cpu.next = instOpRTI4
}

// instOpRTI4 restores the most significant byte of the program counter from the stack and updates the next instruction.
func instOpRTI4(cpu *CPU) {
	data, ok := cpu.read(uint16(cpu.sp) | stackAddr)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

// instOpBRK handles the BRK (Break) instruction by incrementing the program counter and setting the next execution step.
func instOpBRK(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc++
	cpu.next = instOpBRK1
}

// instOpBRK1 handles the first phase of the BRK (Break) operation by pushing the high byte of the program counter onto the stack.
// The stack pointer is decremented, and the next instruction phase is set to instOpBRK2.
func instOpBRK1(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpBRK2
}

// instOpBRK2 handles the second stage of the BRK instruction by storing the low byte of the program counter to the stack.
// It then decrements the stack pointer and prepares the CPU for the next instruction phase by setting `next` to instOpBRK3.
func instOpBRK2(cpu *CPU) {
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpBRK3
}

// instOpBRK3 handles the BRK (break) instruction sequence by pushing flags, managing stack writes, and setting the next operation.
func instOpBRK3(cpu *CPU) {
	data := cpu.pushFlags(true)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	if cpu.pic.HasNMI() {
		cpu.pic.ClearNMI()    // Simulate an edge-triggered input
		cpu.next = instOpNMI5 // Jump to NMI sequence
	} else {
		cpu.next = instOpBRK4
	}
}

// instOpBRK4 handles the BRK instruction by reading the interrupt vector from memory and updating the program counter.
func instOpBRK4(cpu *CPU) {
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpBRK5
}

// instOpBRK5 handles the BRK instruction for the CPU by updating the program counter and setting the next instruction handler.
func instOpBRK5(cpu *CPU) {
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

// instOpBCS handles the Branch if Carry Set (BCS) instruction, branching if the carry flag is set and advancing the PC.
func instOpBCS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBCC handles the branch if carry flag is clear, updating the program counter or setting the next operation accordingly.
func instOpBCC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.cFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBEQ handles the BEQ (Branch if Equal) instruction, branching if the zero flag is set, or proceeding to the next operation.
func instOpBEQ(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBNE handles the BNE (Branch if Not Equal) instruction, branching if the zero flag (zFlag) is not set.
func instOpBNE(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.zFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBVS handles the BVS (Branch if Overflow Set) instruction, branching if the overflow flag is set.
func instOpBVS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBVC handles the BVC (Branch if Overflow Clear) instruction, branching if the overflow flag (vFlag) is 0.
func instOpBVC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if cpu.overflowBranch != nil && cpu.overflowBranch() {
		cpu.vFlag = 1
	}
	if cpu.vFlag != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBMI executes the BMI (Branch if Minus) instruction, branching if the negative flag (nFlag) is set.
func instOpBMI(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) == 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBPL handles the BPL (Branch on Positive) operation by branching if the negative flag (N) is clear.
func instOpBPL(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	if (cpu.nFlag & 0x80) != 0 {
		cpu.next = instOpINI
	} else {
		cpu.branch(data)
	}
}

// instOpBRAnp handles a branch operation without crossing a page boundary, updating the program counter and setting the next instruction.
func instOpBRAnp(cpu *CPU) {
	// No page crossed
	cpu.opFlags |= references.OpFlagIntDelayed
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpINI
}

// instOpBRAbp sets the program counter to the address register for backward branching after a page crossing.
func instOpBRAbp(cpu *CPU) {
	// Page crossed (branch backwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpBRAbp1
}

// instOpBRAbp1 handles branching logic when a specific overflow condition is met and sets the next instruction to execute.
func instOpBRAbp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc + stackAddr); !ok {
		return
	}
	cpu.next = instOpINI
}

// instOpBRAfp handles the forward branch operation by updating the program counter to the address register (ar).
// If the branch crosses a page, it ensures the next instruction is set accordingly.
func instOpBRAfp(cpu *CPU) {
	// Page crossed (branch forwards)
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.pc = cpu.ar
	cpu.next = instOpBRAfp1
}

// instOpBRAfp1 executes a branch operation if the required memory condition is met; sets the next instruction to instOpINI.
func instOpBRAfp1(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc - stackAddr); !ok {
		return
	}
	cpu.next = instOpINI
}
