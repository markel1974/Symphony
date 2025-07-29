package magicdesk

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeMagicDesk represents a software-implemented version of a Magic Desk cartridge for system emulation.
// It implements the IC64Cartridge interface for handling cartridge-specific functionality within an expansion board.
type CartridgeMagicDesk struct {
	*component.BaseComponent
	loaderId  string
	spec      *references.C64CartridgeSpec
	banks     [][]byte
	bankMask  uint8
	regVal    uint8
	slot      uint8
	expansion references.IC64Expansion
}

// NewMagicDesk creates a new instance of CartridgeMagicDesk, registers it with the provided parent, and sets its initial configuration.
func NewMagicDesk(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeMagicDesk {
	md := &CartridgeMagicDesk{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		spec:          references.C64CartridgeSpec8K,
		bankMask:      0x7f,
		regVal:        0,
		slot:          0,
	}
	md.BaseComponent.Register(factory, parent, Identifier(), instance, md, references.IdIC64Cartridge(md, label, instance))
	return md
}

// New creates a new instance of CartridgeMagicDesk, registers it with the provided parent, and sets its initial configuration.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IC64Cartridge {
	return NewMagicDesk(parent, factory, label, instance)
}

func (c *CartridgeMagicDesk) Setup() error {
	return nil
}

func (c *CartridgeMagicDesk) Bind(board references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	c.expansion = board
	c.loaderId = ldr.Id()
	if catalog.Type(ldr.Type()) == catalog.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr)
}

// Reset reinitializes the cartridge state, clearing registers and resetting to default configurations.
func (c *CartridgeMagicDesk) Reset() {

}

// Connect establishes the link between the cartridge and the expansion board, enabling the cartridge's functionality.
func (c *CartridgeMagicDesk) Connect() error {
	return nil
}

// Internal indicates whether the cartridge is an internal type. Always returns false for CartridgeMagicDesk.
func (c *CartridgeMagicDesk) Internal() bool {
	return false
}

// GetLoaderId retrieves the unique identifier of the CartridgeMagicDesk instance.
func (c *CartridgeMagicDesk) GetLoaderId() string {
	return c.loaderId
}

// HardwareButton handles the system response to a physical button press event, updating cartridge state as necessary.
func (c *CartridgeMagicDesk) HardwareButton(pressed bool, value uint8) {
}

// Write attempts to write data to the cartridge at the specified ROM interval and address, returning true if write-protected.
func (c *CartridgeMagicDesk) Write(addr uint16, data uint8) bool {
	fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.slot, addr, data)
	return true
}

// Read retrieves the byte at the specified address if the interval matches the cartridge's configuration, else returns 0.
func (c *CartridgeMagicDesk) Read(addr uint16) uint8 {
	return c.banks[c.slot][addr&0x1fff]
}

// IORead checks if the given address matches the specific range (0xde00) and returns the register value and a success flag.
// If the address does not match, it returns 0 and a failure flag.
func (c *CartridgeMagicDesk) IORead(addr uint16) (uint8, bool) {
	if (addr & 0xfff0) == 0xde00 {
		return c.regVal, true
	}
	return 0, false
}

// IOWrite writes data to the specified I/O address and updates the bank, register, and configuration settings if needed.
// Returns false to indicate no further memory interaction processing is needed.
func (c *CartridgeMagicDesk) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xde00 {
		c.regVal = data & (0x80 | c.bankMask)
		c.slot = data & c.bankMask
		fmt.Println("magic desk slot", c.slot)
		var spec *references.C64CartridgeSpec
		if (data & 0x80) != 0 {
			spec = references.C64CartridgeSpecOff
		} else {
			spec = references.C64CartridgeSpec8K
		}
		if spec != c.spec {
			fmt.Println("magic desk changing config", c.spec)
			c.spec = spec
			c.expansion.GameExRomConfigChanged()
		}
	}
	return false
}

// IRQ handles the Interrupt Request (IRQ) signal for the CartridgeGeneric, enabling appropriate cartridge-specific behavior.
func (c *CartridgeMagicDesk) IRQ(_ uint32) {
}

// IRQCLear clears the state of any active Interrupt Requests (IRQ) for the CartridgeGeneric.
func (c *CartridgeMagicDesk) IRQCLear(_ uint32) {
}

// Config returns the Game line status, ExROM line state, and a boolean indicating successful configuration retrieval.
func (c *CartridgeMagicDesk) Config() (uint8, uint8, bool) {
	return c.spec.Game, c.spec.ExRom, true
}

// Detach disconnects the cartridge from the expansion board and releases associated resources. Returns an error if any issues occur.
func (c *CartridgeMagicDesk) Detach() error {
	//TODO
	return nil
}

// EmulationRequired determines whether the cartridge requires emulation support in the system.
func (c *CartridgeMagicDesk) EmulationRequired() bool {
	return false
}

// Emulate manages the core emulation process for the CartridgeMagicDesk, integrating with the system's cycle operations.
func (c *CartridgeMagicDesk) Emulate() {

}

// initBin initializes the cartridge by dividing the given binary data into 8KB banks and setting the bank mask based on size.
// Returns an error if the size of the data is unsupported.
func (c *CartridgeMagicDesk) initBin(ldr references.IC64CartridgeLoader) error {
	data := ldr.Data()
	c.banks = [][]byte{}
	c.bankMask = 0x7f
	c.regVal = 0
	c.slot = 0
	switch len(data) {
	case 0x100000:
		c.bankMask = 0x3f
	case 0x80000:
		c.bankMask = 0x1f
	case 0x40000:
		c.bankMask = 0x0f
	case 0x20000:
		c.bankMask = 0x07
	case 0x10000:
		c.bankMask = 0x03
	default:
		return fmt.Errorf("unsupported size")
	}
	start := 0
	for start < len(data) {
		end := start + 0x2000
		c.banks = append(c.banks, data[start:end])
		start += end
	}
	return nil
}

// initCrt initializes the cartridge by reading chip headers and populating chip banks from a CRTLoader.
// It validates bank numbers, chip sizes, and addresses and determines the appropriate bank mask. Returns an error if invalid.
func (c *CartridgeMagicDesk) initCrt(ldr references.IC64CartridgeLoader) error {
	c.banks = [][]byte{}
	lastBank := uint16(0)
	c.bankMask = 0x7f
	c.regVal = 0
	c.slot = 0
	for {
		chip, err := ldr.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if (chip.Bank() > 128) || ((chip.Start() != 0x8000) && (chip.Start() != 0xa000)) || (chip.Size() != 0x2000) {
			return fmt.Errorf("invalid chip bank")
		}
		c.banks = append(c.banks, chip.Data())
		if chip.Bank() > lastBank {
			lastBank = chip.Bank()
		}
	}
	if lastBank >= 128 {
		return fmt.Errorf("chip has more than 128 banks")
	}
	if lastBank >= 64 {
		c.bankMask = 0x7f // min 65, max 128 banks
	} else if lastBank >= 32 {
		c.bankMask = 0x3f // min 33, max 64 banks
	} else if lastBank >= 16 {
		c.bankMask = 0x1f // min 17, max 32 banks
	} else if lastBank >= 8 {
		c.bankMask = 0x0f // min 9, max 16 banks
	} else if lastBank >= 4 {
		c.bankMask = 0x07 // min 5, max 8 banks
	} else {
		c.bankMask = 0x03 // max 4 banks
	}
	return nil
}
