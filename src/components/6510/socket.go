package mos6510

// IBanks represents an interface for managing read and write access to memory banks in a system or CPU simulation.
// Provides methods to read an 8-bit value from a specified 16-bit address and to write an 8-bit value to a 16-bit address.
type IBanks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

// IPic defines an interface for handling programmable interrupt controllers with methods for reset, IRQ, and NMI operations.
type IPic interface {
	Reset()
	VerifyIrq(uint8, uint8) uint8
	ClearNMI()
	HasNMI() bool
}

// ISocket represents an interface for components requiring access to banks (IBanks) and a PIC (IPic).
type ISocket interface {
	GetBanks() IBanks
	GetPic() IPic
}
