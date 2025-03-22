package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIROMLoaderC1541(_ IROMLoaderC1541) string {
	return "IROMLoaderC1541"
}

// IROMLoaderC1541 is an interface for handling ROM loading functionality specific to the C1541 drive emulation.
// Setup configures the ROM loader using the provided configuration.
// Load retrieves the raw byte data of the ROM.
type IROMLoaderC1541 interface {
	Setup(cfg *config.Config) error

	Load() []byte
}

func ComponentToIROMLoaderC1541(component IComponent, err error) (IROMLoaderC1541, error) {
	if err = ComponentValidate(component, err); err != nil {
		return nil, err
	}
	v, ok := component.(IROMLoaderC1541)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdIROMLoaderC1541(v))
	}
	return v, nil
}
