package references

import (
	"fmt"
	"io"
)

// RomInterval represents a type used to define distinct memory interval mappings in the ROM of a cartridge system.
type RomInterval int

// ROM_LO represents the low ROM interval value in the RomInterval type.
// ROM_HI_1 represents the first high ROM interval value in the RomInterval type.
// ROM_HI_2 represents the second high ROM interval value in the RomInterval type.
const (
	ROM_LO   = RomInterval(1)
	ROM_HI_1 = RomInterval(2)
	ROM_HI_2 = RomInterval(4)
)

// CartridgeMode represents the operational mode of a cartridge, typically defining its memory layout and behavior.
type CartridgeMode int

// CartridgeMode16K represents the 16KB cartridge mode.
// CartridgeMode8K represents the 8KB cartridge mode.
// CartridgeModeUltimax represents the Ultimax cartridge mode.
// CartridgeModeOff represents a state where the cartridge is turned off.
const (
	CartridgeMode16K = CartridgeMode(iota)
	CartridgeMode8K
	CartridgeModeUltimax
	CartridgeModeOff
)

// CartridgeSpec defines the properties of a cartridge, including its Game/ExRom flags and ROM interval configuration.
type CartridgeSpec struct {
	Game         uint8
	ExRom        uint8
	IntervalLow  RomInterval
	IntervalHigh RomInterval
}

// Data returns the Game, ExRom flags, and merged interval configuration for the cartridge specification.
func (c *CartridgeSpec) Data() (uint8, uint8, RomInterval) {
	return c.Game, c.ExRom, c.IntervalLow | c.IntervalHigh
}

// _cartridgesSpec is a slice of pointers to CartridgeSpec, providing configuration details for various cartridge modes.
var _cartridgesSpec []*CartridgeSpec

// init initializes the _cartridgesSpec array with predefined CartridgeSpec values for different cartridge modes.
func init() {
	_cartridgesSpec = make([]*CartridgeSpec, CartridgeModeOff+1)
	_cartridgesSpec[CartridgeMode16K] = &CartridgeSpec{Game: 0, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_1}
	_cartridgesSpec[CartridgeMode8K] = &CartridgeSpec{Game: 0, ExRom: 1, IntervalLow: ROM_LO, IntervalHigh: 0}
	_cartridgesSpec[CartridgeModeUltimax] = &CartridgeSpec{Game: 1, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_2}
	_cartridgesSpec[CartridgeModeOff] = &CartridgeSpec{Game: 1, ExRom: 1, IntervalLow: 0, IntervalHigh: 0}
}

// GetCartridgeSpecFromDetails retrieves a CartridgeSpec matching the provided game and exRom values. Returns nil if not found.
func GetCartridgeSpecFromDetails(game uint8, exRom uint8) *CartridgeSpec {
	for _, spec := range _cartridgesSpec {
		if spec.Game == game && spec.ExRom == exRom {
			return spec
		}
	}
	return nil
}

// GetCartridgeSpec retrieves the CartridgeSpec configuration for the given CartridgeMode.
// The function returns a pointer to a CartridgeSpec instance corresponding to the specified mode parameter.
// It accesses the global _cartridgesSpec array to fetch the relevant configuration.
func GetCartridgeSpec(ct CartridgeMode) *CartridgeSpec {
	return _cartridgesSpec[ct]
}

// ICartridgeC64 represents the interface for cartridges used in the Commodore 64 system, providing core operations and behavior.
type ICartridgeC64 interface {
	Setup() error

	Bind(board IExpansionC64, loader ICartridgeLoaderC64) error

	GetLoaderId() string

	Reset()

	HardwareButton(pressed bool, value uint8)

	Write(i RomInterval, addr uint16, data uint8) bool

	Read(i RomInterval, addr uint16) (uint8, bool)

	IORead(addr uint16) (uint8, bool)

	IOWrite(addr uint16, data uint8) bool

	IRQ(d uint32)

	IRQCLear(d uint32)

	GetExRom() uint8

	GetGame() uint8

	EmulationRequired() bool

	Emulate()

	Detach() error
}

// ICartridgeLoaderC64 provides methods to handle loading and managing Commodore 64 cartridge data and metadata.
// Setup initializes the loader with the given cartridge id and data.
// GetId retrieves the unique identifier for the cartridge.
// GetType returns the type of the cartridge as an integer representation.
// GetData returns the raw binary data of the cartridge.
// Game fetches the Game signal value of the cartridge.
// ExRom fetches the ExRom signal value of the cartridge.
// Name retrieves the name of the cartridge file.
// ReadChipHeader reads and returns the cartridge chip header information.
type ICartridgeLoaderC64 interface {
	Setup(id string, data []byte) error

	GetId() string

	GetType() int

	GetData() []byte

	Game() int

	ExRom() int

	Name() string

	ReadChipHeader() (ICartridgeChipHeaderC64, error)
}

// ICartridgeChipHeaderC64 represents the interface for accessing cartridge chip header information on C64 cartridges.
// Skip returns the number of bytes to skip in the header.
// Kind returns the type of the cartridge chip.
// Bank returns the ROM bank index of the chip.
// Start returns the start address of the ROM chip in memory.
// Size returns the size of the chip in bytes.
// Data returns the raw data bytes of the ROM chip.
// Write writes the chip header and its data to the given writer and returns an error if writing fails.
type ICartridgeChipHeaderC64 interface {
	Skip() uint32

	Kind() uint16

	Bank() uint16

	Start() uint16

	Size() uint16

	Data() []byte

	Write(w io.Writer) error
}

// IdICartridgeC64 generates a unique identifier for an ICartridgeC64 instance using the provided label and instance number.
func IdICartridgeC64(v ICartridgeC64, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToICartridgeC64 converts an IComponent to an ICartridgeC64 implementation or returns an error if the cast fails.
func ComponentToICartridgeC64(component IComponent) (ICartridgeC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component ICartridgeC64 is nil")
	}
	v, ok := component.(ICartridgeC64)
	if !ok {
		return nil, fmt.Errorf("component is not a ICartridgeC64")
	}
	return v, nil
}
