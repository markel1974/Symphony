package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// irqRasterBit represents the interrupt request bit for raster events.
// irqSpriteToGraphicBit represents the interrupt request bit for sprite-to-graphic collisions.
// irqSpriteToSpriteBit represents the interrupt request bit for sprite-to-sprite collisions.
// irqLightPenBit represents the interrupt request bit for light pen events.
// irqMasterBit represents the master interrupt request bit.
const (
	irqRasterBit          = uint8(0x01)
	irqSpriteToGraphicBit = uint8(0x02)
	irqSpriteToSpriteBit  = uint8(0x04)
	irqLightPenBit        = uint8(0x08)
	irqMasterBit          = uint8(0x80)
)

// irqUnsetMasterBit is the bitwise negation of irqMasterBit, used to clear the master IRQ bit in irqLatch.
const (
	irqUnsetMasterBit = ^irqMasterBit
)

// Interrupts represent the structure for managing interrupt requests (IRQs) with raster and IRQ-related properties.
type Interrupts struct {
	*component.BaseComponent
	irqLatch              uint8  // irqLatch holds an 8-bit value that latches the IRQ (Interrupt Request) configuration.
	irqMask               uint8  // irqMask represents an 8-bit mask used for interrupt request (IRQ) management.
	irqRaster             uint16 // Interrupt raster line
	socketIRQTrigger      func()
	socketIRQClearTrigger func()
}

// NewInterrupts creates and initializes a new Interrupts instance with the specified parameters and configuration.
func NewInterrupts(parent references.IComponent, factory references.IComponentFactory, label string, instance int, socketIRQTrigger func(), socketIRQClearTrigger func()) *Interrupts {
	i := &Interrupts{
		BaseComponent:         component.NewBaseComponent(),
		irqRaster:             0,
		irqLatch:              0,
		irqMask:               0,
		socketIRQTrigger:      socketIRQTrigger,
		socketIRQClearTrigger: socketIRQClearTrigger,
	}
	i.BaseComponent.Register(factory, parent, "interrupts", i, references.IdInternalComponent(label, instance, "Interrupts"))
	return i
}

// Setup initializes the Interrupts component and prepares it for operation.
func (i *Interrupts) Setup() error {
	return nil
}

// Connect establishes the necessary linkages or triggers to initialize the Interrupts component for proper operation.
func (i *Interrupts) Connect() error {
	return nil
}

// EmulationRequired checks if emulation is necessary for the current interrupt component state and returns false by default.
func (i *Interrupts) EmulationRequired() bool {
	return false
}

// Emulate processes a single step in the emulation cycle for the interrupt system, managing raster operations and IRQ states.
func (i *Interrupts) Emulate() {
}

// Internal checks and returns the internal interrupt status as a boolean value.
func (i *Interrupts) Internal() bool {
	return true
}

// Reset clears the internal state, ensuring all interrupt-related settings are reset to their default values.
func (i *Interrupts) Reset() {
}

// Emit handles triggering of an interrupt by setting the IRQ latch and activating the IRQ if it matches the mask.
func (i *Interrupts) Emit(irq uint8) {
	i.irqLatch |= irq
	if (i.irqMask & irq) != 0 {
		i.irqLatch |= irqMasterBit
		i.socketIRQTrigger()
	}
}

// WriteLatch updates the IRQ latch by clearing specific bits based on the provided data and verifies IRQ conditions.
func (i *Interrupts) WriteLatch(data uint8) {
	// Verify implementation
	i.irqLatch &= ^((data & 0xf) | irqMasterBit)
	i.irqVerify() //can emit irq
	//old
	//vic.irqLatch &= ^(data & 0xf)
	//if (vic.irqLatch & vic.irqMask) != 0 {
	//	vic.irqLatch |= irqMasterBit // Set master bit if allowed interrupt still pending
	//} else {
	//	vic.socket.IRQClearTrigger()
	//}
}

// WriteMask updates the IRQ mask by applying a bitwise AND operation with the given data and 0xF.
// It also invokes irqVerify to manage the interrupt state.
func (i *Interrupts) WriteMask(data uint8) {
	i.irqMask = data & 0xf
	i.irqVerify()
}

// ReadMask returns the current interrupt mask with the upper 4 bits set to 1 (0xf0).
func (i *Interrupts) ReadMask() uint8 {
	return i.irqMask | 0xf0
}

// ReadLatch returns the current IRQ latch value combined with a constant offset of 0x70 for additional bits' configuration.
func (i *Interrupts) ReadLatch() uint8 {
	return i.irqLatch | 0x70
}

// WriteRasterLow updates the low byte of the IRQ raster line with the given 8-bit data and triggers raster settings.
func (i *Interrupts) WriteRasterLow(rasterY uint16, data uint8) {
	irqRaster := (i.irqRaster & 0xff00) | uint16(data)
	i.irqRasterSet(rasterY, irqRaster)
}

// WriteRasterHigh sets the high 8 bits of the IRQ raster line by combining the current high bits with the input data.
func (i *Interrupts) WriteRasterHigh(rasterY uint16, data uint16) {
	irqRaster := (i.irqRaster & 0xff) | data
	i.irqRasterSet(rasterY, irqRaster)
}

// irqRasterSet updates the IRQ raster line if it differs from the current value and triggers the interrupt if necessary.
func (i *Interrupts) irqRasterSet(rasterY uint16, irqRaster uint16) {
	if irqRaster != i.irqRaster {
		if rasterY == irqRaster {
			i.Emit(irqRasterBit)
		}
		i.irqRaster = irqRaster
	}
}

// irqVerify verifies the IRQ configuration by checking the current irqLatch against irqMask and updates irqLatch accordingly.
func (i *Interrupts) irqVerify() {
	if (i.irqLatch & i.irqMask) != 0 {
		i.irqLatch |= irqMasterBit
		i.socketIRQTrigger() // Trigger interrupt if pending (now allowed)
	} else {
		i.irqLatch &= irqUnsetMasterBit
		i.socketIRQClearTrigger()
	}
}
