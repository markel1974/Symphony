package references

import (
	"fmt"
)

// IROMLoaderSocket defines an interface for components interacting with ROM loader instances.
type IROMLoaderSocket interface {
}

// IROMLoaderC64 represents the interface for managing ROM data loading processes for the C64 system.
// Setup initializes the ROM loader, preparing it for use.
// Bind connects the ROM loader to the provided IROMLoaderSocket interface.
// Connect establishes any necessary runtime connections for the ROM loader.
// Reset reinitializes the ROM loader to its default state.
// LoadKernal retrieves the bytes of the Kernal ROM.
// LoadBasic retrieves the bytes of the Basic ROM.
// LoadChar retrieves the bytes of the Character ROM.
type IROMLoaderC64 interface {
	Setup() error

	Bind(socket IROMLoaderSocket) error

	Connect() error

	Reset()

	LoadKernal() []byte

	LoadBasic() []byte

	LoadChar() []byte
}

// IdIROMLoaderC64 generates a unique identifier for an IROMLoaderC64 instance based on a label and instance number.
func IdIROMLoaderC64(v IROMLoaderC64, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIROMLoaderC64 attempts to cast the provided IComponent to an IROMLoaderC64. Returns an error if the cast fails.
func ComponentToIROMLoaderC64(component IComponent) (IROMLoaderC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component IROMLoaderC64 is nil")
	}
	v, ok := component.(IROMLoaderC64)
	if !ok {
		return nil, fmt.Errorf("component is not a IROMLoaderC64")
	}
	return v, nil
}
