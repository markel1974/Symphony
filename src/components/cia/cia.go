package mos6526

const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOrEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

//https://emudev.de/q00-c64/cias-timers-keyboard-and-more/

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

func NewCIA(id string) *CIA {
	m := &CIA{
		id:     id,
		tod:    NewTOD(id + "_TOD"),
		timerA: NewTimer(id + "_TIMER_A"),
		timerB: NewTimer(id + "_TIMER_B"),
	}
	return m
}

func (m *CIA) Setup(conn ISocket) {
	m.socket = conn
}

func (m *CIA) Update() {
	if m.tod.Update(m.timerA.GetRTC()) {
		m.icr |= IRQTODAlarmEqual
		m.irqTrigger()
	}
}

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
		m.irqTrigger()
	}

	underFlowB := m.timerB.Emulate(underflowA)
	if underFlowB {
		m.timerBIrqCycle = true
		m.icr |= IRQUnderflowTimerB
		m.irqTrigger()
	}
}

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
		}
		return icr
	case 0x0e:
		return m.timerA.GetCR()
	case 0x0f:
		return m.timerB.GetCR()
	}
	return 0 // Can't happen
}

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
		countMode := uint8(0)
		if (data & crBitInMode) != 0 {
			countMode = 1
		}
		m.timerA.SetControlRegister(data, countMode)
	case 0x0f:
		countMode := (data >> 5) & 0x3
		m.timerB.SetControlRegister(data, countMode)
	}
}

func (m *CIA) GetLastByte() uint8 {
	return 0
}

func (m *CIA) irqTrigger() {
	mask := m.irqMask & 0x1f
	if (m.icr & mask) != 0 {
		m.icr |= IRQOccurred
		m.socket.IRQTrigger()
	}
}

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
