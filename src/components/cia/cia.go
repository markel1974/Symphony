package mos6526

// IRQUnderflowTimerA represents the IRQ flag for Timer A underflow.
// IRQUnderflowTimerB represents the IRQ flag for Timer B underflow.
// IRQTODAlarmEqual represents the IRQ flag for Time-of-Day alarm equality.
// IRQSDRFullOrEmpty represents the IRQ flag for SDR full or empty buffer.
// IRQFlagPin represents the IRQ flag for the state of the FLAG pin.
// IRQOccurred represents the IRQ flag indicating an interrupt occurred.
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

// CIA represents a Complex Interface Adapter providing I/O ports, timers, and a clock for peripheral devices.
// id is the identifier for the CIA instance.
// prA and prB are the Peripheral Registers for ports A and B.
// ddrA and ddrB are the Data Direction Registers for ports A and B.
// sdr is the Serial Data Register for serial communication.
// icr holds the pending interrupt control register state.
// irqMask specifies the enabled interrupts mask for the CIA.
// timerAIrqCycle is a flag indicating if Timer A triggers an IRQ in the next cycle.
// timerBIrqCycle is a flag indicating if Timer B triggers an IRQ in the next cycle.
// tod is a pointer to the Time-of-Day clock functionality managed by the CIA.
// timerA and timerB point to Timer A and Timer B functionality respectively.
// socket represents the interface socket connected to peripheral and communication devices.
type CIA struct {
	id             string
	prA            uint8
	prB            uint8
	ddrA           uint8
	ddrB           uint8
	sdr            uint8
	icr            uint8 // Pending interrupts
	irqMask        uint8 // Enabled interrupts
	timerAIrqCycle bool  // Flag: Trigger Timer A IRQ in next cycle
	timerBIrqCycle bool  // Flag: Trigger Timer B IRQ in next cycle
	tod            *TOD
	timerA         *Timer
	timerB         *Timer
	socket         ISocket
}

// NewCIA initializes and returns a new instance of the CIA struct with the specified ID and associated components.
func NewCIA(id string) *CIA {
	m := &CIA{
		id:     id,
		tod:    NewTOD(id + "_TOD"),
		timerA: NewTimer(id + "_TIMER_A"),
		timerB: NewTimer(id + "_TIMER_B"),
	}
	return m
}

// Setup initializes the CIA instance by assigning a provided ISocket connection to its socket field.
func (m *CIA) Setup(conn ISocket) {
	m.socket = conn
}

// Update checks the TOD alarm condition and triggers an IRQ if the alarm matches the timer.
func (m *CIA) Update() {
	if m.tod.Update(m.timerA.GetRTC()) {
		m.icr |= IRQTODAlarmEqual
		m.irqTrigger()
	}
}

// Emulate processes a single emulation cycle for the CIA timers, handling timer underflows and triggering necessary IRQs.
func (m *CIA) Emulate() {
	if m.timerAIrqCycle {
		m.timerAIrqCycle = false
		//trigger in next cycle
		//m.irqTrigger()
	}

	if m.timerBIrqCycle {
		m.timerBIrqCycle = false
		//trigger in next cycle
		//m.irqTrigger()
	}

	underflowA := m.timerA.Emulate(false)
	if underflowA {
		m.timerAIrqCycle = true
		m.icr |= IRQUnderflowTimerA
		//fmt.Println("EMITTING TIMER A", m.id, m.timerA.timerLatch)
		m.irqTrigger()
	}

	underFlowB := m.timerB.Emulate(underflowA)
	if underFlowB {
		m.timerBIrqCycle = true
		m.icr |= IRQUnderflowTimerB
		//fmt.Println("EMITTING TIMER B", m.id, m.timerA.timerLatch)
		m.irqTrigger()
	}
}

// Reset initializes all registers and flags of the CIA to their default states and resets the timers and TOD.
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
	m.timerA.Reset()
	m.timerB.Reset()
	m.tod.Reset()
}

// ReadRegister reads the data from the specified register address within the CIA component and returns the corresponding value.
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
			m.socket.IRQClear()
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

// WriteRegister writes data to a specified register address in the CIA and updates the appropriate internal state or behavior.
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
		m.icr |= IRQSDRFullOrEmpty
		m.irqTrigger()
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

// GetLastByte returns the last byte of data processed by the CIA instance.
func (m *CIA) GetLastByte() uint8 {
	return 0
}

// irqTrigger checks if any enabled interrupts are pending and triggers an IRQ signal if conditions are met.
func (m *CIA) irqTrigger() {
	mask := m.irqMask & 0x1f
	if (m.icr & mask) != 0 {
		m.icr |= IRQOccurred
		m.socket.IRQTrigger()
	}
}

// irqUpdateMask updates the interrupt mask based on the input data.
// Bits 0-4 determine the conditions for enabling or disabling interrupts.
// Bit 7 determines if the conditions set or clear the corresponding mask bits.
// If all condition bits (0-4) are unset, the mask remains unchanged.
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
