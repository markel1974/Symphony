package mos6510

import (
	"github.com/markel1974/c64emu/src/bits"
	"github.com/markel1974/c64emu/src/components/quartz"
)

const (
	intrRstBit = 0
	intrNmiBit = 1
	intrIrqBit = 2
)

const (
	opFlagIrqDisabled = 0x01
	opFlagIrqEnabled  = 0x02
	opFlagIntDelayed  = 0x04
)

// http://6502.org/tutorials/interrupts.html

type Pic struct {
	quartz        *quartz.Quartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	nmiExec       bool
}

func NewPic() *Pic {
	return &Pic{
		quartz:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		all:           bits.Bits(0),
		irq:           bits.Bits(0),
		nmiExec:       false,
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
	i.nmiExec = false
	//i.irqExec = false
}

func (i *Pic) TriggerReset() {
	i.all.BitSet(intrRstBit)
}

func (i *Pic) ClearReset() {
	i.all.BitClear(intrRstBit)
}

//func (i *Pic) HasReset() bool {
//	return i.all.BitCheck(i.intRstBit)
//}

func (i *Pic) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.irq.BitSet(intr)
	i.all.BitSet(intrIrqBit)
}

func (i *Pic) ClearIRQ(intr uint32) {
	//i.irqExec = false
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(intrIrqBit)
	}
}

func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(intrNmiBit) {
		i.firstNMICycle = i.quartz.Cycle()
		i.all.BitSet(intrNmiBit)
	}
}

func (i *Pic) ClearNMI() {
	i.nmiExec = false
	i.all.BitClear(intrNmiBit)
}

func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(intrNmiBit)
}

func (i *Pic) VerifyIrq(iFlag uint8, opFlags uint8) uint8 {
	const minNMIDistance = 2
	const minIrqDistance = 2
	if i.all != 0 {
		if i.all.BitCheck(intrRstBit) {
			// Edge-triggered
			i.ClearReset()
			return 1
		}
		if i.all.BitCheck(intrNmiBit) && !i.nmiExec {
			if i.computeDistance(i.firstNMICycle, (opFlags&opFlagIntDelayed) != 0, minNMIDistance) {
				// Edge-triggered
				i.nmiExec = true
				return 2
			}
		}
		if i.all.BitCheck(intrIrqBit) /* && !i.irqExec */ {
			if ((iFlag == 0) || ((opFlags & opFlagIrqDisabled) != 0)) && ((opFlags & opFlagIrqEnabled) == 0) {
				if i.computeDistance(i.firstIrqCycle, (opFlags&opFlagIntDelayed) != 0, minIrqDistance) {
					// Level-triggered
					//i.irqExec = true
					return 3
				}
			}
		}
	}
	return 0
}

func (i *Pic) computeDistance(base uint64, hasDelay bool, distance uint64) bool {
	total := base + distance
	if hasDelay {
		total += 1
	}
	cycle := i.quartz.Cycle()
	if cycle >= total {
		return true
	}
	return false
}
