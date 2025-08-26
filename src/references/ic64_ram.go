package references

import (
	"fmt"
)

// IC64RamSocket defines an interface representing a socket interface for C64 RAM-related operations or communication.
type IC64RamSocket interface {
}

// IC64Ram provides an interface for implementing a memory module for the C64 system with read, write, and setup functions.
// Setup initializes the memory module.
// Bind associates the memory module with an IC64RamSocket instance.
// Connect finalizes the connection of the memory module.
// Reset reinitializes the memory module to its default state.
// Read retrieves a byte of data from the specified address.
// ReadInterval returns a slice of data representing a read interval.
// Write stores a byte of data at the specified address.
// WriteInterval returns a slice of data representing a write interval.
type IC64Ram interface {
	Setup() error

	Bind(socket IC64RamSocket) error

	Connect() error

	Reset()

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	Size() int
}

// IdIC64Ram generates a unique identifier string for an IC64Ram component using its label, instance Id, and a fixed identifier.
func IdIC64Ram(v IC64Ram, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Ram converts an IComponent to an IC64Ram instance, returning an error if the conversion is invalid.
func ComponentToIC64Ram(component IComponent) (IC64Ram, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Ram is nil")
	}
	v, ok := component.(IC64Ram)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Ram")
	}
	return v, nil
}
