package mos6510

import (
	"github.com/markel1974/c64emu/src/conversion"
	"log"
	"os"
)

// stackAddr defines the base address for the stack memory operations used in various CPU instructions.
const (
	stackAddr = 0x100
)

// instOpINI handles the initial opcode fetch and subsequent CPU instruction cycle logic based on the current CPU state.
// It considers interrupt conditions, updates the program counter, and sets the next instruction handler.
// If the RDY line is low, the CPU execution is halted by setting the `stop` flag.
func instOpINI(cpu *CPU) {
	// https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt
	// Interrupts are only recognized if the RDY line is high
	if !cpu.rdyLow {
		if !cpu.irqBreaker {
			opFlag := cpu.opFlags
			cpu.opFlags = 0
			switch cpu.pic.VerifyIrq(cpu.iFlag, opFlag) {
			case 1:
				cpu.Reset()
				return
			case 2:
				cpu.irqBreaker = true
				cpu.next = instOpNMI
				cpu.next(cpu)
				return
			case 3:
				cpu.irqBreaker = true
				cpu.next = instOpIRQ
				cpu.next(cpu)
				return
			}
		}
	} else {
		cpu.stop = true
		return
	}
	cpu.irqBreaker = false
	cpu.op = cpu.banks.Read(cpu.pc)
	cpu.pc++
	cpu.next = _modeTable[cpu.op]
}

// instOpIRQ handles the IRQ (Interrupt Request) operation by checking the read status and transitioning to the next state.
func instOpIRQ(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpIRQ1
}

// instOpIRQ1 handles the initial stage of processing an IRQ (Interrupt Request) and sets up the next instruction state.
func instOpIRQ1(cpu *CPU) {
	//internal operation
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpIRQ2
}

// instOpIRQ2 pushes the high byte of the program counter onto the stack and prepares the next interrupt handler step.
func instOpIRQ2(cpu *CPU) {
	//push return address high byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc>>8))
	cpu.sp--
	cpu.next = instOpIRQ3
}

// instOpIRQ3 handles the intermediate step of the IRQ handler by pushing the low byte of the return address onto the stack.
// Updates the stack pointer and sets the next operation to `instOpIRQ4`.
func instOpIRQ3(cpu *CPU) {
	//push return address low byte onto stack
	cpu.banks.Write(uint16(cpu.sp)|stackAddr, uint8(cpu.pc))
	cpu.sp--
	cpu.next = instOpIRQ4
}

// instOpIRQ4 handles the IRQ interrupt by pushing the status register onto the stack and updating the CPU state.
func instOpIRQ4(cpu *CPU) {
	//push status register onto stack
	data := cpu.pushFlags(false)
	cpu.banks.Write((uint16(cpu.sp)&0xff)|stackAddr, data)
	cpu.sp--
	cpu.iFlag = 1
	cpu.next = instOpIRQ5
}

// instOpIRQ5 handles the IRQ interrupt vector by fetching the address from memory at 0xfffe, updating the program counter,
// and setting the next instruction to instOpIRQ6.
func instOpIRQ5(cpu *CPU) {
	//get irq vector from 0xfffe
	data, ok := cpu.read(0xfffe)
	if !ok {
		return
	}
	cpu.pc = uint16(data)
	cpu.next = instOpIRQ6
}

// instOpIRQ6 handles the IRQ interrupt by reading the vector at address 0xFFFF, updating the program counter, and setting the next instruction.
func instOpIRQ6(cpu *CPU) {
	//get irq vector from 0xffff
	data, ok := cpu.read(0xffff)
	if !ok {
		return
	}
	cpu.pc |= uint16(data) << 8
	cpu.next = instOpINI
}

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

// instApZER loads a zero-page address into the address register and sets the next instruction handler.
func instApZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = _opTable[cpu.op]
}

// instApZERx performs a zero-page addressing operation by reading a byte at the program counter and updating the address register.
// It increments the program counter and sets the next instruction to instApZERx1.
func instApZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERx1
}

// instApZERx1 performs a zero-page indexed addressing mode operation using the CPU's address register and X register.
// It updates the address register by adding the X register value, ensuring the result wraps around at 8-bit boundaries.
// The next instruction handler is set based on the current opcode. No operation occurs if the read fails.
func instApZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = _opTable[cpu.op]
}

// instApZERy loads a byte from memory at the program counter into the address register and sets the next instruction.
func instApZERy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApZERy1
}

// instApZERy1 adjusts the address register by adding the Y register value and wraps it within a byte boundary.
// Proceeds to the next operation if the address read was successful.
func instApZERy1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
	cpu.next = _opTable[cpu.op]
}

// instApABS loads the next byte from memory into the address register and advances the program counter.
func instApABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABS1
}

// instApABS1 reads a byte from memory at the program counter, increments the PC, and updates the address register (AR).
// Then, it fetches the next instruction from the operation table based on the current opcode.
func instApABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

// instApABSx fetches a byte from the program counter, increments the PC, and stores the value in the address register.
// Sets the next instruction handler to instApABSx1.
func instApABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSx1
}

// instApABSx1 executes the first step of an absolute addressing mode with X offset, updating CPU registers and next instruction.
func instApABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSx2
	} else {
		cpu.next = instApABSx3
	}
}

// instApABSx2 retrieves data from the address specified by the address register and sets the next operation if successful.
func instApABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApABSx3 performs absolute addressing with an additional stack address adjustment and updates the next instruction.
// If the page is crossed, the function ensures proper handling by checking the address read operation for success.
func instApABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instApABSy loads a byte from memory at the current program counter into the address register and advances the program counter.
func instApABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instApABSy1
}

// instApABSy1 reads a byte from memory, updates address registers, and determines the next instruction to execute.
func instApABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApABSy2
	} else {
		cpu.next = instApABSy3
	}
}

// instApABSy2 handles the execution flow for a specific CPU instruction without crossing a memory page.
// If the memory read operation fails, it terminates further execution for this step.
// Updates the CPU's next instruction handler based on the opcode.
func instApABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApABSy3 handles the execution of an operation, performs a page cross check, and sets the next instruction to execute.
func instApABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instApINDx performs the indirect indexed addressing mode operation, updating CPU state and setting the next instruction.
func instApINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDx1
}

// instApINDx1 executes the first stage of the Indexed Indirect (IND,X) addressing mode operation.
// It reads the memory at the address in cpu.ar2 and checks its availability. If unavailable, the CPU halts.
// Updates cpu.ar2 by adding the X register value, masking the result to fit within 8 bits.
// Sets the next instruction handler to instApINDx2.
func instApINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instApINDx2
}

// instApINDx2 reads a value from the address in ar2, sets it to ar if successful, and updates the next handler to instApINDx3.
func instApINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDx3
}

// instApINDx3 performs an indirect indexed addressing operation, updating the address register and setting the next instruction.
func instApINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = _opTable[cpu.op]
}

// instApINDy fetches a byte from memory at the program counter, increments the PC, and updates the AR2 register.
// It sets the next CPU instruction to instApINDy1.
// If memory reading fails, the function exits early.
func instApINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instApINDy1
}

// instApINDy1 reads a value from the address in `ar2`, sets `ar` with the value, and transitions the CPU to `instApINDy2`.
func instApINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instApINDy2
}

// instApINDy2 performs indirect indexed addressing by updating `ar` using `ar2` and `y`, and sets the next instruction.
// The function reads a byte from memory, updates `ar2`, adjusts `ar`, and determines the next handler based on conditions.
func instApINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar2 = uint16(data)
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (cpu.ar2 << 8)
	if z < stackAddr {
		cpu.next = instApINDy3
	} else {
		cpu.next = instApINDy4
	}
}

// instApINDy3 handles indirect indexed addressing with Y-register offset without page crossing for the CPU.
// It reads from the address stored in the AR register, updates the instruction handler, and checks RDY state.
func instApINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = _opTable[cpu.op]
}

// instApINDy4 performs an indirect indexed addressing mode operation with Y register and updates the CPU state.
// If a page boundary is crossed during execution, it ensures proper handling and advances to the next instruction.
func instApINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeABSx fetches an 8-bit immediate value, increments the program counter, and stores it in the address register (ar).
func instAeABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSx1
}

// instAeABSx1 handles the first step of the absolute-indexed addressing mode with the X register adjustment.
// It combines the X register with the address register and determines the next instruction to execute.
// Updates the address register (ar) and increments the program counter (pc).
func instAeABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSx2
	}
}

// instAeABSx2 handles the Absolute Indexed X addressing mode with page crossing, updating the address and next operation.
func instAeABSx2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeABSy initializes the address register with the value read from the program counter and sets the next instruction handler.
func instAeABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instAeABSy1
}

// instAeABSy1 fetches data, updates the address register, and determines the next instruction based on the address range.
func instAeABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeABSy2
	}
}

// instAeABSy2 adjusts the address register by adding a stack offset and sets the next instruction from the operation table.
// It ensures instruction execution respects memory boundary crossing conditions.
func instAeABSy2(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instAeINDy executes the AE INDY instruction, updating the program counter and secondary address register, then sets the next operation.
func instAeINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instAeINDy1
}

// instAeINDy1 reads data from the memory address specified by `ar2` and stores it in `ar`. Sets the next instruction handler.
func instAeINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instAeINDy2
}

// instAeINDy2 executes the second phase of the AE indirect Y-indexed addressing mode.
// It reads a value from memory, combines it with the address register and Y register, and sets the next instruction.
func instAeINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = _opTable[cpu.op]
	} else {
		cpu.next = instAeINDy3
	}
}

// instAeINDy3 handles the addressing mode operation for instructions that cross a page boundary and updates the CPU state.
func instAeINDy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = _opTable[cpu.op]
}

// instMpZER is a memory operation instruction that reads data from memory into the address register (ar) and sets the next operation.
func instMpZER(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instOpRMW
}

// instMpZERx executes a zero-page read operation, increments the program counter, and sets the next instruction.
func instMpZERx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpZERx1
}

// instMpZERx1 performs a memory read operation at the address in `cpu.ar` and modifies the address by adding `cpu.x`.
// If the memory read fails, the function returns immediately.
// The adjusted address is constrained to an 8-bit boundary, and control proceeds to the `instOpRMW` handler.
func instMpZERx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar = (cpu.ar + uint16(cpu.x)) & 0xff
	cpu.next = instOpRMW
}

//func instMpZERy(cpu *CPU) {
//	data, ok := cpu.read(cpu.pc)
//	if !ok {
//		return
//	}
//	cpu.pc++
//	cpu.ar = uint16(data)
//	cpu.next = instMpZERy1
//}

//func instMpZERy1(cpu *CPU) {
//	data, ok := cpu.read(cpu.ar)
//	if !ok {
//		return
//	}
//	cpu.ar = (cpu.ar + uint16(cpu.y)) & 0xff
//	cpu.next = instOpRMW
//}

// instMpABS loads a byte from memory at the program counter into the address register and updates the next instruction.
func instMpABS(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABS1
}

// instMpABS1 represents an instruction handler that modifies the address register (AR) with data fetched from memory.
// It reads a byte from memory located at the program counter (PC), shifts it 8 bits left, and ORs it with the current AR.
// The program counter is incremented, and the next instruction handler is set to a read-modify-write (RMW) operation.
func instMpABS1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

// instMpABSx fetches the next byte from memory, increments the program counter, and sets it in the address register.
// It updates the CPU's next state to instMpABSx1.
func instMpABSx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSx1
}

// instMpABSx1 is a CPU instruction handler for addressing mode manipulation.
// It reads the next byte from memory and uses it as part of the absolute address computation.
// Updates the address register (ar) with the computed address and determines the next instruction handler.
func instMpABSx1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.x)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSx2
	} else {
		cpu.next = instMpABSx3
	}
}

// instMpABSx2 performs an absolute memory read operation and checks for page crossing issues.
// It sets the next instruction to a read-modify-write operation if the memory read is successful.
func instMpABSx2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpABSx3 performs an operation using absolute addressing with additional adjustments and transitions to the next instruction.
func instMpABSx3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instMpABSy sets the address register (ar) based on the byte at the program counter, then sets the next instruction handler.
func instMpABSy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(data)
	cpu.next = instMpABSy1
}

// instMpABSy1 performs memory page addressing mode with Y register offset, updating `ar` and determining the next instruction.
func instMpABSy1(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpABSy2
	} else {
		cpu.next = instMpABSy3
	}
}

// instMpABSy2 handles the zero-page no-cross memory access and assigns the next instruction to a read-modify-write operation.
// It reads from the address register; if unsuccessful, it stops. Otherwise, it sets the next handler to instOpRMW.
func instMpABSy2(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpABSy3 adjusts the address register and sets up the next operation, handling page crossing scenarios.
func instMpABSy3(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instMpINDx loads an operand from the program counter into the ar2 register and advances the program counter.
// Sets the next CPU instruction handler to instMpINDx1. Returns immediately if the read operation fails.
func instMpINDx(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDx1
}

// instMpINDx1 performs indexed indirect addressing mode update. Adjusts ar2 with x, wraps it, and sets the next operation.
func instMpINDx1(cpu *CPU) {
	if _, ok := cpu.read(cpu.ar2); !ok {
		return
	}
	cpu.ar2 = (cpu.ar2 + uint16(cpu.x)) & 0xff
	cpu.next = instMpINDx2
}

// instMpINDx2 reads data from the memory address in ar2, updates ar with this value if successful, and sets the next instruction.
func instMpINDx2(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDx3
}

// instMpINDx3 performs an indexed memory read, updates the address register, and sets the next instruction handler.
func instMpINDx3(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	cpu.ar = cpu.ar | (uint16(data) << 8)
	cpu.next = instOpRMW
}

// instMpINDy reads a byte from memory at the program counter, increments the program counter, and updates ar2 and next state.
func instMpINDy(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar2 = uint16(data)
	cpu.next = instMpINDy1
}

// instMpINDy1 reads data using the CPU's ar2 register, updates the ar register, and sets the next operation to instMpINDy2.
func instMpINDy1(cpu *CPU) {
	data, ok := cpu.read(cpu.ar2)
	if !ok {
		return
	}
	cpu.ar = uint16(data)
	cpu.next = instMpINDy2
}

// instMpINDy2 handles the indirect indexed addressing mode logic by adjusting the address register and setting the next instruction.
func instMpINDy2(cpu *CPU) {
	data, ok := cpu.read((cpu.ar2 + 1) & 0xff)
	if !ok {
		return
	}
	z := cpu.ar + uint16(cpu.y)
	cpu.ar = (z & 0xff) | (uint16(data) << 8)
	if z < stackAddr {
		cpu.next = instMpINDy3
	} else {
		cpu.next = instMpINDy4
	}
}

// instMpINDy3 handles the memory instruction with indirect addressing, updating the CPU's next operation if successful.
func instMpINDy3(cpu *CPU) {
	// No page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.next = instOpRMW
}

// instMpINDy4 handles indirect indexed addressing with Y register offset and handles potential page crossing.
// If a page boundary is crossed, the function adjusts the address register and sets the next operation to instOpRMW.
func instMpINDy4(cpu *CPU) {
	// Page crossed
	if _, ok := cpu.read(cpu.ar); !ok {
		return
	}
	cpu.ar += stackAddr
	cpu.next = instOpRMW
}

// instOpRMW reads data from the address in the CPU's address register, stores it in the `rmw` buffer, and sets the next operation.
func instOpRMW(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.rmw = data
	cpu.next = instOpRMW1
}

// instOpRMW1 executes the second phase of a read-modify-write operation by writing the modified value and updating the next instruction.
func instOpRMW1(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.rmw)
	cpu.next = _opTable[cpu.op]
}

// instOpLDA loads a byte from memory at the address in the address register (AR) into the accumulator (A).
// Updates the negative (N) and zero (Z) flags based on the loaded value.
// Sets the next instruction to `instOpINI` if reading from memory is successful; does nothing otherwise.
func instOpLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDA loads a byte from memory into the accumulator, updates the negative and zero flags, and sets the next instruction.
func instOiLDA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOpLDX loads a value from memory into the X register and updates the negative and zero flags.
func instOpLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDX loads a byte from memory into the X register, updating the negative and zero flags, and sets the next instruction.
func instOiLDX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.x = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

func instOpLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// instOiLDY loads a value into the Y register, updates the negative and zero flags, and sets the next instruction handler.
func instOiLDY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.y = data
	cpu.nFlag = data
	cpu.zFlag = data
	cpu.next = instOpINI
}

// Store

// instOpSTA stores the value in the accumulator (A register) into memory at the address stored in the address register (AR).
// It updates the next CPU instruction handler to instOpINI.
func instOpSTA(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.a)
	cpu.next = instOpINI
}

// instOpSTX stores the value of the X register into memory at the address specified by the AR register and sets the next instruction.
func instOpSTX(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.x)
	cpu.next = instOpINI
}

// instOpSTY stores the value of the Y register into memory at the address in the address register and sets the next instruction.
func instOpSTY(cpu *CPU) {
	cpu.banks.Write(cpu.ar, cpu.y)
	cpu.next = instOpINI
}

// Transfer

// instOpTAX performs the TAX instruction, transferring the value of the accumulator (A) to the X register.
// Updates the negative (nFlag) and zero (zFlag) flags based on the value of A. Sets the next instruction to instOpINI.
func instOpTAX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpTXA transfers the X register to the A register, updating the negative and zero flags based on the value of X.
func instOpTXA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.x
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpTAY transfers the value of the accumulator (A) to the Y register and updates the negative and zero flags.
func instOpTAY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y = cpu.a
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpTYA transfers the value of the Y register to the A register and updates the negative and zero flags accordingly.
func instOpTYA(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.a = cpu.y
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpTSX loads the stack pointer into the X register, updates the negative and zero flags, and sets the next instruction.
func instOpTSX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x = cpu.sp
	cpu.nFlag = cpu.sp
	cpu.zFlag = cpu.sp
	cpu.next = instOpINI
}

// instOpTXS transfers the value from the X register to the stack pointer and sets the next instruction to instOpINI.
func instOpTXS(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.sp = cpu.x
	cpu.next = instOpINI
}

// Arithmetic

// instOpADC performs the ADC (Add with Carry) operation by reading data from memory and calling the ADC handler.
func instOpADC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doADC(data)
	cpu.next = instOpINI
}

// instOiADC executes the ADC instruction at the current program counter, updating the accumulator and advancing the PC.
// Invokes the next instruction handler after performing the operation.
func instOiADC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doADC(data)
	cpu.next = instOpINI
}

// instOpSBC executes the SBC (Subtract with Carry) operation using data from the address register and updates the next instruction.
func instOpSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.doSBC(data)
	cpu.next = instOpINI
}

// instOiSBC handles the SBC (Subtract with Carry) instruction. It reads data at the program counter, increments the PC,
// performs the SBC operation using the read data, and sets the next instruction to instOpINI.
func instOiSBC(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.doSBC(data)
	cpu.next = instOpINI
}

// Increment, decrement

// instOpINX increments the X register, updates the negative and zero flags, and sets the next instruction to instOpINI.
func instOpINX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x++
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpDEX decrements the X register, setting the negative and zero flags, and sets the next instruction to instOpINI.
func instOpDEX(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.x--
	cpu.nFlag = cpu.x
	cpu.zFlag = cpu.x
	cpu.next = instOpINI
}

// instOpINY increments the Y register, updates the negative and zero flags, and sets the next instruction to instOpINI.
func instOpINY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y++
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpDEY decrements the Y register, updating the negative and zero flags based on the resulting value.
func instOpDEY(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.y--
	cpu.nFlag = cpu.y
	cpu.zFlag = cpu.y
	cpu.next = instOpINI
}

// instOpINC increments a value, updates the N and Z flags, writes the result to memory, and sets the next operation.
func instOpINC(cpu *CPU) {
	v := cpu.rmw + 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOpDEC performs a decrement operation on the value in the RMW register, updating CPU flags and writing the result.
func instOpDEC(cpu *CPU) {
	v := cpu.rmw - 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOpAND performs a bitwise AND operation between the accumulator and memory at the address in the address register.
// Updates the negative and zero flags based on the result. Sets the next instruction handler.
func instOpAND(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiAND performs the AND operation between the accumulator and fetched data. Updates N and Z flags accordingly.
func instOiAND(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a &= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpORA performs the ORA (logical OR with accumulator) operation, updates the negative/zero flags, and sets the next instruction.
func instOpORA(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiOPA is an instruction that performs a bitwise OR operation between the accumulator and a fetched operand.
// Updates negative and zero flags based on the result. Sets the next instruction to instOpINI.
func instOiOPA(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a |= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpEOR performs the EOR (Exclusive OR) operation on the accumulator with a value from memory and updates CPU flags.
func instOpEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOiEOR executes the Exclusive OR (EOR) operation on the accumulator with a fetched memory value and updates flags.
func instOiEOR(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.a ^= data
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpCMP performs a comparison between the accumulator and a memory value, updating CPU flags based on the result.
func instOpCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCMP performs a comparison between the CPU's accumulator and the operand, updating CPU flags accordingly.
func instOiCMP(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.a) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpCPX executes the CPX (Compare with X Register) instruction, updating the CPU's flags based on the result.
func instOpCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCPX executes the OiCPX instruction by reading a value from memory and performing a subtraction with register X.
// It updates the CPU's flags (negative, zero, carry) and the address register (AR) as per the operation's result.
// Finally, it sets the next instruction handler to instOpINI.
func instOiCPX(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.x) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOpCPY handles the CPY (Compare Y Register) instruction, updating flags based on comparison with memory data.
func instOpCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// instOiCPY handles the CPY (Compare Y Register with Memory) instruction in immediate addressing mode.
// It updates the negative, zero, and carry flags based on the comparison result.
// The next instruction to execute is set to instOpINI.
func instOiCPY(cpu *CPU) {
	data, ok := cpu.read(cpu.pc)
	if !ok {
		return
	}
	cpu.pc++
	cpu.ar = uint16(cpu.y) - uint16(data)
	cpu.nFlag = uint8(cpu.ar)
	cpu.zFlag = uint8(cpu.ar)
	cpu.cFlag = conversion.BoolToUint8(cpu.ar < stackAddr)
	cpu.next = instOpINI
}

// Bit-test

// instOpBIT performs the BIT instruction, updating the CPU flags based on a bitwise AND between the accumulator and memory data.
func instOpBIT(cpu *CPU) {
	data, ok := cpu.read(cpu.ar)
	if !ok {
		return
	}
	cpu.zFlag = cpu.a & data
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.next = instOpINI
}

// instOpASL performs the ASL (Arithmetic Shift Left) operation on the CPU's RMW buffer and updates relevant CPU flags.
// It shifts the RMW buffer left by one bit, sets the carry flag to the original top bit, and updates zero/negative flags.
// The result is written back to memory at the address pointed to by the address register (ar).
// Finally, sets the next CPU instruction to instOpINI for subsequent execution.
func instOpASL(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x80
	v := cpu.rmw << 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOaASL executes the ASL (Arithmetic Shift Left) operation on the accumulator, updating flags and the next instruction.
func instOaASL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x80
	cpu.a <<= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpLSR performs the Logical Shift Right (LSR) operation on the CPU's RMW register, updating flags and writing the result.
func instOpLSR(cpu *CPU) {
	cpu.cFlag = cpu.rmw & 0x1
	v := cpu.rmw >> 1
	cpu.nFlag = v
	cpu.zFlag = v
	cpu.banks.Write(cpu.ar, v)
	cpu.next = instOpINI
}

// instOaLSR performs a logical shift right (LSR) on the A register, updating the carry, zero, and negative flags accordingly.
func instOaLSR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = cpu.a & 0x1
	cpu.a >>= 1
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.next = instOpINI
}

// instOpROL performs the ROL (Rotate Left) operation on the `rmw` register, updating CPU flags and memory state.
func instOpROL(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw << 1) | 0x1
	} else {
		t = cpu.rmw << 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x80
	cpu.next = instOpINI
}

// instOaROL performs a rotate left operation on the accumulator and updates the CPU flags accordingly.
func instOaROL(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x80
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a << 1) | 0x1
	} else {
		cpu.a = cpu.a << 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = instOpINI
}

// instOpROR performs a Rotate Right operation on the CPU's RMW register, updates flags, and writes the result to memory.
func instOpROR(cpu *CPU) {
	var t uint8
	if cpu.cFlag != 0 {
		t = (cpu.rmw >> 1) | 0x80
	} else {
		t = cpu.rmw >> 1
	}
	cpu.nFlag = t
	cpu.zFlag = t
	cpu.banks.Write(cpu.ar, t)
	cpu.cFlag = cpu.rmw & 0x1
	cpu.next = instOpINI
}

// instOaROR performs the ROR (Rotate Right) operation on the accumulator,
// updating the negative, zero, and carry flags and setting the next instruction.
func instOaROR(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	data := cpu.a & 0x1
	if cpu.cFlag != 0 {
		cpu.a = (cpu.a >> 1) | 0x80
	} else {
		cpu.a = cpu.a >> 1
	}
	cpu.nFlag = cpu.a
	cpu.zFlag = cpu.a
	cpu.cFlag = data
	cpu.next = instOpINI
}

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
		cpu.opFlags |= opFlagIrqDisabled
	} else if iFlagPrev != 0 && cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqEnabled
	}
	cpu.next = instOpINI
}

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
	cpu.opFlags |= opFlagIntDelayed
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

// Flag

// instOpSEC sets the Carry flag (cFlag) to 1 and moves execution to the next instruction handler (instOpINI).
func instOpSEC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 1
	cpu.next = instOpINI
}

// instOpCLC clears the carry flag in the CPU and sets the next instruction to instOpINI. It halts if the current PC read fails.
func instOpCLC(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.cFlag = 0
	cpu.next = instOpINI
}

// instOpSED sets the decimal mode flag (dFlag) to 1 and assigns the next instruction handler to instOpINI.
func instOpSED(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 1
	cpu.next = instOpINI
}

// instOpCLD clears the decimal mode flag (dFlag) and sets the next instruction handler to instOpINI.
func instOpCLD(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.dFlag = 0
	cpu.next = instOpINI
}

// instOpSEI sets the interrupt disable flag and updates CPU state to handle the next instruction cycle.
func instOpSEI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqDisabled
	}
	cpu.iFlag = 1
	cpu.next = instOpINI
}

// instOpCLI clears the interrupt disable flag and sets the next instruction to instOpINI if the current opcode is valid.
func instOpCLI(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	if cpu.iFlag == 0 {
		cpu.opFlags |= opFlagIrqEnabled
	}
	cpu.iFlag = 0
	cpu.next = instOpINI
}

// instOpCLV clears the overflow flag in the CPU and sets the next instruction handler to instOpINI.
// If the current PC address cannot be read, the operation is aborted.
func instOpCLV(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.vFlag = 0
	cpu.next = instOpINI
}

// instOpNOP is a no-operation function for the CPU that progresses the state to the next instruction without modifying registers.
func instOpNOP(cpu *CPU) {
	if _, ok := cpu.read(cpu.pc); !ok {
		return
	}
	cpu.next = instOpINI
}

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
