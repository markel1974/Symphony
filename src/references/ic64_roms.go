package references

import (
	"fmt"
)

// IC64RomsSocket defines an interface for components interacting with ROM loader instances.
type IC64RomsSocket interface {
}

// IC64Roms represents the interface for managing ROM data loading processes for the C64 system.
// Setup initializes the ROM loader, preparing it for use.
// Bind connects the ROM loader to the provided IC64RomsSocket interface.
// Connect establishes any necessary runtime connections for the ROM loader.
// Reset reinitializes the ROM loader to its default state.
// LoadKernal retrieves the bytes of the Kernal ROM.
// LoadBasic retrieves the bytes of the Basic ROM.
// LoadChar retrieves the bytes of the Character ROM.
type IC64Roms interface {
	Setup() error

	Bind(socket IC64RomsSocket) error

	Connect() error

	Reset()

	KernalRead(addr uint16) uint8

	BasicRead(addr uint16) uint8

	CharRead(addr uint16) uint8
}

// IdIC64Roms generates a unique identifier for an IC64Roms instance based on a label and instance number.
func IdIC64Roms(v IC64Roms, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Roms attempts to cast the provided IComponent to an IC64Roms. Returns an error if the cast fails.
func ComponentToIC64Roms(component IComponent) (IC64Roms, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Roms is nil")
	}
	v, ok := component.(IC64Roms)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Roms")
	}
	return v, nil
}
