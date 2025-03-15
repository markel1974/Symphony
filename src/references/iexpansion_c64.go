package references

//OFF
//GAME = 1, EXROM = 1

//8K Cartridge, $8000-$9FFF (ROML).
//GAME = 1, EXROM = 0
//ROML is read only. Basic ROM and Kernal ROM are available.

//16K Cartridge, $8000-$9FFF / $A000-$BFFF (ROML / ROMH).
//GAME = 0, EXROM = 0
//ROML/ROMH are read only, Basic ROM is overwritten by ROMH.

//16K Cartridge, $8000-$9FFF / $E000-$FFFF (ROML / ROMH). Ultimax mode.
//GAME = 0, EXROM = 1
//Ultimax mode is an emulation of the Japanese CBM machine called “MAX”. It is a predecessor of the C64 with less RAM. In Ultimax mode ROMH replaces the kernal at $E000. You do not need ROML for a cartridge to function and can be left out.

// https://www.c64-wiki.com/wiki/Expansion_Port

// IExpansionC64 defines an interface for C64 expansion boards, including memory access and hardware interaction methods.
// Cycle returns the current cycle count.
// CycleAlarm registers a cycle-based alarm with a callback and returns an alarm interface.
// GameExRomConfigChanged notifies the board of a Game/ExRom configuration change.
// Read reads a byte from a given memory address.
// Write writes a byte to a given memory address.
// RamSetWriteTrigger assigns a write trigger function to a memory address, returning the trigger ID.
// RamRemoveWriteTrigger removes a write trigger from a memory address by ID.
// ResetTrigger sends a hardware reset signal to the board.
// NMITrigger triggers a non-maskable interrupt.
// IRQTrigger triggers an interrupt request.
// IRQClear clears an active interrupt request.
// IRQTriggerBind binds a callback function to IRQ triggering events.
// IRQClearBind binds a callback function to IRQ clearing events.
// SetDMALow sets the DMA (Direct Memory Access) low state.
// BusAvailable checks whether the system bus is available for access.
// AECAvailable checks whether an AEC (Advanced Expansion Controller) is available.
// RmwFlags retrieves flags for read-modify-write operations (non-standard).
type IExpansionC64 interface {
	Cycle() uint64

	CycleAlarm(string, QuartzAlarmCallback) IQuartzAlarm

	GameExRomConfigChanged()

	Read(uint16) uint8

	Write(uint16, uint8)

	RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int

	RamRemoveWriteTrigger(addr uint16, id int)

	ResetTrigger()

	NMITrigger()

	IRQTrigger()

	IRQClear()

	IRQTriggerBind(fn func(uint32))

	IRQClearBind(fn func(uint32))

	SetDMALow(bool)

	BusAvailable() bool

	AECAvailable() bool //TODO NOT STANDARD

	RmwFlags() uint8 //TODO NOT STANDARD
}

// IExpansionSocketC64 defines an interface for managing expansion socket operations and interactions in a given system.
// Config retrieves the configuration data of the expansion socket and additional state details.
// Read fetches a value from a specified ROM interval and address, returning the value and success state.
// IORead performs a read operation from an I/O address, returning the value and success state.
// IOWrite executes a write operation to an I/O address with the specified data, indicating success.
type IExpansionSocketC64 interface {
	Config() (uint8, uint8, bool)

	Read(interval RomInterval, addr uint16) (uint8, bool)

	IORead(addr uint16) (uint8, bool)

	IOWrite(addr uint16, data uint8) bool
}
