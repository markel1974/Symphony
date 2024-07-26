package board

import (
	"github.com/markel1974/c64emu/src/Interrupt"
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/config"
)

// Interrupts
const (
	IntVic = 0x1
	IntCia = 0x2
	IntNmi = 0x4
	IntRst = 0x8
	IntBus = 0x10
)

type Interrupts struct {
	quartz        *quartz.Quartz
	prefs         *config.Config
	intr          Interrupt.Interrupt
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

func (i *Interrupts) HasAny() bool {
	return i.intr != 0
}

func (i *Interrupts) TriggerReset() {
	i.intr.BitSet(IntRst)
}

func (i *Interrupts) HasReset() bool {
	return i.intr.BitCheck(IntRst)
}

func (i *Interrupts) HasIRQ() bool {
	return i.intr.BitCheck(IntVic) || i.intr.BitCheck(IntCia)
}

func (i *Interrupts) TriggerBUSIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntBus)
}

func (i *Interrupts) ClearBUSIRQ() {
	i.intr.BitClear(IntBus)
}

func (i *Interrupts) TriggerVICIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntVic)
}

func (i *Interrupts) ClearVICIRQ() {
	i.intr.BitClear(IntVic)
}

func (i *Interrupts) TriggerCIAIRQ() {
	i.updateIrqCycle()
	i.intr.BitSet(IntCia)
}

func (i *Interrupts) ClearCIAIRQ() {
	i.intr.BitClear(IntCia)
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

func (i *Interrupts) updateIrqCycle() {
	vic := i.intr.BitCheck(IntVic)
	cia := i.intr.BitCheck(IntCia)
	bus := i.intr.BitCheck(IntBus)
	if !(vic || cia || bus) {
		i.firstIrqCycle = i.quartz.Cycle()
	}
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
