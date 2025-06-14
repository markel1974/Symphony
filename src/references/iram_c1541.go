package references

import (
	"fmt"
)

// IdIRamC1541 generates a unique identifier for an IRamC1541 component based on a label, instance, and constant ID string.
func IdIRamC1541(_ IRamC1541, label string, instance int) string {
	return IdInternalComponent(label, instance, "IdIRamC1541")
}

// IRamC1541Socket represents an interface for interacting with a RAM C1541 socket structure or component.
type IRamC1541Socket interface {
}

// IRamC1541 defines the interface for interacting with a RAM module compatible with the C1541 system.
// Setup initializes the RAM module.
// Bind associates the RAM module with an IRamC1541Socket instance for communication.
// Connect establishes the connection for the RAM module with other components.
// Reset clears or prepares the RAM to its initial state.
// Read retrieves the data stored at the specified memory address.
// Write stores data at the specified memory address in the RAM.
// Size provides the total memory size of the RAM module.
type IRamC1541 interface {
	Setup() error

	Bind(socket IRamC1541Socket) error

	Connect() error

	Reset()

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	Size() int
}

// ComponentToIRamC1541 converts an IComponent to an IRamC1541 and returns an error if the conversion is invalid or unsupported.
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
