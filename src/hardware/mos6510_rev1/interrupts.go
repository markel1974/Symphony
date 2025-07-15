package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	OpFlagIrqDisabled = 0x01
	OpFlagIrqEnabled  = 0x02
	OpFlagIntDelayed  = 0x04
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

// Interrupts represents a programmable interrupt controller managing IRQ and NMI signals within a system.
type Interrupts struct {
	*component.BaseComponent
	//quartz        references.IQuartz
	cycles        func() uint64
	all           bits.Bits
	irq           bits.Bits
	firstIrqCycle uint64
	firstNMICycle uint64
	nmiExec       bool
}

// NewInterrupts creates and initializes a new instance of Interrupts with default values and registers it with the specified parent and factory.
func NewInterrupts(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *Interrupts {
	p := &Interrupts{
		BaseComponent: component.NewBaseComponent(),
		cycles:        nil,
		firstIrqCycle: 0,
		firstNMICycle: 0,
		all:           bits.Bits(0),
		irq:           bits.Bits(0),
		nmiExec:       false,
	}
	p.BaseComponent.Register(factory, parent, "interrupts", p, references.IdInternalComponent(label, instance, "INTERRUPTS"))
	return p
}

// Setup initializes the Interrupts instance with the provided quartz reference and other dependencies. Returns an error on failure.
func (i *Interrupts) Setup() error {
	return nil
}

func (i *Interrupts) Bind(q references.IQuartz) error {
	i.cycles = q.Cycle
	return nil
}

// Connect establishes necessary bindings or configurations between the PIC and other components, returning an error if it fails.
func (i *Interrupts) Connect() error {
	return nil
}

// Internal indicates if the `Interrupts` is set as an internal device. Always returns false in this implementation.
func (i *Interrupts) Internal() bool {
	return true
}

// Emulate processes the internal state of the Interrupts instance during an emulation cycle.
func (i *Interrupts) Emulate() {
}

// EmulationRequired determines if emulation is required for the current state of the Interrupts component. Always returns false.
func (i *Interrupts) EmulationRequired() bool {
	return false
}

// Reset reinitializes the Interrupts by clearing all internal state and resetting signals to their default values.
func (i *Interrupts) Reset() {
	i.all = 0
	i.irq = 0
	i.firstIrqCycle = 0
	i.firstNMICycle = 0
	i.nmiExec = false
}

// TriggerReset sets the reset interrupt flag in the internal bitfield to signal a reset condition.
func (i *Interrupts) TriggerReset() {
	i.all.BitSet(intrRstBit)
}

// ClearReset clears the reset interrupt bit by resetting the bit at the position defined by intrRstBit in i.all.
func (i *Interrupts) ClearReset() {
	i.all.BitClear(intrRstBit)
}

// TriggerIRQ sets the specified IRQ bit and records the cycle if no IRQ has been set previously. It also updates interrupt flags.
func (i *Interrupts) TriggerIRQ(intr uint32) {
	if i.irq == 0 {
		i.firstIrqCycle = i.cycles()
	}
	i.irq.BitSet(intr)
	i.all.BitSet(intrIrqBit)
}

// ClearIRQ clears the specified interrupt request (IRQ) bit and updates related signals if necessary.
func (i *Interrupts) ClearIRQ(intr uint32) {
	//i.irqExec = false
	i.irq.BitClear(intr)
	if i.irq == 0 {
		i.all.BitClear(intrIrqBit)
	}
}

// TriggerNMI sets the NMI (Non-Maskable Interrupt) bit and captures the cycle when the NMI was triggered if not already set.
func (i *Interrupts) TriggerNMI() {
	if !i.all.BitCheck(intrNmiBit) {
		i.firstNMICycle = i.cycles()
		i.all.BitSet(intrNmiBit)
	}
}

// ClearNMI clears the NMI (Non-Maskable Interrupt) state, disabling its execution and clearing the associated bit.
func (i *Interrupts) ClearNMI() {
	i.nmiExec = false
	i.all.BitClear(intrNmiBit)
}

// HasNMI checks if the Non-Maskable Interrupt (NMI) bit is set and returns true if active, otherwise false.
func (i *Interrupts) HasNMI() bool {
	return i.all.BitCheck(intrNmiBit)
}

// VerifyIrq evaluates and handles interrupt requests (IRQ, NMI, RESET) based on internal state and input flags.
func (i *Interrupts) VerifyIrq(iFlag uint8, opFlags uint8) uint8 {
	const minNMIDistance = 2
	const minIrqDistance = 2
	if i.all != 0 {
		if i.all.BitCheck(intrRstBit) {
			// Edge-triggered
			i.ClearReset()
			return 1
		}
		if i.all.BitCheck(intrNmiBit) && !i.nmiExec {
			if i.computeDistance(i.firstNMICycle, (opFlags&OpFlagIntDelayed) != 0, minNMIDistance) {
				// Edge-triggered
				i.nmiExec = true
				return 2
			}
		}
		if i.all.BitCheck(intrIrqBit) /* && !i.irqExec */ {
			if ((iFlag == 0) || ((opFlags & OpFlagIrqDisabled) != 0)) && ((opFlags & OpFlagIrqEnabled) == 0) {
				if i.computeDistance(i.firstIrqCycle, (opFlags&OpFlagIntDelayed) != 0, minIrqDistance) {
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
func (i *Interrupts) computeDistance(base uint64, hasDelay bool, distance uint64) bool {
	total := base + distance
	if hasDelay {
		total += 1
	}
	cycle := i.cycles()
	if cycle >= total {
		return true
	}
	return false
}

//func (i *Interrupts) HasReset() bool {
//	return i.all.BitCheck(i.intRstBit)
//}
