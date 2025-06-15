package references

import (
	"fmt"
)

// IRomsC64Socket defines an interface for components interacting with ROM loader instances.
type IRomsC64Socket interface {
}

// IRomsC64 represents the interface for managing ROM data loading processes for the C64 system.
// Setup initializes the ROM loader, preparing it for use.
// Bind connects the ROM loader to the provided IRomsC64Socket interface.
// Connect establishes any necessary runtime connections for the ROM loader.
// Reset reinitializes the ROM loader to its default state.
// LoadKernal retrieves the bytes of the Kernal ROM.
// LoadBasic retrieves the bytes of the Basic ROM.
// LoadChar retrieves the bytes of the Character ROM.
type IRomsC64 interface {
	Setup() error

	Bind(socket IRomsC64Socket) error

	Connect() error

	Reset()

	KernalRead(addr uint16) uint8

	BasicRead(addr uint16) uint8

	CharRead(addr uint16) uint8
}

// IdIRomsC64 generates a unique identifier for an IRomsC64 instance based on a label and instance number.
func IdIRomsC64(v IRomsC64, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIRomsC64 attempts to cast the provided IComponent to an IRomsC64. Returns an error if the cast fails.
func ComponentToIRomsC64(component IComponent) (IRomsC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component IRomsC64 is nil")
	}
	v, ok := component.(IRomsC64)
	if !ok {
		return nil, fmt.Errorf("component is not a IRomsC64")
	}
	return v, nil
}
