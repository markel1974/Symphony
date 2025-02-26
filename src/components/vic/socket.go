package mos6569

// IBanks defines an interface for handling memory read operations in different regions.
// It provides methods to read from character ROM, direct memory, and color memory.
// The specific implementation determines how the data is retrieved from each region.
type IBanks interface {
	ReadCharRom(uint16) uint8
	ReadDirect(uint16) uint8
	ReadColor(uint16) uint8
}

// IDisplayBuffer is an interface for managing display buffer operations in graphical rendering systems.
// Set sets a single value at the specified index in the display buffer.
// SetMulti8 sets a multi-bit value at the specified index in the display buffer.
// Set8 sets an array of 8-bit values starting at the specified index in the display buffer.
type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}

// ISocket represents an interface for handling system cycle interactions and peripheral communications.
type ISocket interface {
	Cycle() uint64
	GetDisplayBuffer() IDisplayBuffer
	GetBanks() IBanks
	IRQTrigger()
	IRQClear()
	BALow(d bool)
	AECLow(d bool)
	VBlank()
	LastCycle()
}
