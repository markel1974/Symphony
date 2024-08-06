package board

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

const (
	intRstId = 0x1
	intNmiId = 0x2
)

type Pic struct {
	quartz        *quartz.Quartz
	prefs         *config.Config
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
}

func NewPic() *Pic {
	return &Pic{
		quartz:        nil,
		prefs:         nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		all:           0,
		irq:           0,
	}
}

func (i *Pic) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

func (i *Pic) Reset() {
	i.all = 0
	i.irq = 0
}

func (i *Pic) HasAny() bool {
	return i.all != 0
}

func (i *Pic) TriggerReset() {
	i.all.BitSet(intRstId)
}

func (i *Pic) HasReset() bool {
	return i.all.BitCheck(intRstId)
}

func (i *Pic) HasIRQ() bool {
	return i.irq != 0
}

func (i *Pic) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.all.BitSet(intr)
	i.irq.BitSet(intr)
}

func (i *Pic) ClearIRQ(intr uint32) {
	i.all.BitClear(intr)
	i.irq.BitClear(intr)
}

func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(intNmiId) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.all.BitSet(intNmiId)
}

func (i *Pic) ClearNMI() {
	i.all.BitClear(intNmiId)
}

func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(intNmiId)
}

func (i *Pic) GetNMICycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstNMICycle, uint64(delay))
}

func (i *Pic) GetIrqCycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstIrqCycle, uint64(delay))
}

func (i *Pic) computeDistance(base uint64, delay uint64) uint64 {
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
