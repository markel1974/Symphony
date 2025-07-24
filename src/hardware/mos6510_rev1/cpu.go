package mos6510_rev1

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// IControlUnit defines the interface for managing and configuring operations and modes within a CPU control unit.
// Setup initializes the control unit and prepares it for operation.
// CreateOpTable generates the operational function table for the control unit.
// CreateModeTable generates the mode-handling function table for the control unit.
// GetOpFn retrieves the operation function associated with a given string identifier.
// GetOpId returns the string identifier associated with a given operation function.
// GetInstOpINI provides the operation function for the "INI" instruction.
// GetInstOpHalt provides the operation function for the "Halt" instruction.
// GetInstOpNMI provides the operation function for the "NMI" (Non-Maskable Interrupt) instruction.
// GetInstOpIRQ provides the operation function for the "IRQ" (Interrupt Request) instruction.
// GetInstOpBRAbp provides the operation function for a branch instruction on breakpoint.
// GetInstOpBRAfp provides the operation function for a branch instruction on function pointer.
// GetInstOpBRAnp provides the operation function for a branch instruction on no pointer.
type IControlUnit interface {
	Setup() error

	CreateOpTable() []func(cpu *CPU)

	CreateModeTable() []func(cpu *CPU)

	GetOpFn(v string) (func(cpu *CPU), bool)

	GetOpId(func(cpu *CPU)) (string, bool)

	GetInstOpINI() func(cpu *CPU)

	GetInstOpHalt() func(cpu *CPU)

	GetInstOpNMI() func(cpu *CPU)

	GetInstOpIRQ() func(cpu *CPU)

	GetInstOpBRAbp() func(cpu *CPU)

	GetInstOpBRAfp() func(cpu *CPU)

	GetInstOpBRAnp() func(cpu *CPU)
}

//https://web.archive.org/web/20221112220344if_/http://archive.6502.org/datasheets/synertek_programming_manual.pdf
//https://dustlayer.com/c64-architecture
//https://www.zimmers.net/cbmpics/cbm/c64/vic-ii.txt

//Notes
//https://codebase64.org/lib/exe/fetch.php?media=base:safely_freezing_the_c64.pdf

// CPU represents a simulated central processing unit with registers, flags, and associated helper components.
type CPU struct {
	*component.BaseComponent
	next           func(cpu *CPU) // next is a function pointer that executes the next CPU instruction or operation during emulation.
	overflowBranch func() bool    // overflowBranch determines if the CPU should branch based on the overflow condition.
	nFlag          uint8          // Negative flag - Only the highest bit of the nFlag variable is used
	zFlag          uint8          // Zero flag - The zFlag variable has the inverse meaning of the 6510 Z flag
	vFlag          uint8          // Overflow flag
	dFlag          uint8          // Decimal mode flag
	iFlag          uint8          // Interrupt disable flag
	cFlag          uint8          // Carry flag
	a              uint8          // Register
	x              uint8          // Register
	y              uint8          // Register
	sp             uint8          // Stack pointer
	pc             uint16         // Program counter
	op             uint8          // Current opcode
	ar             uint16         // Address register
	ar2            uint16         // Address register 2
	rmw            uint8          // Data buffer for RMW instructions
	rdyLow         bool           // current RDY state
	aecLow         bool           // current AEC state
	opFlags        uint8          // opFlags is a uint8 value used to store operational flags for the CPU's current instruction state.
	savedNext      func(cpu *CPU)
	modeTable      []func(*CPU)
	opTable        []func(*CPU)
	interrupts     *Interrupts
	bus            *Bus
	control        IControlUnit
	label          string

	instOpINI   func(cpu *CPU)
	instOpHalt  func(cpu *CPU)
	instOpNMI   func(cpu *CPU)
	instOpIRQ   func(cpu *CPU)
	instOpBRAbp func(cpu *CPU)
	instOpBRAfp func(cpu *CPU)
	instOpBRAnp func(cpu *CPU)
}

// NewCPU initializes and returns a new CPU instance with the provided id.
func NewCPU(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CPU {
	cpu := &CPU{
		BaseComponent: component.NewBaseComponent(),
		label:         label,
	}
	cpu.BaseComponent.Register(factory, parent, Identifier(), cpu, references.IdIMos6510(cpu, label, instance))
	return cpu
}

// Setup initializes the CPU by configuring its PIC and banks from the given socket.
func (cpu *CPU) Setup() error {
	cpu.interrupts = NewInterrupts(cpu, cpu.GetFactory(), cpu.label, 0)
	cpu.bus = NewBus(cpu, cpu.GetFactory(), cpu.label, 0)
	cpu.control = NewControlUnit(cpu, cpu.GetFactory(), cpu.label, 0)
	if err := cpu.interrupts.Setup(); err != nil {
		return err
	}
	if err := cpu.bus.Setup(); err != nil {
		return err
	}
	if err := cpu.control.Setup(); err != nil {
		return err
	}
	cpu.opTable = cpu.control.CreateOpTable()
	cpu.modeTable = cpu.control.CreateModeTable()

	cpu.instOpINI = cpu.control.GetInstOpINI()
	cpu.instOpHalt = cpu.control.GetInstOpHalt()
	cpu.instOpNMI = cpu.control.GetInstOpNMI()
	cpu.instOpIRQ = cpu.control.GetInstOpIRQ()
	cpu.instOpBRAbp = cpu.control.GetInstOpBRAbp()
	cpu.instOpBRAfp = cpu.control.GetInstOpBRAfp()
	cpu.instOpBRAnp = cpu.control.GetInstOpBRAnp()
	return nil
}

// Bind initializes the connections between the CPU and its socket, PIC, and memory banks, configuring necessary handlers.
func (cpu *CPU) Bind(_ references.IMos6510Socket, q references.IQuartz, banks references.IMos6510Banks) error {
	if err := cpu.interrupts.Bind(q, cpu.Reset, cpu.NMI, cpu.IRQ); err != nil {
		return err
	}
	if err := cpu.bus.Bind(cpu.setModeHalt, banks); err != nil {
		return err
	}
	return nil
}

// Connect establishes necessary connections for the CPU, preparing it for operation.
func (cpu *CPU) Connect() error {
	return nil
}

// Internal determines whether the CPU is in an internal state, returning false as the default implementation.
func (cpu *CPU) Internal() bool {
	return false
}

// Reset initializes or restores the CPU to a default state by resetting internal flags, registers, and setting the program counter.
func (cpu *CPU) Reset() {
	cpu.interrupts.Reset()
	cpu.bus.Reset()
	cpu.pc = uint16(cpu.bus.ReadDirect(0xfffc)) | (uint16(cpu.bus.ReadDirect(0xfffd)) << 8) // Read reset vector
	cpu.opFlags = 0
	cpu.next = cpu.instOpINI
}

// SetOverflowBranch sets the function used to handle conditional overflow branching for the CPU.
func (cpu *CPU) SetOverflowBranch(sob func() bool) {
	cpu.overflowBranch = sob
}

// TriggerReset initiates a reset operation on the CPU by activating the corresponding reset mechanism in the PIC.
func (cpu *CPU) TriggerReset() {
	cpu.interrupts.TriggerReset()
}

// ClearIRQ clears the interrupt request (IRQ) for the given IRQ vector by invoking the associated PIC method.
func (cpu *CPU) ClearIRQ(v uint32) {
	cpu.interrupts.ClearIRQ(v)
}

// TriggerIRQ triggers an interrupt request (IRQ) by delegating the request to the programmable interrupt controller (PIC).
func (cpu *CPU) TriggerIRQ(v uint32) {
	cpu.interrupts.TriggerIRQ(v)
}

// TriggerNMI forces a Non-Maskable Interrupt (NMI) to occur on the CPU by signaling the programmable interrupt controller.
func (cpu *CPU) TriggerNMI() {
	cpu.interrupts.TriggerNMI()
}

// ClearNMI clears the Non-Maskable Interrupt (NMI) state by delegating the operation to the programmable interrupt controller.
func (cpu *CPU) ClearNMI() {
	cpu.interrupts.ClearNMI()
}

// SetAECLow sets the AEC line state to low or high based on the provided boolean value and updates the CPU stop state accordingly.
func (cpu *CPU) SetAECLow(aecLow bool) {
	cpu.aecLow = aecLow
	if cpu.aecLow {
		cpu.bus.SetAECLowMode()
	}
}

// SetRDYLow sets the RDY line state to low or high based on the provided boolean value and updates the CPU stop state accordingly.
func (cpu *CPU) SetRDYLow(rdyLow bool) {
	cpu.rdyLow = rdyLow
	if cpu.rdyLow {
		cpu.bus.SetRDYLowMode()
	} else {
		cpu.bus.SetNormalMode()
		if cpu.savedNext != nil {
			cpu.next = cpu.savedNext
			cpu.savedNext = nil
		}
	}
}

// setModeHalt transitions the CPU into halt mode by saving the current state and setting the next operation to halt.
func (cpu *CPU) setModeHalt() {
	if cpu.savedNext == nil {
		cpu.savedNext = cpu.next
		cpu.next = cpu.instOpHalt
	}
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

//go:nosplit
func (cpu *CPU) NMI() {
	cpu.next = cpu.instOpNMI
	cpu.next(cpu)
}

//go:nosplit
func (cpu *CPU) IRQ() {
	cpu.next = cpu.instOpIRQ
	cpu.next(cpu)
}

// popFlags updates the CPU state flags based on the input data, setting various flags like nFlag, vFlag, dFlag, etc.
func (cpu *CPU) popFlags(data uint8) {
	cpu.nFlag = data
	cpu.vFlag = data & 0x40
	cpu.dFlag = data & 0x08
	cpu.iFlag = data & 0x04
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
func (cpu *CPU) computeBranch(data uint8) func(*CPU) {
	cpu.ar = cpu.pc + uint16(int8(data))
	if (cpu.ar >> 8) != (cpu.pc >> 8) {
		if (data & 0x80) != 0 {
			return cpu.instOpBRAbp
		}
		return cpu.instOpBRAfp
	}
	return cpu.instOpBRAnp
}

// doADC performs the add with Carry (ADC) operation using the given operand and CPU state.
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
		p1 := (uint16(cpu.a) ^ uint16(data)) & 0x80
		p2 := (uint16(cpu.a) ^ tmp) & 0x80
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
	if (p1 != 0) && (p2 == 0) {
		cpu.vFlag = 1
	} else {
		cpu.vFlag = 0
	}
	if ah > 9 {
		ah += 6
	}
	// BCD fixup for upper nybble
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
		if tmp < 0x100 {
			cpu.cFlag = 1
		} else {
			cpu.cFlag = 0
		}
		p1 := (uint16(cpu.a) ^ tmp) & 0x80
		p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
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
	if uint16(tmp) < 0x100 {
		cpu.cFlag = 1
	} else {
		cpu.cFlag = 0
	}
	p1 := (uint16(cpu.a) ^ tmp) & 0x80
	p2 := (uint16(cpu.a) ^ uint16(data)) & 0x80
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
