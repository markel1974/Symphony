package ocean

import (
	"fmt"

	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeOcean represents a cartridge model with bank switching and IO interaction mechanisms for ROM emulation.
type CartridgeOcean struct {
	*component.BaseComponent
	loaderId string
	lastData uint8
	banks    [][]byte
	ioMask   uint8
	currBank uint8
	spec     *references.C64CartridgeSpec
	board    references.IC64Expansion
}

// NewCartridgeOcean creates and returns a new instance of the Ocean Cartridge conforming to the IC64Cartridge interface.
func NewCartridgeOcean(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeOcean {
	co := &CartridgeOcean{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		spec:          references.C64CartridgeSpec16K,
	}
	co.BaseComponent.Register(factory, parent, Identifier(), instance, co, references.IdIC64Cartridge(co, label, instance))
	return co
}

// New creates and returns a new instance of the Ocean Cartridge conforming to the IC64Cartridge interface.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IC64Cartridge {
	return NewCartridgeOcean(parent, factory, label, instance)
}

func (c *CartridgeOcean) reset(hard bool) {
	c.spec = references.C64CartridgeSpec16K
	c.lastData = 0
	c.currBank = 0
	if hard {
		c.ioMask = 0
		c.banks = nil
	}
}

// Setup initializes the cartridge with the specified expansion board and CRT loader, setting up necessary configurations.
func (c *CartridgeOcean) Setup() error {
	c.reset(true)
	return nil
}

// Bind associates the CartridgeOcean instance with the provided expansion board and cartridge loader, initializing it accordingly.
func (c *CartridgeOcean) Bind(board references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	c.board = board
	c.loaderId = ldr.Id()
	if catalog.Type(ldr.Type()) == catalog.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr)
}

// Connect establishes a connection or initializes the cartridge within the system context, returning an error if unsuccessful.
func (c *CartridgeOcean) Connect() error {
	return nil
}

func (c *CartridgeOcean) Internal() bool {
	return false
}

// Reset restores the CartridgeOcean state to its initial default, clearing any active configurations or settings.
func (c *CartridgeOcean) Reset() {
	c.reset(false)
}

// GetLoaderId returns the unique identifier of the CartridgeOcean instance.
func (c *CartridgeOcean) GetLoaderId() string {
	return c.loaderId
}

// HardwareButton handles the system response to a physical button press event, updating cartridge state as necessary.
func (c *CartridgeOcean) HardwareButton(pressed bool, value uint8) {
}

// Read reads a byte from the current bank, based on the specified ROM interval and address. Returns the byte and success status.
func (c *CartridgeOcean) Read(addr uint16) uint8 {
	return c.banks[c.currBank][addr&0x1fff]
}

// IORead reads from a memory-mapped I/O address and returns the data and a success flag if the address is valid.
func (c *CartridgeOcean) IORead(addr uint16) (uint8, bool) {
	if (addr & 0xfff0) == 0xde00 {
		return c.lastData, true
	}
	return 0, false
}

// IOWrite processes writes to the cartridge's I/O address space, handling bank switching and configuration updates.
func (c *CartridgeOcean) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xde00 {
		currBank := (data & c.ioMask) & 0x3f
		c.currBank = currBank
		c.lastData = data
		//fmt.Printf("[OCEAN] Bank switching %x => %d, %d\n", addr, data, c.currBank)
		return true
	}
	return false
}

// IRQ handles the Interrupt Request (IRQ) signal for the CartridgeGeneric, enabling appropriate cartridge-specific behavior.
func (c *CartridgeOcean) IRQ(_ uint32) {
}

// IRQCLear clears the state of any active Interrupt Requests (IRQ) for the CartridgeGeneric.
func (c *CartridgeOcean) IRQCLear(_ uint32) {
}

// Config returns the Game line status, ExROM line state, and a boolean indicating successful configuration retrieval.
func (c *CartridgeOcean) Config() (uint8, uint8, bool) {
	return c.spec.Game, c.spec.ExRom, true
}

// Detach removes the cartridge from the system, performing any necessary cleanup or state reset operations.
func (c *CartridgeOcean) Detach() error {
	//TODO
	return nil
}

// EmulationRequired indicates if additional emulation logic is needed for the cartridge. Always returns false for this type.
func (c *CartridgeOcean) EmulationRequired() bool {
	return false
}

// Emulate handles the execution of the cartridge's emulation logic during each cycle of the system's operation.
func (c *CartridgeOcean) Emulate() {

}

// initBin initializes the cartridge by parsing binary data, validating it and segmenting it into fixed-size memory banks.
// It calculates the I/O mask and sets initial values for `lastData` and `currBank`. Returns an error if validation fails.
func (c *CartridgeOcean) initBin(ldr references.IC64CartridgeLoader) error {
	data := ldr.Data()
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	const cSize = 0x2000
	lCartridge := len(data)
	c.ioMask = uint8((lCartridge >> 13) - 1)
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
	c.lastData = 0
	c.currBank = 0
	return nil
}

// initCrt initializes the cartridge by reading chip headers from the provided CRTLoader and validates the chip data.
func (c *CartridgeOcean) initCrt(ldr references.IC64CartridgeLoader) error {
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
		if (chip.Bank() > 63) || ((chip.Start() != 0x8000) && (chip.Start() != 0xa000)) || (chip.Size() != 0x2000) {
			return fmt.Errorf("invalid chip bank")
		}
		c.banks = append(c.banks, chip.Data())
		romSize += int(chip.Size())
	}
	c.ioMask = uint8((romSize >> 13) - 1)
	c.lastData = 0
	c.currBank = 0
	return nil
}
