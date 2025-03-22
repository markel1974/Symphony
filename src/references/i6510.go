package references

import "fmt"

// I6510Banks represents an interface for managing and accessing memory banks in a system.
// Read retrieve the value from a specified memory address within the bank.
// Write sets the specified value at a given memory address in the bank.
type I6510Banks interface {
	Read(uint16) uint8

	Write(uint16, uint8)
}

func IdI6510(_ I6510) string {
	return "I6510"
}

// I6510 represents the 6510 CPU interface for emulation and interaction with hardware components.
// Reset reinitializes the CPU state to default values.
// Emulate performs a single emulation step of the 6510 CPU.
// Setup configures the CPU to interact with the specified I6510Socket.
// SetRDYLow sets the RDY line state to low or high.
// SetAECLow sets the AEC line state to low or high.
// SetOverflowBranch assigns a callback function for signaling overflow during branch instructions.
type I6510 interface {
	Reset()

	Emulate()

	Setup(socket I6510Socket) error

	SetRDYLow(rdyLow bool)

	SetAECLow(aecLow bool)

	SetOverflowBranch(sob func() bool)
}

// I6510Socket defines an interface for interacting with 6510 PIC and memory banks.
// Provides methods to retrieve the programmable interrupt controller (PIC) and banked memory interface.
// GetBanks retrieves the memory bank interface for read/write operations.
// GetPic retrieves the programmable interrupt controller for managing IRQ/NMI and resets.
type I6510Socket interface {
	GetBanks() I6510Banks

	GetPic() IPIC6510
}

func ComponentToI6510(component IComponent, err error) (I6510, error) {
	if err = ComponentValidate(component, err); err != nil {
		return nil, err
	}
	v, ok := component.(I6510)
	if !ok {
		return nil, fmt.Errorf("component is not a %s", IdI6510(v))
	}
	return v, nil
}
