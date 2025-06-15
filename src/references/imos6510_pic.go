package references

import (
	"fmt"
)

// OpFlagIrqDisabled represents a flag indicating that interrupts are disabled.
// OpFlagIrqEnabled represents a flag indicating that interrupts are enabled.
// OpFlagIntDelayed represents a flag indicating that interrupts are delayed.
const (
	OpFlagIrqDisabled = 0x01
	OpFlagIrqEnabled  = 0x02
	OpFlagIntDelayed  = 0x04
)

// IMos6510PicSocket is an interface used to establish communication between a PIC and its corresponding socket.
type IMos6510PicSocket interface {
}

// IMos6510Pic defines an interface for managing a programmable interrupt controller (PIC) in a 6510 CPU simulation.
// Reset reinitializes the PIC to its default state.
// Setup initializes the PIC, returning an error if setup fails.
// Bind binds the PIC to a socket and quartz clock, returning an error if binding fails.
// Connect establishes necessary connections for the PIC, making it operational.
// ClearIRQ clears an interrupt request associated with the specified IRQ vector.
// TriggerIRQ generates an interrupt request for the specified IRQ vector.
// TriggerReset forces a reset operation in the PIC.
// TriggerNMI forces a non-maskable interrupt (NMI) in the PIC.
// VerifyIrq verifies the specified IRQ vector states, returning the result.
// ClearNMI clears the non-maskable interrupt (NMI) state.
// HasNMI checks if a non-maskable interrupt (NMI) is currently active, returning a boolean.
type IMos6510Pic interface {
	Reset()

	Setup() error

	Bind(socket IMos6510PicSocket, quartz IQuartz) error

	Connect() error

	ClearIRQ(uint32)

	TriggerIRQ(uint32)

	TriggerReset()

	TriggerNMI()

	VerifyIrq(uint8, uint8) uint8

	ClearNMI()

	HasNMI() bool
}

// IdIMos6510Pic generates a unique identifier for an IMos6510Pic component using the provided label, instance, and interface name.
func IdIMos6510Pic(v IMos6510Pic, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6510Pic attempts to cast an IComponent to an IMos6510Pic interface and returns an error if the cast fails.
func ComponentToIMos6510Pic(component IComponent) (IMos6510Pic, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6510Pic is nil")
	}
	v, ok := component.(IMos6510Pic)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6510Pic")
	}
	return v, nil
}
