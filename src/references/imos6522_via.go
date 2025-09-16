package references

import (
	"fmt"
)

// IMos6522Socket defines the interface for interacting with a MOS 6522 VIA socket in a system.
// ReadPortA reads the current state of port A.
// ReadPortB reads the current state of port B.
// ReadCA1 returns the logic state of control line CA1.
// ReadCB1 returns the logic state of control line CB1.
// ReadCB2 returns the logic state of control line CB2.
// ReadPB6 returns the logic state of bit 6 of port B.
// SignalPRA sends an updated value to the port A register.
// SignalPRB sends an updated value to the port B register.
// SignalDDRA sends an updated value to the data direction register A.
// SignalDDRB sends an updated value to the data direction register B.
// SignalPCR updates the peripheral control register with a new value.
// IRQClearTrigger clears the current IRQ trigger state.
// IRQTrigger asserts an IRQ trigger to signal an interrupt request.
type IMos6522Socket interface {
	ReadPortA() uint8

	ReadPortB() uint8

	ReadCA1() bool

	ReadCB1() bool

	ReadCB2() bool

	WriteCB2(bool)

	ReadPB6() bool

	SignalPRA(uint8)

	SignalPRB(uint8)

	SignalDDRA(uint8)

	SignalDDRB(uint8)

	SignalPCR(uint8)

	IRQClearTrigger()

	IRQTrigger()
}

// IMos6522 defines an interface for the MOS 6522 VIA (Versatile Interface Adapter) chip emulation.
// Setup initializes the VIA for operation, preparing its components.
// Bind associates the VIA with a given socket for external interactions and signal handling.
// Connect establishes necessary connections for the VIA within the system.
// Reset reinitializes the VIA, clearing registers and resetting its state.
// Emulate performs a single step of VIA emulation logic, processing registers and signals.
// ReadByte retrieves an 8-bit value from a specified register address in the VIA memory map.
// WriteByte writes an 8-bit value to a specified register address in the VIA memory map.
// ReadDDRA returns the current state of the Code Direction Register for Port A (DDRA).
// ReadDDRB returns the current state of the Code Direction Register for Port B (DDRB).
// ReadPRA retrieves the current value of the Port A register (PRA).
// ReadPRB retrieves the current value of the Port B register (PRB).
// ReadACR retrieves the current value of the Auxiliary Control Register (ACR).
// ReadPCR retrieves the current value of the Peripheral Control Register (PCR).
type IMos6522 interface {
	Setup() error

	Bind(conn IMos6522Socket) error

	Connect() error

	Reset()

	Emulate()

	ReadByte(addr uint16) uint8

	WriteByte(addr uint16, data uint8)

	ReadDDRA() uint8

	ReadDDRB() uint8

	ReadPRA() uint8

	ReadPRB() uint8

	ReadACR() uint8

	ReadPCR() uint8
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
