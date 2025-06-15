package references

import (
	"fmt"
)

// IMos6581Socket is an interface representing a socket for communication with an IMos6581 implementation.
type IMos6581Socket interface {
}

// IMos6581 defines an interface for interacting with a SID (Sound Interface Device) component.
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
type IMos6581 interface {
	Setup() error

	Bind(socket IMos6581Socket, fragFreq int, rasters int) error

	Connect() error

	Reset()

	SetPotX(x uint8)

	SetPotY(y uint8)

	Prepare()

	Update()

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8
}

// IdIMos6581 generates a unique identifier string for an IMos6581 component using its label, instance number, and interface name.
func IdIMos6581(v IMos6581, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6581 converts an IComponent into an IMos6581 if possible, returning an error if the conversion fails.
func ComponentToIMos6581(component IComponent) (IMos6581, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6581 is nil")
	}
	v, ok := component.(IMos6581)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6581")
	}
	return v, nil
}
