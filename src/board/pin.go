package board

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

const (
	IntNmi = 0x4
	IntRst = 0x8
)

type Pin struct {
	quartz        *quartz.Quartz
	prefs         *config.Config
	pinOut        bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
}

func NewPin() *Pin {
	return &Pin{
		quartz:        nil,
		prefs:         nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		pinOut:        0,
		irq:           0,
	}
}

func (i *Pin) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

func (i *Pin) Reset() {
	i.pinOut = 0
	i.irq = 0
}

func (i *Pin) HasAny() bool {
	return i.pinOut != 0
}

func (i *Pin) TriggerReset() {
	i.pinOut.BitSet(IntRst)
}

func (i *Pin) HasReset() bool {
	return i.pinOut.BitCheck(IntRst)
}

func (i *Pin) HasIRQ() bool {
	return i.irq != 0
}

func (i *Pin) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.pinOut.BitSet(intr)
	i.irq.BitSet(intr)
}

func (i *Pin) ClearIRQ(intr uint32) {
	i.pinOut.BitClear(intr)
	i.irq.BitClear(intr)
}

func (i *Pin) TriggerNMI() {
	if !i.pinOut.BitCheck(IntNmi) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.pinOut.BitSet(IntNmi)
}

func (i *Pin) ClearNMI() {
	i.pinOut.BitClear(IntNmi)
}

func (i *Pin) HasNMI() bool {
	return i.pinOut.BitCheck(IntNmi)
}

func (i *Pin) GetNMICycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstNMICycle, uint64(delay))
}

func (i *Pin) GetIrqCycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstIrqCycle, uint64(delay))
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

//func (i *Pin) AsyncNMI() {
//	i.pinOut.BitSet(IntNmi)
//}
