package mos6510_pic_rev1

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// intrRstBit represents the bit position for the reset interrupt.
// intrNmiBit represents the bit position for the non-maskable interrupt.
// intrIrqBit represents the bit position for the IRQ interrupt.
const (
	intrRstBit = 0
	intrNmiBit = 1
	intrIrqBit = 2
)

// http://6502.org/tutorials/interrupts.html

// Pic represents a programmable interrupt controller managing IRQ and NMI signals within a system.
type Pic struct {
	*component.BaseComponent
	quartz        references.IQuartz
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	nmiExec       bool
}

// NewPIC creates and initializes a new instance of Pic with default values and registers it with the specified parent and factory.
func NewPIC(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Pic {
	p := &Pic{
		BaseComponent: component.NewBaseComponent(),
		quartz:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		all:           bits.Bits(0),
		irq:           bits.Bits(0),
		nmiExec:       false,
	}
	p.BaseComponent.Register(factory, parent, Identifier(), p, references.IdIMos6510Pic(p, label, instance))
	return p
}

// Setup initializes the Pic instance with the provided quartz reference and other dependencies. Returns an error on failure.
func (i *Pic) Setup() error {
	return nil
}

func (i *Pic) Bind(_ references.IMos6510PicSocket, q references.IQuartz) error {
	i.quartz = q
	return nil
}

// Connect establishes necessary bindings or configurations between the PIC and other components, returning an error if it fails.
func (i *Pic) Connect() error {
	return nil
}

// Internal indicates if the `Pic` is set as an internal device. Always returns false in this implementation.
func (i *Pic) Internal() bool {
	return false
}

// Emulate processes the internal state of the Pic instance during an emulation cycle.
func (i *Pic) Emulate() {
}

// EmulationRequired determines if emulation is required for the current state of the Pic component. Always returns false.
func (i *Pic) EmulationRequired() bool {
	return false
}

// Reset reinitializes the Pic by clearing all internal state and resetting signals to their default values.
func (i *Pic) Reset() {
	i.all = 0
	i.irq = 0
	i.firstIrqCycle = 0
	i.firstNMICycle = 0
	i.nmiExec = false
}

// TriggerReset sets the reset interrupt flag in the internal bitfield to signal a reset condition.
func (i *Pic) TriggerReset() {
	i.all.BitSet(intrRstBit)
}

// ClearReset clears the reset interrupt bit by resetting the bit at the position defined by intrRstBit in i.all.
func (i *Pic) ClearReset() {
	i.all.BitClear(intrRstBit)
}

func (i *Pic) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.quartz.Cycle()
	}
	i.irq.BitSet(intr)
	i.all.BitSet(intrIrqBit)
}

// ClearIRQ clears the specified interrupt request (IRQ) bit and updates related signals if necessary.
func (i *Pic) ClearIRQ(intr uint32) {
	//i.irqExec = false
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(intrIrqBit)
	}
}

// TriggerNMI sets the NMI (Non-Maskable Interrupt) bit and captures the cycle when the NMI was triggered if not already set.
func (i *Pic) TriggerNMI() {
	if !i.all.BitCheck(intrNmiBit) {
		i.firstNMICycle = i.quartz.Cycle()
		i.all.BitSet(intrNmiBit)
	}
}

// ClearNMI clears the NMI (Non-Maskable Interrupt) state, disabling its execution and clearing the associated bit.
func (i *Pic) ClearNMI() {
	i.nmiExec = false
	i.all.BitClear(intrNmiBit)
}

// HasNMI checks if the Non-Maskable Interrupt (NMI) bit is set and returns true if active, otherwise false.
func (i *Pic) HasNMI() bool {
	return i.all.BitCheck(intrNmiBit)
}

// VerifyIrq evaluates and handles interrupt requests (IRQ, NMI, RESET) based on internal state and input flags.
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
			if i.computeDistance(i.firstNMICycle, (opFlags&references.OpFlagIntDelayed) != 0, minNMIDistance) {
				// Edge-triggered
				i.nmiExec = true
				return 2
			}
		}
		if i.all.BitCheck(intrIrqBit) /* && !i.irqExec */ {
			if ((iFlag == 0) || ((opFlags & references.OpFlagIrqDisabled) != 0)) && ((opFlags & references.OpFlagIrqEnabled) == 0) {
				if i.computeDistance(i.firstIrqCycle, (opFlags&references.OpFlagIntDelayed) != 0, minIrqDistance) {
					// Level-triggered
					//i.irqExec = true
					return 3
				}
			}
		}
	}
	return 0
}

// computeDistance calculates if the current cycle has met or exceeded the distance from a given base value.
// It optionally accounts for a delay by adding 1 to the total distance if hasDelay is true.
// Returns true if the current cycle is greater than or equal to the total distance, false otherwise.
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

//func (i *Pic) HasReset() bool {
//	return i.all.BitCheck(i.intRstBit)
//}
