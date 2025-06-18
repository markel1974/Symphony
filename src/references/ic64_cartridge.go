package references

import (
	"fmt"
	"io"
)

// C64RomInterval represents a type used to define distinct memory interval mappings in the ROM of a cartridge system.

// ROM_LO represents the low ROM interval value in the C64RomInterval type.
// ROM_HI_1 represents the first high ROM interval value in the C64RomInterval type.
// ROM_HI_2 represents the second high ROM interval value in the C64RomInterval type.
const (
	ROM_UNK  = uint8(0)
	ROM_LO   = uint8(1)
	ROM_HI_1 = uint8(2)
	ROM_HI_2 = uint8(4)
)

// C64CartridgeSpec defines the properties of a cartridge, including its Game/ExRom flags and ROM interval configuration.
type C64CartridgeSpec struct {
	Game         uint8
	ExRom        uint8
	Intervals    uint8
	IntervalLow  uint8
	IntervalHigh uint8
	Mode         uint8
}

func NewC64CartridgeSpec(game uint8, exRom uint8, intervalLow uint8, intervalHigh uint8) *C64CartridgeSpec {
	return &C64CartridgeSpec{
		Game:         game,
		ExRom:        exRom,
		Intervals:    intervalLow | intervalHigh,
		IntervalLow:  intervalLow,
		IntervalHigh: intervalHigh,
		Mode:         c64CartridgeComputeMode(game, exRom),
	}
}

// Data returns the Game, ExRom flags, and merged interval configuration for the cartridge specification.
func (c *C64CartridgeSpec) Data() (uint8, uint8, uint8) {
	return c.Game, c.ExRom, c.IntervalLow | c.IntervalHigh
}

var C64CartridgeSpec16K = NewC64CartridgeSpec(0, 0, ROM_LO, ROM_HI_1)
var C64CartridgeSpec8K = NewC64CartridgeSpec(0, 1, ROM_LO, 0)
var C64CartridgeSpecUltimax = NewC64CartridgeSpec(1, 0, ROM_LO, ROM_HI_2)
var C64CartridgeSpecOff = NewC64CartridgeSpec(1, 1, 0, 0)

// _c64CartridgesSpec is a slice of pointers to C64CartridgeSpec, providing configuration details for various cartridge modes.
var _c64CrtSpec []*C64CartridgeSpec

var _c64BanksType []uint8

func init() {
	_c64BanksType = make([]uint8, 0xf+1)
	for idx := range _c64BanksType {
		_c64BanksType[idx] = ROM_UNK
	}
	_c64BanksType[0x8] = ROM_LO
	_c64BanksType[0x9] = ROM_LO
	_c64BanksType[0xa] = ROM_HI_1
	_c64BanksType[0xb] = ROM_HI_1
	_c64BanksType[0xe] = ROM_HI_2
	_c64BanksType[0xf] = ROM_HI_2

	_c64CrtSpec = make([]*C64CartridgeSpec, 0x3+1)
	for idx := range _c64CrtSpec {
		_c64CrtSpec[idx] = C64CartridgeSpecOff
	}
	_c64CrtSpec[C64CartridgeSpec16K.Mode] = C64CartridgeSpec16K
	_c64CrtSpec[C64CartridgeSpec8K.Mode] = C64CartridgeSpec8K
	_c64CrtSpec[C64CartridgeSpecUltimax.Mode] = C64CartridgeSpecUltimax
}

func c64CartridgeComputeMode(game uint8, exRom uint8) uint8 {
	return (game<<1 | exRom) & 0x3
}

// GetCartridgeSpec retrieves a C64CartridgeSpec matching the provided game and exRom values. Returns nil if not found.
func GetCartridgeSpec(game uint8, exRom uint8) *C64CartridgeSpec {
	m := c64CartridgeComputeMode(game, exRom)
	return _c64CrtSpec[m]
}

func C64CartridgeBank(addr uint16) uint8 {
	v := addr >> 12
	return _c64BanksType[v]
}

// IC64Cartridge represents the interface for cartridges used in the Commodore 64 system, providing core operations and behavior.
type IC64Cartridge interface {
	Setup() error

	Bind(board IC64Expansion, loader IC64CartridgeLoader) error

	GetLoaderId() string

	Reset()

	HardwareButton(pressed bool, value uint8)

	Read(addr uint16) uint8

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

// IC64CartridgeLoader provides methods to handle loading and managing Commodore 64 cartridge data and metadata.
// Setup initializes the loader with the given cartridge id and data.
// GetId retrieves the unique identifier for the cartridge.
// GetType returns the type of the cartridge as an integer representation.
// GetData returns the raw binary data of the cartridge.
// Game fetches the Game signal value of the cartridge.
// ExRom fetches the ExRom signal value of the cartridge.
// Name retrieves the name of the cartridge file.
// ReadChipHeader reads and returns the cartridge chip header information.
type IC64CartridgeLoader interface {
	Setup(id string, data []byte) error

	GetId() string

	GetType() int

	GetData() []byte

	Game() int

	ExRom() int

	Name() string

	ReadChipHeader() (IC64CartridgeChipHeader, error)
}

// IC64CartridgeChipHeader represents the interface for accessing cartridge chip header information on C64 cartridges.
// Skip returns the number of bytes to skip in the header.
// Kind returns the type of the cartridge chip.
// Bank returns the ROM bank index of the chip.
// Start returns the start address of the ROM chip in memory.
// Size returns the size of the chip in bytes.
// Data returns the raw data bytes of the ROM chip.
// Write writes the chip header and its data to the given writer and returns an error if writing fails.
type IC64CartridgeChipHeader interface {
	Skip() uint32

	Kind() uint16

	Bank() uint16

	Start() uint16

	Size() uint16

	Data() []byte

	Write(w io.Writer) error
}

// IdIC64Cartridge generates a unique identifier for an IC64Cartridge instance using the provided label and instance number.
func IdIC64Cartridge(v IC64Cartridge, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Cartridge converts an IComponent to an IC64Cartridge implementation or returns an error if the cast fails.
func ComponentToIC64Cartridge(component IComponent) (IC64Cartridge, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Cartridge is nil")
	}
	v, ok := component.(IC64Cartridge)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Cartridge")
	}
	return v, nil
}
