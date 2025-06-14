package references

import (
	"fmt"
)

// ISIDSocket is an interface representing a socket for communication with an ISID implementation.
type ISIDSocket interface {
}

// ISID defines an interface for interacting with a SID (Sound Interface Device) component.
// Setup initializes the SID component, preparing it for operation.
// Bind attaches the SID to a socket with specified frame and raster frequencies.
// Connect establishes any required connections for the SID to function in a system.
// Reset reinitializes the SID to its default state.
// SetPotX adjusts the X-axis potentiometer for the SID to a specific 8-bit value.
// SetPotY adjusts the Y-axis potentiometer for the SID to a specific 8-bit value.
// Prepare prepares the SID for updates or operational adjustments.
// Update processes any changes or updates required for the SID's functioning.
// WriteRegister writes an 8-bit value to a specified 16-bit memory address in the SID.
// ReadRegister reads an 8-bit value from a specified 16-bit memory address in the SID.
type ISID interface {
	Setup() error

	Bind(socket ISIDSocket, fragFreq int, rasters int) error

	Connect() error

	Reset()

	SetPotX(x uint8)

	SetPotY(y uint8)

	Prepare()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

// IdISID generates a unique identifier string for an ISID component using its label, instance number, and interface name.
func IdISID(v ISID, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToISID converts an IComponent into an ISID if possible, returning an error if the conversion fails.
func ComponentToISID(component IComponent) (ISID, error) {
	if component == nil {
		return nil, fmt.Errorf("component ISID is nil")
	}
	v, ok := component.(ISID)
	if !ok {
		return nil, fmt.Errorf("component is not a ISID")
	}
	return v, nil
}
