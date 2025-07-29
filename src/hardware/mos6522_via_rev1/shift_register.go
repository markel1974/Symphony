package mos6522

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// ShiftRegister represents a shift register utilizing CB2 (Control Bus 2) for read and write operations.
type ShiftRegister struct {
	*component.BaseComponent
	shiftCounter uint8 // symphony:export shiftCounter tracks the number of shifts performed, resetting after completing an 8-bit operation.
	sr           uint8 // symphony:export sr represents the internal 8-bit storage for the shift register operations.
	reflect      *ShiftRegisterReflect
	readCB2      func() bool
	writeCB2     func(bool)
}

// NewShiftRegister creates and initializes a new instance of ShiftRegister and registers it with the provided factory and parent.
func NewShiftRegister(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *ShiftRegister {
	s := &ShiftRegister{
		BaseComponent: component.NewBaseComponent(),
	}
	s.BaseComponent.Register(factory, parent, "shiftRegister", instance, s, references.IdInternalComponent(label, instance, "ShiftRegister"))
	s.reflect = NewShiftRegisterReflect(s)
	return s
}

// Setup initializes the shift register component and prepares it for use. It returns an error if initialization fails.
func (v *ShiftRegister) Setup() error {
	return nil
}

// Connect establishes the necessary connections for the ShiftRegister's operation and ensures integration with other components.
func (v *ShiftRegister) Connect() error {
	return nil
}

// EmulationRequired determines if emulation functionality is required for this ShiftRegister. Always returns false.
func (v *ShiftRegister) EmulationRequired() bool {
	return false
}

// Emulate triggers the emulation process for the shift register, updating its internal state as required.
func (v *ShiftRegister) Emulate() {
}

// Internal returns a boolean indicating whether the shift register is in its internal mode or state.
func (v *ShiftRegister) Internal() bool {
	return true
}

// Initialize sets up the read and write callback functions for the shift register.
func (v *ShiftRegister) Initialize(readCB2 func() bool, writeCB2 func(bool)) {
	v.readCB2 = readCB2
	v.writeCB2 = writeCB2
}

// Reset resets the shift register and shift counter to their initial state (0).
func (v *ShiftRegister) Reset() {
	v.shiftCounter = 0
	v.sr = 0
}

// Get retrieves the current value stored in the shift register.
func (v *ShiftRegister) Get() uint8 {
	return v.sr
}

// Set updates the internal shift register with the provided 8-bit data.
func (v *ShiftRegister) Set(data uint8) {
	v.sr = data
}

// Handle processes a single shift operation in the shift register, either shifting in or out based on isShiftIn flag.
// It increments the shift counter and returns true when a full 8-bit operation is completed.
func (v *ShiftRegister) Handle(shiftIn bool) bool {
	if v.shiftCounter < 8 {
		if shiftIn {
			inBit := uint8(0)
			if v.readCB2() {
				inBit = 1
			}
			v.sr = (v.sr << 1) | inBit
		} else {
			outBit := (v.sr & 0x80) != 0
			v.writeCB2(outBit)
			v.sr <<= 1
		}
		v.shiftCounter++
	}
	if v.shiftCounter == 8 {
		v.shiftCounter = 9
		return true
	}
	return false
}
