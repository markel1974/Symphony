package references

import (
	"fmt"
)

// IMos6526Socket represents an interface to manage input/output operations and IRQ handling for a CIA component.
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
type IMos6526Socket interface {
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

// IMos6526 defines the interface for a CIA (Complex Interface Adapter) component in an emulator or system.
// Setup initializes the CIA for operation, returning an error if initialization fails.
// Bind links the CIA to a socket interface implementing the IMos6526Socket interface.
// Connect establishes the necessary connections for the CIA within the system.
// Reset reinitializes the CIA to its default state.
// Emulate allows the CIA to execute its required operations in synchronization with the overall system.
// Update performs periodic updates or state adjustments for the CIA during emulation.
// WriteRegister writes the specified 8-bit data value to a given 16-bit register address.
// ReadRegister fetches and returns the 8-bit data value from a given 16-bit register address.
type IMos6526 interface {
	Setup() error

	Bind(conn IMos6526Socket) error

	Connect() error

	Reset()

	Emulate()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

// IdIMos6526 generates a unique identifier string for the given IMos6526 interface, label, and instance number.
func IdIMos6526(v IMos6526, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6526 converts an IComponent to an IMos6526 interface if compatible, or returns an error if not.
func ComponentToIMos6526(component IComponent) (IMos6526, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6526 is nil")
	}
	v, ok := component.(IMos6526)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6526")
	}
	return v, nil
}
