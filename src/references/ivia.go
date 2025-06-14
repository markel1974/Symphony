package references

import (
	"fmt"
)

// IVIASocket represents an interface for socket communication and control in a system with read/write and IRQ operations.
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
type IVIASocket interface {
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

// IVIA defines an interface for a Versatile Interface Adapter (VIA) with setup, binding, and signaling functionalities.
// Setup initializes the VIA and prepares it for operation.
// Bind connects the VIA to a given IVIASocket for external interactions.
// Connect establishes any required connections after initialization.
// Reset reinitializes the VIA to its default state.
// Emulate processes emulation cycles for the VIA.
// ReadByte reads an 8-bit value from the specified memory address.
// WriteByte writes an 8-bit value to the specified memory address.
// SignalPRA triggers a signal specific to Port A.
// SignalPRB triggers a signal specific to Port B.
type IVIA interface {
	Setup() error

	Bind(conn IVIASocket) error

	Connect() error

	Reset()

	Emulate()

	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)

	SignalPRA()

	SignalPRB()

	//ByteReady() bool
}

// IdIVIA generates a unique identifier for a VIA component based on the provided label, instance, and interface name.
func IdIVIA(v IVIA, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIVIA converts an IComponent instance to an IVIA if possible or returns an error if the conversion fails.
func ComponentToIVIA(component IComponent) (IVIA, error) {
	if component == nil {
		return nil, fmt.Errorf("component IVIA is nil")
	}
	v, ok := component.(IVIA)
	if !ok {
		return nil, fmt.Errorf("component is not a IVIA")
	}
	return v, nil
}
