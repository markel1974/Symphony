package references

import (
	"fmt"
)

// ICIASocket represents an interface to manage input/output operations and IRQ handling for a CIA component.
// ReadPortA reads and returns port A data based on port and data direction registers.
// ReadPortB reads and returns port B data based on port and data direction registers.
// WritePortA writes to port A using the provided port and data direction register values.
// WritePortB writes to port B using the provided port and data direction register values.
// WriteDdrA writes data to the data direction register for port A.
// WriteDdrB writes data to the data direction register for port B.
// ReadSP retrieves the current level of the SP line.
// WriteSP sets the level of the SP line.
// IRQTrigger triggers an interrupt request on the system.
// IRQClearTrigger clears any active interrupt requests.
type ICIASocket interface {
	ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8

	ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8

	WritePortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WritePortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WriteDdrA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	WriteDdrB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)

	ReadSP() bool

	WriteSP(level bool)

	IRQTrigger()

	IRQClearTrigger()
}

// ICIA defines the interface for a CIA (Complex Interface Adapter) component in an emulator or system.
// Setup initializes the CIA for operation, returning an error if initialization fails.
// Bind links the CIA to a socket interface implementing the ICIASocket interface.
// Connect establishes the necessary connections for the CIA within the system.
// Reset reinitializes the CIA to its default state.
// Emulate allows the CIA to execute its required operations in synchronization with the overall system.
// Update performs periodic updates or state adjustments for the CIA during emulation.
// WriteRegister writes the specified 8-bit data value to a given 16-bit register address.
// ReadRegister fetches and returns the 8-bit data value from a given 16-bit register address.
type ICIA interface {
	Setup() error

	Bind(conn ICIASocket) error

	Connect() error

	Reset()

	Emulate()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

// IdICIA generates a unique identifier string for the given ICIA interface, label, and instance number.
func IdICIA(v ICIA, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToICIA converts an IComponent to an ICIA interface if compatible, or returns an error if not.
func ComponentToICIA(component IComponent) (ICIA, error) {
	if component == nil {
		return nil, fmt.Errorf("component ICIA is nil")
	}
	v, ok := component.(ICIA)
	if !ok {
		return nil, fmt.Errorf("component is not a ICIA")
	}
	return v, nil
}
