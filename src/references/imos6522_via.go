package references

import (
	"fmt"
)

// IMos6522Socket represents an interface for socket communication and control in a system with read/write and IRQ operations.
// ReadPRA reads the value of port A register at the specified indexes.
// ReadPRB reads the value of port B register at the specified indexes.
// WritePRA writes a value to the port A register at the specified indexes.
// WritePRB writes a value to the port B register at the specified indexes.
// WriteDDRA writes a value to the data direction register A at the specified indexes.
// WriteDDRB writes a value to the data direction register B at the specified indexes.
// WriteCA2 writes a boolean control value for the CA2 pin.
// WriteCB2 writes a boolean control value for the CB2 pin.
// IRQClearTrigger clears the IRQ trigger state.
// IRQTrigger triggers an IRQ signal.
type IMos6522Socket interface {
	ReadPRA(uint8, uint8) uint8

	ReadPRB(uint8, uint8) uint8

	WritePRA(uint8, uint8)

	WritePRB(uint8, uint8)

	WriteDDRA(uint8, uint8)

	WriteDDRB(uint8, uint8)

	WriteCA2(bool)

	WriteCB2(bool)

	IRQClearTrigger()

	IRQTrigger()
}

// IMos6522 defines an interface for a Versatile Interface Adapter (VIA) with setup, binding, and signaling functionalities.
// Setup initializes the VIA and prepares it for operation.
// Bind connects the VIA to a given IMos6522Socket for external interactions.
// Connect establishes any required connections after initialization.
// Reset reinitializes the VIA to its default state.
// Emulate processes emulation cycles for the VIA.
// ReadByte reads an 8-bit value from the specified memory address.
// WriteByte writes an 8-bit value to the specified memory address.
// SignalPRA triggers a signal specific to Port A.
// SignalPRB triggers a signal specific to Port B.
type IMos6522 interface {
	Setup() error

	Bind(conn IMos6522Socket) error

	Connect() error

	Reset()

	Emulate()

	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)

	SignalPRA()

	SignalPRB()

	//ByteReady() bool
}

// IdIMos6522 generates a unique identifier for a VIA component based on the provided label, instance, and interface name.
func IdIMos6522(v IMos6522, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6522 converts an IComponent instance to an IMos6522 if possible or returns an error if the conversion fails.
func ComponentToIMos6522(component IComponent) (IMos6522, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6522 is nil")
	}
	v, ok := component.(IMos6522)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6522")
	}
	return v, nil
}
