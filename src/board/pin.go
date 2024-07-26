package board

import (
	"github.com/markel1974/c64emu/src/Interrupt"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

// Pin
const (
	IntVic = 0x1
	IntCia = 0x2
	IntNmi = 0x4
	IntRst = 0x8
	IntBus = 0x10
)

type Pin struct {
	quartz        *quartz.Quartz
	prefs         *config.Config
	intr          Interrupt.Interrupt
	firstIrqCycle uint64
	firstNMICycle uint64
}

func NewPin() *Pin {
	return &Pin{
		quartz:        nil,
		prefs:         nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
	}
}

func (i *Pin) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

func (i *Pin) Reset() {
	i.intr = 0
}

func (i *Pin) HasAny() bool {
	return i.intr != 0
}

func (i *Pin) TriggerReset() {
	i.intr.BitSet(IntRst)
}

func (i *Pin) HasReset() bool {
	return i.intr.BitCheck(IntRst)
}

func (i *Pin) HasIRQ() bool {
	return i.intr.BitCheck(IntVic) || i.intr.BitCheck(IntCia)
}

func (i *Pin) TriggerBUSIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntBus)
}

func (i *Pin) ClearBUSIRQ() {
	i.intr.BitClear(IntBus)
}

func (i *Pin) TriggerVICIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntVic)
}

func (i *Pin) ClearVICIRQ() {
	i.intr.BitClear(IntVic)
}

func (i *Pin) TriggerCIAIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntCia)
}

func (i *Pin) ClearCIAIRQ() {
	i.intr.BitClear(IntCia)
}

func (i *Pin) AsyncNMI() {
	i.intr.BitSet(IntNmi)
}

func (i *Pin) TriggerNMI() {
	if !i.intr.BitCheck(IntNmi) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.intr.BitSet(IntNmi)
}

func (i *Pin) ClearNMI() {
	i.intr.BitClear(IntNmi)
}

func (i *Pin) HasNMI() bool {
	return i.intr.BitCheck(IntNmi)
}

func (i *Pin) GetNMICycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstNMICycle, uint64(delay))
}

func (i *Pin) GetIrqCycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstIrqCycle, uint64(delay))
}

func (i *Pin) updateIrqCycle() {
	vic := i.intr.BitCheck(IntVic)
	cia := i.intr.BitCheck(IntCia)
	bus := i.intr.BitCheck(IntBus)
	if !(vic || cia || bus) {
		i.firstIrqCycle = i.quartz.Cycle()
	}
}

func (i *Pin) computeDistance(base uint64, delay uint64) uint64 {
	cycle := i.quartz.Cycle()
	if base > cycle {
		return 0
	}
	v := cycle - base
	if v < delay {
		return 0
	}
	v -= delay
	return v
}
