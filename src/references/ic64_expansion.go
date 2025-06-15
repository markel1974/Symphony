package references

import (
	"fmt"
)

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

// IC64Expansion defines an interface for managing expansion hardware behavior and state in the C64 system.
type IC64Expansion interface {
	Cycle() uint64

	CycleAlarm(string, QuartzAlarmCallback) IQuartzAlarm

	Read(uint16) uint8

	Write(uint16, uint8)

	RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int

	RamRemoveWriteTrigger(addr uint16, id int)

	ResetTrigger()

	NMITrigger()

	IRQTrigger()

	IRQClearTrigger()

	SetDMALow(bool)

	BusAvailable() bool

	AECAvailable() bool //TODO NOT STANDARD

	GameExRomConfigChanged()

	RmwFlags() uint8 //TODO NOT STANDARD

	LedActivity(uint8, bool)
}

// IC64CartridgeManagerSocket represents an interface for managing socket interactions in a C64 cartridge system.
type IC64CartridgeManagerSocket interface {
}

// IC64CartridgeManager defines an interface for managing C64 cartridges and their interactions within an emulator system.
type IC64CartridgeManager interface {
	Setup() error

	Bind(socket IC64CartridgeManagerSocket, expansion IC64Expansion) error

	Config() (uint8, uint8, bool)

	Connect() error

	CreateCartridges() error

	HardwareButton(pressed bool, value uint8)

	Read(interval C64RomInterval, addr uint16) (uint8, bool)

	Write(interval C64RomInterval, addr uint16, data uint8) bool

	IORead(addr uint16) (uint8, bool)

	IOWrite(addr uint16, data uint8) bool

	IRQ(d uint32)

	IRQClear(d uint32)

	Reset()

	Emulate()

	Add(kind string, name string, data []uint8) (string, error)
}

// IdIC64CartridgeManager generates a unique identifier for an IC64CartridgeManager interface using the given label and instance.
func IdIC64CartridgeManager(v IC64CartridgeManager, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64CartridgeManager converts an IComponent to an IC64CartridgeManager, returning an error if the cast fails.
func ComponentToIC64CartridgeManager(component IComponent) (IC64CartridgeManager, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64CartridgeManager is nil")
	}
	v, ok := component.(IC64CartridgeManager)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64CartridgeManager")
	}
	return v, nil
}
