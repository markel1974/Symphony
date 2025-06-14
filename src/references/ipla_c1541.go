package references

import (
	"fmt"
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
	Setup() error

	Bind(socket IPLAc1541Socket, via1 IVIA, via2 IVIA, ram IRamC1541, romLoader IROMLoaderC1541) error

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
