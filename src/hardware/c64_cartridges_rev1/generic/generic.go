package generic

import (
	"fmt"

	"github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/symphony/src/kernel/component"
	"github.com/markel1974/symphony/src/references"
)

// cSize16K defines the constant value for 16KB, equivalent to 0x4000 in hex.
const cSize16K = 0x4000

// cSize8K defines the constant value for 8KB, equivalent to 0x2000 in hex.
const cSize8K = 0x2000

// CartridgeGeneric represents a generic cartridge component used in the system.
// It contains memory banks and associated cartridge specifications.
// The type leverages base component functionality for shared behaviors.
type CartridgeGeneric struct {
	*component.BaseComponent
	loaderId  string
	bank0     []uint8
	bank1     []uint8
	spec      *references.C64CartridgeSpec
	expansion references.IC64Expansion
}

// GetType returns the integer type identifier for a cartridge, specifically catalog.CartridgeCRT.
func GetType() int {
	return catalog.CartridgeCRT
}

// NewCartridgeGeneric initializes and returns a new instance of CartridgeGeneric with the specified parent, factory, label, and instance.
func NewCartridgeGeneric(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeGeneric {
	g := &CartridgeGeneric{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		spec:          references.C64CartridgeSpec8K,
	}
	g.BaseComponent.Register(factory, parent, Identifier(), instance, g, references.IdIC64Cartridge(g, label, instance))
	return g
}

// New creates and initializes a new IC64Cartridge instance using the provided parent component, factory, label, and instance.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IC64Cartridge {
	return NewCartridgeGeneric(parent, factory, label, instance)
}

// Setup initializes the CartridgeGeneric component. Returns an error if initialization fails.
func (c *CartridgeGeneric) Setup() error {
	return nil
}

// Bind initializes the cartridge by associating it with an expansion and a cartridge loader, configuring its type and data.
func (c *CartridgeGeneric) Bind(expansion references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	c.expansion = expansion
	c.loaderId = ldr.Id()
	if catalog.Type(ldr.Type()) == catalog.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initRaw(ldr.Data())
}

// Reset restores the cartridge state to its default initialization, clearing any active configurations or states.
func (c *CartridgeGeneric) Reset() {

}

// Connect establishes the connection of the cartridge to the expansion port, preparing it for operation.
func (c *CartridgeGeneric) Connect() error {
	return nil
}

// Internal indicates whether the cartridge operates in internal mode. Returns false as the default implementation.
func (c *CartridgeGeneric) Internal() bool {
	return false
}

// GetLoaderId returns the Id of the cartridge loader associated with the CartridgeGeneric instance.
func (c *CartridgeGeneric) GetLoaderId() string {
	return c.loaderId
}

// initCrt initializes the cartridge banks and specification based on the data and headers provided by the loader.
// The method supports 8K, 16K, and Ultimax specifications and validates compatibility.
// Returns an error if the chip configuration is unsupported or if the header data is invalid.
func (c *CartridgeGeneric) initCrt(ldr references.IC64CartridgeLoader) error {
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
				c.applySpec(references.C64CartridgeSpec8K)
				return nil
			} else if chip2.Size() == cSize8K {
				if chip2.Start() == 0x8000 {
					copy(c.bank0, chip1.Data())
					copy(c.bank1, chip2.Data())
					c.applySpec(references.C64CartridgeSpec16K)
					return nil
				} else if chip2.Start() == 0xe000 {
					copy(c.bank0, chip1.Data())
					copy(c.bank1, chip2.Data())
					c.applySpec(references.C64CartridgeSpecUltimax)
					return nil
				}
			}
		} else if chip1.Size() == cSize16K {
			copy(c.bank0, chip1.Data()[:cSize8K])
			copy(c.bank1, chip1.Data()[cSize8K:])
			c.applySpec(references.C64CartridgeSpec16K)
			return nil
		}
	}
	return fmt.Errorf("unsupported crt")
}

// initRaw initializes cartridge memory banks with the provided raw data, validating its size and content.
// Returns an error if the data is invalid or the size is not supported (8KB or 16KB).
func (c *CartridgeGeneric) initRaw(data []byte) error {
	c.bank0 = make([]uint8, cSize8K)
	c.bank1 = make([]uint8, cSize8K)
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	if len(data) == cSize8K {
		copy(c.bank0, data)
		c.spec = references.C64CartridgeSpec8K
		return nil
	}
	if len(data) == cSize16K {
		copy(c.bank0, data[:cSize8K])
		copy(c.bank1, data[cSize8K:])
		c.spec = references.C64CartridgeSpec16K
		return nil
	}
	return fmt.Errorf("invalid size")
}

// HardwareButton handles hardware button interactions, determining actions based on whether the button is pressed and its value.
func (c *CartridgeGeneric) HardwareButton(_ bool, _ uint8) {
}

// Read retrieves a value from the cartridge memory based on the provided address and the cartridge bank mapping.
func (c *CartridgeGeneric) Read(addr uint16) uint8 {
	i := references.C64CartridgeBank(addr)
	if c.spec.IntervalLow == i {
		return c.bank0[(addr & 0x1fff)]
	}
	if c.spec.IntervalHigh == i {
		return c.bank1[(addr & 0x1fff)]
	}
	return 0
}

// IRQ triggers an interrupt request for the cartridge with the provided cycle count.
func (c *CartridgeGeneric) IRQ(_ uint32) {
}

// IRQCLear clears or deactivates the IRQ signal for the cartridge.
func (c *CartridgeGeneric) IRQCLear(_ uint32) {
}

// IORead reads data from the I/O address space
// and returns the value along with a boolean indicating success or failure.
func (c *CartridgeGeneric) IORead(_ uint16) (uint8, bool) {
	return 0, false
}

// IOWrite handles writing to the cartridge I/O space and returns a boolean indicating if the operation was successful.
func (c *CartridgeGeneric) IOWrite(_ uint16, _ uint8) bool {
	return false
}

// Config returns the Game line status, ExROM line state, and a boolean indicating successful configuration retrieval.
func (c *CartridgeGeneric) Config() (uint8, uint8, bool) {
	return c.spec.Game, c.spec.ExRom, true
}

// EmulationRequired checks if the cartridge requires emulation and returns false to indicate no emulation is needed.
func (c *CartridgeGeneric) EmulationRequired() bool {
	return false
}

// Emulate handles the execution cycle logic of the cartridge emulation for the connected hardware or simulation framework.
func (c *CartridgeGeneric) Emulate() {
}

// Detach disconnects the cartridge from the expansion system and performs necessary cleanup operations.
func (c *CartridgeGeneric) Detach() error {
	//TODO
	return nil
}

func (c *CartridgeGeneric) applySpec(spec *references.C64CartridgeSpec) {
	c.spec = spec
	c.expansion.GameExRomConfigChanged()
}
