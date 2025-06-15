package references

import (
	"fmt"
)

// IC1541PlaSocket defines an interface that serves as a connection point for integrating components with the PLA system.
type IC1541PlaSocket interface {
}

// IC1541Pla represents an interface for managing and simulating the Programmable Logic Array (PLA) in a C1541 system.
// Setup initializes the PLA and prepares it for operation.
// Bind connects the PLA to compatible components, including VIAs, RAM, and ROM loader, for functional integration.
// Connect establishes necessary runtime connections for the PLA to interact within the system.
// Read retrieves an 8-bit value from the specified memory address in the PLA's memory map.
// Write stores an 8-bit value at the specified memory address in the PLA's memory map.
type IC1541Pla interface {
	Setup() error

	Bind(socket IC1541PlaSocket, via1 IMos6522, via2 IMos6522, ram IC1541Ram, romLoader IC1541Roms) error

	Connect() error

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)
}

// IdIC1541Pla generates a unique identifier string for an IC1541Pla component based on its label, instance, and interface name.
func IdIC1541Pla(v IC1541Pla, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC1541Pla converts an IComponent to an IC1541Pla instance, returning an error if the conversion fails.
func ComponentToIC1541Pla(component IComponent) (IC1541Pla, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC1541Pla is nil")
	}
	v, ok := component.(IC1541Pla)
	if !ok {
		return nil, fmt.Errorf("component is not a IC1541Pla")
	}
	return v, nil
}
