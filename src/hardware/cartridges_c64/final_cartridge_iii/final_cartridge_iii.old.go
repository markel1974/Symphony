package final_cartridge_iii

/*
import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/catalog"
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeFinalCartridgeIII represents a cartridge model with bank switching and IO interaction mechanisms for ROM emulation.
type CartridgeFinalCartridgeIII struct {
	*component.BaseComponent
	loaderId     string
	intervals    references.RomInterval
	banks        [][]byte
	currBank     uint8
	game         uint8
	exRom        uint8
	board        references.IExpansionC64
	lastData     uint8
	freezeStatus uint8
}

// GetType returns the type identifier of the Ocean cartridge as an integer constant.
func GetType() int {
	return catalog.CartridgeFinalIII
}

// NewCartridgeFinalCartridgeIII creates and returns a new instance of the Ocean Cartridge conforming to the ICartridgeC64 interface.
func NewCartridgeFinalCartridgeIII(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeFinalCartridgeIII {
	v := references.GetCartridgeSpec(references.CartridgeModeUltimax)
	co := &CartridgeFinalCartridgeIII{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		game:          v.Game,
		exRom:         v.ExRom,
		intervals:     v.IntervalLow | v.IntervalHigh,
		lastData:      0,
		freezeStatus:  0,
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
	c.lastData = 0
	c.freezeStatus = 0
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

// ReadOld reads a byte from the current bank, based on the specified ROM interval and address. Returns the byte and success status.
func (c *CartridgeFinalCartridgeIII) ReadOld(i references.RomInterval, addr uint16) (uint8, bool) {
	if i&c.intervals == 0 {
		return 0, false
	}
	var effectiveBank uint8
	var offset uint16
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		// Usa il banco corrente per l'area bassa
		effectiveBank = c.currBank
		offset = addr - 0x8000
	case addr >= 0xA000 && addr <= 0xBFFF:
		// Usa il banco successivo per l'area alta
		effectiveBank = c.currBank + 1
		offset = addr - 0xA000
	default:
		return 0, false
	}
	// Controlla che il banco calcolato sia valido
	if int(effectiveBank) >= len(c.banks) {
		return 0, false
	}
	//v := c.banks[effectiveBank][offset]
	//fmt.Printf("CartridgeFinalCartridgeIII read [bank %d] %x => %d\n", effectiveBank, addr, v)
	return c.banks[effectiveBank][offset], true
}

//var _counter = 0

func (c *CartridgeFinalCartridgeIII) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	if i&c.intervals == 0 {
		return 0, false
	}

	var effectiveBank uint8
	var offset uint16

	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		effectiveBank = c.currBank
		offset = addr - 0x8000

	case addr >= 0xA000 && addr <= 0xBFFF:
		// In modalità non-freezer (game=1), questa parte non viene mai raggiunta.
		// Viene raggiunta solo in modalità freezer (game=0).
		effectiveBank = c.currBank + 1
		offset = addr - 0xA000

	// Nei dump da 64KB, i banchi 14 e 15 sono mappati ai banchi 6 e 7.
	case addr >= 0xC000 && addr <= 0xDFFF:
		effectiveBank = 0 // 6 Banco fisso per il Toolkit
		offset = addr - 0xC000

	case addr >= 0xE000 && addr <= 0xFFFF:
		effectiveBank = 1 // 7 Banco fisso per Freezer e Reset Vector
		offset = addr - 0xE000

	//case addr >= 0xE000 && addr <= 0xFFFF:
	// Mappatura Ultimax: E000-FFFF
	//	effectiveBank = c.currBank & 0x03 // Solo 4 banchi in modalità Ultimax
	//	fmt.Println(effectiveBank)
	//	offset = addr - 0xE000
	default:
		return 0, false
	}

	// --- ECCO LA CORREZIONE ---
	// Nei dump da 64KB, i banchi 14 e 15 sono mappati ai banchi 6 e 7.
	//case addr >= 0xC000 && addr <= 0xDFFF:
	//	effectiveBank = 6 // Banco fisso per il Toolkit
	//	offset = addr - 0xC000

	//case addr >= 0xE000 && addr <= 0xFFFF:
	//	effectiveBank = 7 // Banco fisso per Freezer e Reset Vector
	//	offset = addr - 0xE000

	// Controllo di sicurezza
	if int(effectiveBank) >= len(c.banks) {
		// Potresti voler loggare un avviso qui, ma non restituire un errore fatale
		return 0, true
	}

	v := c.banks[effectiveBank][offset]
	//fmt.Printf("CartridgeFinalCartridgeIII read [bank %d] %x => %d, counter %d\n", effectiveBank, addr, v, _counter)
	//_counter++
	return v, true
}

// IORead reads from a memory-mapped I/O address and returns the data and a success flag if the address is valid.
func (c *CartridgeFinalCartridgeIII) IORead(addr uint16) (uint8, bool) {
	switch addr & 0xFFF0 {
	case 0xDF00:
		return c.lastData, true
	case 0xDF80:
		return c.freezeStatus, true
	}
	return 0, false
}

// IOWrite processes writes to the cartridge's I/O address space, handling bank switching and configuration updates.
func (c *CartridgeFinalCartridgeIII) IOWrite(addr uint16, data uint8) bool {
	fmt.Printf("CartridgeFinalCartridgeIII IOWrite %x, %x\n", addr, data)
	if (addr & 0xFF00) != 0xDF00 {
		return false
	}
	c.lastData = data
	c.currBank = data & 0x1F
	// Modalità standard (ExRom=0, Game=1)
	spec := references.GetCartridgeSpec(references.CartridgeModeUltimax)
	switch {
	case c.currBank == 0x0C:
		// Disabilita cartuccia (ExRom=1, Game=1)
		spec = references.GetCartridgeSpec(references.CartridgeModeOff)
	case c.currBank >= 0x10 && c.currBank <= 0x1F:
		// Modalità freezer (ExRom=0, Game=0)
		spec = references.GetCartridgeSpec(references.CartridgeMode16K)
	}
	c.exRom = spec.ExRom
	c.game = spec.Game
	c.intervals = spec.IntervalLow | spec.IntervalHigh
	c.board.GameExRomConfigChanged()
	return true
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

// initCrt initializes the cartridge by reading chip headers from the provided CRTLoader and validates the chip data.
func (c *CartridgeFinalCartridgeIII) initCrt(ldr references.ICartridgeLoaderC64) error {
	const totalBanks = 32
	c.banks = make([][]byte, totalBanks)
	for i := range c.banks {
		c.banks[i] = make([]byte, 0x2000)
	}
	for {
		chip, err := ldr.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if chip.Start() != 0x8000 {
			return fmt.Errorf("invalid chip start")
		}
		if chip.Size() != 0x4000 {
			return fmt.Errorf("invalid chip size")
		}
		index := chip.Bank() * 2
		bankIndex8kLow := index
		bankIndex8kHigh := index + 1
		if int(bankIndex8kHigh) >= len(c.banks) {
			return fmt.Errorf("bank index %d out of range (0-%d)", bankIndex8kHigh, len(c.banks))
		}
		copy(c.banks[bankIndex8kLow], chip.Data()[0x0000:0x2000])
		copy(c.banks[bankIndex8kHigh], chip.Data()[0x2000:0x4000])

		fmt.Printf("Bank Low %d: High %d\n", bankIndex8kLow, bankIndex8kHigh)
	}
	c.currBank = 0
	return nil
}

// initBin initializes the cartridge by parsing binary data, validating it and segmenting it into fixed-size memory banks.
// It calculates the I/O mask and sets initial values for `lastData` and `currBank`. Returns an error if validation fails.
func (c *CartridgeFinalCartridgeIII) initBin(ldr references.ICartridgeLoaderC64) error {
	data := ldr.GetData()
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	const cSize = 0x2000
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
	c.lastData = 0
	return nil
}

*/

/*
func (c *CartridgeFinalCartridgeIII) isAddressActive(addr uint16) bool {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		return c.game == 1 && c.exRom == 0
	case addr >= 0xA000 && addr <= 0xBFFF:
		return c.game == 0 && c.exRom == 0
	}
	return false
}

func (c *CartridgeFinalCartridgeIII) calculateBankOffset(addr uint16) (uint8, uint16, bool) {
	switch {
	case addr >= 0x8000 && addr <= 0x9FFF:
		return c.currBank, addr - 0x8000, true
	case addr >= 0xA000 && addr <= 0xBFFF:
		return c.currBank + 1, addr - 0xA000, true
	}
	return 0, 0, false
}


func (c *CartridgeFinalCartridgeIII) updateBankMapping() {
	switch {
	case c.currBank == 0x0C:
		// Disabilita cartuccia (ExRom=1, Game=1)
		c.exRom = 1
		c.game = 1
	case c.currBank >= 0x10 && c.currBank <= 0x1F:
		// Modalità freezer (ExRom=0, Game=0)
		c.exRom = 0
		c.game = 0
	default:
		// Modalità standard (ExRom=0, Game=1)
		c.exRom = 0
		c.game = 1
	}
}
*/
