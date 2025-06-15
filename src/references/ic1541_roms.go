package references

import (
	"fmt"
)

// IC1541RomsSocket defines an interface for connecting and interacting with a ROM loader mechanism.
type IC1541RomsSocket interface {
}

// IC1541Roms represents an interface for loading and managing ROM data in a C1541 drive emulation system.
// Setup initializes the ROM loader, preparing it for operation within the system.
// Bind associates the ROM loader with a specified socket, enabling integration with the system.
// Connect establishes connections required for the ROM loader to function properly.
// Load retrieves the ROM data as a byte slice from the associated storage medium.
type IC1541Roms interface {
	Setup() error

	Bind(rom IC1541RomsSocket) error

	Connect() error

	KernalRead(addr uint16) uint8
}

// IdIC1541Roms generates a unique hardware identifier for an IC1541Roms instance based on the label and instance number.
func IdIC1541Roms(v IC1541Roms, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC1541Roms converts an IComponent into an IC1541Roms if possible, returning an error otherwise.
func ComponentToIC1541Roms(component IComponent) (IC1541Roms, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC1541Roms is nil")
	}
	v, ok := component.(IC1541Roms)
	if !ok {
		return nil, fmt.Errorf("component is not a IC1541Roms")
	}
	return v, nil
}
