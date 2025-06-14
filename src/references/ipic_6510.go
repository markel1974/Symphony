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

// IPIC6510Socket is an interface used to establish communication between a PIC and its corresponding socket.
type IPIC6510Socket interface {
}

// IPIC6510 defines an interface for managing a programmable interrupt controller (PIC) in a 6510 CPU simulation.
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
type IPIC6510 interface {
	Reset()

	Setup() error

	Bind(socket IPIC6510Socket, quartz IQuartz) error

	Connect() error

	ClearIRQ(uint32)

	TriggerIRQ(uint32)

	TriggerReset()

	TriggerNMI()

	VerifyIrq(uint8, uint8) uint8

	ClearNMI()

	HasNMI() bool
}

// IdIPIC6510 generates a unique identifier for an IPIC6510 component using the provided label, instance, and interface name.
func IdIPIC6510(v IPIC6510, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIPIC6510 attempts to cast an IComponent to an IPIC6510 interface and returns an error if the cast fails.
func ComponentToIPIC6510(component IComponent) (IPIC6510, error) {
	if component == nil {
		return nil, fmt.Errorf("component IPIC6510 is nil")
	}
	v, ok := component.(IPIC6510)
	if !ok {
		return nil, fmt.Errorf("component is not a IPIC6510")
	}
	return v, nil
}
