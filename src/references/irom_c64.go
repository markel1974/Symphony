package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIROMLoaderC64(_ IROMLoaderC64, instance int) string {
	return IdInternalComponent("IROMLoaderC64", instance)
}

// IROMLoaderC64 is an interface that provides methods to load various ROM sections, including Kernal, Basic, and Char ROMs.
// LoadKernal loads the Kernal ROM bytes and returns the data as a slice of bytes.
// LoadBasic loads the Basic ROM bytes and returns the data as a slice of bytes.
// LoadChar loads the Character ROM bytes and returns the data as a slice of bytes.
type IROMLoaderC64 interface {
	Setup(cfg *config.Config) error

	Reset()

	LoadKernal() []byte

	LoadBasic() []byte

	LoadChar() []byte
}

func ComponentToIROMLoaderC64(component IComponent) (IROMLoaderC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component is nil")
	}
	v, ok := component.(IROMLoaderC64)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdIROMLoaderC64(v, 0))
	}
	return v, nil
}

func ComponentsToIROMLoaderC64(cc map[string]IComponent, instance int) (IROMLoaderC64, error) {
	id := IdIROMLoaderC64(nil, instance)
	c, err := ComponentToIROMLoaderC64(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
