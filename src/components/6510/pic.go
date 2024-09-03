package mos6510

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/components/quartz"
)

const (
	opFlagIrqDisabled = 0x01
	opFlagIrqEnabled  = 0x02
	opFlagIntDelayed  = 0x04
)

const minIrqCycleDistance = 32

type Pic struct {
	intRstBit     uint32
	intNmiBit     uint32
	intIrqBit     uint32
	quartz        *quartz.Quartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	lastIrqCycle  uint64
	opFlags       uint8
}

func NewPic(rstBit uint32, nmiBit uint32, irqBit uint32) *Pic {
	return &Pic{
		intRstBit:     rstBit,
		intNmiBit:     nmiBit,
		intIrqBit:     irqBit,
		quartz:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		lastIrqCycle:  0,
		all:           0,
		irq:           0,
		opFlags:       0,
	}
}

func (i *Pic) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

func (i *Pic) Reset() {
	i.all = 0
	i.irq = 0
	i.firstIrqCycle = 0
	i.firstNMICycle = 0
	i.lastIrqCycle = 0
	i.opFlags = 0
}

func (i *Pic) TriggerReset() {
	i.all.BitSet(i.intRstBit)
}

func (i *Pic) ClearReset() {
	i.all.BitClear(i.intRstBit)
}

//func (i *Pic) HasReset() bool {
//	return i.all.BitCheck(i.intRstBit)
//}

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

//func (i *Pic) HasIRQ() bool {
//	return i.irq != 0
//}

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

func (i *Pic) SetOpFlagIrqDisabled() {
	i.opFlags |= opFlagIrqEnabled
}

func (i *Pic) SetOpFlagIrqEnabled() {
	i.opFlags |= opFlagIrqEnabled
}

func (i *Pic) SetOpFlagIntDelayed() {
	i.opFlags |= opFlagIntDelayed
}

func (i *Pic) ClearOPFlags() {
	i.opFlags = 0
}

func (i *Pic) VerifyIrq(iFlag uint8) uint8 {
	if i.all != 0 && (i.quartz.Cycle()-i.lastIrqCycle) > minIrqCycleDistance {
		if i.all.BitCheck(i.intRstBit) {
			i.lastIrqCycle = i.quartz.Cycle()
			i.opFlags = 0
			// Edge-triggered
			i.ClearReset()
			return 1
		} else if i.all.BitCheck(i.intNmiBit) {
			if (i.computeDistance(i.firstNMICycle)) >= 2 {
				i.lastIrqCycle = i.quartz.Cycle()
				i.opFlags = 0
				// Edge-triggered
				i.ClearNMI()
				return 2
			}
		} else if i.irq != 0 {
			if ((iFlag == 0) || ((i.opFlags & opFlagIrqDisabled) != 0)) && ((i.opFlags & opFlagIrqEnabled) == 0) {
				if (i.computeDistance(i.firstIrqCycle)) >= 2 {
					// Level-triggered
					i.lastIrqCycle = i.quartz.Cycle()
					i.opFlags = 0
					return 3
				}
			}
		}
	}
	return 0
}

func (i *Pic) computeDistance(base uint64) uint64 {
	delay := uint64(0)
	if (i.opFlags & opFlagIntDelayed) != 0 {
		delay = 1
	}
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
