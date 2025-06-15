package references

import (
	"fmt"
)

// IPLAc1541Socket defines an interface that serves as a connection point for integrating components with the PLA system.
type IPLAc1541Socket interface {
}

// IPLAc1541 represents an interface for managing and simulating the Programmable Logic Array (PLA) in a C1541 system.
// Setup initializes the PLA and prepares it for operation.
// Bind connects the PLA to compatible components, including VIAs, RAM, and ROM loader, for functional integration.
// Connect establishes necessary runtime connections for the PLA to interact within the system.
// Read retrieves an 8-bit value from the specified memory address in the PLA's memory map.
// Write stores an 8-bit value at the specified memory address in the PLA's memory map.
type IPLAc1541 interface {
	Setup() error

	Bind(socket IPLAc1541Socket, via1 IVIA, via2 IVIA, ram IRamC1541, romLoader IRomsC1541) error

	Connect() error

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)
}

// IdIPLAc1541 generates a unique identifier string for an IPLAc1541 component based on its label, instance, and interface name.
func IdIPLAc1541(v IPLAc1541, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIPLAc1541 converts an IComponent to an IPLAc1541 instance, returning an error if the conversion fails.
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
