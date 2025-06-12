package final_cartridge_iii

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/catalog"
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeFinalCartridgeIII represents a cartridge model with bank switching and IO interaction mechanisms for ROM emulation.
type CartridgeFinalCartridgeIII struct {
	*component.BaseComponent
	loaderId  string
	intervals references.RomInterval
	//lastData  uint8
	banks [][]byte
	//ioMask    uint8
	currBank uint8
	game     uint8
	exRom    uint8
	board    references.IExpansionC64
}

// GetType returns the type identifier of the Ocean cartridge as an integer constant.
func GetType() int {
	return catalog.CartridgeFinalIII
}

// NewCartridgeFinalCartridgeIII creates and returns a new instance of the Ocean Cartridge conforming to the ICartridgeC64 interface.
func NewCartridgeFinalCartridgeIII(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeFinalCartridgeIII {
	v := references.GetCartridgeSpec(references.CartridgeMode16K)
	co := &CartridgeFinalCartridgeIII{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		game:          v.Game,
		exRom:         v.ExRom,
		intervals:     v.IntervalLow | v.IntervalHigh,
		//lastData:      0,
	}
	co.BaseComponent.Register(factory, parent, Identifier(), co, references.IdICartridgeC64(co, label, instance))
	return co
}

// New creates and returns a new instance of the Ocean Cartridge conforming to the ICartridgeC64 interface.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.ICartridgeC64 {
	return NewCartridgeFinalCartridgeIII(parent, factory, label, instance)
}

// Setup initializes the cartridge with the specified expansion board and CRT loader, setting up necessary configurations.
func (c *CartridgeFinalCartridgeIII) Setup() error {
	return nil
}

func (c *CartridgeFinalCartridgeIII) Bind(board references.IExpansionC64, ldr references.ICartridgeLoaderC64) error {
	c.board = board
	c.loaderId = ldr.GetId()
	if catalog.Type(ldr.GetType()) == catalog.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr)
}

func (c *CartridgeFinalCartridgeIII) Connect() error {
	return nil
}

func (c *CartridgeFinalCartridgeIII) Internal() bool {
	return false
}

// Reset restores the CartridgeOcean state to its initial default, clearing any active configurations or settings.
func (c *CartridgeFinalCartridgeIII) Reset() {

}

// GetLoaderId returns the unique identifier of the CartridgeOcean instance.
func (c *CartridgeFinalCartridgeIII) GetLoaderId() string {
	return c.loaderId
}

// Write attempts to write data to the cartridge at the specified address and interval. Returns true if the write is blocked.
func (c *CartridgeFinalCartridgeIII) Write(i references.RomInterval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("CartridgeFinalCartridgeIII can't be write [bank %d] %x => %d\n", c.currBank, addr, data)
		return true
	}
	return false
}

// Read reads a byte from the current bank, based on the specified ROM interval and address. Returns the byte and success status.
func (c *CartridgeFinalCartridgeIII) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	if i&c.intervals != 0 {
		//return c.banks[c.currBank][addr&0x1fff], true
		bank := c.banks[c.currBank]
		return bank[addr&0x1fff], true
	}
	return 0, false
}

// IORead reads from a memory-mapped I/O address and returns the data and a success flag if the address is valid.
func (c *CartridgeFinalCartridgeIII) IORead(addr uint16) (uint8, bool) {
	//if (addr & 0xfff0) == 0xde00 {
	//	return c.lastData, true
	//}
	return 0, false
}

// IOWrite processes writes to the cartridge's I/O address space, handling bank switching and configuration updates.
func (c *CartridgeFinalCartridgeIII) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xdf00 {
		c.currBank = data & 0x1F
		//currBank := (data & c.ioMask) & 0x3f
		//c.currBank = currBank
		//c.lastData = data
		fmt.Printf("[CartridgeFinalCartridgeIII] Bank switching %x => %d, %d\n", addr, data, c.currBank)
	}
	return false
}

// IRQ handles the Interrupt Request (IRQ) signal for the CartridgeGeneric, enabling appropriate cartridge-specific behavior.
func (c *CartridgeFinalCartridgeIII) IRQ(_ uint32) {
}

// IRQCLear clears the state of any active Interrupt Requests (IRQ) for the CartridgeGeneric.
func (c *CartridgeFinalCartridgeIII) IRQCLear(_ uint32) {
}

// GetExRom returns the value of the ExROM line, which indicates the configuration of the cartridge in the memory map.
func (c *CartridgeFinalCartridgeIII) GetExRom() uint8 {
	return c.exRom
}

// GetGame returns the current game state identifier as a uint8 value.
func (c *CartridgeFinalCartridgeIII) GetGame() uint8 {
	return c.game
}

// Detach removes the cartridge from the system, performing any necessary cleanup or state reset operations.
func (c *CartridgeFinalCartridgeIII) Detach() error {
	//TODO
	return nil
}

// EmulationRequired indicates if additional emulation logic is needed for the cartridge. Always returns false for this type.
func (c *CartridgeFinalCartridgeIII) EmulationRequired() bool {
	return false
}

// Emulate handles the execution of the cartridge's emulation logic during each cycle of the system's operation.
func (c *CartridgeFinalCartridgeIII) Emulate() {

}

// initBin initializes the cartridge by parsing binary data, validating it and segmenting it into fixed-size memory banks.
// It calculates the I/O mask and sets initial values for `lastData` and `currBank`. Returns an error if validation fails.
func (c *CartridgeFinalCartridgeIII) initBin(ldr references.ICartridgeLoaderC64) error {
	data := ldr.GetData()
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	const cSize = 0x2000
	//lCartridge := len(data)
	//c.ioMask = uint8((lCartridge >> 13) - 1)
	totalBanks := len(data) / cSize
	c.banks = make([][]byte, totalBanks)
	for bankIdx := 0; bankIdx < totalBanks; bankIdx++ {
		bank := make([]byte, cSize)
		offset := bankIdx * cSize
		for y := 0; y < cSize; y++ {
			bank[y] = data[offset+y]
		}
		c.banks[bankIdx] = bank
	}
	c.currBank = 0
	return nil
}

// initCrt initializes the cartridge by reading chip headers from the provided CRTLoader and validates the chip data.
func (c *CartridgeFinalCartridgeIII) initCrt(ldr references.ICartridgeLoaderC64) error {
	c.banks = [][]byte{}
	romSize := 0
	for {
		chip, err := ldr.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if chip.Bank() > 31 {
			return fmt.Errorf("invalid chip bank")
		}
		if chip.Start() != 0x8000 {
			return fmt.Errorf("invalid chip start")
		}
		if chip.Size() != 0x4000 {
			return fmt.Errorf("invalid chip size")
		}
		c.banks = append(c.banks, chip.Data())
		romSize += int(chip.Size())
	}
	c.currBank = 0
	return nil
}
