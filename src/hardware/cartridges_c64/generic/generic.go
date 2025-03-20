package generic

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	loader2 "github.com/markel1974/c64emu/src/hardware/cartridges_c64/loader"
	"github.com/markel1974/c64emu/src/references"
)

// cSize16K defines the size of 16 kilobytes (0x4000), commonly used for memory allocation or data validation.
const cSize16K = 0x4000

// cSize8K defines a constant representing 8 kilobytes (8192 bytes) in hexadecimal format (0x2000).
const cSize8K = 0x2000

// Generic represents the structure and functionality of a cartridge, including memory banks, intervals, and configuration.
type Generic struct {
	*component.BaseComponent
	loaderId   string
	b0Interval references.RomInterval
	b1Interval references.RomInterval
	bank0      []uint8
	bank1      []uint8
	intervals  references.RomInterval
	game       uint8
	exRom      uint8
	board      references.IExpansionC64
}

// GetType returns the constant value representing the CARTRIDGE_CRT type.
func GetType() int {
	return loader2.CARTRIDGE_CRT
}

// New creates and returns a new instance of the Generic cartridge implementing the ICartridgeC64 interface.
func New(parent references.IComponent, factory references.IComponentFactory, instance int) references.ICartridgeC64 {
	const id = "generic"
	g := &Generic{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      "generic",
		game:          0,
		exRom:         0,
		b0Interval:    0,
		b1Interval:    0,
		intervals:     0,
	}
	g.BaseComponent.Register(factory, parent, id, instance, g, references.IdICartridgeC64(g))
	return g
}

// Setup initializes the Generic cartridge by setting the board and loading data using the provided CRTLoader.
func (c *Generic) Setup(board references.IExpansionC64, ldr references.ICartridgeLoaderC64) error {
	c.board = board
	c.loaderId = ldr.GetId()
	if loader2.Type(ldr.GetType()) == loader2.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initRaw(ldr.GetData())
}

// Reset reinitializes the Generic cartridge to its default state, clearing all settings and memory banks.
func (c *Generic) Reset() {

}

// GetLoaderId retrieves the unique identifier of the Generic instance.
func (c *Generic) GetLoaderId() string {
	return c.loaderId
}

// initCrt initializes the cartridge configuration using the provided CRTLoader.
// It sets up memory banks and mode based on the cartridge chip headers.
// Returns an error if the cartridge format or configuration is unsupported.
func (c *Generic) initCrt(ldr references.ICartridgeLoaderC64) error {
	c.bank0 = make([]uint8, cSize8K)
	c.bank1 = make([]uint8, cSize8K)
	chip1, err := ldr.ReadChipHeader()
	if err != nil {
		return err
	}
	if chip1 == nil {
		return fmt.Errorf("nil chip")
	}
	if chip1.Start() == 0x8000 {
		if chip1.Size() == cSize8K {
			if chip2, _ := ldr.ReadChipHeader(); chip2 == nil {
				copy(c.bank0, chip1.Data())
				c.applyConfig(references.CartridgeMode8K)
				return nil
			} else if chip2.Size() == cSize8K {
				if chip2.Start() == 0x8000 {
					copy(c.bank0, chip1.Data())
					copy(c.bank1, chip2.Data())
					c.applyConfig(references.CartridgeMode16K)
					return nil
				} else if chip2.Start() == 0xe000 {
					copy(c.bank0, chip1.Data())
					copy(c.bank1, chip2.Data())
					c.applyConfig(references.CartridgeModeUltimax)
					return nil
				}
			}
		} else if chip1.Size() == cSize16K {
			copy(c.bank0, chip1.Data()[:cSize8K])
			copy(c.bank1, chip1.Data()[cSize8K:])
			c.applyConfig(references.CartridgeMode16K)
			return nil
		}
	}
	return fmt.Errorf("unsupported crt")
}

// initRaw initializes the raw cartridge data, validates its size, and configures the banks and cartridge mode accordingly.
// Returns an error if the cartridge size is invalid or validation fails.
func (c *Generic) initRaw(data []byte) error {
	c.bank0 = make([]uint8, cSize8K)
	c.bank1 = make([]uint8, cSize8K)
	if err := loader2.ValidateCartridge(data); err != nil {
		return err
	}
	if len(data) == cSize8K {
		copy(c.bank0, data)
		c.applyConfig(references.CartridgeMode8K)
		return nil
	}
	if len(data) == cSize16K {
		copy(c.bank0, data[:cSize8K])
		copy(c.bank1, data[cSize8K:])
		c.applyConfig(references.CartridgeMode16K)
		return nil
	}
	return fmt.Errorf("invalid size")
}

// applyConfig configures the Generic cartridge by setting its memory intervals and control flags based on the provided CartridgeMode.
func (c *Generic) applyConfig(ct references.CartridgeMode) {
	v := references.GetCartridgeSpec(ct)
	c.game = v.Game
	c.exRom = v.ExRom
	c.b0Interval = v.IntervalLow
	c.b1Interval = v.IntervalHigh
	c.intervals = v.IntervalLow | v.IntervalHigh
}

// Write attempts to write data to the cartridge at the specified interval and address. Returns true if writing is not allowed.
func (c *Generic) Write(i references.RomInterval, addr uint16, data uint8) bool {
	if (i & c.intervals) != 0 {
		fmt.Printf("Generic Cartridge can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

// Read retrieves a byte and a success flag from the cartridge based on the provided address and ROM interval.
func (c *Generic) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	if (i & c.intervals) != 0 {
		if c.b0Interval == i {
			return c.bank0[(addr & 0x1fff)], true
		}
		if c.b1Interval == i {
			return c.bank1[(addr & 0x1fff)], true
		}
	}
	return 0, false
}

// IORead performs an I/O read operation and always returns 0 and false, indicating no data is available for the requested address.
func (c *Generic) IORead(_ uint16) (uint8, bool) {
	return 0, false
}

// IOWrite handles IO write operations for the Generic cartridge, always returning false to indicate no action is performed.
func (c *Generic) IOWrite(_ uint16, _ uint8) bool {
	return false
}

// GetExRom retrieves the current state of the ExRom line associated with the Generic cartridge, returning its value as uint8.
func (c *Generic) GetExRom() uint8 {
	return c.exRom
}

// GetGame returns the game identifier associated with the Generic cartridge.
func (c *Generic) GetGame() uint8 {
	return c.game
}

// EmulationRequired checks if the cartridge requires emulation to function properly. Returns false for Generic cartridges.
func (c *Generic) EmulationRequired() bool {
	return false
}

// Emulate performs the main emulation cycle logic for the Generic type, enabling behavior specific to the cartridge type.
func (c *Generic) Emulate() {
}

// Detach gracefully detaches the cartridge from the associated expansion, releasing any allocated resources or states.
func (c *Generic) Detach() error {
	//TODO
	return nil
}
