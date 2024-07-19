package cia

import (
	"github.com/markel1974/c64emu/src/board/cpu"
	"github.com/markel1974/c64emu/src/board/flag"
	"github.com/markel1974/c64emu/src/board/vic"
	"github.com/markel1974/c64emu/src/preferences"
)

type MOS6526_1 struct {
	*MOS6526
	intr        *cpu.Interrupts
	vic         *vic.MOS6569
	prevLPState uint8    // Previous state of LP line (bit 4
	KeyMatrix   [8]uint8 // C64 keyboard matrix, 1 bit/key (0: key down, 1: key up)
	RevMatrix   [8]uint8 // Reversed keyboard matrix
	Joystick1   uint8    // Joystick 1 AND value
	Joystick2   uint8    // Joystick 2 AND value
	prefs       *preferences.Prefs
}

func NewMOS6526_1() *MOS6526_1 {
	m := &MOS6526_1{}
	m.MOS6526 = NewMOS6526(m.TriggerInterrupt)
	return m
}

func (cia1 *MOS6526_1) Setup(intr *cpu.Interrupts, vic *vic.MOS6569, prefs *preferences.Prefs) {
	cia1.intr = intr
	cia1.vic = vic
	cia1.prefs = prefs
}

func (cia1 *MOS6526_1) Reset() {
	cia1.MOS6526.Reset()

	for i := 0; i < 8; i++ {
		cia1.KeyMatrix[i] = 0xff
		cia1.RevMatrix[i] = 0xff
	}
	cia1.Joystick1 = 0xff
	cia1.Joystick2 = 0xff
	cia1.prevLPState = 0x10
}

func (cia1 *MOS6526_1) SetKeyUp(c64Byte int, c64Bit int, shifted bool, joyKey1 uint8, joyKey2 uint8) {
	if shifted {
		cia1.KeyMatrix[6] |= 0x10
		cia1.RevMatrix[4] |= 0x40
	}
	cia1.KeyMatrix[c64Byte] |= 1 << c64Bit
	cia1.RevMatrix[c64Bit] |= 1 << c64Byte
	cia1.Joystick1 = joyKey1
	cia1.Joystick2 = joyKey2
}

func (cia1 *MOS6526_1) SetKeyDown(c64Byte int, c64Bit int, shifted bool, joyKey1 uint8, joyKey2 uint8) {
	if shifted {
		cia1.KeyMatrix[6] &= 0xef
		cia1.RevMatrix[4] &= 0xbf
	}
	cia1.KeyMatrix[c64Byte] &= ^(1 << c64Bit)
	cia1.RevMatrix[c64Bit] &= ^(1 << c64Byte)
	cia1.Joystick1 = joyKey1
	cia1.Joystick2 = joyKey2
}

func (cia1 *MOS6526_1) SetJoystick1(port1 uint8) {
	cia1.Joystick1 = port1
}

func (cia1 *MOS6526_1) SetJoystick2(port2 uint8) {
	cia1.Joystick2 = port2
}

func (cia1 *MOS6526_1) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		ret := cia1.prA | ^cia1.ddrA
		tst := (cia1.prB | ^cia1.ddrB) & cia1.Joystick1
		// AND all active columns
		if !flag.Uint8ToBool(tst & 0x01) {
			ret &= cia1.RevMatrix[0]
		}
		if !flag.Uint8ToBool(tst & 0x02) {
			ret &= cia1.RevMatrix[1]
		}
		if !flag.Uint8ToBool(tst & 0x04) {
			ret &= cia1.RevMatrix[2]
		}
		if !flag.Uint8ToBool(tst & 0x08) {
			ret &= cia1.RevMatrix[3]
		}
		if !flag.Uint8ToBool(tst & 0x10) {
			ret &= cia1.RevMatrix[4]
		}
		if !flag.Uint8ToBool(tst & 0x20) {
			ret &= cia1.RevMatrix[5]
		}
		if !flag.Uint8ToBool(tst & 0x40) {
			ret &= cia1.RevMatrix[6]
		}
		if !flag.Uint8ToBool(tst & 0x80) {
			ret &= cia1.RevMatrix[7]
		}
		return ret & cia1.Joystick2

	case 0x01:
		ret := ^cia1.ddrB
		tst := (cia1.prA | ^cia1.ddrA) & cia1.Joystick2
		if !flag.Uint8ToBool(tst & 0x01) {
			ret &= cia1.KeyMatrix[0]
		} // AND all active rows
		if !flag.Uint8ToBool(tst & 0x02) {
			ret &= cia1.KeyMatrix[1]
		}
		if !flag.Uint8ToBool(tst & 0x04) {
			ret &= cia1.KeyMatrix[2]
		}
		if !flag.Uint8ToBool(tst & 0x08) {
			ret &= cia1.KeyMatrix[3]
		}
		if !flag.Uint8ToBool(tst & 0x10) {
			ret &= cia1.KeyMatrix[4]
		}
		if !flag.Uint8ToBool(tst & 0x20) {
			ret &= cia1.KeyMatrix[5]
		}
		if !flag.Uint8ToBool(tst & 0x40) {
			ret &= cia1.KeyMatrix[6]
		}
		if !flag.Uint8ToBool(tst & 0x80) {
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
		cia1.intr.ClearCIAIRQ() // Clear IRQ
		//ClearIRQSignal.Emit();
		return ret

	case 0x0e:
		return cia1.crA

	case 0x0f:
		return cia1.crB
	}
	return 0 // Can't happen
}

func (cia1 *MOS6526_1) WriteRegister(addr uint16, data uint8) {
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
		if !flag.Uint8ToBool(cia1.crA & 1) {
			// Reload timer if stopped
			cia1.timerA = cia1.latchA
		}

	case 0x6:
		cia1.latchB = (cia1.latchB & 0xff00) | uint16(data)

	case 0x7:
		cia1.latchB = (cia1.latchB & 0xff) | (uint16(data) << 8)
		if !flag.Uint8ToBool(cia1.crB & 1) {
			// Reload timer if stopped
			cia1.timerB = cia1.latchB
		}

	case 0x8:
		if flag.Uint8ToBool(cia1.crB & 0x80) {
			cia1.alm10ths = data & 0x0f
		} else {
			cia1.tod10ths = data & 0x0f
		}

	case 0x9:
		if flag.Uint8ToBool(cia1.crB & 0x80) {
			cia1.almSec = data & 0x7f
		} else {
			cia1.todSec = data & 0x7f
		}

	case 0xa:
		if flag.Uint8ToBool(cia1.crB & 0x80) {
			cia1.almMin = data & 0x7f
		} else {
			cia1.todMin = data & 0x7f
		}

	case 0xb:
		if flag.Uint8ToBool(cia1.crB & 0x80) {
			cia1.almHr = data & 0x9f
		} else {
			cia1.todHr = data & 0x9f
		}

	case 0xc:
		cia1.sdr = data
		cia1.TriggerInterrupt(8) // Fake SDR interrupt for programs that need it

	case 0xd:
		if flag.Uint8ToBool(data & 0x80) {
			cia1.intMask |= data & 0x7f
		} else {
			cia1.intMask &= ^data
		}
		if flag.Uint8ToBool(cia1.icr & cia1.intMask & 0x1f) {
			// Trigger IRQ if pending
			cia1.icr |= 0x80
			cia1.intr.TriggerCIAIRQ()
		}

	case 0xe:
		cia1.hasNewCrA = true // Delay write by 1 cycle
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

func (cia1 *MOS6526_1) TriggerInterrupt(bit uint8) {
	cia1.icr |= bit
	if flag.Uint8ToBool(cia1.intMask & bit) {
		cia1.icr |= 0x80
		//_cpu->TriggerCIAIRQ();
		//TriggerIRQSignal.Emit();
		cia1.intr.TriggerCIAIRQ()
	}
}

// Write to port B, check for lightPen interrupt
func (cia1 *MOS6526_1) checkLightPen() {
	if ((cia1.prB | ^cia1.ddrB) & 0x10) != cia1.prevLPState {
		cia1.vic.TriggerLightPen()
	}
	cia1.prevLPState = (cia1.prB | ^cia1.ddrB) & 0x10
}
