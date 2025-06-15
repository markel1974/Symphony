package mos6526

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
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
	prA              uint8
	prB              uint8
	ddrA             uint8
	ddrB             uint8
	sdr              uint8
	icr              uint8 // Pending interrupts
	irqMask          uint8 // Enabled interrupts
	timerAIrqCycle   bool  // Flag: Trigger Timer A IRQ in next cycle
	timerBIrqCycle   bool  // Flag: Trigger Timer B IRQ in next cycle
	tod              *TOD
	timerA           *Timer
	timerB           *Timer
	sdrShiftRegister uint8 // Lo shift register interno
	sdrShiftCounter  uint8 // Contatore per i bit (da 8 a 0)
	todClockDivider  int
	socket           references.IMos6526Socket
	label            string
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
	m.BaseComponent.Register(factory, parent, Identifier(), m, references.IdIMos6526(m, label, instance))
	return m
}

// Setup initializes the CIA component, creating TOD and Timer instances, binding timer underflow signals, and setting the socket.
func (m *CIA) Setup() error {
	m.tod = NewTOD(m, m.GetFactory(), m.label, 0)
	m.timerA = NewTimer(m, m.GetFactory(), m.label, 0)
	m.timerA.UnderflowSignal().Bind(m.timerAUnderflowSlot)
	m.timerB = NewTimer(m, m.GetFactory(), m.label, 1)
	m.timerB.UnderflowSignal().Bind(m.timerBUnderflowSlot)
	return nil
}

func (m *CIA) Bind(socket references.IMos6526Socket) error {
	m.socket = socket
	return nil
}

// Connect establishes the necessary connections for the CIA component and prepares it for operation.
func (m *CIA) Connect() error {
	return nil
}

func (m *CIA) Internal() bool {
	return false
}

// Update checks the TOD alarm condition against the RTC, triggers an IRQ if conditions match, and sets the alarm flag.
func (m *CIA) Update() {
	//if m.tod.Update(m.timerA.GetRTC()) {
	//	m.icr |= IRQTODAlarmEqual
	//	m.irqTrigger()
	//}
}

// Emulate performs one emulation cycle, updating internal state, timers, and their interactions without triggering IRQs.
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
			freq = 20000 // 50 Hz (PAL) 1,000,000 / 50 = 20,000
		} else {
			freq = 16667 // 60 Hz (NTSC) 1,000,000 / 60 = 16,667
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
	// If a serial bit shift is in progress...
	if m.sdrShiftCounter > 0 {
		// Check Timer A mode (input or output)
		if (m.timerA.GetCR() & crBitSPMode) != 0 {
			// Send the most significant bit (MSB) to SP pin
			// (bit 7 of our shift register)
			msbIsSet := (m.sdrShiftRegister & 0x80) != 0
			m.socket.WriteSP(msbIsSet)
			// Shift bits left to prepare for next one
			m.sdrShiftRegister <<= 1
		} else {
			// Read bit from SP pin
			bit := m.socket.ReadSP()
			// Shift bits left
			m.sdrShiftRegister <<= 1
			// Insert new bit at the end (at LSB, bit 0)
			if bit {
				m.sdrShiftRegister |= 1
			}
		}
		m.sdrShiftCounter--
		// If we have finished shifting all 8 bits...
		if m.sdrShiftCounter == 0 {
			// In INPUT mode, copy received byte to visible SDR register
			if (m.timerA.GetCR() & crBitSPMode) == 0 {
				m.sdr = m.sdrShiftRegister
			}
			m.icr |= IRQSDRFullOrEmpty
			m.irqTrigger()
		}
	}
	m.timerAIrqCycle = true
	m.icr |= IRQUnderflowTimerA
	//fmt.Println("EMITTING TIMER A", m.id, m.timerA.timerLatch)
	//m.irqTrigger()
}

// timerBUnderflowSlot handles the underflow event for Timer B, sets the IRQ flag, and triggers an interrupt request.
func (m *CIA) timerBUnderflowSlot() {
	m.timerBIrqCycle = true
	m.icr |= IRQUnderflowTimerB
	//fmt.Println("EMITTING TIMER B", m.id, m.timerA.timerLatch)
	//m.irqTrigger()
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
	m.sdrShiftRegister = 0
	m.sdrShiftCounter = 0
	m.timerA.Reset()
	m.timerB.Reset()
	m.tod.Reset()
}

// ReadRegister reads the value of the specified register from the CIA based on the provided address.
func (m *CIA) ReadRegister(addr uint16) uint8 {
	reg := addr & 0x0f
	switch reg {
	case 0x00:
		return m.socket.ReadPortA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x01:
		ret := m.socket.ReadPortB(m.prA, m.ddrA, m.prB, m.ddrB)
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
	case 0x02:
		return m.ddrA
	case 0x03:
		return m.ddrB
	case 0x04:
		return m.timerA.GetTimerLow()
	case 0x05:
		return m.timerA.GetTimerHigh()
	case 0x06:
		return m.timerB.GetTimerLow()
	case 0x07:
		return m.timerB.GetTimerHigh()
	case 0x08:
		return m.tod.Get10ths()
	case 0x09:
		return m.tod.GetSec()
	case 0x0a:
		return m.tod.GetMin()
	case 0x0b:
		return m.tod.GetHour()
	case 0x0c:
		return m.sdr
	case 0x0d:
		icr := m.icr
		m.icr = 0
		if icr != 0 {
			m.socket.IRQClearTrigger()
			//fmt.Println("CLEARING ", m.id)
		}
		return icr
	case 0x0e:
		return m.timerA.GetCR()
	case 0x0f:
		return m.timerB.GetCR()
	}
	return 0 // Impossible
}

// WriteRegister writes data to the specified address within the CIA, updating registers or triggering system operations.
func (m *CIA) WriteRegister(addr uint16, data uint8) {
	reg := addr & 0x0f
	switch reg {
	case 0x00:
		m.prA = data
		m.socket.WritePortA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x01:
		m.prB = data
		m.socket.WritePortB(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x02:
		m.ddrA = data
		m.socket.WriteDdrA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x03:
		m.ddrB = data
		m.socket.WriteDdrB(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x04:
		m.timerA.SetTimerLow(data)
	case 0x05:
		m.timerA.SetTimerHigh(data)
	case 0x06:
		m.timerB.SetTimerLow(data)
	case 0x07:
		m.timerB.SetTimerHigh(data)
	case 0x08:
		m.tod.Set10ths(m.timerB.GetRTC(), data)
	case 0x09:
		m.tod.SetSec(m.timerB.GetRTC(), data)
	case 0x0a:
		m.tod.SetMin(m.timerB.GetRTC(), data)
	case 0x0b:
		m.tod.SetHour(m.timerB.GetRTC(), data)
	case 0x0c:
		m.sdr = data
		if (m.timerA.GetCR() & crBitSPMode) != 0 {
			m.sdrShiftRegister = data
			m.sdrShiftCounter = 8
			//sdr interrupt at the end of the transmission
		}
	case 0x0d:
		m.irqUpdateMask(data)
		m.irqTrigger()
	case 0x0e:
		//00 = Timer counts system cycles
		//01 = Timer counts positive slope at CNT-pin
		countMode := uint8(0)
		if (data & crBitInMode) != 0 {
			countMode = 1
		}
		m.timerA.SetControlRegister(data, countMode)
	case 0x0f:
		//00 = Timer counts System cycle
		//01 = Timer counts positive slope on CNT-pin
		//10 = Timer counts underflow of timer A
		//11 = Timer counts underflow of timer A if the CNT-pin is high
		//crBitInMode | crBitSPMode
		countMode := (data >> 5) & 0x3
		m.timerB.SetControlRegister(data, countMode)
	}
}

// GetLastByte returns the last byte that was operated on or updated in the CIA component.
func (m *CIA) GetLastByte() uint8 {
	return 0
}

// SetCNTLevel sets the CNT pin state for both Timer A and Timer B to the specified level.
func (m *CIA) SetCNTLevel(level bool) {
	m.timerA.SetCNTLevel(level)
	m.timerB.SetCNTLevel(level)
}

// SetCNTPulse sends a pulse signal to Timer A's CNT line, triggering it to perform a defined operation.
func (m *CIA) SetCNTPulse() {
	m.timerA.SetCNTPulse()
}

// irqTrigger checks active interrupts against the current IRQ mask and triggers an IRQ if conditions are met.
func (m *CIA) irqTrigger() {
	mask := m.irqMask & 0x1f
	if (m.icr & mask) != 0 {
		m.icr |= IRQOccurred
		m.socket.IRQTrigger()
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
