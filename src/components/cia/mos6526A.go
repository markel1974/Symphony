package cia

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
)

//https://emudev.de/q00-c64/cias-timers-keyboard-and-more/

type MOS6526A struct {
	id                    string
	prA                   uint8
	prB                   uint8
	ddrA                  uint8
	ddrB                  uint8
	sdr                   uint8
	icr                   uint8 // Pending interrupts
	intMask               uint8 // Enabled interrupts
	timerAIrqNextCycle    bool  // Flag: Trigger TA IRQ in next cycle
	timerBIrqNextCycle    bool  // Flag: Trigger Timer B IRQ in next cycle
	signalIRQTrigger      *signals.SignalUint32
	signalIRQClear        *signals.SignalUint32
	signalLightPenTrigger *signals.Signal
	prevLPState           uint8    // Previous state of LP line (bit 4
	keyMatrix             [8]uint8 // C64 keyboard matrix, 1 bit/key (0: key down, 1: key up)
	revMatrix             [8]uint8 // Reversed keyboard matrix
	joy1                  uint8    // Joystick 1 AND value
	joy2                  uint8    // Joystick 2 AND value
	tod                   *TOD
	timerA                *Timer
	timerB                *Timer
}

func NewMOS6526A(id string) *MOS6526A {
	m := &MOS6526A{
		id:                    id,
		signalIRQTrigger:      signals.NewSignalUint32(),
		signalIRQClear:        signals.NewSignalUint32(),
		signalLightPenTrigger: signals.NewSignal(),
		tod:                   NewTOD(id + "_TOD"),
		timerA:                NewTimer(id+"_TIMER_A", false),
		timerB:                NewTimer(id+"_TIMER_B", true),
	}
	m.timerA.SignalUnderflowBind(func() { m.icr |= IRQUnderflowTimerA; m.timerAIrqNextCycle = true })
	m.timerB.SignalUnderflowBind(func() { m.icr |= IRQUnderflowTimerB; m.timerBIrqNextCycle = true })
	return m
}

func (cia1 *MOS6526A) Setup(_ *config.Config) {
}

func (cia1 *MOS6526A) CheckIRQs() {
	if cia1.timerAIrqNextCycle {
		cia1.timerAIrqNextCycle = false
		cia1.triggerIrq()
	}
	if cia1.timerBIrqNextCycle {
		cia1.timerBIrqNextCycle = false
		cia1.triggerIrq()
	}
}

func (cia1 *MOS6526A) Emulate() {
	underflow := cia1.timerA.Emulate(false)
	cia1.timerB.Emulate(underflow)
}

func (cia1 *MOS6526A) SignalTriggerIRQBind(fn func(uint32)) {
	cia1.signalIRQTrigger.Bind(fn)
}

func (cia1 *MOS6526A) SignalClearIRQBind(fn func(uint32)) {
	cia1.signalIRQClear.Bind(fn)
}

func (cia1 *MOS6526A) SignalLightPenTriggerBind(fn func()) {
	cia1.signalLightPenTrigger.Bind(fn)
}

func (cia1 *MOS6526A) Update() {
	if cia1.tod.Update(cia1.timerA.GetRTC()) {
		cia1.icr |= IRQTODAlarmEqual
		cia1.triggerIrq()
	}
}

func (cia1 *MOS6526A) Reset() {
	cia1.prA = 0
	cia1.prB = 0
	cia1.ddrA = 0
	cia1.ddrB = 0
	cia1.sdr = 0
	cia1.icr = 0
	cia1.intMask = 0
	cia1.timerAIrqNextCycle = false
	cia1.timerBIrqNextCycle = false
	cia1.timerA.Reset()
	cia1.timerB.Reset()
	cia1.tod.Reset()

	//External
	for i := 0; i < 8; i++ {
		cia1.keyMatrix[i] = 0xff
		cia1.revMatrix[i] = 0xff
	}
	cia1.joy1 = 0xff
	cia1.joy2 = 0xff
	cia1.prevLPState = 0x10
}

func (cia1 *MOS6526A) SetKeyUp(keyM int, revM int, shifted bool) {
	if shifted {
		cia1.keyMatrix[6] |= 0x10
		cia1.revMatrix[4] |= 0x40
	}
	cia1.keyMatrix[keyM] |= 1 << revM
	cia1.revMatrix[revM] |= 1 << keyM
}

func (cia1 *MOS6526A) SetKeyDown(keyM int, revM int, shifted bool) {
	if shifted {
		cia1.keyMatrix[6] &= 0xef
		cia1.revMatrix[4] &= 0xbf
	}
	cia1.keyMatrix[keyM] &= ^(1 << revM)
	cia1.revMatrix[revM] &= ^(1 << keyM)
}

func (cia1 *MOS6526A) SetJoystick1(port1 uint8) {
	cia1.joy1 = port1
}

func (cia1 *MOS6526A) SetJoystick2(port2 uint8) {
	cia1.joy2 = port2
}

func (cia1 *MOS6526A) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		return cia1.readPortA()
	case 0x01:
		return cia1.readPortB()
	case 0x02:
		return cia1.ddrA
	case 0x03:
		return cia1.ddrB
	case 0x04:
		return cia1.timerA.GetTimerLow()
	case 0x05:
		return cia1.timerA.GetTimerHigh()
	case 0x06:
		return cia1.timerB.GetTimerLow()
	case 0x07:
		return cia1.timerB.GetTimerHigh()
	case 0x08:
		return cia1.tod.Get10ths()
	case 0x09:
		return cia1.tod.GetSec()
	case 0x0a:
		return cia1.tod.GetMin()
	case 0x0b:
		return cia1.tod.GetHour()
	case 0x0c:
		return cia1.sdr
	case 0x0d:
		icr := cia1.icr
		cia1.icr = 0
		if icr != 0 {
			cia1.signalIRQClear.Emit(intrCia1Id)
		}
		return icr
	case 0x0e:
		return cia1.timerA.GetCR()
	case 0x0f:
		return cia1.timerB.GetCR()
	}
	return 0 // Can't happen
}

func (cia1 *MOS6526A) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		cia1.prA = data
	case 0x01:
		cia1.prB = data
		cia1.updateLightPen()
	case 0x02:
		cia1.ddrA = data
	case 0x03:
		cia1.ddrB = data
		cia1.updateLightPen()
	case 0x04:
		cia1.timerA.SetTimerLow(data)
	case 0x05:
		cia1.timerA.SetTimerHigh(data)
	case 0x06:
		cia1.timerB.SetTimerLow(data)
	case 0x07:
		cia1.timerB.SetTimerHigh(data)
	case 0x08:
		cia1.tod.Set10ths(cia1.timerB.GetRTC(), data)
	case 0x09:
		cia1.tod.SetSec(cia1.timerB.GetRTC(), data)
	case 0x0a:
		cia1.tod.SetMin(cia1.timerB.GetRTC(), data)
	case 0x0b:
		cia1.tod.SetHour(cia1.timerB.GetRTC(), data)
	case 0x0c:
		cia1.sdr = data
		cia1.icr |= IRQSDRFullOrEmpty
		cia1.triggerIrq()
	case 0x0d:
		cia1.updateIntMask(data)
		cia1.triggerIrq()
	case 0x0e:
		cia1.timerA.SetControlRegister(data)
	case 0x0f:
		cia1.timerB.SetControlRegister(data)
	}
}

func (cia1 *MOS6526A) updateIntMask(data uint8) {
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
			cia1.intMask |= bits //data & 0x7f
		} else {
			//set bits 0..4 are clearing the according mask bit.
			cia1.intMask &= ^bits //^data
		}
	}
}

func (cia1 *MOS6526A) triggerIrq() {
	mask := cia1.intMask & 0x1f
	if (cia1.icr & mask) != 0 {
		cia1.icr |= IRQOccurred
		cia1.signalIRQTrigger.Emit(intrCia1Id)
	}
}

func (cia1 *MOS6526A) readPortA() uint8 {
	//Joy port 2
	ret := cia1.prA | ^cia1.ddrA
	tst := (cia1.prB | ^cia1.ddrB) & cia1.joy1
	if (tst & 0x01) == 0 {
		ret &= cia1.revMatrix[0]
	}
	if (tst & 0x02) == 0 {
		ret &= cia1.revMatrix[1]
	}
	if (tst & 0x04) == 0 {
		ret &= cia1.revMatrix[2]
	}
	if (tst & 0x08) == 0 {
		ret &= cia1.revMatrix[3]
	}
	if (tst & 0x10) == 0 {
		ret &= cia1.revMatrix[4]
	}
	if (tst & 0x20) == 0 {
		ret &= cia1.revMatrix[5]
	}
	if (tst & 0x40) == 0 {
		ret &= cia1.revMatrix[6]
	}
	if (tst & 0x80) == 0 {
		ret &= cia1.revMatrix[7]
	}
	return ret & cia1.joy2
}

func (cia1 *MOS6526A) readPortB() uint8 {
	//joy port 1
	ret := ^cia1.ddrB
	tst := (cia1.prA | ^cia1.ddrA) & cia1.joy2
	if (tst & 0x01) == 0 {
		ret &= cia1.keyMatrix[0]
	}
	if (tst & 0x02) == 0 {
		ret &= cia1.keyMatrix[1]
	}
	if (tst & 0x04) == 0 {
		ret &= cia1.keyMatrix[2]
	}
	if (tst & 0x08) == 0 {
		ret &= cia1.keyMatrix[3]
	}
	if (tst & 0x10) == 0 {
		ret &= cia1.keyMatrix[4]
	}
	if (tst & 0x20) == 0 {
		ret &= cia1.keyMatrix[5]
	}
	if (tst & 0x40) == 0 {
		ret &= cia1.keyMatrix[6]
	}
	if (tst & 0x80) == 0 {
		ret &= cia1.keyMatrix[7]
	}
	return (ret | (cia1.prB & cia1.ddrB)) & cia1.joy1
}

func (cia1 *MOS6526A) updateLightPen() {
	if ((cia1.prB | ^cia1.ddrB) & 0x10) != cia1.prevLPState {
		cia1.signalLightPenTrigger.Emit()
	}
	cia1.prevLPState = (cia1.prB | ^cia1.ddrB) & 0x10
}
