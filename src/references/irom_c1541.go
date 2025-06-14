package references

import (
	"fmt"
)

// IROMLoaderC1541Socket defines an interface for connecting and interacting with a ROM loader mechanism.
type IROMLoaderC1541Socket interface {
}

// IROMLoaderC1541 represents an interface for loading and managing ROM data in a C1541 drive emulation system.
// Setup initializes the ROM loader, preparing it for operation within the system.
// Bind associates the ROM loader with a specified socket, enabling integration with the system.
// Connect establishes connections required for the ROM loader to function properly.
// Load retrieves the ROM data as a byte slice from the associated storage medium.
type IROMLoaderC1541 interface {
	Setup() error

	Bind(rom IROMLoaderC1541Socket) error

	Connect() error

	Load() []byte
}

// IdIROMLoaderC1541 generates a unique hardware identifier for an IROMLoaderC1541 instance based on the label and instance number.
func IdIROMLoaderC1541(v IROMLoaderC1541, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIROMLoaderC1541 converts an IComponent into an IROMLoaderC1541 if possible, returning an error otherwise.
func ComponentToIROMLoaderC1541(component IComponent) (IROMLoaderC1541, error) {
	if component == nil {
		return nil, fmt.Errorf("component IROMLoaderC1541 is nil")
	}
	v, ok := component.(IROMLoaderC1541)
	if !ok {
		return nil, fmt.Errorf("component is not a IROMLoaderC1541")
	}
	return v, nil
}
