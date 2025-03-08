package magicdesk

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/components/board"
)

// CartridgeMagicDesk represents a software-implemented version of a Magic Desk cartridge for system emulation.
// It implements the ICartridge interface for handling cartridge-specific functionality within an expansion board.
type CartridgeMagicDesk struct {
	*board.BaseComponent
	loaderId string
	spec     *icartridge.CartridgeSpec
	banks    [][]byte
	bankMask uint8
	regVal   uint8
	slot     uint8
	board    icartridge.IExpansion
}

// GetType returns an integer identifier representing the type of the Magic Desk cartridge.
func GetType() int {
	return loader.CARTRIDGE_MAGIC_DESK
}

// New creates and returns a new instance of a CartridgeMagicDesk implementing the icartridge.ICartridge interface.
func New(parenNode *board.Node, suffix string) icartridge.ICartridge {
	md := &CartridgeMagicDesk{
		BaseComponent: board.NewBaseComponent(),
		loaderId:      "magicDesk",
		spec:          icartridge.GetCartridgeSpec(icartridge.CartridgeMode8K),
		bankMask:      0x7f,
		regVal:        0,
		slot:          0,
	}
	md.Register(md, "magicDesk", suffix, parenNode, nil)
	return md
}

// Setup initializes the CartridgeMagicDesk by configuring its board and loading data via the provided CRTLoader.
func (c *CartridgeMagicDesk) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.loaderId = ldr.GetId()
	if ldr.GetType() == loader.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

// Reset reinitializes the cartridge state, clearing registers and resetting to default configurations.
func (c *CartridgeMagicDesk) Reset() {

}

// GetLoaderId retrieves the unique identifier of the CartridgeMagicDesk instance.
func (c *CartridgeMagicDesk) GetLoaderId() string {
	return c.loaderId
}

// Write attempts to write data to the cartridge at the specified ROM interval and address, returning true if write-protected.
func (c *CartridgeMagicDesk) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if (i & (c.spec.IntervalLow | c.spec.IntervalHigh)) != 0 {
		fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.slot, addr, data)
		return true
	}
	return false
}

// Read retrieves the byte at the specified address if the interval matches the cartridge's configuration, else returns 0.
func (c *CartridgeMagicDesk) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if (i & (c.spec.IntervalLow | c.spec.IntervalHigh)) != 0 {
		//if c.b0Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		//if c.b1Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		return c.banks[c.slot][addr&0x1fff], true
	}
	return 0, false
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
		var spec *icartridge.CartridgeSpec
		if (data & 0x80) != 0 {
			spec = icartridge.GetCartridgeSpec(icartridge.CartridgeModeOff)
		} else {
			spec = icartridge.GetCartridgeSpec(icartridge.CartridgeMode8K)
		}
		if spec != c.spec {
			fmt.Println("magic desk changing config", c.spec)
			c.spec = spec
			c.board.GameExRomConfigChanged()
		}
	}
	return false
}

// GetExRom retrieves the ExRom value from the cartridge specification.
func (c *CartridgeMagicDesk) GetExRom() uint8 {
	return c.spec.ExRom
}

// GetGame retrieves the value of the "Game" field from the cartridge specification associated with the instance.
func (c *CartridgeMagicDesk) GetGame() uint8 {
	return c.spec.Game
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
func (c *CartridgeMagicDesk) initBin(data []byte) error {
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
func (c *CartridgeMagicDesk) initCrt(loader *loader.CRTLoader) error {
	c.banks = [][]byte{}
	lastBank := uint16(0)
	c.bankMask = 0x7f
	c.regVal = 0
	c.slot = 0
	for {
		chip, err := loader.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if (chip.Bank > 128) || ((chip.Start != 0x8000) && (chip.Start != 0xa000)) || (chip.Size != 0x2000) {
			return fmt.Errorf("invalid chip bank")
		}
		c.banks = append(c.banks, chip.Data)
		if chip.Bank > lastBank {
			lastBank = chip.Bank
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
