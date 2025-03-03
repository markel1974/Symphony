package mos6569

import "github.com/markel1974/c64emu/src/components/board"

// IBanks defines an interface for handling memory read operations in different regions.
// It provides methods to read from character ROM, direct memory, and color memory.
// The specific implementation determines how the data is retrieved from each region.
type IBanks interface {
	ReadCharRom(uint16) uint8
	ReadDirect(uint16) uint8
	ReadColor(uint16) uint8
}

// ISocket represents an interface for handling system cycle interactions and peripheral communications.
type ISocket interface {
	Cycle() uint64
	GetDisplayBuffer() board.IDisplayBuffer
	GetBanks() IBanks
	IRQTrigger()
	IRQClear()
	BALow(d bool)
	AECLow(d bool)
	VBlank()
	LastCycle()
}
