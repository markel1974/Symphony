package references

import (
	"fmt"
	"io"
)

// RomInterval represents a type for defining ROM address ranges or intervals in memory mapping and cartridge operations.
type RomInterval int

// ROM_LO indicates the low ROM interval.
// ROM_HI_1 represents the first high ROM interval.
// ROM_HI_2 represents the second high ROM interval.
const (
	ROM_LO   = RomInterval(1)
	ROM_HI_1 = RomInterval(2)
	ROM_HI_2 = RomInterval(4)
)

// CartridgeMode represents the mode or configuration type of a cartridge, often determining its behavior or memory layout.
type CartridgeMode int

// CartridgeMode16K represents a 16KB cartridge mode.
// CartridgeMode8K represents an 8KB cartridge mode.
// CartridgeModeUltimax represents the Ultimax cartridge mode.
// CartridgeModeOff represents the cartridge mode being off.
const (
	CartridgeMode16K = CartridgeMode(iota)
	CartridgeMode8K
	CartridgeModeUltimax
	CartridgeModeOff
)

// CartridgeSpec defines the configuration parameters for a cartridge, including ROM intervals and hardware signals.
type CartridgeSpec struct {
	Game         uint8
	ExRom        uint8
	IntervalLow  RomInterval
	IntervalHigh RomInterval
}

// _cartridgesSpec stores specifications for each cartridge mode as a slice of pointers to CartridgeSpec.
var _cartridgesSpec []*CartridgeSpec

// init initializes the _cartridgesSpec array with predefined CartridgeSpec configurations for different cartridge modes.
func init() {
	_cartridgesSpec = make([]*CartridgeSpec, CartridgeModeOff+1)
	_cartridgesSpec[CartridgeMode16K] = &CartridgeSpec{Game: 0, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_1}
	_cartridgesSpec[CartridgeMode8K] = &CartridgeSpec{Game: 0, ExRom: 1, IntervalLow: ROM_LO, IntervalHigh: 0}
	_cartridgesSpec[CartridgeModeUltimax] = &CartridgeSpec{Game: 1, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_2}
	_cartridgesSpec[CartridgeModeOff] = &CartridgeSpec{Game: 1, ExRom: 1, IntervalLow: 0, IntervalHigh: 0}
}

func GetCartridgeSpecFromDetails(game uint8, exRom uint8) *CartridgeSpec {
	for _, spec := range _cartridgesSpec {
		if spec.Game == game && spec.ExRom == exRom {
			return spec
		}
	}
	return nil
}

// GetCartridgeSpec retrieves the CartridgeSpec associated with a given CartridgeMode from the predefined set of specifications.
func GetCartridgeSpec(ct CartridgeMode) *CartridgeSpec {
	return _cartridgesSpec[ct]
}

func IdICartridgeC64(_ ICartridgeC64, label string, instance int) string {
	return IdInternalComponent(label, instance, "ICartridgeC64")
}

// ICartridgeC64 represents the interface for a C64-compatible cartridge, defining methods for setup, memory operations, and emulation.
// Setup initializes the cartridge with the provided expansion board and loader.
// GetLoaderId retrieves the unique identifier for the cartridge loader.
// Reset resets the cartridge to its initial state.
// HardwareButton handles the response to a physical button press event
// Write writes data to a specified memory address within a ROM interval and returns success status.
// Read reads data from a specified memory address within a ROM interval and returns the data and success status.
// IORead performs an I/O read from a specified address and returns the data and success status.
// IOWrite performs an I/O write to a specified address with data and returns success status.
// GetExRom retrieves the current ExROM configuration value.
// GetGame retrieves the current Game configuration value.
// EmulationRequired checks if the cartridge requires emulation.
// Emulate initiates the emulation process for the cartridge if required.
// Detach detaches the cartridge, releasing any associated resources.
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

// ICartridgeLoaderC64 defines an interface for loading and managing C64 cartridge data.
// Setup initializes the cartridge loader with the given identifier and data.
// GetId retrieves the unique identifier for the cartridge.
// GetType returns the type of the cartridge as an integer.
// GetData provides the binary data associated with the cartridge.
// Game retrieves the value representing the game configuration.
// ExRom retrieves the value representing the ExROM configuration.
// Name returns the name associated with the cartridge.
// ReadChipHeader reads and parses a chip header from the cartridge data, returning the header and any encountered error.
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

// ICartridgeChipHeaderC64 defines an interface for handling cartridge chip headers in CRT files.
// Skip returns the number of bytes to skip after the chip data.
// Kind retrieves the type of the chip as a 16-bit value.
// Bank retrieves the chip's assigned bank number as a 16-bit value.
// Start provides the start address of the chip within the address space.
// Size returns the size of the chip data in bytes.
// Data retrieves the binary data associated with the chip.
// Write writes the chip data to the provided io.Writer and returns an error if unsuccessful.
type ICartridgeChipHeaderC64 interface {
	Skip() uint32

	Kind() uint16

	Bank() uint16

	Start() uint16

	Size() uint16

	Data() []byte

	Write(w io.Writer) error
}

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
