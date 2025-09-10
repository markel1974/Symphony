package final_cartridge_iii

import (
	"fmt"

	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	banksMin     = 4
	banksMax     = 16
	banksSize16k = 0x4000
	banksSize8k  = 0x2000
)

// CartridgeFinalCartridgeIII represents the Final Cartridge III hardware component for the system.
// It includes memory banks, control register, and operational states for emulation of cartridge behavior.
// banksL contains all "low" 8KB banks mapped to $8000-$9FFF memory addresses.
// banksH contains all "high" 8KB banks mapped to $A000-$BFFF memory addresses.
// banksTotal indicates the conceptual number of 16KB banks (can be 4 or 16).
// regEnabled denotes whether the control register is currently enabled.
// loaderId stores the identifier for the ROM being loaded.
// intervals defines the memory intervals for ROM mapping.
// banksCurrent tracks the currently active bank in use.
// game sets the GAME line status for cartridge interaction.
// exRom sets the EXROM line status for cartridge interaction.
// board represents the associated expansion interface of the C64.
// lastData stores the last accessed data byte from the cartridge.
type CartridgeFinalCartridgeIII struct {
	*component.BaseComponent
	board    references.IC64Expansion
	loaderId string
	spec     *references.C64CartridgeSpec
	//game          uint8
	//exRom         uint8
	reg           uint8
	regEnabled    bool
	freezeCounter int
	banksL        [][]byte
	banksH        [][]byte
	banksCurrent  uint8
	banksTotal    int
}

// NewCartridgeFinalCartridgeIII creates and initializes a new Final Cartridge III component for use in the emulation system.
// It registers the cartridge with the provided parent component, factory, label, and instance identifier.
func NewCartridgeFinalCartridgeIII(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *CartridgeFinalCartridgeIII {
	co := &CartridgeFinalCartridgeIII{
		BaseComponent: component.NewBaseComponent(),
		loaderId:      Identifier(),
		spec:          references.C64CartridgeSpec16K,
	}
	co.BaseComponent.Register(factory, parent, Identifier(), instance, co, references.IdIC64Cartridge(co, label, instance))
	return co
}

// New creates a new instance of IC64Cartridge using the provided parent component, factory, label, and instance identifier.
func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IC64Cartridge {
	return NewCartridgeFinalCartridgeIII(parent, factory, label, instance)
}

func (c *CartridgeFinalCartridgeIII) reset(hard bool) {
	//c.game, c.exRom, _ = references.C64CartridgeSpec16K.Data()
	c.spec = references.C64CartridgeSpec16K
	c.banksCurrent = 0
	c.regEnabled = true
	c.reg = 0
	c.freezeCounter = 0
	if hard {
		c.banksH = nil
		c.banksL = nil
		c.banksTotal = 0
	}
}

// Setup initializes the CartridgeFinalCartridgeIII instance, preparing it for use.
func (c *CartridgeFinalCartridgeIII) Setup() error {
	c.reset(true)
	return nil
}

// Reset initializes or reinitializes the cartridge state, setting core properties to their default values.
func (c *CartridgeFinalCartridgeIII) Reset() {
	c.reset(false)
}

// Bind initializes the cartridge by associating it with the provided board and loader, loading ROM data in the process.
func (c *CartridgeFinalCartridgeIII) Bind(board references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	c.board = board
	c.loaderId = ldr.Id()
	if catalog.Type(ldr.Type()) == catalog.TypeCrt {
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

// IRQ is a placeholder method for handling IRQ-related logic for the CartridgeFinalCartridgeIII.
func (c *CartridgeFinalCartridgeIII) IRQ(_ uint32) {
}

// IRQCLear clears the IRQ signal for the Final Cartridge III.
func (c *CartridgeFinalCartridgeIII) IRQCLear(_ uint32) {
}

// Config returns the Game line status, ExROM line state, and a boolean indicating successful configuration retrieval.
func (c *CartridgeFinalCartridgeIII) Config() (uint8, uint8, bool) {
	return c.spec.Game, c.spec.ExRom, true
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
			c.doFreeze()
			return
		}
	}
}

// Read recupera un byte dalla ROM della cartuccia.
// MODIFICATA per gestire correttamente la modalità Ultimax/Freezer.
func (c *CartridgeFinalCartridgeIII) Read(addr uint16) uint8 {
	// IOWrite must have set `c.intervals` to cover the entire $8000-$FFFF area when entering Freezer mode.
	bank := c.banksCurrent
	if int(bank) >= c.banksTotal {
		fmt.Printf("Operand: invalid bank %d >= %d", bank, c.banksTotal)
		return 0
	}
	//if addr >= 0xda00 {
	//	fmt.Printf("[%d][Operand] Ultimax mode read addr: 0x%x\n", c.board.Cycle(), addr)
	//}
	// The FCIII logic in this mode performs a mirroring
	// of the 16KB bank on the entire $8000-$FFFF area.
	// The mask & 0x3FFF handles this mirroring automatically.
	// Example: a read at $FFFA becomes a read at offset $3FFA of the bank.
	// A read at $BFFC becomes a read at offset $3FFC of the bank.
	offset16k := addr & 0x3FFF
	if offset16k < banksSize8k {
		return c.banksL[bank][offset16k]
	}
	offset8k := offset16k & 0x1FFF
	return c.banksH[bank][offset8k]
}

// IORead reads data from the cartridge memory at the specified address and returns the value with a success flag.
// It handles mirroring logic for $DE00 (second-to-last page) and $DFxx (last page), with special handling for $DFFF.
// Returns 0 and false if the address is invalid or out of bounds.
func (c *CartridgeFinalCartridgeIII) IORead(addr uint16) (uint8, bool) {
	target := addr & 0xff00
	if target == 0xde00 {
		// mirroring logic for $DE00 (second-to-last page).
		if int(c.banksCurrent) >= len(c.banksL) {
			return 0, false
		}
		offset := 0x1e00 + (addr & 0x00ff)
		return c.banksL[c.banksCurrent][offset], true
	}
	if target == 0xdf00 {
		if addr == 0xdfff {
			if (c.reg & 0x80) != 0 {
				return 0xff, true
			}
			val := ((c.reg - 1) & 2) / 2 * 0xFF
			return val, true
		}
		// for other addresses in $DFxx ($DF00-$DFFE), use standard mirroring of the last page
		if int(c.banksCurrent) >= len(c.banksL) {
			return 0, false
		}
		offset := 0x1f00 + (addr & 0x00ff)
		return c.banksL[c.banksCurrent][offset], true
	}
	return 0, false
}

// IOWrite writes data to the cartridge control register if enabled and the address is within the valid range (0xDF00-0xDFFF).
// Updates internal state, including the configuration of EXROM, GAME, and memory intervals. Returns true if the write was successful.
func (c *CartridgeFinalCartridgeIII) IOWrite(addr uint16, data uint8) bool {
	if !c.regEnabled || (addr&0xff) != 0xff {
		return false
	}
	c.reg = data
	c.regEnabled = ((data >> 7) & 1) == 0
	if (data & 0x80) != 0 {
		c.banksCurrent = 0
	} else {
		c.banksCurrent = data & (uint8(c.banksTotal) - 1)
	}
	command := (data >> 4) & 0x03
	switch command {
	case 0b00:
		//c.game, c.exRom, _ = references.C64CartridgeSpec16K.Data()
		c.spec = references.C64CartridgeSpec16K
		c.board.GameExRomConfigChanged()
	case 0b01:
		c.doFreeze()
		return true
	case 0b10:
		c.spec = references.C64CartridgeSpec8K
		//c.game, c.exRom, _ = references.C64CartridgeSpec8K.Data()
		c.board.GameExRomConfigChanged()
	case 0b11:
		c.spec = references.C64CartridgeSpecOff
		//c.game, c.exRom, _ = references.C64CartridgeSpecOff.Data()
		c.board.GameExRomConfigChanged()
	}
	return true
}

// doFreeze increments the freeze counter, modifies bank switching, and triggers configuration updates for cartridge mode.
func (c *CartridgeFinalCartridgeIII) doFreeze() {
	c.freezeCounter++
	if c.freezeCounter != 1 {
		return
	}
	c.banksCurrent = uint8(c.banksTotal - 1)
	c.spec = references.C64CartridgeSpecUltimax
	//c.game, c.exRom, _ = spec.Data()
	c.board.GameExRomConfigChanged()
	c.board.NMITrigger()

	t := c.board.CycleAlarm("Freezer", func(mainCpuClk uint64, offset uint64) {
		c.freezeCounter = 0
	})
	_ = t.Set(5000000)
}

// initCrt initializes the CRT by loading cartridge data via the provided loader and validates the chip structure and sizes.
func (c *CartridgeFinalCartridgeIII) initCrt(ldr references.IC64CartridgeLoader) error {
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
		if chip.Size() != banksSize16k {
			return fmt.Errorf("invalid chip size")
		}
		offset := int(chip.Bank()) * banksSize16k
		if (offset + int(chip.Size())) > len(rawCart) {
			return fmt.Errorf("bank data %d out of bounds", chip.Bank())
		}
		copy(rawCart[offset:], chip.Data())
		banksLoaded++
	}
	if banksLoaded != banksMin && banksLoaded != banksMax {
		return fmt.Errorf("unsupported banks number in crt file (%d)", banksLoaded)
	}
	data := rawCart[:(banksLoaded * banksSize16k)]
	return c.loadData(data)
}

// loadData loads the provided ROM data into the cartridge while validating its size and formatting it into banks.
// It ensures the ROM data size is a multiple of 16KB and supports only specific numbers of banks (4 or 16).
// Returns an error if the ROM size is invalid or unsupported.
func (c *CartridgeFinalCartridgeIII) loadData(data []byte) error {
	size := len(data)
	if (size % banksSize16k) != 0 {
		return fmt.Errorf("rom size %d is not a multiple of 16KB", size)
	}
	total := size / banksSize16k
	if total != banksMin && total != banksMax {
		return fmt.Errorf("number of banks (%d) not supported", total)
	}
	c.banksCurrent = 0
	c.banksTotal = total
	c.banksL = make([][]byte, c.banksTotal)
	c.banksH = make([][]byte, c.banksTotal)
	for i := 0; i < c.banksTotal; i++ {
		start := i * banksSize16k
		c.banksL[i] = data[start:(start + banksSize8k)]
		c.banksH[i] = data[(start + banksSize8k):(start + banksSize16k)]
	}
	return nil
}

// initBin initializes the cartridge using binary file data from the provided loader.
// Returns an error if the BIN file size is not 64KB or 256KB or if validation fails.
func (c *CartridgeFinalCartridgeIII) initBin(ldr references.IC64CartridgeLoader) error {
	data := ldr.Data()
	if err := catalog.ValidateCartridge(data); err != nil {
		return err
	}
	size := len(data)
	if size != 0x10000 && size != 0x40000 {
		return fmt.Errorf("bin file size (%d) is not 64KB or 256KB", size)
	}
	return c.loadData(data)
}
