package cia

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
)

type MOS6526A struct {
	*MOS6526
	signalIRQTrigger      *signals.SignalUint32
	signalIRQClear        *signals.SignalUint32
	signalLightPenTrigger *signals.Signal
	prevLPState           uint8    // Previous state of LP line (bit 4
	KeyMatrix             [8]uint8 // C64 keyboard matrix, 1 bit/key (0: key down, 1: key up)
	RevMatrix             [8]uint8 // Reversed keyboard matrix
	Joystick1             uint8    // Joystick 1 AND value
	Joystick2             uint8    // Joystick 2 AND value
}

func NewMOS6526A() *MOS6526A {
	m := &MOS6526A{
		signalIRQTrigger:      signals.NewSignalUint32(),
		signalIRQClear:        signals.NewSignalUint32(),
		signalLightPenTrigger: signals.NewSignal(),
	}
	m.MOS6526 = NewMOS6526(m.TriggerInterrupt)
	return m
}

func (cia1 *MOS6526A) Setup(prefs *config.Config) {
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

func (cia1 *MOS6526A) Reset() {
	cia1.MOS6526.Reset()

	for i := 0; i < 8; i++ {
		cia1.KeyMatrix[i] = 0xff
		cia1.RevMatrix[i] = 0xff
	}
	cia1.Joystick1 = 0xff
	cia1.Joystick2 = 0xff
	cia1.prevLPState = 0x10
}

func (cia1 *MOS6526A) SetKeyUp(keyM int, revM int, shifted bool, joyKey1 uint8, joyKey2 uint8) {
	if shifted {
		cia1.KeyMatrix[6] |= 0x10
		cia1.RevMatrix[4] |= 0x40
	}
	cia1.KeyMatrix[keyM] |= 1 << revM
	cia1.RevMatrix[revM] |= 1 << keyM
	cia1.Joystick1 = joyKey1
	cia1.Joystick2 = joyKey2
}

func (cia1 *MOS6526A) SetKeyDown(keyM int, revM int, shifted bool, joyKey1 uint8, joyKey2 uint8) {
	if shifted {
		cia1.KeyMatrix[6] &= 0xef
		cia1.RevMatrix[4] &= 0xbf
	}
	cia1.KeyMatrix[keyM] &= ^(1 << revM)
	cia1.RevMatrix[revM] &= ^(1 << keyM)
	cia1.Joystick1 = joyKey1
	cia1.Joystick2 = joyKey2
}

func (cia1 *MOS6526A) SetJoystick1(port1 uint8) {
	cia1.Joystick1 = port1
}

func (cia1 *MOS6526A) SetJoystick2(port2 uint8) {
	cia1.Joystick2 = port2
}

func (cia1 *MOS6526A) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		ret := cia1.prA | ^cia1.ddrA
		tst := (cia1.prB | ^cia1.ddrB) & cia1.Joystick1
		if (tst & 0x01) == 0 {
			ret &= cia1.RevMatrix[0]
		}
		if (tst & 0x02) == 0 {
			ret &= cia1.RevMatrix[1]
		}
		if (tst & 0x04) == 0 {
			ret &= cia1.RevMatrix[2]
		}
		if (tst & 0x08) == 0 {
			ret &= cia1.RevMatrix[3]
		}
		if (tst & 0x10) == 0 {
			ret &= cia1.RevMatrix[4]
		}
		if (tst & 0x20) == 0 {
			ret &= cia1.RevMatrix[5]
		}
		if (tst & 0x40) == 0 {
			ret &= cia1.RevMatrix[6]
		}
		if (tst & 0x80) == 0 {
			ret &= cia1.RevMatrix[7]
		}
		return ret & cia1.Joystick2

	case 0x01:
		ret := ^cia1.ddrB
		tst := (cia1.prA | ^cia1.ddrA) & cia1.Joystick2
		if (tst & 0x01) == 0 {
			ret &= cia1.KeyMatrix[0]
		}
		if (tst & 0x02) == 0 {
			ret &= cia1.KeyMatrix[1]
		}
		if (tst & 0x04) == 0 {
			ret &= cia1.KeyMatrix[2]
		}
		if (tst & 0x08) == 0 {
			ret &= cia1.KeyMatrix[3]
		}
		if (tst & 0x10) == 0 {
			ret &= cia1.KeyMatrix[4]
		}
		if (tst & 0x20) == 0 {
			ret &= cia1.KeyMatrix[5]
		}
		if (tst & 0x40) == 0 {
			ret &= cia1.KeyMatrix[6]
		}
		if (tst & 0x80) == 0 {
			ret &= cia1.KeyMatrix[7]
		}
		return (ret | (cia1.prB & cia1.ddrB)) & cia1.Joystick1

	case 0x02:
		return cia1.ddrA

	case 0x03:
		return cia1.ddrB

	case 0x04:
		return uint8(cia1.timerA)

	case 0x05:
		return uint8(cia1.timerA >> 8)

	case 0x06:
		return uint8(cia1.timerB)

	case 0x07:
		return uint8(cia1.timerB >> 8)

	case 0x08:
		cia1.todHalt = false
		return cia1.tod10ths

	case 0x09:
		return cia1.todSec

	case 0x0a:
		return cia1.todMin

	case 0x0b:
		cia1.todHalt = true
		return cia1.todHr

	case 0x0c:
		return cia1.sdr

	case 0x0d:
		ret := cia1.icr // Read and clear ICR
		cia1.icr = 0
		cia1.signalIRQClear.Emit(intrCia1Id)
		return ret

	case 0x0e:
		return cia1.crA

	case 0x0f:
		return cia1.crB
	}
	return 0 // Can't happen
}

func (cia1 *MOS6526A) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x0:
		cia1.prA = data

	case 0x1:
		cia1.prB = data
		cia1.checkLightPen()

	case 0x2:
		cia1.ddrA = data

	case 0x3:
		cia1.ddrB = data
		cia1.checkLightPen()

	case 0x4:
		cia1.latchA = (cia1.latchA & 0xff00) | uint16(data)

	case 0x5:
		cia1.latchA = (cia1.latchA & 0xff) | (uint16(data) << 8)
		if (cia1.crA & 1) == 0 {
			// Reload timer if stopped
			cia1.timerA = cia1.latchA
		}

	case 0x6:
		cia1.latchB = (cia1.latchB & 0xff00) | uint16(data)

	case 0x7:
		cia1.latchB = (cia1.latchB & 0xff) | (uint16(data) << 8)
		if (cia1.crB & 1) == 0 {
			// Reload timer if stopped
			cia1.timerB = cia1.latchB
		}

	case 0x8:
		if (cia1.crB & 0x80) != 0 {
			cia1.alm10ths = data & 0x0f
		} else {
			cia1.tod10ths = data & 0x0f
		}

	case 0x9:
		if (cia1.crB & 0x80) != 0 {
			cia1.almSec = data & 0x7f
		} else {
			cia1.todSec = data & 0x7f
		}

	case 0xa:
		if (cia1.crB & 0x80) != 0 {
			cia1.almMin = data & 0x7f
		} else {
			cia1.todMin = data & 0x7f
		}

	case 0xb:
		if (cia1.crB & 0x80) != 0 {
			cia1.almHr = data & 0x9f
		} else {
			cia1.todHr = data & 0x9f
		}

	case 0xc:
		cia1.sdr = data
		// Fake SDR interrupt for programs that need it
		cia1.TriggerInterrupt(8)

	case 0xd:
		if (data & 0x80) != 0 {
			cia1.intMask |= data & 0x7f
		} else {
			cia1.intMask &= ^data
		}
		if (cia1.icr & cia1.intMask & 0x1f) != 0 {
			// Trigger IRQ if pending
			cia1.icr |= 0x80
			cia1.signalIRQTrigger.Emit(intrCia1Id)
		}

	case 0xe:
		// Delay write by 1 cycle
		cia1.hasNewCrA = true
		cia1.newCrA = data
		cia1.timerACntPhi2 = (data & 0x20) == 0x00

	case 0xf:
		// Delay write by 1 cycle
		cia1.hasNewCrB = true
		cia1.newCrB = data
		cia1.timerBCntPhi2 = (data & 0x60) == 0x00
		cia1.timerBCntTimerA = (data & 0x60) == 0x40
	}
}

func (cia1 *MOS6526A) TriggerInterrupt(bit uint8) {
	cia1.icr |= bit
	if (cia1.intMask & bit) != 0 {
		cia1.icr |= 0x80
		cia1.signalIRQTrigger.Emit(intrCia1Id)
	}
}

// Write to port B, check for lightPen interrupt
func (cia1 *MOS6526A) checkLightPen() {
	if ((cia1.prB | ^cia1.ddrB) & 0x10) != cia1.prevLPState {
		cia1.signalLightPenTrigger.Emit()
	}
	cia1.prevLPState = (cia1.prB | ^cia1.ddrB) & 0x10
}
