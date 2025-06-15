package references

import (
	"fmt"
)

// IRamC64Socket defines an interface representing a socket interface for C64 RAM-related operations or communication.
type IRamC64Socket interface {
}

// IRamC64 provides an interface for implementing a memory module for the C64 system with read, write, and setup functions.
// Setup initializes the memory module.
// Bind associates the memory module with an IRamC64Socket instance.
// Connect finalizes the connection of the memory module.
// Reset reinitializes the memory module to its default state.
// Read retrieves a byte of data from the specified address.
// ReadInterval returns a slice of data representing a read interval.
// Write stores a byte of data at the specified address.
// WriteInterval returns a slice of data representing a write interval.
type IRamC64 interface {
	Setup() error

	Bind(socket IRamC64Socket) error

	Connect() error

	Reset()

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	Size() int

	ReadColor(addr uint16) uint8

	WriteColor(addr uint16, data uint8)
}

// IdIRamC64 generates a unique identifier string for an IRamC64 component using its label, instance ID, and a fixed identifier.
func IdIRamC64(v IRamC64, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIRamC64 converts an IComponent to an IRamC64 instance, returning an error if the conversion is invalid.
func ComponentToIRamC64(component IComponent) (IRamC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component IRamC64 is nil")
	}
	v, ok := component.(IRamC64)
	if !ok {
		return nil, fmt.Errorf("component is not a IRamC64")
	}
	return v, nil
}
