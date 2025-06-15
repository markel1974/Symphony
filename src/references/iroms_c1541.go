package references

import (
	"fmt"
)

// IRomsC1541Socket defines an interface for connecting and interacting with a ROM loader mechanism.
type IRomsC1541Socket interface {
}

// IRomsC1541 represents an interface for loading and managing ROM data in a C1541 drive emulation system.
// Setup initializes the ROM loader, preparing it for operation within the system.
// Bind associates the ROM loader with a specified socket, enabling integration with the system.
// Connect establishes connections required for the ROM loader to function properly.
// Load retrieves the ROM data as a byte slice from the associated storage medium.
type IRomsC1541 interface {
	Setup() error

	Bind(rom IRomsC1541Socket) error

	Connect() error

	KernalRead(addr uint16) uint8
}

// IdIRomsC1541 generates a unique hardware identifier for an IRomsC1541 instance based on the label and instance number.
func IdIRomsC1541(v IRomsC1541, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIRomsC1541 converts an IComponent into an IRomsC1541 if possible, returning an error otherwise.
func ComponentToIRomsC1541(component IComponent) (IRomsC1541, error) {
	if component == nil {
		return nil, fmt.Errorf("component IRomsC1541 is nil")
	}
	v, ok := component.(IRomsC1541)
	if !ok {
		return nil, fmt.Errorf("component is not a IRomsC1541")
	}
	return v, nil
}
