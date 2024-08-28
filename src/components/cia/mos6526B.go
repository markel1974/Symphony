package cia

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/signals"
)

//https://emudev.de/q00-c64/cias-timers-keyboard-and-more/

type MOS6526B struct {
	prA                uint8
	prB                uint8
	ddrB               uint8
	ddrA               uint8
	sdr                uint8
	icr                uint8 // Pending interrupts
	intMask            uint8 // Enabled interrupts
	timerAIrqNextCycle bool  // Flag: Trigger Timer A IRQ in next cycle
	timerBIrqNextCycle bool  // Flag: Trigger Timer B IRQ in next cycle
	bus                IBus
	cfg                *config.Config
	signalNMITrigger   *signals.Signal
	signalNMIClear     *signals.Signal
	signalChangedVA    *signals.SignalByte
	tod                *TOD
	timerA             *TimerA
	timerB             *TimerB
}

func NewMOS6526B() *MOS6526B {
	m := &MOS6526B{
		signalNMITrigger: signals.NewSignal(),
		signalNMIClear:   signals.NewSignal(),
		signalChangedVA:  signals.NewSignalByte(),
		tod:              NewTOD(),
	}
	m.timerA = NewTimerA()
	m.timerB = NewTimerB()
	m.timerA.SignalTimerAUnderflowBind(m.timerAUnderflowSlot)
	m.timerB.SignalTimerBUnderflowBind(m.timerBUnderflowSlot)
	return m
}

func (cia2 *MOS6526B) Setup(bus IBus, cfg *config.Config) {
	cia2.bus = bus
	cia2.cfg = cfg
	//TODO IMPLEMENT
}

func (cia2 *MOS6526B) CheckIRQs() {
	// Trigger pending interrupts
	if cia2.timerAIrqNextCycle {
		cia2.timerAIrqNextCycle = false
		cia2.triggerNMI(IRQUnderflowTimerA)
	}
	if cia2.timerBIrqNextCycle {
		cia2.timerBIrqNextCycle = false
		cia2.triggerNMI(IRQUnderflowTimerB)
	}
}

func (cia2 *MOS6526B) Emulate() {
	taUnderflow := cia2.timerA.Emulate()
	cia2.timerB.Emulate(taUnderflow)
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

func (cia2 *MOS6526B) Update() {
	if cia2.tod.Update(cia2.timerA.crA & 0x80) {
		cia2.triggerNMI(IRQTODAlarmEqual)
	}
}

func (cia2 *MOS6526B) Reset() {
	cia2.prA = 0
	cia2.prB = 0
	cia2.ddrB = 0
	cia2.ddrA = 0
	cia2.sdr = 0
	cia2.icr = 0
	cia2.intMask = 0

	cia2.timerAIrqNextCycle = false
	cia2.timerBIrqNextCycle = false

	cia2.timerA.Reset()
	cia2.timerB.Reset()
	cia2.tod.Reset()
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
		ret := uint8(cia2.timerA.timerA)
		return ret

	case 0x05:
		ret := uint8(cia2.timerA.timerA >> 8)
		return ret

	case 0x06:
		ret := uint8(cia2.timerB.timerB)
		return ret

	case 0x07:
		ret := uint8(cia2.timerB.timerB >> 8)
		return ret

	case 0x08:
		v := cia2.tod.Get10ths()
		cia2.tod.Unfreeze()
		return v

	case 0x09:
		return cia2.tod.GetSec()

	case 0x0a:
		return cia2.tod.GetMin()

	case 0x0b:
		cia2.tod.Freeze()
		return cia2.tod.GetHour()

	case 0x0c:
		return cia2.sdr

	case 0x0d:
		ret := cia2.icr
		cia2.icr = 0
		if ret != 0 {
			cia2.signalNMIClear.Emit()
		}
		return ret

	case 0x0e:
		return cia2.timerA.crA

	case 0x0f:
		return cia2.timerB.crB
	}
	return 0 // Can't happen
}

func (cia2 *MOS6526B) WriteRegister(addr uint16, data uint8) {
	addr = addr & 0x0f
	switch addr {
	case 0x0:
		//Bit 0..1: Select the position of the VIC-memory
		//Bit 2: RS-232: TXD Output, userport: Data PA 2 (pin M)
		//Bit 3..5: serial bus Output (0=High/Inactive, 1=Low/Active)
		//Bit 6..7: serial bus Input (0=Low/Active, 1=High/Inactive)
		cia2.prA = data
		cia2.UpdateVA()
		cia2.bus.CpuWrite(data)

	case 0x1:
		cia2.prB = data

	case 0x2:
		cia2.ddrA = data
		cia2.UpdateVA()

	case 0x3:
		cia2.ddrB = data

	case 0x4:
		cia2.timerA.latchA = (cia2.timerA.latchA & 0xff00) | uint16(data)

	case 0x5:
		cia2.timerA.latchA = (cia2.timerA.latchA & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if (cia2.timerA.crA & 1) == 0 {
			cia2.timerA.timerA = cia2.timerA.latchA
		}

	case 0x6:
		cia2.timerB.latchB = (cia2.timerB.latchB & 0xff00) | uint16(data)

	case 0x7:
		cia2.timerB.latchB = (cia2.timerB.latchB & 0xff) | (uint16(data) << 8)
		// Reload timer if stopped
		if (cia2.timerB.crB & 1) == 0 {
			cia2.timerB.timerB = cia2.timerB.latchB
		}

	case 0x08:
		if (cia2.timerB.crB & 0x80) != 0 {
			cia2.tod.SetAlarm10ths(data & 0x0f)
		} else {
			cia2.tod.Set10ths(data & 0x0f)
		}

	case 0x09:
		if (cia2.timerB.crB & 0x80) != 0 {
			cia2.tod.SetAlarmSec(data & 0x7f)
		} else {
			cia2.tod.SetSec(data & 0x7f)
		}

	case 0x0a:
		if (cia2.timerB.crB & 0x80) != 0 {
			cia2.tod.SetAlarmMin(data & 0x7f)
		} else {
			cia2.tod.SetMin(data & 0x7f)
		}

	case 0x0b:
		if (cia2.timerB.crB & 0x80) != 0 {
			cia2.tod.SetAlarmHour(data & 0x9f)
		} else {
			cia2.tod.SetHour(data & 0x9f)
		}

	case 0xc:
		cia2.sdr = data
		// Fake SDR interrupt for programs that need it
		cia2.triggerNMI(IRQSDRFullOtEmpty)

	case 0xd:
		if bits := data & 0x1f; bits != 0 {
			if (data & 0x80) != 0 {
				cia2.intMask |= bits //data & 0x7f
			} else {
				cia2.intMask &= ^bits //^data
			}
		}
		// Trigger NMI if pending
		mask := cia2.intMask & 0x1f
		if (cia2.icr & mask) != 0 {
			cia2.icr |= IRQOccurred
			cia2.signalNMITrigger.Emit()
		}

	case 0xe:
		// Delay write by 1 cycle
		cia2.timerA.hasNewCrA = true
		cia2.timerA.newCrA = data
		cia2.timerA.timerACntPhi2 = (data & 0x20) == 0x00

	case 0xf:
		// Delay write by 1 cycle
		cia2.timerB.hasNewCrB = true
		cia2.timerB.newCrB = data
		cia2.timerB.timerBCntPhi2 = (data & 0x60) == 0x00
		cia2.timerB.timerBCntTimerA = (data & 0x60) == 0x40
	}
}

func (cia2 *MOS6526B) timerAUnderflowSlot() {
	cia2.icr |= IRQUnderflowTimerA
	cia2.timerAIrqNextCycle = true
}

func (cia2 *MOS6526B) timerBUnderflowSlot() {
	cia2.icr |= IRQUnderflowTimerB
	cia2.timerBIrqNextCycle = true
}

func (cia2 *MOS6526B) triggerNMI(bit uint8) {
	cia2.icr |= bit
	if (cia2.intMask & bit) != 0 {
		cia2.icr |= IRQOccurred
		cia2.signalNMITrigger.Emit()
	}
}

func (cia2 *MOS6526B) UpdateVA() {
	//%00, 0: Bank 3: $C000-$FFFF, 49152-65535
	//%01, 1: Bank 2: $8000-$BFFF, 32768-49151
	//%10, 2: Bank 1: $4000-$7FFF, 16384-32767
	//%11, 3: Bank 0: $0000-$3FFF, 0-16383 (standard)
	va := (^(cia2.prA | (^cia2.ddrA))) & 3
	cia2.signalChangedVA.Emit(va)
}
