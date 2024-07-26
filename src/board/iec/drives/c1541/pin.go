package c1541

import (
	"github.com/markel1974/c64emu/src/Interrupt"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

// Pin
const (
	IntVia1 = 0x1
	IntVia2 = 0x2
	IntNmi  = 0x4
	IntRst  = 0x8
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
	return i.intr.BitCheck(IntVia1) || i.intr.BitCheck(IntVia2)
}

func (i *Pin) TriggerVIA1IRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntVia1)
}

func (i *Pin) ClearVIA1IRQ() {
	i.intr.BitClear(IntVia1)
}

func (i *Pin) TriggerVIA2IRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntVia2)
}

func (i *Pin) ClearVIA2IRQ() {
	i.intr.BitClear(IntVia2)
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
	via1 := i.intr.BitCheck(IntVia1)
	via2 := i.intr.BitCheck(IntVia2)
	if !(via1 || via2) {
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
