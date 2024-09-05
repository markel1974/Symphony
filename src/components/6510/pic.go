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

//const MinIrqCycleDistance = 32

type Pic struct {
	//irqCycleDistance uint64
	quartz        *quartz.Quartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	//lastIrqCycle     uint64
	opFlags uint8
}

func NewPic() *Pic {
	return &Pic{
		//irqCycleDistance: uint64(irqCycleDistance),
		quartz:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		//lastIrqCycle:     0,
		all:     0,
		irq:     0,
		opFlags: 0,
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
	//i.lastIrqCycle = 0
	i.opFlags = 0
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
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(intrIrqBit)
	}
}

//func (i *Pic) HasIRQ() bool {
//	return i.irq != 0
//}

func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(intrNmiBit) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.all.BitSet(intrNmiBit)
}

func (i *Pic) ClearNMI() {
	i.all.BitClear(intrNmiBit)
}

func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(intrNmiBit)
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
	const minIrqDistance = 2
	if i.all != 0 {
		if i.all.BitCheck(intrRstBit) {
			i.opFlags = 0
			// Edge-triggered
			i.ClearReset()
			return 1
		} else if i.all.BitCheck(intrNmiBit) {
			if (i.computeDistance(i.firstNMICycle)) >= minIrqDistance {
				i.opFlags = 0
				// Edge-triggered
				i.ClearNMI()
				return 2
			}
		} else if i.all.BitCheck(intrIrqBit) {
			if ((iFlag == 0) || ((i.opFlags & opFlagIrqDisabled) != 0)) && ((i.opFlags & opFlagIrqEnabled) == 0) {
				if (i.computeDistance(i.firstIrqCycle)) >= minIrqDistance {
					// Level-triggered
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
