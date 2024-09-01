package mos6510

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/components/quartz"
)

type Pic struct {
	intRstBit     uint32
	intNmiBit     uint32
	intIrqBit     uint32
	quartz        *quartz.Quartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
}

func NewPic(rstBit uint32, nmiBit uint32, irqBit uint32) *Pic {
	return &Pic{
		intRstBit:     rstBit,
		intNmiBit:     nmiBit,
		intIrqBit:     irqBit,
		quartz:        nil,
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
	i.all.BitSet(i.intRstBit)
}

func (i *Pic) HasReset() bool {
	return i.all.BitCheck(i.intRstBit)
}

func (i *Pic) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.all.BitSet(i.intIrqBit)
	i.irq.BitSet(intr)
}

func (i *Pic) ClearIRQ(intr uint32) {
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(i.intIrqBit)
	}
}

func (i *Pic) HasIRQ() bool {
	return i.irq != 0
}

func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(i.intNmiBit) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.all.BitSet(i.intNmiBit)
}

func (i *Pic) ClearNMI() {
	i.all.BitClear(i.intNmiBit)
}

func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(i.intNmiBit)
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
