package references

import (
	"fmt"
)

// IColorRamC64Socket represents an interface for managing the Color RAM socket in Commodore 64 systems.
type IColorRamC64Socket interface {
}

// IC64ColorRam defines an interface representing the color RAM component for the C64 with memory and control operations.
// Setup initializes the color RAM component, preparing it for use.
// Bind associates the color RAM with a given RAM socket interface for communication.
// Connect establishes any necessary internal connections for the color RAM.
// Reset reinitializes the state of the color RAM to default settings.
// Size returns the size of the color RAM in bytes.
// Read fetches the data at a given address in the color RAM.
// Write stores data at a specified address in the color RAM.
type IC64ColorRam interface {
	Setup() error

	Bind(socket IColorRamC64Socket) error

	Connect() error

	Reset()

	Size() int

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8

	WriteRegister(addr uint16, data uint8)
}

// IdIC64ColorRam generates an identifier string for a C64 RAM color component using its label, instance, and interface name.
func IdIC64ColorRam(v IC64ColorRam, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64ColorRam converts the provided IComponent to an IC64ColorRam if possible, or returns an error if not.
func ComponentToIC64ColorRam(component IComponent) (IC64ColorRam, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64ColorRam is nil")
	}
	v, ok := component.(IC64ColorRam)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64ColorRam")
	}
	return v, nil
}
