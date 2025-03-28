package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIPLAc1541(_ IPLAc1541, label string, instance int) string {
	return IdInternalComponent(label, instance, "IPLAc1541")
}

type IPLAc1541Socket interface {
}

// IPLAc1541 represents an interface for handling PLA logic in a 1541 disk drive emulation.
// Setup initializes the interface by linking it to VIA components, ROM loader, and configuration data.
// Read retrieves the value from the specified memory address.
// Write writes a value to the specified memory address.
type IPLAc1541 interface {
	Setup(cc map[string]IComponent, cfg *config.Config) error

	Bind(socket IPLAc1541Socket, via1 IVIA, via2 IVIA, romLoader IROMLoaderC1541) error

	Connect() error

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)
}

func ComponentToIPLAc1541(component IComponent) (IPLAc1541, error) {
	if component == nil {
		return nil, fmt.Errorf("component IPLAc1541 is nil")
	}
	v, ok := component.(IPLAc1541)
	if !ok {
		return nil, fmt.Errorf("component is not a IPLAc1541")
	}
	return v, nil
}

func ComponentsToIPLAc1541(cc map[string]IComponent, label string, instance int) (IPLAc1541, error) {
	id := IdIPLAc1541(nil, label, instance)
	c, err := ComponentToIPLAc1541(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
