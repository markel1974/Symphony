package references

// I6510Banks represents an interface for managing and accessing memory banks in a system.
// Read retrieve the value from a specified memory address within the bank.
// Write sets the specified value at a given memory address in the bank.
type I6510Banks interface {
	Read(uint16) uint8

	Write(uint16, uint8)
}

// I6510Pic provides an interface for managing the programmable interrupt controller (PIC) in a 6510 CPU simulation.
// Reset reinitializes the state of the PIC to default values.
// VerifyIrq determines and returns the type of interrupt request (IRQ) based on specified input conditions.
// ClearNMI clears the non-maskable interrupt (NMI) state within the PIC.
// HasNMI checks if a non-maskable interrupt (NMI) has been triggered and returns a boolean result.
type I6510Pic interface {
	Reset()

	VerifyIrq(uint8, uint8) uint8

	ClearNMI()

	HasNMI() bool
}

// I6510 represents an interface for a 6510 CPU implementation with methods for reset, emulation, and signal control.
type I6510 interface {
	Reset()

	Emulate()

	Setup(socket I6510Socket)

	SetRDYLow(rdyLow bool)

	SetAECLow(aecLow bool)
}

// I6510Socket defines an interface for interacting with 6510 PIC and memory banks.
// Provides methods to retrieve the programmable interrupt controller (PIC) and banked memory interface.
// GetBanks retrieves the memory bank interface for read/write operations.
// GetPic retrieves the programmable interrupt controller for managing IRQ/NMI and resets.
type I6510Socket interface {
	GetBanks() I6510Banks

	GetPic() I6510Pic
}
