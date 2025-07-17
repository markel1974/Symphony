package references

import (
	"fmt"
)

// IMos6526Socket defines the interface for socket communication with the MOS 6526 CIA chip for emulation or integration.
// ReadPortA reads and returns the current state of Port A (8-bit).
// ReadPortB reads and returns the current state of Port B (8-bit).
// ReadSP reads and returns the state of the Serial Port input line.
// SignalSP sets the level of the Serial Port input line based on the provided boolean value.
// SignalPRA sets the value of Port A (8-bit) to the specified value.
// SignalPRB sets the value of Port B (8-bit) to the specified value.
// SignalDDRA sets the Data Direction Register for Port A (8-bit).
// SignalDDRB sets the Data Direction Register for Port B (8-bit).
// IRQTrigger signals an Interrupt Request (IRQ) to be triggered.
// IRQClearTrigger clears any triggered Interrupt Request (IRQ).
type IMos6526Socket interface {
	ReadPortA(prA uint8, prB uint8, ddrA uint8, ddrB uint8) uint8

	ReadPortB(prA uint8, prB uint8, ddrA uint8, ddrB uint8) uint8

	ReadSP() bool

	SignalSP(level bool)

	SignalPRA(prA uint8)

	SignalPRB(prB uint8)

	SignalDDRA(ddrA uint8)

	SignalDDRB(ddrB uint8)

	IRQTrigger()

	IRQClearTrigger()
}

// IMos6526 represents the interface for the MOS 6526 CIA (Complex Interface Adapter) chip used in emulation contexts.
// Setup initializes the MOS 6526 instance for operation, returning an error if initialization fails.
// Bind links the MOS 6526 instance to a socket connection, facilitating communication and data exchange.
// Connect establishes the runtime connections necessary for the MOS 6526 to interact with other system components.
// Reset reinitializes the MOS 6526, clearing its internal state and returning it to its default configuration.
// Emulate performs the simulation or emulation of the chip's behavior based on its current state and inputs.
// Update triggers an internal update of the chip's state and outputs in response to the latest inputs.
// WriteRegister writes an 8-bit value to a specified 16-bit address in the chip's register space.
// ReadRegister retrieves an 8-bit value from a specified 16-bit address in the chip's register space.
// ReadPRA reads and returns the 8-bit value stored in Port A Data Register (PRA).
// ReadPRB reads and returns the 8-bit value stored in Port B Data Register (PRB).
// ReadDDRA retrieves the 8-bit value from the Data Direction Register for Port A (DDRA).
// ReadDDRB retrieves the 8-bit value from the Data Direction Register for Port B (DDRB).
type IMos6526 interface {
	Setup() error

	Bind(conn IMos6526Socket) error

	Connect() error

	Reset()

	Emulate()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8

	ReadPRA() uint8

	ReadPRB() uint8

	ReadDDRA() uint8

	ReadDDRB() uint8
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
