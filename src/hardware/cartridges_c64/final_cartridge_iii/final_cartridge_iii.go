package final_cartridge_iii

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/catalog"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

const (
	minBanks = 4
	maxBanks = 16

	bankSize16k = 0x4000
	bankSize8k  = 0x2000
)

// CartridgeFinalCartridgeIII represents the Final Cartridge III hardware component for the system.
// It includes memory banks, control register, and operational states for emulation of cartridge behavior.
// romLBanks contains all "low" 8KB banks mapped to $8000-$9FFF memory addresses.
// romHBanks contains all "high" 8KB banks mapped to $A000-$BFFF memory addresses.
// numBanks indicates the conceptual number of 16KB banks (can be 4 or 16).
// regEnabled denotes whether the control register is currently enabled.
// loaderId stores the identifier for the ROM being loaded.
// intervals defines the memory intervals for ROM mapping.
// currBank tracks the currently active bank in use.
// game sets the GAME line status for cartridge interaction.
// exRom sets the EXROM line status for cartridge interaction.
// board represents the associated expansion interface of the C64.
// lastData stores the last accessed data byte from the cartridge.
type CartridgeFinalCartridgeIII struct {
	*component.BaseComponent
	romLBanks  [][]byte // Tutti i banchi "bassi" da 8KB ($8000-$9FFF)
	romHBanks  [][]byte // Tutti i banchi "alti" da 8KB ($A000-$BFFF)
	numBanks   int      // Numero di banchi concettuali da 16KB (4 o 16)
	reg        uint8
	regEnabled bool
	loaderId   string
	intervals  references.RomInterval
	currBank   uint8
	game       uint8
	exRom      uint8
	freeze     bool
	board      references.IExpansionC64
}

// NewCartridgeFinalCartridgeIII creates and initializes a new Final Cartridge III component for use in the emulation system.
// It registers the cartridge with the provided parent component, factory, label, and instance identifier.
func NewCartridgeFinalCartridgeIII(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeFinalCartridgeIII {
	v := references.GetCartridgeSpec(references.CartridgeMode16K)
	co := &CartridgeFinalCartridgeIII{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		reg:           0,
		game:          v.Game,
		exRom:         v.ExRom,
		intervals:     v.IntervalLow | v.IntervalHigh,
		currBank:      0,
		regEnabled:    true,
		freeze:        false,
	}
	co.BaseComponent.Register(factory, parent, Identifier(), co, references.IdICartridgeC64(co, label, instance))
	return co
}

// New creates a new instance of ICartridgeC64 using the provided parent component, factory, label, and instance identifier.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.ICartridgeC64 {
	return NewCartridgeFinalCartridgeIII(parent, factory, label, instance)
}

// Setup initializes the CartridgeFinalCartridgeIII instance, preparing it for use.
func (c *CartridgeFinalCartridgeIII) Setup() error {
	return nil
}

// Bind initializes the cartridge by associating it with the provided board and loader, loading ROM data in the process.
func (c *CartridgeFinalCartridgeIII) Bind(board references.IExpansionC64, ldr references.ICartridgeLoaderC64) error {
	c.board = board
	c.loaderId = ldr.GetId()
	if catalog.Type(ldr.GetType()) == catalog.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr)
}

// Connect initializes the cartridge and integrates it with the system, ensuring readiness for operation.
func (c *CartridgeFinalCartridgeIII) Connect() error {
	return nil
}

// Internal checks if the cartridge operates internally or externally.
func (c *CartridgeFinalCartridgeIII) Internal() bool {
	return false
}

func (c *CartridgeFinalCartridgeIII) Reset() {
	v := references.GetCartridgeSpec(references.CartridgeMode16K)
	c.game = v.Game
	c.exRom = v.ExRom
	c.intervals = v.IntervalLow | v.IntervalHigh
	c.currBank = 0
	c.regEnabled = true
	c.reg = 0
	c.freeze = false
}

// IRQ is a placeholder method for handling IRQ-related logic for the CartridgeFinalCartridgeIII.
func (c *CartridgeFinalCartridgeIII) IRQ(_ uint32) {
}

// IRQCLear clears the IRQ signal for the Final Cartridge III.
func (c *CartridgeFinalCartridgeIII) IRQCLear(_ uint32) {
}

// GetExRom returns the current value of the ExRom signal.
func (c *CartridgeFinalCartridgeIII) GetExRom() uint8 {
	return c.exRom
}

// GetGame retrieves the current `game` state of the cartridge as a uint8.
func (c *CartridgeFinalCartridgeIII) GetGame() uint8 {
	return c.game
}

// EmulationRequired checks if emulation is required for the cartridge and returns false, indicating no emulation is needed.
func (c *CartridgeFinalCartridgeIII) EmulationRequired() bool {
	return false
}

// Emulate executes the emulation process for the Final Cartridge III by managing ROM operations and hardware interaction.
func (c *CartridgeFinalCartridgeIII) Emulate() {
}

// Detach disengages the cartridge from the connected system and prepares it for cleanup or removal.
func (c *CartridgeFinalCartridgeIII) Detach() error {
	return nil
}

// GetLoaderId retrieves the identifier of the loader associated with the cartridge.
func (c *CartridgeFinalCartridgeIII) GetLoaderId() string {
	return c.loaderId
}

// HardwareButton triggers a non-maskable interrupt (NMI) on the board if the button is pressed.
func (c *CartridgeFinalCartridgeIII) HardwareButton(pressed bool, value uint8) {
	if pressed {
		if c.board.AECAvailable() && c.board.BusAvailable() {
			c.freezer()
			/*
				var funcId int
				fmt.Println("HardwareButton: NMI triggered")
				c.board.NMITrigger()
				funcId = c.board.RamSetWriteTrigger(0xfffa, func(addr uint16, _ uint8) {
					spec := references.GetCartridgeSpec(references.CartridgeModeUltimax)
					c.exRom = spec.ExRom
					c.game = spec.Game
					c.intervals = spec.IntervalHigh | spec.IntervalLow
					c.board.GameExRomConfigChanged()
					c.board.RamRemoveWriteTrigger(0xfffa, funcId)
				})
			*/

			//t := c.board.CycleAlarm("Freezer", func(mainCpuClk uint64, offset uint64) {
			//  spec := references.GetCartridgeSpec(references.CartridgeMode8K)
			//	c.exRom = spec.ExRom
			//	c.game = spec.Game
			//	c.intervals = spec.IntervalHigh | spec.IntervalLow
			//	c.board.GameExRomConfigChanged()
			//})
			//_ = t.Set(10)
			return
		}
	}
}

// Write attempts to write data to the cartridge at the specified address within the given ROM interval.
// Returns true if the operation is invalid for the cartridge, otherwise false.
func (c *CartridgeFinalCartridgeIII) Write(i references.RomInterval, addr uint16, data uint8) bool {
	if (i & c.intervals) != 0 {
		fmt.Printf("Write: can't write [bank %d] %x => %d\n", c.currBank, addr, data)
		return true
	}
	return false
}

// Read recupera un byte dalla ROM della cartuccia.
// MODIFICATA per gestire correttamente la modalità Ultimax/Freezer.
func (c *CartridgeFinalCartridgeIII) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	// IOWrite must have set `c.intervals` to cover the entire $8000-$FFFF area when entering Freezer mode.
	if (i & c.intervals) == 0 {
		fmt.Printf("Read: invalid interval %d != %d", c.intervals, i)
		return 0, false
	}
	if int(c.currBank) >= c.numBanks {
		log.Printf("Read: invalid bank %d >= %d", c.currBank, c.numBanks)
		return 0, false
	}
	if addr >= 0xda00 {
		fmt.Printf("[%d][Read] Ultimax mode read addr: 0x%x\n", c.board.Cycle(), addr)
	}
	// The FCIII logic in this mode performs a mirroring
	// of the 16KB bank on the entire $8000-$FFFF area.
	// The mask & 0x3FFF handles this mirroring automatically.
	// Example: a read at $FFFA becomes a read at offset $3FFA of the bank.
	// A read at $BFFC becomes a read at offset $3FFC of the bank.
	offset16k := addr & 0x3FFF
	if offset16k < bankSize8k {
		return c.romLBanks[c.currBank][offset16k], true
	} else {
		offset8k := offset16k & 0x1FFF
		return c.romHBanks[c.currBank][offset8k], true
	}
}

// IORead handles I/O read requests by decoding the given address and fetching the corresponding data from ROM banks.
// Returns the data byte and a boolean indicating success or failure based on the address validity and current state.
func (c *CartridgeFinalCartridgeIII) IORead(addr uint16) (uint8, bool) {
	var offset uint16
	switch addr & 0xff00 {
	case 0xde00:
		offset = 0x1e00 + (addr & 0x00ff)
	case 0xdf00:
		offset = 0x1f00 + (addr & 0x00ff)
	default:
		return 0, false
	}
	if int(c.currBank) >= len(c.romLBanks) {
		return 0, false
	}
	return c.romLBanks[c.currBank][offset], true
}

// IOWrite writes data to the cartridge control register if enabled and the address is within the valid range (0xDF00-0xDFFF).
// Updates internal state, including the configuration of EXROM, GAME, and memory intervals. Returns true if the write was successful.
func (c *CartridgeFinalCartridgeIII) IOWrite(addr uint16, data uint8) bool {
	if !c.regEnabled || (addr&0xff00) != 0xdf00 {
		return false
	}
	freeze := c.freeze
	c.reg = data
	c.freeze = (data & 0x40) == 0
	c.regEnabled = ((data >> 7) & 1) == 0
	c.exRom = (data >> 4) & 1
	c.game = (data >> 5) & 1
	c.currBank = data & (uint8(c.numBanks) - 1)

	if v := references.GetCartridgeSpecFromDetails(c.game, c.exRom); v != nil {
		if v.IntervalLow|v.IntervalHigh != c.intervals {
			fmt.Printf("[%d][IOWrite] configuration changed: game %d, exRom %d\n", c.board.Cycle(), c.game, c.exRom)
		}
		c.intervals = v.IntervalHigh | v.IntervalLow
	}
	c.board.GameExRomConfigChanged()
	if !freeze && c.freeze {
		c.freezer()
		fmt.Printf("[%d][IOWrite] current data before NMI trigger %08b (0x%x) addr: 0x%x - game: %d, exRom: %d\n", c.board.Cycle(), data, data, addr, c.game, c.exRom)
	}
	return true
}

func (c *CartridgeFinalCartridgeIII) freezer() {
	spec := references.GetCartridgeSpec(references.CartridgeModeUltimax)
	c.currBank = uint8(c.numBanks - 1)
	c.exRom = spec.ExRom
	c.game = spec.Game
	c.intervals = spec.IntervalHigh | spec.IntervalLow
	c.board.GameExRomConfigChanged()
	c.board.NMITrigger()
}

// initCrt initializes the CRT by loading cartridge data via the provided loader and validates the chip structure and sizes.
func (c *CartridgeFinalCartridgeIII) initCrt(ldr references.ICartridgeLoaderC64) error {
	rawCart := make([]byte, 256*1024)
	banksLoaded := 0
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
		if chip.Size() != bankSize16k {
			return fmt.Errorf("invalid chip size")
		}
		offset := int(chip.Bank()) * bankSize16k
		if (offset + int(chip.Size())) > len(rawCart) {
			return fmt.Errorf("bank data %d out of bounds", chip.Bank())
		}
		copy(rawCart[offset:], chip.Data())
		banksLoaded++
	}
	if banksLoaded != minBanks && banksLoaded != maxBanks {
		return fmt.Errorf("unsupported banks number in crt file (%d)", banksLoaded)
	}
	data := rawCart[:(banksLoaded * bankSize16k)]
	return c.loadData(data)
}

// loadData loads the provided ROM data into the cartridge while validating its size and formatting it into banks.
// It ensures the ROM data size is a multiple of 16KB and supports only specific numbers of banks (4 or 16).
// Returns an error if the ROM size is invalid or unsupported.
func (c *CartridgeFinalCartridgeIII) loadData(data []byte) error {
	size := len(data)
	if (size % bankSize16k) != 0 {
		return fmt.Errorf("rom size %d is not a multiple of 16KB", size)
	}
	numBanks := size / bankSize16k
	if numBanks != minBanks && numBanks != maxBanks {
		return fmt.Errorf("number of banks (%d) not supported", numBanks)
	}
	c.romLBanks = make([][]byte, numBanks)
	c.romHBanks = make([][]byte, numBanks)
	c.numBanks = numBanks
	for i := 0; i < numBanks; i++ {
		start := i * bankSize16k
		c.romLBanks[i] = data[start:(start + bankSize8k)]
		c.romHBanks[i] = data[(start + bankSize8k):(start + bankSize16k)]
	}
	return nil
}

// initBin initializes the cartridge using binary file data from the provided loader.
// Returns an error if the BIN file size is not 64KB or 256KB or if validation fails.
func (c *CartridgeFinalCartridgeIII) initBin(ldr references.ICartridgeLoaderC64) error {
	data := ldr.GetData()
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	size := len(data)
	if size != 0x10000 && size != 0x40000 {
		return fmt.Errorf("bin file size (%d) is not 64KB or 256KB", size)
	}
	return c.loadData(data)
}

/*
// ReadOld retrieves a byte from the specified address and interval. Returns the byte and a bool indicating success.
func (c *CartridgeFinalCartridgeIII) ReadOld(i references.RomInterval, addr uint16) (uint8, bool) {
	if (i & c.intervals) == 0 {
		return 0, false
	}
	offset := addr & 0x1fff
	if addr >= 0x8000 && addr <= 0x9fff {
		if int(c.currBank) >= len(c.romLBanks) {
			return 0, false
		}
		return c.romLBanks[c.currBank][offset], true
	} else if addr >= 0xa000 && addr <= 0xBfff {
		if int(c.currBank) >= len(c.romHBanks) {
			return 0, false
		}
		return c.romHBanks[c.currBank][offset], true
	}
	return 0, false
}
*/
