package mos6526

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// ShiftRegister represents an 8-bit shift register component with bit counter and custom socket signaling logic.
type ShiftRegister struct {
	*component.BaseComponent
	register       uint8
	counter        uint8
	socketReadSP   func() bool
	socketSignalSP func(bool)
}

// NewShiftRegister creates and initializes a new ShiftRegister instance, registering it with the provided factory and parent component.
func NewShiftRegister(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *ShiftRegister {
	s := &ShiftRegister{
		BaseComponent: component.NewBaseComponent(),
	}
	s.BaseComponent.Register(factory, parent, "shiftRegister", s, references.IdInternalComponent(label, instance, "ShiftRegister"))
	return s
}

// Setup initializes the ShiftRegister component and prepares it for use. It returns an error if initialization fails.
func (v *ShiftRegister) Setup() error {
	return nil
}

// Connect establishes the necessary connections for the ShiftRegister component and prepares it for operation.
func (v *ShiftRegister) Connect() error {
	return nil
}

// EmulationRequired determines if the ShiftRegister requires emulation, always returning false as no emulation is needed.
func (v *ShiftRegister) EmulationRequired() bool {
	return false
}

// Emulate performs the core emulation logic for the shift register, managing bit shifts and invoking communication sockets.
func (v *ShiftRegister) Emulate() {
}

// Internal determines whether the ShiftRegister operates in internal mode, returning true if it does.
func (v *ShiftRegister) Internal() bool {
	return true
}

// Initialize initializes the ShiftRegister by assigning the provided socketReadSP and socketSignalSP functions.
func (v *ShiftRegister) Initialize(socketReadSP func() bool, socketSignalSP func(bool)) {
	v.socketReadSP = socketReadSP
	v.socketSignalSP = socketSignalSP
}

// Get returns the current value of the shift register stored in register.
func (v *ShiftRegister) Get() uint8 {
	return v.register
}

// Counter returns the current value of the shift counter (counter).
func (v *ShiftRegister) Counter() uint8 {
	return v.counter
}

// Set updates the shift register with the provided data and resets the shift counter to 8 bits.
func (v *ShiftRegister) Set(data uint8) {
	v.register = data
	v.counter = 8
}

// Reset resets the shift register and its counter to their initial state (zero).
func (v *ShiftRegister) Reset() {
	v.register = 0
	v.counter = 0
}

// Handle processes a single bit in the shift register based on the SP mode and decrements the shift counter.
// It returns true if all bits have been processed; otherwise, it returns false.
func (v *ShiftRegister) Handle(spMode bool) bool {
	if v.counter == 0 {
		return false
	}

	// Check Timer A mode (input or output)
	if spMode {
		// Send the most significant bit (MSB) to SP pin
		// (bit 7 of our shift register)
		msbIsSet := (v.register & 0x80) != 0
		v.socketSignalSP(msbIsSet)
		// Shift bits left to prepare for next one
		v.register <<= 1
	} else {
		// Read bit from SP pin
		bit := v.socketReadSP()
		// Shift bits left
		v.register <<= 1
		// Insert new bit at the end (at LSB, bit 0)
		if bit {
			v.register |= 1
		}
	}
	v.counter--
	// If we have finished shifting all 8 bits...
	if v.counter == 0 {
		return true
	}
	return false
}
