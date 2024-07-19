package cia

import (
	"github.com/markel1974/c64emu/src/board/flag"
	"github.com/markel1974/c64emu/src/preferences"
)

type MOS6526_2 struct {
	*MOS6526
	intr  IInterrupts
	vic   IVic
	bus   IBus
	prefs *preferences.Prefs
}

func NewMOS6526_2() *MOS6526_2 {
	m := &MOS6526_2{}
	m.MOS6526 = NewMOS6526(m.TriggerInterrupt)
	return m
}

func (cia2 *MOS6526_2) Setup(intr IInterrupts, vic IVic, bus IBus, prefs *preferences.Prefs) {
	cia2.intr = intr
	cia2.vic = vic
	cia2.bus = bus
	cia2.prefs = prefs

	//TODO IMPLEMENT
}

func (cia2 *MOS6526_2) Reset() {
	cia2.MOS6526.Reset()
	// VA14/15 = 0
	cia2.vic.ChangedVA(0)
}

func (cia2 *MOS6526_2) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		return ((cia2.prA | (^cia2.ddrA)) & 0x3f) | cia2.bus.CpuRead()

	case 0x01:
		return cia2.prB | (^cia2.ddrB)

	case 0x02:
		return cia2.ddrA

	case 0x03:
		return cia2.ddrB

	case 0x04:
		return uint8(cia2.timerA)

	case 0x05:
		return uint8(cia2.timerA >> 8)

	case 0x06:
		return uint8(cia2.timerB)

	case 0x07:
		return uint8(cia2.timerB >> 8)

	case 0x08:
		cia2.todHalt = false
		return cia2.tod10ths

	case 0x09:
		return cia2.todSec

	case 0x0a:
		return cia2.todMin

	case 0x0b:
		cia2.todHalt = true
		return cia2.todHr

	case 0x0c:
		return cia2.sdr

	case 0x0d:
		ret := cia2.icr // Read and clear ICR
		cia2.icr = 0
		cia2.intr.ClearNMI()
		return ret

	case 0x0e:
		return cia2.crA

	case 0x0f:
		return cia2.crB
	}
	return 0 // Can't happen
}

func (cia2 *MOS6526_2) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x0:
		cia2.prA = data
		cia2.vic.ChangedVA((^(cia2.prA | (^cia2.ddrA))) & 3)
		cia2.bus.CpuWrite(data)

	case 0x1:
		cia2.prB = data

	case 0x2:
		cia2.ddrA = data
		cia2.vic.ChangedVA((^(cia2.prA | (^cia2.ddrA))) & 3)

	case 0x3:
		cia2.ddrB = data

	case 0x4:
		cia2.latchA = (cia2.latchA & 0xff00) | uint16(data)

	case 0x5:
		cia2.latchA = (cia2.latchA & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if !flag.Uint8ToBool(cia2.crA & 1) {
			cia2.timerA = cia2.latchA
		}

	case 0x6:
		cia2.latchB = (cia2.latchB & 0xff00) | uint16(data)

	case 0x7:
		cia2.latchB = (cia2.latchB & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if !flag.Uint8ToBool(cia2.crB & 1) {
			cia2.timerB = cia2.latchB
		}

	case 0x8:
		if flag.Uint8ToBool(cia2.crB & 0x80) {
			cia2.alm10ths = data & 0x0f
		} else {
			cia2.tod10ths = data & 0x0f
		}

	case 0x9:
		if flag.Uint8ToBool(cia2.crB & 0x80) {
			cia2.almSec = data & 0x7f
		} else {
			cia2.todSec = data & 0x7f
		}

	case 0xa:
		if flag.Uint8ToBool(cia2.crB & 0x80) {
			cia2.almMin = data & 0x7f
		} else {
			cia2.todMin = data & 0x7f
		}

	case 0xb:
		if flag.Uint8ToBool(cia2.crB & 0x80) {
			cia2.almHr = data & 0x9f
		} else {
			cia2.todHr = data & 0x9f
		}

	case 0xc:
		cia2.sdr = data
		// Fake SDR interrupt for programs that need it
		cia2.TriggerInterrupt(8)

	case 0xd:
		if flag.Uint8ToBool(data & 0x80) {
			cia2.intMask |= data & 0x7f
		} else {
			cia2.intMask &= ^data
		}
		// Trigger NMI if pending
		if flag.Uint8ToBool(cia2.icr & cia2.intMask & 0x1f) {
			cia2.icr |= 0x80
			cia2.intr.TriggerNMI()
		}

	case 0xe:
		// Delay write by 1 cycle
		cia2.hasNewCrA = true
		cia2.newCrA = data
		cia2.timerACntPhi2 = (data & 0x20) == 0x00

	case 0xf:
		// Delay write by 1 cycle
		cia2.hasNewCrB = true
		cia2.newCrB = data
		cia2.timerBCntPhi2 = (data & 0x60) == 0x00
		cia2.timerBCntTimerA = (data & 0x60) == 0x40
	}
}

func (cia2 *MOS6526_2) TriggerInterrupt(bit uint8) {
	cia2.icr |= bit
	if flag.Uint8ToBool(cia2.intMask & bit) {
		cia2.icr |= 0x80
		cia2.intr.TriggerNMI()
	}
}
