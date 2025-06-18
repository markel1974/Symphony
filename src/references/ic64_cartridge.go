package references

import (
	"fmt"
	"io"
)

// C64RomInterval represents a type used to define distinct memory interval mappings in the ROM of a cartridge system.

// ROM_UNK represents an unknown ROM type.
// ROM_LO represents a low-priority ROM type.
// ROM_HI_1 represents a high-priority ROM type (level 1).
// ROM_HI_2 represents a high-priority ROM type (level 2).
const (
	ROM_UNK  = uint8(0)
	ROM_LO   = uint8(1)
	ROM_HI_1 = uint8(2)
	ROM_HI_2 = uint8(4)
)

// C64CartridgeSpec defines the specifications for a Commodore 64 cartridge, including its mode, memory intervals, and signals.
type C64CartridgeSpec struct {
	Game         uint8
	ExRom        uint8
	Intervals    uint8
	IntervalLow  uint8
	IntervalHigh uint8
	Mode         uint8
}

// NewC64CartridgeSpec creates a new C64CartridgeSpec with specified game, exRom, intervalLow, and intervalHigh values.
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

// C64CartridgeSpec16K defines a 16KB cartridge specification for C64 with GAME=0, EXROM=0, and ROM ranges $8000-$9FFF, $A000-$BFFF.
var C64CartridgeSpec16K = NewC64CartridgeSpec(0, 0, ROM_LO, ROM_HI_1)

// C64CartridgeSpec8K represents the 8K cartridge specification for the Commodore 64, with GAME=0, EXROM=1, and ROM_LO interval.
var C64CartridgeSpec8K = NewC64CartridgeSpec(0, 1, ROM_LO, 0)

// C64CartridgeSpecUltimax defines the specification for a Commodore 64 Ultimax cartridge mode with specific ROM intervals.
var C64CartridgeSpecUltimax = NewC64CartridgeSpec(1, 0, ROM_LO, ROM_HI_2)

// C64CartridgeSpecOff represents the default cartridge configuration with GAME and EXROM set to 1 and intervals disabled.
var C64CartridgeSpecOff = NewC64CartridgeSpec(1, 1, 0, 0)

// _c64CrtSpec is a slice of pointers to C64CartridgeSpec that defines cartridge specifications for different modes.
var _c64CrtSpec []*C64CartridgeSpec

// _c64BanksType is an array of uint8 used to map memory bank addresses to corresponding ROM bank types.
var _c64BanksType []uint8

// init initializes the `_c64BanksType` and `_c64CrtSpec` arrays with default ROM types and cartridge specifications.
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

// c64CartridgeComputeMode determines the cartridge mode based on the GAME and EXROM signals provided as inputs.
// The function combines GAME and EXROM using bitwise operations and returns the resulting 2-bit mode value.
func c64CartridgeComputeMode(game uint8, exRom uint8) uint8 {
	return (game<<1 | exRom) & 0x3
}

// GetCartridgeSpec returns a pointer to a C64CartridgeSpec based on the provided game and exRom signals.
func GetCartridgeSpec(game uint8, exRom uint8) *C64CartridgeSpec {
	m := c64CartridgeComputeMode(game, exRom)
	return _c64CrtSpec[m]
}

// C64CartridgeBank computes the cartridge memory bank index based on the given address.
func C64CartridgeBank(addr uint16) uint8 {
	v := addr >> 12
	return _c64BanksType[v]
}

// IC64Cartridge defines an interface for managing Commodore 64 cartridge functionality and operations.
// Setup initializes the cartridge for use.
// Bind associates the cartridge with a provided board and loader.
// GetLoaderId retrieves the unique identifier of the associated loader.
// Reset resets the cartridge state to default.
// HardwareButton handles hardware button press events with given state and value.
// Read reads data from the specified memory address.
// IORead performs an I/O read operation at the specified address, returning data and validity.
// IOWrite performs an I/O write operation to the specified address with given data, returning success.
// IRQ triggers an interrupt request with the given delay.
// IRQCLear clears interrupt requests after the given delay.
// GetExRom retrieves the current ExRom value.
// GetGame retrieves the current Game line value.
// EmulationRequired checks if emulation is mandatory for the cartridge.
// Emulate executes the emulation process if required.
// Detach handles cartridge removal and clean-up operations.
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

// IC64CartridgeLoader defines the interface for loading and managing Commodore 64 cartridge data and metadata.
// Setup initializes the cartridge loader with a unique ID and binary data.
// GetId retrieves the loader's unique identifier.
// GetType returns the type identifier of the cartridge being loaded.
// GetData retrieves the raw binary data loaded into the cartridge.
// Game returns the GAME line status for the cartridge.
// ExRom returns the EXROM line status for the cartridge.
// Name fetches the name or label of the cartridge.
// ReadChipHeader parses and retrieves the cartridge chip header for detailed configuration.
type IC64CartridgeLoader interface {
	Setup(id string, data []byte) error

	Id() string

	Type() int

	Data() []byte

	Game() int

	ExRom() int

	Name() string

	ReadChipHeader() (IC64CartridgeChipHeader, error)
}

// IC64CartridgeChipHeader defines the structure for handling C64 cartridge chip headers.
// Skip returns the offset to skip to the next chip.
// Kind returns the type identifier of the chip.
// Bank returns the chip's associated memory bank number.
// Start returns the starting address of the chip in the memory.
// Size returns the size of the chip data in bytes.
// Data returns the raw chip data as a byte slice.
// Write writes the chip data to an io.Writer and returns an error if the operation fails.
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

// ComponentToIC64Cartridge converts an IComponent to an IC64Cartridge, returning an error if the conversion fails or input is nil.
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
