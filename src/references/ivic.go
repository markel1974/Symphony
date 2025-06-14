package references

import (
	"fmt"
)

// IVIC defines an interface for emulating a VIC component with methods for setup, interaction, state management, and operations.
type IVIC interface {
	Setup() error

	Bind(socket IVICSocket) error

	Connect() error

	Reset()

	Emulate()

	GetBALow() bool

	GetAECLow() bool

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8

	LightPenTrigger()

	GetText() []byte

	ChangedVA(newVA uint8)

	GetLastByte() uint8
}

// IVICBanks provides methods to interact with various VIC memory regions, including character ROM and color memory.
// ReadCharRom reads a byte from the character ROM at the specified address.
// ReadDirect reads a byte from direct memory at the specified address.
// ReadColor retrieves a color byte from VIC memory at the given address.
type IVICBanks interface {
	ReadCharRom(uint16) uint8

	ReadDirect(uint16) uint8

	ReadColor(uint16) uint8
}

// IVICSocket defines the interface for VIC's connection socket, providing methods for synchronization and signal control.
// Cycle returns the current cycle count of the VIC socket.
// GetBanks retrieves the associated IVICBanks for memory access operations.
// IRQTrigger signals an interrupt request trigger on the VIC socket.
// IRQClearTrigger clears an active interrupt request signal on the VIC socket.
// BALow sets or unsets the state of the BA (Bus Available) line.
// AECLow sets or unsets the state of the AEC (Address Enable Control) line.
// VBlank notifies the VIC to enter vertical blanking interval.
// LastCycle is invoked to finalize operations at the last cycle of the frame.
type IVICSocket interface {
	Cycle() uint64

	GetBanks() IVICBanks

	IRQTrigger()

	IRQClearTrigger()

	BALow(d bool)

	AECLow(d bool)

	VBlank()

	LastCycle()
}

// IdIVIC generates a unique identifier for an IVIC interface based on the given label, instance, and interface name.
func IdIVIC(v IVIC, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIVIC converts an IComponent to an IVIC if the provided component matches the IVIC interface type.
// Returns an IVIC instance and nil error if successful, or nil and an error if the conversion fails.
// An error is also returned if the input component is nil.
func ComponentToIVIC(component IComponent) (IVIC, error) {
	if component == nil {
		return nil, fmt.Errorf("component IVIC is nil")
	}
	v, ok := component.(IVIC)
	if !ok {
		return nil, fmt.Errorf("component is not a IVIC")
	}
	return v, nil
}
