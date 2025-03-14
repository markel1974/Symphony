package references

// I6510Banks represents an interface for managing read and write access to memory banks in a system or CPU simulation.
// Provides methods to read an 8-bit value from a specified 16-bit address and to write an 8-bit value to a 16-bit address.
type I6510Banks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

// I6510Pic defines an interface for handling programmable interrupt controllers with methods for reset, IRQ, and NMI operations.
type I6510Pic interface {
	Reset()
	VerifyIrq(uint8, uint8) uint8
	ClearNMI()
	HasNMI() bool
}

// I6510Socket represents an interface for components requiring access to banks (I6510Banks) and a PIC (I6510Pic).
type I6510Socket interface {
	GetBanks() I6510Banks
	GetPic() I6510Pic
}
