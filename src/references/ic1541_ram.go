package references

import (
	"fmt"
)

// IC1541RamSocket represents an interface for binding and integrating a RAM instance specific to C1541 functionality.
type IC1541RamSocket interface {
}

// IC1541Ram represents an interface for managing a RAM module in a C1541 drive emulation system.
// Setup initializes the RAM module and prepares it for interaction.
// Bind associates the RAM module with a provided IC1541RamSocket.
// Connect establishes the necessary connections for the RAM module to operate within the system.
// Reset reinitializes the RAM module to its default state.
// Read retrieves an 8-bit value from the specified memory address in the RAM module.
// Write stores an 8-bit value at the specified memory address in the RAM module.
// Size returns the total size of the RAM module in bytes.
type IC1541Ram interface {
	Setup() error

	Bind(socket IC1541RamSocket) error

	Connect() error

	Reset()

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	Size() int
}

// IdIC1541Ram generates a unique hardware Id for an IC1541Ram component using a label, instance number, and reference.
func IdIC1541Ram(v IC1541Ram, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC1541Ram converts an IComponent to an IC1541Ram, ensuring type compatibility and returning an error if invalid.
func ComponentToIC1541Ram(component IComponent) (IC1541Ram, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC1541Ram is nil")
	}
	v, ok := component.(IC1541Ram)
	if !ok {
		return nil, fmt.Errorf("component is not a IC1541Ram")
	}
	return v, nil
}
