package cia

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
)

type MOS6526B struct {
	*MOS6526
	bus              IBus
	prefs            *config.Config
	signalNMITrigger *signals.Signal
	signalNMIClear   *signals.Signal
	signalChangedVA  *signals.SignalByte
}

func NewMOS6526B() *MOS6526B {
	m := &MOS6526B{
		signalNMITrigger: signals.NewSignal(),
		signalNMIClear:   signals.NewSignal(),
		signalChangedVA:  signals.NewSignalByte(),
	}
	m.MOS6526 = NewMOS6526(m.TriggerInterrupt)
	return m
}

func (cia2 *MOS6526B) Setup(bus IBus, prefs *config.Config) {
	cia2.bus = bus
	cia2.prefs = prefs

	//TODO IMPLEMENT
}

func (cia2 *MOS6526B) SignalTriggerNMIBind(fn func()) {
	cia2.signalNMITrigger.Bind(fn)
}

func (cia2 *MOS6526B) SignalClearNMIBind(fn func()) {
	cia2.signalNMIClear.Bind(fn)
}

func (cia2 *MOS6526B) SignalChangedVABind(fn func(uint8)) {
	cia2.signalChangedVA.Bind(fn)
}

func (cia2 *MOS6526B) Reset() {
	cia2.MOS6526.Reset()
	// VA14/15 = 0
	cia2.signalChangedVA.Emit(0)
}

func (cia2 *MOS6526B) ReadRegister(addr uint16) uint8 {
	addr = addr & 0x0f
	switch addr {
	case 0x00:
		data := cia2.bus.CpuRead()
		ret := ((cia2.prA | (^cia2.ddrA)) & 0x3f) | data
		return ret

	case 0x01:
		ret := cia2.prB | (^cia2.ddrB)
		return ret

	case 0x02:
		return cia2.ddrA

	case 0x03:
		return cia2.ddrB

	case 0x04:
		ret := uint8(cia2.timerA)
		return ret

	case 0x05:
		ret := uint8(cia2.timerA >> 8)
		return ret

	case 0x06:
		ret := uint8(cia2.timerB)
		return ret

	case 0x07:
		ret := uint8(cia2.timerB >> 8)
		return ret

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
		ret := cia2.icr
		cia2.icr = 0
		cia2.signalNMIClear.Emit()
		return ret

	case 0x0e:
		return cia2.crA

	case 0x0f:
		return cia2.crB
	}
	return 0 // Can't happen
}

func (cia2 *MOS6526B) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x0:
		cia2.prA = data
		va := (^(cia2.prA | (^cia2.ddrA))) & 3
		cia2.signalChangedVA.Emit(va)
		cia2.bus.CpuWrite(data)

	case 0x1:
		cia2.prB = data

	case 0x2:
		cia2.ddrA = data
		va := (^(cia2.prA | (^cia2.ddrA))) & 3
		cia2.signalChangedVA.Emit(va)

	case 0x3:
		cia2.ddrB = data

	case 0x4:
		cia2.latchA = (cia2.latchA & 0xff00) | uint16(data)

	case 0x5:
		cia2.latchA = (cia2.latchA & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if (cia2.crA & 1) == 0 {
			cia2.timerA = cia2.latchA
		}

	case 0x6:
		cia2.latchB = (cia2.latchB & 0xff00) | uint16(data)

	case 0x7:
		cia2.latchB = (cia2.latchB & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if (cia2.crB & 1) == 0 {
			cia2.timerB = cia2.latchB
		}

	case 0x8:
		if (cia2.crB & 0x80) != 0 {
			cia2.alm10ths = data & 0x0f
		} else {
			cia2.tod10ths = data & 0x0f
		}

	case 0x9:
		if (cia2.crB & 0x80) != 0 {
			cia2.almSec = data & 0x7f
		} else {
			cia2.todSec = data & 0x7f
		}

	case 0xa:
		if (cia2.crB & 0x80) != 0 {
			cia2.almMin = data & 0x7f
		} else {
			cia2.todMin = data & 0x7f
		}

	case 0xb:
		if (cia2.crB & 0x80) != 0 {
			cia2.almHr = data & 0x9f
		} else {
			cia2.todHr = data & 0x9f
		}

	case 0xc:
		cia2.sdr = data
		// Fake SDR interrupt for programs that need it
		cia2.TriggerInterrupt(8)

	case 0xd:
		if (data & 0x80) != 0 {
			cia2.intMask |= data & 0x7f
		} else {
			cia2.intMask &= ^data
		}
		// Trigger NMI if pending
		if (cia2.icr & cia2.intMask & 0x1f) != 0 {
			cia2.icr |= 0x80
			cia2.signalNMITrigger.Emit()
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

func (cia2 *MOS6526B) TriggerInterrupt(bit uint8) {
	cia2.icr |= bit
	if (cia2.intMask & bit) != 0 {
		cia2.icr |= 0x80
		cia2.signalNMITrigger.Emit()
	}
}
