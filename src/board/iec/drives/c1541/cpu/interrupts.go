package cpu

import (
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/preferences"
)

// Interrupts
const (
	IntVia1 = 0x1
	IntVia2 = 0x1
	IntNmi  = 0x4
	IntRst  = 0x8
)

type interrupt uint32

func (i *interrupt) Clear() {
	*i = 0
}

func (i *interrupt) BitSet(n uint32) {
	*i = *i | (1 << n)
}

func (i *interrupt) BitClear(n uint32) {
	*i = *i & ^(1 << n)
}

func (i *interrupt) BitCheck(n uint32) bool {
	v := (*i >> n) & 1
	return v != 0
}

type Interrupts struct {
	quartz        *quartz.Quartz
	prefs         *preferences.Prefs
	intr          interrupt
	firstIrqCycle uint64
	firstNMICycle uint64
}

func NewInterrupts() *Interrupts {
	return &Interrupts{
		quartz:        nil,
		prefs:         nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
	}
}

func (i *Interrupts) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

func (i *Interrupts) Reset() {
	i.intr = 0
}

func (i *Interrupts) HasInterrupt() bool {
	return i.intr != 0
}

func (i *Interrupts) AsyncReset() {
	i.intr.BitSet(IntRst)
}

func (i *Interrupts) HasReset() bool {
	return i.intr.BitCheck(IntRst)
}

func (i *Interrupts) ClearVIA1IRQ() {
	i.intr.BitClear(IntVia1) //INT_VIA1IRQ
}

func (i *Interrupts) ClearVIA2IRQ() {
	i.intr.BitClear(IntVia2) //INT_VIA1IRQ
}

func (i *Interrupts) TriggerVIA1() {
	via1 := i.intr.BitCheck(IntVia1)
	via2 := i.intr.BitCheck(IntVia2)
	if !(via1 || via2) {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.intr.BitSet(IntVia1)
}

func (i *Interrupts) TriggerVIA2() {
	via1 := i.intr.BitCheck(IntVia1)
	via2 := i.intr.BitCheck(IntVia2)
	if !(via1 || via2) {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.intr.BitSet(IntVia2)
}

func (i *Interrupts) HasVIA1() bool {
	return i.intr.BitCheck(IntVia1)
}

func (i *Interrupts) HasVIA2() bool {
	return i.intr.BitCheck(IntVia2)
}

func (i *Interrupts) AsyncNMI() {
	i.intr.BitSet(IntNmi)
}

func (i *Interrupts) TriggerNMI() {
	if !i.intr.BitCheck(IntNmi) {
		i.firstNMICycle = i.quartz.Cycle()
	}
	i.intr.BitSet(IntNmi)
}

func (i *Interrupts) ClearNMI() {
	i.intr.BitClear(IntNmi)
}

func (i *Interrupts) HasNMI() bool {
	return i.intr.BitCheck(IntNmi)
}

func (i *Interrupts) GetNMICycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstNMICycle, uint64(delay))
}

func (i *Interrupts) GetIrqCycleDistance(delay int) uint64 {
	return i.computeDistance(i.firstIrqCycle, uint64(delay))
}

func (i *Interrupts) computeDistance(base uint64, delay uint64) uint64 {
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
