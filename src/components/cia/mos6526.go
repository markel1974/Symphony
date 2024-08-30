package cia

import (
	"github.com/markel1974/c64emu/src/signals"
)

const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOrEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

//https://emudev.de/q00-c64/cias-timers-keyboard-and-more/

type MOS6526 struct {
	id                 string
	prA                uint8
	prB                uint8
	ddrA               uint8
	ddrB               uint8
	sdr                uint8
	icr                uint8 // Pending interrupts
	irqMask            uint8 // Enabled interrupts
	timerAIrqNextCycle bool  // Flag: Trigger TA IRQ in next cycle
	timerBIrqNextCycle bool  // Flag: Trigger Timer B IRQ in next cycle
	signalIRQTrigger   *signals.SignalUint32
	signalIRQClear     *signals.SignalUint32
	tod                *TOD
	timerA             *Timer
	timerB             *Timer
	conn               IWiring
}

func NewMOS6526(id string) *MOS6526 {
	m := &MOS6526{
		id:               id,
		signalIRQTrigger: signals.NewSignalUint32(),
		signalIRQClear:   signals.NewSignalUint32(),
		tod:              NewTOD(id + "_TOD"),
		timerA:           NewTimer(id+"_TIMER_A", false),
		timerB:           NewTimer(id+"_TIMER_B", true),
	}
	m.timerA.SignalUnderflowBind(func() { m.icr |= IRQUnderflowTimerA; m.timerAIrqNextCycle = true })
	m.timerB.SignalUnderflowBind(func() { m.icr |= IRQUnderflowTimerB; m.timerBIrqNextCycle = true })
	return m
}

func (m *MOS6526) Setup(conn IWiring, trigger func(uint32), clear func(uint32)) {
	m.conn = conn
	m.signalIRQTrigger.Bind(trigger)
	m.signalIRQClear.Bind(clear)
}

func (m *MOS6526) CheckIRQs() {
	if m.timerAIrqNextCycle {
		m.timerAIrqNextCycle = false
		m.triggerIrq()
	}
	if m.timerBIrqNextCycle {
		m.timerBIrqNextCycle = false
		m.triggerIrq()
	}
}

func (m *MOS6526) Emulate() {
	underflow := m.timerA.Emulate(false)
	m.timerB.Emulate(underflow)
}

func (m *MOS6526) Update() {
	if m.tod.Update(m.timerA.GetRTC()) {
		m.icr |= IRQTODAlarmEqual
		m.triggerIrq()
	}
}

func (m *MOS6526) Reset() {
	m.prA = 0
	m.prB = 0
	m.ddrA = 0
	m.ddrB = 0
	m.sdr = 0
	m.icr = 0
	m.irqMask = 0
	m.timerAIrqNextCycle = false
	m.timerBIrqNextCycle = false
	m.timerA.Reset()
	m.timerB.Reset()
	m.tod.Reset()
	m.conn.Reset()
}

func (m *MOS6526) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		return m.conn.ReadPortA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x01:
		return m.conn.ReadPortB(m.prA, m.ddrA, m.prB, m.ddrB)
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
			m.signalIRQClear.Emit(intrCia1Id)
		}
		return icr
	case 0x0e:
		return m.timerA.GetCR()
	case 0x0f:
		return m.timerB.GetCR()
	}
	return 0 // Can't happen
}

func (m *MOS6526) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		m.prA = data
		m.conn.WritePortA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x01:
		m.prB = data
		m.conn.WritePortB(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x02:
		m.ddrA = data
		m.conn.WriteDdrA(m.prA, m.ddrA, m.prB, m.ddrB)
	case 0x03:
		m.ddrB = data
		m.conn.WriteDdrB(m.prA, m.ddrA, m.prB, m.ddrB)
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
		m.triggerIrq()
	case 0x0d:
		m.updateIrqMask(data)
		m.triggerIrq()
	case 0x0e:
		m.timerA.SetControlRegister(data)
	case 0x0f:
		m.timerB.SetControlRegister(data)
	}
}

func (m *MOS6526) updateIrqMask(data uint8) {
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

func (m *MOS6526) triggerIrq() {
	mask := m.irqMask & 0x1f
	if (m.icr & mask) != 0 {
		m.icr |= IRQOccurred
		m.signalIRQTrigger.Emit(intrCia1Id)
	}
}
