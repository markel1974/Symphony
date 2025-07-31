package mos6526

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	RegisterSize  = 0xf
	RegisterCount = RegisterSize + 1
)

const (
	dividerPAL  = 20000 // 50 Hz (PAL) 1,000,000 / 50 = 20,000
	dividerNTSC = 16667 // 60 Hz (NTSC) 1,000,000 / 60 = 16,667
)

// IRQUnderflowTimerA represents the IRQ flag for underflow on Timer A.
// IRQUnderflowTimerB represents the IRQ flag for underflow on Timer B.
// IRQTODAlarmEqual represents the IRQ flag for Time-of-Day alarm match.
// IRQSDRFullOrEmpty represents the IRQ flag for SDR full or empty condition.
// IRQFlagPin represents the IRQ flag for the external flag pin.
// IRQOccurred represents the IRQ flag indicating an interrupt has occurred.
const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOrEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

//https://web.archive.org/web/20181126000922if_/http://archive.6502.org/datasheets/mos_6526_cia_recreated.pdf
//https://emudev.de/q00-c64/cias-timers-keyboard-and-more/

// CIA represents the Complex Interface Adapter, a chip used for I/O operations and timers.
type CIA struct {
	*component.BaseComponent
	reflect               *CIAReflect
	tod                   *TOD
	timerA                *Timer
	timerB                *Timer
	shiftRegister         *ShiftRegister
	reads                 [RegisterCount]func() uint8
	writes                [RegisterCount]func(uint8)
	prA                   uint8 // symphony:export prA represents the current state of the port A register for the CIA component.
	prB                   uint8 // symphony:export prB represents the current state of the port B register for the CIA component.
	ddrA                  uint8 // symphony:export ddrA represents the Data Direction Register for port A, used to configure input/output direction of each bit.
	ddrB                  uint8 // symphony:export ddrB represents the Data Direction Register for Port B, used to configure input/output for the CIA port B pins.
	sdr                   uint8 // symphony:export sdr holds the contents of the Serial Data Register (SDR), used for serial communication operations in the CIA chip.
	icr                   uint8 // symphony:export icr holds the Interrupt Control Register value, managing interrupt flags and statuses for the CIA component.
	irqMask               uint8 // symphony:export irqMask defines the interrupt control mask used to enable or disable specific interrupt sources within the CIA.
	timerAIrqCycle        bool  // symphony:export timerAIrqCycle indicates whether an interrupt request should be triggered in the next cycle for Timer A.
	timerBIrqCycle        bool  // symphony:export timerBIrqCycle indicates if Timer B's interrupt request will be triggered in the next emulation cycle.
	todClockDivider       int   // symphony:export todClockDivider determines the number of cycles before updating the internal Time of Day (TOD) clock.
	label                 string
	socketReadPortA       func(uint8, uint8, uint8, uint8) uint8
	socketReadPortB       func(uint8, uint8, uint8, uint8) uint8
	socketSignalPRA       func(uint8)
	socketSignalPRB       func(uint8)
	socketSignalDDRA      func(uint8)
	socketSignalDDRB      func(uint8)
	socketIRQTrigger      func()
	socketIRQClearTrigger func()
}

// NewCIA creates and initializes a new instance of CIA, registering it with the provided factory, parent, and instance ID.
func NewCIA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CIA {
	m := &CIA{
		BaseComponent: component.NewBaseComponent(),
		tod:           nil,
		timerA:        nil,
		timerB:        nil,
		label:         label,
	}
	//m.BaseComponent.Register(factory, parent, Identifier(), instance, m, references.IdIMos6526(m, label, instance))
	m.reflect = NewCIAReflect(m, factory, parent, Identifier(), instance, references.IdIMos6526(m, label, instance))
	return m
}

// Setup initializes the CIA component, creating TOD and Timer instances, binding timer underflow signals, and setting the socket.
func (m *CIA) Setup() error {
	m.tod = NewTOD(m, m.GetFactory(), m.label, 0)
	m.timerA = NewTimer(m, m.GetFactory(), m.label, 0)
	m.timerA.UnderflowSignal().Bind(m.timerAUnderflowSlot)
	m.timerB = NewTimer(m, m.GetFactory(), m.label, 1)
	m.timerB.UnderflowSignal().Bind(m.timerBUnderflowSlot)
	m.shiftRegister = NewShiftRegister(m, m.GetFactory(), m.label, 0)
	return nil
}

func (m *CIA) Bind(socket references.IMos6526Socket) error {
	//m.socketSignalSP = socket.SignalSP
	//m.socketReadSP = socket.ReadSP
	m.socketReadPortA = socket.ReadPortA
	m.socketReadPortB = socket.ReadPortB
	m.socketSignalPRA = socket.SignalPRA
	m.socketSignalPRB = socket.SignalPRB
	m.socketSignalDDRA = socket.SignalDDRA
	m.socketSignalDDRB = socket.SignalDDRB
	m.socketIRQTrigger = socket.IRQTrigger
	m.socketIRQClearTrigger = socket.IRQClearTrigger

	m.shiftRegister.Initialize(socket.ReadSP, socket.SignalSP)
	m.reads = m.createReadRegister()
	m.writes = m.createWriteRegister()
	return nil
}

// Connect establishes the necessary connections for the CIA component and prepares it for operation.
func (m *CIA) Connect() error {
	return nil
}

// Internal returns a fixed boolean value, used as a placeholder or representation of internal state.
func (m *CIA) Internal() bool {
	return false
}

// Emulate performs one emulation cycle, updating internal state, timers, and their interactions without triggering IRQs.
//
//go:nosplit
func (m *CIA) Emulate() {
	if m.timerAIrqCycle {
		m.timerAIrqCycle = false //next cycle trigger
		m.irqTrigger()
	}
	if m.timerBIrqCycle {
		m.timerBIrqCycle = false //next cycle trigger
		m.irqTrigger()
	}
	//m.timerAIrqCycle = false
	//m.timerBIrqCycle = false
	m.timerA.Emulate()
	m.timerB.SetUnderflowIn(m.timerA.GetUnderflowOut())
	m.timerB.Emulate()

	m.todClockDivider--
	if m.todClockDivider <= 0 {
		var freq int
		if m.timerA.GetRTC() {
			freq = dividerPAL
		} else {
			freq = dividerNTSC
		}
		m.todClockDivider = freq
		if m.tod.Update() {
			m.icr |= IRQTODAlarmEqual
			m.irqTrigger()
		}
	}
}

// EmulationRequired returns true indicating that emulation is necessary for this CIA instance.
func (m *CIA) EmulationRequired() bool {
	return true
}

// timerAUnderflowSlot handles the underflow event of Timer A by setting the respective interrupt cycle flag and triggering IRQ.
func (m *CIA) timerAUnderflowSlot() {
	if m.shiftRegister.Counter() > 0 {
		if m.shiftRegister.Handle((m.timerA.ReadCR() & crBitSPMode) != 0) {
			if (m.timerA.ReadCR() & crBitSPMode) == 0 {
				m.sdr = m.shiftRegister.Get()
			}
			m.icr |= IRQSDRFullOrEmpty
			m.irqTrigger()
		}
	}
	m.timerAIrqCycle = true
	m.icr |= IRQUnderflowTimerA
}

// timerBUnderflowSlot handles the underflow event for Timer B, sets the IRQ flag, and triggers an interrupt request.
func (m *CIA) timerBUnderflowSlot() {
	m.timerBIrqCycle = true
	m.icr |= IRQUnderflowTimerB
}

// Reset reinitializes the internal state of the CIA by resetting all registers, timers, and the time of day (TOD) clock.
func (m *CIA) Reset() {
	m.prA = 0
	m.prB = 0
	m.ddrA = 0
	m.ddrB = 0
	m.sdr = 0
	m.icr = 0
	m.irqMask = 0
	m.timerAIrqCycle = false
	m.timerBIrqCycle = false
	m.todClockDivider = 0
	m.timerA.Reset()
	m.timerB.Reset()
	m.tod.Reset()
}

// ReadPRA returns the current value of the port A register (PRA) of the CIA.
func (m *CIA) ReadPRA() uint8 {
	return m.prA
}

// ReadPRB returns the current value of the port B register (PRB) of the CIA.
func (m *CIA) ReadPRB() uint8 {
	return m.prB
}

// ReadDDRA reads the value of the Data Direction Register A (DDRA), which determines input/output configuration for port A.
func (m *CIA) ReadDDRA() uint8 {
	return m.ddrA
}

// ReadDDRB reads and returns the current value of the data direction register B (DDRB) in the CIA.
func (m *CIA) ReadDDRB() uint8 {
	return m.ddrB
}

// ReadSDR retrieves and returns the current value of the sdr field from the CIA instance.
func (m *CIA) ReadSDR() uint8 {
	return m.sdr
}

// ReadRegister reads the value of the specified register from the CIA based on the provided address.
func (m *CIA) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x0f
	return m.reads[reg]()
}

// WriteRegister writes data to the specified address within the CIA, updating registers or triggering system operations.
func (m *CIA) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x0f
	m.writes[reg](data)
}

// WriteCNTLevel sets the CNT pin state for both Timer A and Timer B to the specified level.
func (m *CIA) WriteCNTLevel(level bool) {
	m.timerA.SetCNTLevel(level)
	m.timerB.SetCNTLevel(level)
}

// WriteCNTPulse sends a pulse signal to Timer A's CNT line, triggering it to perform a defined operation.
func (m *CIA) WriteCNTPulse() {
	m.timerA.SetCNTPulse()
}

// irqTrigger checks active interrupts against the current IRQ mask and triggers an IRQ if conditions are met.
func (m *CIA) irqTrigger() {
	mask := m.irqMask & 0x1f
	if (m.icr & mask) != 0 {
		m.icr |= IRQOccurred
		m.socketIRQTrigger()
	}
}

// irqUpdateMask updates the interrupt mask based on the provided data byte.
// Bits [0..4] determine the mask, and bit 7 specifies whether to set or clear the mask bits.
func (m *CIA) irqUpdateMask(data uint8) {
	//Bit 0: 1 = Interrupt release through timer A underflow
	//Bit 1: 1 = Interrupt release through timer B underflow
	//Bit 2: 1 = Interrupt release if clock=alarm
	//Bit 3: 1 = Interrupt release if a complete byte has been received/sent.
	//Bit 4: 1 = Interrupt release if a positive slope occurs at the FLAG-Pin.
	//Bit 5..6: unused
	//Bit 7: Source bit.
	//     0 = set bits 0..4 are clearing the according mask bit.
	//     1 = set bits 0..4 are setting the according mask bit.
	//If all 5 bits [0..4] are cleared, there will be no change to the mask.
	if bits := data & 0x1f; bits != 0 {
		//Bit 7: Source bit.
		// 1 = set bits 0..4 are setting the according mask bit.
		// 0 = set bits 0..4 are clearing the according mask bit.
		if (data & 0x80) != 0 {
			//set bits 0..4 are setting the according mask bit.
			m.irqMask |= bits //data & 0x7f
		} else {
			//set bits 0..4 are clearing the according mask bit.
			m.irqMask &= ^bits //^data
		}
	}
}

// ReadPortA reads the current state of port A by considering its data direction and peripheral registers.
func (m *CIA) ReadPortA() uint8 {
	return m.socketReadPortA(m.prA, m.prB, m.ddrA, m.ddrB)
}

// ReadPortB reads and returns the current value of Port B, including modifications based on timers with PB output enabled.
func (m *CIA) ReadPortB() uint8 {
	ret := m.socketReadPortB(m.prA, m.prB, m.ddrA, m.ddrB)
	// TA/TB output to PB enabled
	if m.timerA.HasPBOn() {
		if m.timerA.ToggleModeApply(m.timerAIrqCycle) {
			ret |= 0x40
		} else {
			ret &= 0xbf
		}
	}
	if m.timerB.HasPBOn() {
		if m.timerB.ToggleModeApply(m.timerBIrqCycle) {
			ret |= 0x80
		} else {
			ret &= 0x7f
		}
	}
	return ret
}

// ReadICR reads the current ICR value, resets it, and clears the IRQ trigger if ICR is non-zero. Returns the ICR value.
func (m *CIA) ReadICR() uint8 {
	icr := m.icr
	m.icr = 0
	if icr != 0 {
		m.socketIRQClearTrigger()
	}
	return icr
}

// WritePRA writes the given data to the prA register and triggers the socketSignalPRA with the updated value.
func (m *CIA) WritePRA(data uint8) {
	m.prA = data
	m.socketSignalPRA(m.prA)
}

// WritePRB sets the prB field to the provided data value and triggers socketSignalPRB with the updated prB value.
func (m *CIA) WritePRB(data uint8) {
	m.prB = data
	m.socketSignalPRB(m.prB)
}

// WriteDDRA updates the data direction register A and triggers the associated socket signal with the new value.
func (m *CIA) WriteDDRA(data uint8) {
	m.ddrA = data
	m.socketSignalDDRA(m.ddrA)
}

// WriteDDRB writes an 8-bit data value to the DDRB register and triggers the socketSignalDDRB function with the updated value.
func (m *CIA) WriteDDRB(data uint8) {
	m.ddrB = data
	m.socketSignalDDRB(m.ddrB)
}

// WriteSDR writes an 8-bit value to the SDR (Serial Data Register) and updates the shift register if serial mode is active.
func (m *CIA) WriteSDR(data uint8) {
	m.sdr = data
	if (m.timerA.ReadCR() & crBitSPMode) != 0 {
		//sdr interrupt at the end of the transmission
		m.shiftRegister.Set(data)
	}
}

// WriteICR writes the given data to the Interrupt Control Register (ICR) and updates the IRQ state accordingly.
func (m *CIA) WriteICR(data uint8) {
	m.irqUpdateMask(data)
	m.irqTrigger()
}

// WriteControlRegisterTimerA sets the control register for Timer A with the given data and determines the counting mode.
func (m *CIA) WriteControlRegisterTimerA(data uint8) {
	//00 = Timer counts system cycles
	//01 = Timer counts positive slope at CNT-pin
	countMode := uint8(0)
	if (data & crBitInMode) != 0 {
		countMode = 1
	}
	m.timerA.SetControlRegister(data, countMode)
}

// WriteControlRegisterTimerB sets the control register for timer B based on the provided data and count mode settings.
func (m *CIA) WriteControlRegisterTimerB(data uint8) {
	//00 = Timer counts System cycle
	//01 = Timer counts positive slope on CNT-pin
	//10 = Timer counts underflow of timer A
	//11 = Timer counts underflow of timer A if the CNT-pin is high
	//crBitInMode | crBitSPMode
	countMode := (data >> 5) & 0x3
	m.timerB.SetControlRegister(data, countMode)
}

// WriteTOD10ths writes the 10ths component of the Time-of-Day (TOD) clock using the provided data.
func (m *CIA) WriteTOD10ths(data uint8) {
	m.tod.Set10ths(m.timerB.GetRTC(), data)
}

// WriteTODSec sets the seconds value of the Time of Day (TOD) clock using the provided data parameter.
func (m *CIA) WriteTODSec(data uint8) {
	m.tod.SetSec(m.timerB.GetRTC(), data)
}

// WriteTODMin updates the minute value of the TOD clock using the provided data and RTC configuration.
func (m *CIA) WriteTODMin(data uint8) {
	m.tod.SetMin(m.timerB.GetRTC(), data)
}

// WriteTODHour updates the TOD hour value using the provided data and the RTC configuration from Timer B.
func (m *CIA) WriteTODHour(data uint8) {
	m.tod.SetHour(m.timerB.GetRTC(), data)
}

// createReadRegister initializes and returns an array of functions to read various registers of the CIA chip.
func (m *CIA) createReadRegister() [RegisterCount]func() uint8 {
	var reads [RegisterCount]func() uint8
	reads[0x00] = m.ReadPortA
	reads[0x01] = m.ReadPortB
	reads[0x02] = m.ReadDDRA
	reads[0x03] = m.ReadDDRB
	reads[0x04] = m.timerA.ReadTimerLow
	reads[0x05] = m.timerA.ReadTimerHigh
	reads[0x06] = m.timerB.ReadTimerLow
	reads[0x07] = m.timerB.ReadTimerHigh
	reads[0x08] = m.tod.Read10ths
	reads[0x09] = m.tod.ReadSec
	reads[0x0a] = m.tod.ReadMin
	reads[0x0b] = m.tod.ReadHour
	reads[0x0c] = m.ReadSDR
	reads[0x0d] = m.ReadICR
	reads[0x0e] = m.timerA.ReadCR
	reads[0x0f] = m.timerB.ReadCR
	return reads
}

// createWriteRegister initializes an array of functions to handle writing to specific VIA registers.
func (m *CIA) createWriteRegister() [RegisterCount]func(uint8) {
	var writes [RegisterCount]func(uint8)
	writes[0x00] = m.WritePRA
	writes[0x01] = m.WritePRB
	writes[0x02] = m.WriteDDRA
	writes[0x03] = m.WriteDDRB
	writes[0x04] = m.timerA.WriteTimerLow
	writes[0x05] = m.timerA.WriteTimerHigh
	writes[0x06] = m.timerB.WriteTimerLow
	writes[0x07] = m.timerB.WriteTimerHigh
	writes[0x08] = m.WriteTOD10ths
	writes[0x09] = m.WriteTODSec
	writes[0x0a] = m.WriteTODMin
	writes[0x0b] = m.WriteTODHour
	writes[0x0c] = m.WriteSDR
	writes[0x0d] = m.WriteICR
	writes[0x0e] = m.WriteControlRegisterTimerA
	writes[0x0f] = m.WriteControlRegisterTimerB
	return writes
}
