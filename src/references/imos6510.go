package references

import (
	"fmt"
)

// IMos6510Banks represents an interface for reading from and writing to memory banks in the 6510 CPU simulation.
// Read retrieves an 8-bit unsigned value from the specified 16-bit memory address.
// Write stores an 8-bit unsigned value at the specified 16-bit memory address.
type IMos6510Banks interface {
	Read(uint16) uint8

	Write(uint16, uint8)
}

// IMos6510Socket provides an interface for facilitating communication and interaction with the 6510 CPU.
type IMos6510Socket interface {
}

// IMos6510 defines an interface for emulating a MOS 6510 CPU and managing its operations and interactions.
type IMos6510 interface {
	Setup() error

	Bind(socket IMos6510Socket, quartz IQuartz, banks IMos6510Banks) error

	Connect() error

	Reset()

	Emulate()

	SetRDYLow(rdyLow bool)

	SetAECLow(aecLow bool)

	SetOverflowBranch(sob func() bool)

	TriggerReset()

	TriggerIRQ(uint32)

	ClearIRQ(uint32)

	TriggerNMI()

	ClearNMI()
}

// IdIMos6510 generates a unique identifier for an IMos6510 component using its label, instance index, and reference identity.
func IdIMos6510(v IMos6510, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6510 converts a given IComponent to an IMos6510 type, returning an error if the conversion is invalid or the component is nil.
func ComponentToIMos6510(component IComponent) (IMos6510, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6510 is nil")
	}
	v, ok := component.(IMos6510)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6510")
	}
	return v, nil
}
