package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIROMLoaderC1541(_ IROMLoaderC1541, label string, instance int) string {
	return IdInternalComponent(label, instance, "IROMLoaderC1541")
}

type IROMLoaderC1541Socket interface {
}

// IROMLoaderC1541 is an interface for handling ROM loading functionality specific to the C1541 drive emulation.
// Setup configures the ROM loader using the provided configuration.
// Load retrieves the raw byte data of the ROM.
type IROMLoaderC1541 interface {
	Setup(cc map[string]IComponent, cfg *config.Config) error

	Bind(rom IROMLoaderC1541Socket) error

	Connect() error

	Load() []byte
}

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

func ComponentsToIROMLoaderC1541(cc map[string]IComponent, label string, instance int) (IROMLoaderC1541, error) {
	id := IdIROMLoaderC1541(nil, label, instance)
	c, err := ComponentToIROMLoaderC1541(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
