package references

import (
	"fmt"
)

// I6510Banks represents an interface for reading from and writing to memory banks in the 6510 CPU simulation.
// Read retrieves an 8-bit unsigned value from the specified 16-bit memory address.
// Write stores an 8-bit unsigned value at the specified 16-bit memory address.
type I6510Banks interface {
	Read(uint16) uint8

	Write(uint16, uint8)
}

// I6510Socket provides an interface for facilitating communication and interaction with the 6510 CPU.
type I6510Socket interface {
}

// I6510 defines the interface for a simulation of the 6510 CPU, including setup, binding, connection, and operation methods.
// Setup initializes the 6510 simulation, returning an error if setup fails.
// Bind associates necessary components like socket, PIC, and memory banks to the CPU simulation, returning an error if binding fails.
// Connect finalizes the connection process for the 6510, making it operational in the simulation.
// Reset reinitializes the CPU simulation to its default state.
// Emulate executes one emulation cycle of the 6510 logic, advancing its internal state.
// SetRDYLow alters the RDY line state, taking a boolean to indicate if it should be set low.
// SetAECLow changes the AEC line state, taking a boolean to specify if it should be lowered.
// SetOverflowBranch assigns an overflow branch function, which should return a boolean representing overflow status.
type I6510 interface {
	Setup() error

	Bind(socket I6510Socket, pic IPIC6510, banks I6510Banks) error

	Connect() error

	Reset()

	Emulate()

	SetRDYLow(rdyLow bool)

	SetAECLow(aecLow bool)

	SetOverflowBranch(sob func() bool)
}

// IdI6510 generates a unique identifier for an I6510 component using its label, instance index, and reference identity.
func IdI6510(v I6510, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToI6510 converts a given IComponent to an I6510 type, returning an error if the conversion is invalid or the component is nil.
func ComponentToI6510(component IComponent) (I6510, error) {
	if component == nil {
		return nil, fmt.Errorf("component I6510 is nil")
	}
	v, ok := component.(I6510)
	if !ok {
		return nil, fmt.Errorf("component is not a I6510")
	}
	return v, nil
}
