package references

import (
	"fmt"
)

// IRamC1541Socket represents an interface for binding and integrating a RAM instance specific to C1541 functionality.
type IRamC1541Socket interface {
}

// IRamC1541 represents an interface for managing a RAM module in a C1541 drive emulation system.
// Setup initializes the RAM module and prepares it for interaction.
// Bind associates the RAM module with a provided IRamC1541Socket.
// Connect establishes the necessary connections for the RAM module to operate within the system.
// Reset reinitializes the RAM module to its default state.
// Read retrieves an 8-bit value from the specified memory address in the RAM module.
// Write stores an 8-bit value at the specified memory address in the RAM module.
// Size returns the total size of the RAM module in bytes.
type IRamC1541 interface {
	Setup() error

	Bind(socket IRamC1541Socket) error

	Connect() error

	Reset()

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	Size() int
}

// IdIRamC1541 generates a unique hardware ID for an IRamC1541 component using a label, instance number, and reference.
func IdIRamC1541(v IRamC1541, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIRamC1541 converts an IComponent to an IRamC1541, ensuring type compatibility and returning an error if invalid.
func ComponentToIRamC1541(component IComponent) (IRamC1541, error) {
	if component == nil {
		return nil, fmt.Errorf("component IRamC1541 is nil")
	}
	v, ok := component.(IRamC1541)
	if !ok {
		return nil, fmt.Errorf("component is not a IRamC1541")
	}
	return v, nil
}
