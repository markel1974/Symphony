package mos6510_rev1

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

//https://web.archive.org/web/20221112220344if_/http://archive.6502.org/datasheets/synertek_programming_manual.pdf
//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf

// CPU represents a simulated central processing unit with registers, flags, and associated helper components.
type CPU struct {
	*component.BaseComponent
	lineRead  func(uint16) uint8
	lineWrite func(uint16, uint8)
	busRead   func(uint16) uint8  // busRead is a function that reads a byte from a specified 16-bit memory address in the CPU's memory bank.
	busWrite  func(uint16, uint8) // busWrite is a function that writes a byte to a specified 16-bit memory address in the CPU's memory bank.

	picReset       func()                                 // picReset is a function that resets or reinitializes the Programmable Interrupt Controller (PIC).
	picHasNMI      func() bool                            // picHasNMI checks if the PIC (Programmable Interrupt Controller) has a Non-Maskable Interrupt (NMI) pending.
	picClearNMI    func()                                 // picClearNMI clears the Non-Maskable Interrupt (NMI) signal in the system's Programmable Interrupt Controller (PIC).
	picVerifyIrq   func(iFlag uint8, opFlags uint8) uint8 // picVerifyIrq verifies interrupt conditions by comparing the iFlag and opFlags and returns the updated interrupt state.
	next           func(cpu *CPU)                         // next is a function pointer that executes the next CPU instruction or operation during emulation.
	overflowBranch func() bool                            // overflowBranch determines if the CPU should branch based on the overflow condition.
	nFlag          uint8                                  // Negative flag - Only the highest bit of the nFlag variable is used
	zFlag          uint8                                  // Zero flag - The zFlag variable has the inverse meaning of the 6510 Z flag
	vFlag          uint8                                  // Overflow flag
	dFlag          uint8                                  // Decimal mode flag
	iFlag          uint8                                  // Interrupt disable flag
	cFlag          uint8                                  // Carry flag
	a              uint8                                  // Register
	x              uint8                                  // Register
	y              uint8                                  // Register
	sp             uint8                                  // Stack pointer
	pc             uint16                                 // Program counter
	op             uint8                                  // Current opcode
	ar             uint16                                 // Address register
	ar2            uint16                                 // Address register 2
	rmw            uint8                                  // Data buffer for RMW instructions
	rdyLow         bool                                   // current RDY state
	aecLow         bool                                   // current AEC state
	opFlags        uint8                                  // opFlags is a uint8 value used to store operational flags for the CPU's current instruction state.
	irqBreaker     bool                                   // irqBreaker indicates whether the CPU's interrupt request is currently blocked.
	savedNext      func(cpu *CPU)
	modeTable      []func(*CPU)
	opTable        []func(*CPU)
}

// NewCPU initializes and returns a new CPU instance with the provided id.
func NewCPU(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CPU {
	cpu := &CPU{
		BaseComponent: component.NewBaseComponent(),
	}
	cpu.BaseComponent.Register(factory, parent, Identifier(), cpu, references.IdIMos6510(cpu, label, instance))
	return cpu
}

// Setup initializes the CPU by configuring its PIC and banks from the given socket.
func (cpu *CPU) Setup() error {
	cpu.opTable = CreateOpTable()
	cpu.modeTable = CreateModeTable()
	return nil
}

func (cpu *CPU) Bind(_ references.IMos6510Socket, pic references.IMos6510Pic, banks references.IMos6510Banks) error {
	cpu.picReset = pic.Reset
	cpu.picVerifyIrq = pic.VerifyIrq
	cpu.picHasNMI = pic.HasNMI
	cpu.picClearNMI = pic.ClearNMI

	cpu.lineRead = banks.Read
	cpu.lineWrite = banks.Write

	cpu.busRead = cpu.lineRead
	cpu.busWrite = cpu.lineWrite
	return nil
}

// Connect establishes necessary connections for the CPU, preparing it for operation.
func (cpu *CPU) Connect() error {
	return nil
}

func (cpu *CPU) Internal() bool {
	return false
}

// Reset initializes or restores the CPU to a default state by resetting internal flags, registers, and setting the program counter.
func (cpu *CPU) Reset() {
	cpu.picReset()
	cpu.pc = uint16(cpu.busRead(0xfffc)) | (uint16(cpu.busRead(0xfffd)) << 8) // Read reset vector
	cpu.next = InstOpINI
	cpu.opFlags = 0
	cpu.irqBreaker = false
}

// SetOverflowBranch sets the function used to handle conditional overflow branching for the CPU.
func (cpu *CPU) SetOverflowBranch(sob func() bool) {
	cpu.overflowBranch = sob
}

// SetAECLow sets the AEC line state to low or high based on the provided boolean value and updates the CPU stop state accordingly.
func (cpu *CPU) SetAECLow(aecLow bool) {
	cpu.aecLow = aecLow
	if cpu.aecLow {
		cpu.disconnectBus()
	}
}

// SetRDYLow sets the RDY line to the provided state and resumes the CPU if RDY is not low.
func (cpu *CPU) SetRDYLow(rdyLow bool) {
	cpu.rdyLow = rdyLow
	if !cpu.rdyLow {
		cpu.setModeRun()
	}
}

// setModeHalt transitions the CPU into halt mode by saving the current state and setting the next operation to halt.
func (cpu *CPU) setModeHalt() {
	if cpu.savedNext == nil {
		cpu.savedNext = cpu.next
		cpu.next = halt
	}
}

// setModeRun transitions the CPU into run mode by connecting to the bus and restoring any previously saved state.
func (cpu *CPU) setModeRun() {
	cpu.connectBus()
	if cpu.savedNext != nil {
		cpu.next = cpu.savedNext
		cpu.savedNext = nil
	}
}

// disconnectBus disconnects the CPU from its current bus by setting read and write operations to their disconnected state.
func (cpu *CPU) disconnectBus() {
	cpu.busRead = cpu.busDisconnectedRead
	cpu.busWrite = cpu.busDisconnectedWrite
}

// connectBus initializes the CPU's internal bus by linking busRead and busWrite to corresponding line signals.
func (cpu *CPU) connectBus() {
	cpu.busRead = cpu.lineRead
	cpu.busWrite = cpu.lineWrite
}

// Emulate processes one CPU cycle by invoking the next instruction handler unless the CPU is stopped.
//
//go:nosplit
func (cpu *CPU) Emulate() {
	cpu.next(cpu)
}

// EmulationRequired determines if the CPU requires emulation for the current operation, returning true if necessary.
func (cpu *CPU) EmulationRequired() bool {
	return true
}

// busDisconnectedRead performs a read operation while the CPU bus is disconnected, always returning a default value of 0.
func (cpu *CPU) busDisconnectedRead(_ uint16) uint8 {
	return 0
}

// busDisconnectedWrite is a placeholder method called when an illegal write operation is attempted on a disconnected bus.
func (cpu *CPU) busDisconnectedWrite(_ uint16, _ uint8) {
}

// halt pauses the CPU's operation by acting as a no-op function while the CPU remains in the halted state.
func halt(_ *CPU) {
}

// Read retrieves a byte from the specified memory address.
//
// Returns the byte read and a boolean indicating success.
//
// If the RDY line is low (rdyLow == true), indicating that the VIC-II is currently
// accessing memory, the CPU pauses execution by setting the internal 'stop' flag,
// and the function returns 0, false. This simulates the behavior of the 6510's RDY line,
// which is used by the VIC-II during "bad-lines".  The 'stop' flag is specific
// to this emulator and is NOT part of the real 6510 hardware.
//
// If the RDY line is high (rdyLow == false), the function reads a byte from memory
// using the cpu.banks.Read method and returns the byte and true.
func (cpu *CPU) read(addr uint16) (uint8, bool) {
	if cpu.rdyLow {
		cpu.setModeHalt()
		return 0, false
	}
	return cpu.busRead(addr), true
}

// popFlags updates the CPU state flags based on the input data, setting various flags like nFlag, vFlag, dFlag, etc.
func (cpu *CPU) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
	//cpu.zFlag = conversion.BoolToUint8((data & 0x02) == 0)
	if (data & 0x02) == 0 {
		cpu.zFlag = 1
	} else {
		cpu.zFlag = 0
	}
	cpu.cFlag = data & 0x01
}

// pushFlags generates a status register byte based on the CPU's flag states and the provided B flag parameter.
func (cpu *CPU) pushFlags(bFlags bool) uint8 {
	data := 0x20 | (cpu.nFlag & 0x80)
	if cpu.vFlag != 0 {
		data |= 0x40
	}
	if bFlags {
		data |= 0x10
	}
	if cpu.dFlag != 0 {
		data |= 0x08
	}
	if cpu.iFlag != 0 {
		data |= 0x04
	}
	if cpu.zFlag == 0 {
		data |= 0x02
	}
	if cpu.cFlag != 0 {
		data |= 0x01
	}
	return data
}

// branch handles the logic for branching to a new program counter location based on offset and page crossing.
func (cpu *CPU) branch(data uint8) {
	cpu.ar = cpu.pc + uint16(int8(data))
	if (cpu.ar >> 8) != (cpu.pc >> 8) {
		if (data & 0x80) != 0 {
			cpu.next = InstOpBRAbp
		} else {
			cpu.next = InstOpBRAfp
		}
	} else {
		cpu.next = InstOpBRAnp
	}
}

// doADC performs the Add with Carry (ADC) operation using the given operand and CPU state.
// Supports both binary and decimal modes. Updates processor flags: carry, zero, negative, and overflow.
func (cpu *CPU) doADC(data uint8) {
	k := uint8(0)
	if cpu.cFlag != 0 {
		k = 1
	}
	if cpu.dFlag == 0 {
		// Binary mode
		tmp := uint16(cpu.a) + uint16(data) + uint16(k)
		if tmp > 0xff {
			cpu.cFlag = 1
		} else {
			cpu.cFlag = 0
		}
		//cpu.cFlag = conversion.BoolToUint8(tmp > 0xff)
		p1 := (uint16(cpu.a) ^ uint16(data)) & 0x80
		p2 := (uint16(cpu.a) ^ tmp) & 0x80
		//cpu.vFlag = conversion.BoolToUint8((p1 == 0) && (p2 != 0))
		if (p1 == 0) && (p2 != 0) {
			cpu.vFlag = 1
		} else {
			cpu.vFlag = 0
		}
		cpu.a = uint8(tmp)
		cpu.nFlag = uint8(tmp)
		cpu.zFlag = uint8(tmp)
		return
	}
	// Decimal mode
	al := (cpu.a & 0x0f) + (data & 0x0f) + k
	if al > 9 {
		al += 6 // BCD fixup
	}
	ah := (cpu.a >> 4) + (data >> 4)
	if al > 0x0f {
		ah++
	}
	cpu.zFlag = cpu.a + data + k
	cpu.nFlag = ah << 4 // Only the highest bit used
	p1 := ((ah << 4) ^ cpu.a) & 0x80
	p2 := (cpu.a ^ data) & 0x80
	//cpu.vFlag = conversion.BoolToUint8((p1 != 0) && (p2 == 0))
	if (p1 != 0) && (p2 == 0) {
		cpu.vFlag = 1
	} else {
		cpu.vFlag = 0
	}
	if ah > 9 {
		ah += 6
	}
	// BCD fixup for upper nybble
	//cpu.cFlag = conversion.BoolToUint8(ah > 0x0f) // carry flag
	if ah > 0x0f {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	cpu.a = (ah << 4) | (al & 0x0f) // result
}

// doSBC performs the SBC (Subtract with Carry) instruction on the CPU, handling both binary and BCD (decimal) modes.
// It calculates the result of the subtraction of `data` and the current accumulator with the carry-in value.
// The method updates the accumulator and relevant CPU flags: negative, zero, overflow, and carry.
func (cpu *CPU) doSBC(data uint8) {
	k := uint8(0)
	if cpu.cFlag == 0 {
		k = 1
	}
	tmp := uint16(cpu.a) - uint16(data) - uint16(k)
	if cpu.dFlag == 0 {
		// Binary mode
		//cpu.cFlag = conversion.BoolToUint8(tmp < 0x100)
		if tmp < 0x100 {
			cpu.cFlag = 1
		} else {
			cpu.cFlag = 0
		}
		p1 := (uint16(cpu.a) ^ tmp) & 0x80
		p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
		//cpu.vFlag = conversion.BoolToUint8((p1 != 0) && (p2 != 0))
		if (p1 != 0) && (p2 != 0) {
			cpu.vFlag = 1
		} else {
			cpu.vFlag = 0
		}
		cpu.a = uint8(tmp)
		cpu.nFlag = uint8(tmp)
		cpu.zFlag = uint8(tmp)
		return
	}
	// Decimal mode
	al := (cpu.a & 0x0f) - (data & 0x0f) - k
	ah := (cpu.a >> 4) - (data >> 4)
	if (al & 0x10) != 0 {
		al -= 6 // BCD fixup
		ah--
	}
	if (ah & 0x10) != 0 {
		ah -= 6 // BCD fixup
	}
	//cpu.cFlag = conversion.BoolToUint8(uint16(tmp) < 0x100)
	if uint16(tmp) < 0x100 {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	p1 := (uint16(cpu.a) ^ tmp) & 0x80
	p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
	//cpu.vFlag = conversion.BoolToUint8((p1 != 0) && (p2 != 0))
	if (p1 != 0) && (p2 != 0) {
		cpu.vFlag = 1
	} else {
		cpu.vFlag = 0
	}
	cpu.zFlag = uint8(tmp)
	cpu.nFlag = uint8(tmp)
	cpu.a = (ah << 4) | (al & 0x0f)
}

// printRegisters logs the state of the CPU registers and flags. It outputs a formatted string with relevant values.
func (cpu *CPU) printRegisters(qCycle uint64, baLow bool) {
	bLow := 0
	if baLow {
		bLow = 1
	}
	fmt.Printf("CPU] %d|%d||%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|%d\n",
		qCycle,
		bLow,
		cpu.nFlag,
		cpu.zFlag,
		cpu.vFlag,
		cpu.dFlag,
		cpu.iFlag,
		cpu.cFlag,
		cpu.a,
		cpu.x,
		cpu.y,
		cpu.sp,
		cpu.pc,
		cpu.op,
		cpu.ar,
		cpu.ar2,
		cpu.rmw)
}
