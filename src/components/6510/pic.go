package mos6510

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/quartz"
)

// intrRstBit represents the bit position for the Reset interrupt signal.
// intrNmiBit represents the bit position for the Non-Maskable Interrupt (NMI) signal.
// intrIrqBit represents the bit position for the Maskable Interrupt (IRQ) signal.
const (
	intrRstBit = 0
	intrNmiBit = 1
	intrIrqBit = 2
)

// opFlagIrqDisabled represents the flag for interrupt requests being disabled.
// opFlagIrqEnabled represents the flag for interrupt requests being enabled.
// opFlagIntDelayed represents the flag for delayed interrupt requests.
const (
	opFlagIrqDisabled = 0x01
	opFlagIrqEnabled  = 0x02
	opFlagIntDelayed  = 0x04
)

// http://6502.org/tutorials/interrupts.html

// Pic represents a programmable interrupt controller with IRQ and NMI handling capabilities.
type Pic struct {
	*board.BaseComponent
	quartz        *quartz.Quartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	nmiExec       bool
}

func (i *Pic) Emulate() {
	//TODO implement me
	panic("implement me")
}

// NewPic initializes and returns a pointer to a new Pic instance with default values.
func NewPic(parent board.IComponent, suffix string) *Pic {
	p := &Pic{
		BaseComponent: board.NewBaseComponent("pic", suffix),
		quartz:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		all:           bits.Bits(0),
		irq:           bits.Bits(0),
		nmiExec:       false,
	}
	board.Register(parent, p)
	return p
}

// Setup initializes the Pic instance with a Quartz instance, establishing internal dependencies.
func (i *Pic) Setup(quartz *quartz.Quartz) {
	i.quartz = quartz
}

// Reset reinitializes the Pic instance by clearing all internal state variables and flags.
func (i *Pic) Reset() {
	i.all = 0
	i.irq = 0
	i.firstIrqCycle = 0
	i.firstNMICycle = 0
	i.nmiExec = false
	//i.irqExec = false
}

// TriggerReset sets the reset bit in the internal state to initiate a reset operation.
func (i *Pic) TriggerReset() {
	i.all.BitSet(intrRstBit)
}

// ClearReset clears the reset interrupt bit in the internal state of the Pic instance.
func (i *Pic) ClearReset() {
	i.all.BitClear(intrRstBit)
}

//func (i *Pic) HasReset() bool {
//	return i.all.BitCheck(i.intRstBit)
//}

// TriggerIRQ sets the specified interrupt bit and records the cycle if it's the first IRQ occurrence.
func (i *Pic) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.irq.BitSet(intr)
	i.all.BitSet(intrIrqBit)
}

// ClearIRQ clears the specified IRQ bit in the `irq` field and resets the `intrIrqBit` in `all` if no IRQs remain active.
func (i *Pic) ClearIRQ(intr uint32) {
	//i.irqExec = false
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(intrIrqBit)
	}
}

// TriggerNMI sets the NMI interrupt bit and records the current cycle if the NMI interrupt is not already active.
func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(intrNmiBit) {
		i.firstNMICycle = i.quartz.Cycle()
		i.all.BitSet(intrNmiBit)
	}
}

// ClearNMI clears the NMI (Non-Maskable Interrupt) bit and resets the NMI execution flag to false.
func (i *Pic) ClearNMI() {
	i.nmiExec = false
	i.all.BitClear(intrNmiBit)
}

// HasNMI checks if the Non-Maskable Interrupt (NMI) bit is set in the internal state and returns true if it is set.
func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(intrNmiBit)
}

// VerifyIrq evaluates and handles interrupt requests based on the current state, flags, and timing constraints.
// Returns a uint8 value indicating the type of interrupt serviced: 0 (none), 1 (reset), 2 (NMI), or 3 (IRQ).
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

// computeDistance checks if the given cycle distance is reached, optionally adding a delay.
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
