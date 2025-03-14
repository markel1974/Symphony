package references

// IVicBanks defines an interface for handling memory read operations in different regions.
// It provides methods to read from character ROM, direct memory, and color memory.
// The specific implementation determines how the data is retrieved from each region.
type IVicBanks interface {
	ReadCharRom(uint16) uint8

	ReadDirect(uint16) uint8

	ReadColor(uint16) uint8
}

// IVicSocket represents an interface for handling system cycle interactions and peripheral communications.
type IVicSocket interface {
	Cycle() uint64

	GetDisplayBuffer() IDisplayBuffer

	GetBanks() IVicBanks

	IRQTrigger()

	IRQClear()

	BALow(d bool)

	AECLow(d bool)

	VBlank()

	LastCycle()
}
