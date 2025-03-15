package easyflash

//https://skoe.de/easyflash/files/devdocs/EasyFlash-ProgRef.pdf

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/common/filler"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/easyflash/flash"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/easyflash/snapshot"
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/loader"
	"github.com/markel1974/c64emu/src/references"
	"io"
	"log"
	"os"
)

// CartridgeEasyFlash represents the implementation of an EasyFlash cartridge for execution on supported hardware.
// It contains fields for configuration, state management, memory mapping, flash state, and cartridge-specific options.
type CartridgeEasyFlash struct {
	*component.BaseComponent
	factory         references.IComponentFactory
	loaderId        string
	board           references.IC64Expansion
	intervalLo      references.RomInterval
	intervalHi      references.RomInterval
	memoryConfigIdx int
	game            uint8
	exRom           uint8
	stateLow        *flash.Flash040 /* the 29F040B statemachine */
	stateHigh       *flash.Flash040 /* the 29F040B statemachine */
	jumper          bool
	writeEnabled    bool    /* writing back to crt enabled */
	optimize        bool    /* optimizing crt enabled */
	register00      uint8   /* backup of the registers */
	register02      uint8   /* backup of the registers */
	ram             []uint8 /* extra RAM */
	filename        string  /* filename when attached */
	filetype        loader.Type
	led             bool
	updateEApi      bool
}

// GetType returns the cartridge type constant representing an EasyFlash cartridge.
func GetType() int {
	return loader.CARTRIDGE_EASYFLASH
}

// New creates and returns a new instance of a CartridgeEasyFlash implementing the IC64Cartridge interface.
func New(parent component.IComponent, factory references.IComponentFactory, suffix string) references.IC64Cartridge {
	ef := &CartridgeEasyFlash{
		BaseComponent:   component.NewBaseComponent("easyFlash", suffix),
		factory:         factory,
		loaderId:        "easyFlash",
		game:            1,
		exRom:           1,
		register00:      0,
		register02:      0,
		intervalHi:      references.ROM_HI_1,
		intervalLo:      references.ROM_LO,
		stateLow:        nil,
		stateHigh:       nil,
		filetype:        0,
		jumper:          false,
		led:             false,
		ram:             make([]byte, CartRamSize),
		memoryConfigIdx: -1,
		updateEApi:      true,
	}
	component.Register(parent, ef)
	return ef
}

// Setup initializes the CartridgeEasyFlash instance with the provided board and CRT loader data.
func (c *CartridgeEasyFlash) Setup(board references.IC64Expansion, ldr references.IC64CartridgeLoader) error {
	var rawCart []byte
	c.board = board
	c.loaderId = ldr.GetId()
	c.game = uint8(ldr.Game())
	c.exRom = uint8(ldr.ExRom())
	rp := filler.New(255, 2, 1, 0x100, 255, 0, 0, 0)
	rp.InitWithPattern(c.ram, CartRamSize)
	c.filename = ldr.Name()
	c.filetype = 0
	var err error
	if loader.Type(ldr.GetType()) == loader.TypeCrt {
		c.filetype = loader.TypeCrt
		if rawCart, err = c.crtAttach(ldr); err != nil {
			return err
		}
	} else {
		c.filetype = loader.TypeBin
		if rawCart, err = c.binAttach(ldr); err != nil {
			return err
		}
	}
	c.register00 = 0
	c.controlUpdate(0, false)
	c.initialize(rawCart)
	return nil
}

// Reset reinitializes the CartridgeEasyFlash state, clearing any temporary data and restoring default settings.
func (c *CartridgeEasyFlash) Reset() {

}

// GetLoaderId retrieves the unique identifier of the CartridgeEasyFlash instance.
func (c *CartridgeEasyFlash) GetLoaderId() string {
	return c.loaderId
}

// GetExRom returns the current value of the exRom property for the CartridgeEasyFlash instance as an unsigned 8-bit integer.
func (c *CartridgeEasyFlash) GetExRom() uint8 {
	return c.exRom
}

// GetGame retrieves the currently selected game index from the EasyFlash cartridge.
func (c *CartridgeEasyFlash) GetGame() uint8 {
	return c.game
}

// EmulationRequired determines whether emulation is required for the CartridgeEasyFlash. Always returns false.
func (c *CartridgeEasyFlash) EmulationRequired() bool {
	return false
}

// Emulate starts the emulation process for the CartridgeEasyFlash, executing associated operations and behaviors.
func (c *CartridgeEasyFlash) Emulate() {
}

// SetJumper sets the jumper state for the EasyFlash cartridge to the specified boolean value.
func (c *CartridgeEasyFlash) SetJumper(j bool) {
	c.jumper = j
}

// SetWriteEnabled updates the write-enabled state of the cartridge.
func (c *CartridgeEasyFlash) SetWriteEnabled(e bool) {
	c.writeEnabled = e
}

// SetOptimize sets the optimize flag for the CartridgeEasyFlash instance.
func (c *CartridgeEasyFlash) SetOptimize(o bool) {
	c.optimize = o
}

// SetUpdateEApi configures the updateEApi flag to enable or disable the extended API updating functionality.
func (c *CartridgeEasyFlash) SetUpdateEApi(e bool) {
	c.updateEApi = e
}

// initialize prepares the interleaved low and high memory banks from the raw cartridge data
// and initializes the flash states.
func (c *CartridgeEasyFlash) initialize(rawCart []byte) {
	low := make([]byte, 0x80000)
	high := make([]byte, 0x80000)
	// split interleaved low and high banks
	for i := uint(0); i < NBanks; i++ {
		const size = 0x2000
		start := i * size
		p1 := i * (size * 2)
		p2 := p1 + size
		copy(low[start:start+size], rawCart[p1:p1+size])
		copy(high[start:start+size], rawCart[p2:p2+size])
	}
	c.stateLow = flash.NewFlash040(c.board, flash.KindB, low)
	c.stateHigh = flash.NewFlash040(c.board, flash.KindB, high)

	if c.updateEApi {
		eApiHeader := high[EApiStartAddress : EApiStartAddress+len(EApiHeader)]
		if bytes.Compare(eApiHeader, []byte(EApiHeader)) == 0 {
			//updating eapi
			//eApi := make([]byte, 17)
			//for k := 0; k < 16; k++ {
			//	eApi[k] = c.stateHigh.Peek(uint32(0x1804 + k))
			//}
			_ = c.stateHigh.StoreInterval(0x1800, _eApiAM29f040)
		}
	}
}

// controlUpdate handles updates to the cartridge control register
// and adjusts the cartridge mode and configuration accordingly.
// It manages changes to GAME/EXROM signal settings, LED state, and memory configuration based on input values.
// Logs warnings for unsupported modes and invokes a board configuration change if required.
func (c *CartridgeEasyFlash) controlUpdate(value uint8, update bool) {
	c.register02 = value & 0x87 // led, mode, exrom, game [led 0x80, other 0x07]
	mode := references.CartridgeModeOff
	mxg := value & 0x07
	switch mxg {
	case 0:
		//GAME from jumper, EXROM high (i.e., Ultimax or Off)
		if !c.jumper {
			mode = references.CartridgeModeUltimax
		} else {
			mode = references.CartridgeModeOff
		}
	case 1, 3:
		//Reserved, don’t use these
		log.Printf("CartridgeEasyFlash: unsupported mode %d", mode)
	case 2:
		//GAME from jumper, EXROM low (i.e. 16K or 8K)
		if !c.jumper {
			mode = references.CartridgeMode16K
		} else {
			mode = references.CartridgeMode8K
		}
	case 4:
		//Cartridge ROM off (RAM at $DF00 still available)
		mode = references.CartridgeModeOff
	case 5:
		//Ultimax (Low bank at $8000, high bank at $e000) GAME = 0, EXROM = 1
		mode = references.CartridgeModeUltimax
	case 6:
		//8k Cartridge (Low bank at $8000) GAME = 1, EXROM = 0
		mode = references.CartridgeMode8K
	case 7:
		//16k cartridge (Low bank at $8000, high bank at $a000)
		mode = references.CartridgeMode16K
	}
	if int(mode) != c.memoryConfigIdx {
		c.memoryConfigIdx = int(mode)
		v := references.GetCartridgeSpec(mode)
		c.game = v.Game
		c.exRom = v.ExRom
		c.intervalLo = v.IntervalLow
		c.intervalHi = v.IntervalHigh
		//fmt.Println("EASYFLASH MEMORY CONFIG CHANGED:", mxg, "exrom:", c.exRom, "game:", c.game, "LO", c.intervalLo, "HIGH", c.intervalHi)
		if update {
			c.board.GameExRomConfigChanged()
		}
	}
	if led := value&0x80 == 0x80; led != c.led {
		c.led = led
		fmt.Println("EASYFLASH LED:", c.led)
	}
}

// Write attempts to handle a write operation to the EasyFlash cartridge but currently does not implement writing logic.
func (c *CartridgeEasyFlash) Write(i references.RomInterval, _ uint16, _ uint8) bool {
	if c.intervalLo == i {
		fmt.Printf("EASYFLASH Write LOW NOT DEFINED\n")
	} else if c.intervalHi == i {
		fmt.Printf("EASYFLASH Write HIGH NOT DEFINED\n")
	}
	return false
}

// Read retrieves a byte of data from the cartridge memory based on the provided ROM interval and address.
func (c *CartridgeEasyFlash) Read(i references.RomInterval, addr uint16) (uint8, bool) {
	if c.intervalLo == i {
		return c.romLRead(addr), true
	} else if c.intervalHi == i {
		return c.romHRead(addr), true
	}
	return 0, false
}

// IORead reads a value from the specified address within the EasyFlash cartridge's I/O space.
// Returns the value and a boolean indicating if the read is valid.
func (c *CartridgeEasyFlash) IORead(addr uint16) (uint8, bool) {
	bank := addr & 0xff00
	if bank == 0xdf00 {
		v := c.io2Read(addr)
		return v, true
	} else if bank == 0xde00 {
		// read is never valid, regs are write-only
		fmt.Printf("EASYFLASH WARNING: regs $de00-$deff are write only -> $%x\n", addr)
		return 0, false
	}
	return 0, false
}

// IOWrite handles input/output write operations for the CartridgeEasyFlash.
// It processes writes to specific memory address ranges (0xde00 and 0xdf00).
// Updates internal registers and control states based on the provided data.
// Returns true if the operation is valid and the data is handled, otherwise false.
func (c *CartridgeEasyFlash) IOWrite(addr uint16, data uint8) bool {
	bank := addr & 0xff00
	if bank == 0xde00 {
		switch addr & 0x2 {
		case 0:
			c.register00 = data & BankMask
			return true
		case 1:
			return true
		case 2:
			c.controlUpdate(data, true)
			return true
		case 3:
			return true
		}
	} else if bank == 0xdf00 {
		c.io2Store(addr, data)
		return true
	}
	return false
}

// romLRead reads a byte from the low ROM bank at the given address within the EasyFlash cartridge.
func (c *CartridgeEasyFlash) romLRead(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateLow.Read(v)
}

// romHRead reads a byte of data from the high ROM bank based on the given address and the current register state.
func (c *CartridgeEasyFlash) romHRead(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateHigh.Read(v)
}

// romLWrite writes the provided value to the low bank memory at the calculated address based on the current register state.
func (c *CartridgeEasyFlash) romLWrite(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateLow.Store(v, value)
}

// romHWrite handles writing a value to the high ROM bank at the specified address within the cartridge memory.
func (c *CartridgeEasyFlash) romHWrite(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateHigh.Store(v, value)
}

// io1Peek reads and returns a value from the specified IO1 address within the EasyFlash cartridge's memory map.
// The returned value depends on the address bits; either register02 or register00 is accessed.
func (c *CartridgeEasyFlash) io1Peek(addr uint16) uint8 {
	if addr&2 != 0 {
		return c.register02
	}
	return c.register00
}

// io2Read reads a byte from the I/O bank 2 memory-mapped area at the specified address and returns the value.
func (c *CartridgeEasyFlash) io2Read(addr uint16) uint8 {
	return c.ram[addr&0xff]
}

// io2Store writes a given 8-bit value into the cartridge's RAM at the specified address masked to a 256-byte range.
func (c *CartridgeEasyFlash) io2Store(addr uint16, value uint8) {
	c.ram[addr&0xff] = value
}

// writeChipIfNotEmpty writes the chip data to the given writer
// if the chip data is not empty or optimization is disabled.
func (c *CartridgeEasyFlash) writeChipIfNotEmpty(fd io.Writer, chip references.IC64CartridgeChipHeader) error {
	for i := uint16(0); i < chip.Size(); i++ {
		if (chip.Data()[i] != 0xff) || !c.optimize {
			if err := chip.Write(fd); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

// binAttach initializes a raw cartridge with default data and copies data from the provided CRTLoader into it.
func (c *CartridgeEasyFlash) binAttach(ldr references.IC64CartridgeLoader) ([]byte, error) {
	rawCart := make([]uint8, 0x100000)
	for idx := range rawCart {
		rawCart[idx] = 0xff
	}
	copy(rawCart, ldr.GetData())
	return rawCart, nil
}

// crtAttach processes CRT cartridge data using the provided CRTLoader
// and returns the formatted cartridge data or an error.
func (c *CartridgeEasyFlash) crtAttach(ldr references.IC64CartridgeLoader) ([]byte, error) {
	raw := make([]uint8, 0x100000)
	for idx := range raw {
		raw[idx] = 0xff
	}
	for {
		chip, err := ldr.ReadChipHeader()
		if err != nil {
			return nil, err
		}
		if chip == nil {
			break
		}
		if uint16(len(chip.Data())) != chip.Size() {
			return nil, fmt.Errorf("invalid chip size")
		}
		if chip.Size() == 0x2000 {
			if chip.Bank() >= NBanks || !(chip.Start() == 0x8000 || chip.Start() == 0xa000 || chip.Start() == 0xe000) {
				return nil, fmt.Errorf("invalid start")
			}
			index := uint(chip.Bank()) << 14
			offset := uint(chip.Start()) & uint(chip.Size())
			target := index | offset
			copy(raw[target:target+uint(chip.Size())], chip.Data())
		} else if chip.Size() == 0x4000 {
			if chip.Bank() >= NBanks || chip.Start() != 0x8000 {
				return nil, fmt.Errorf("invalid start")
			}
			target := uint(chip.Bank()) << 14
			copy(raw[target:target+uint(chip.Size())], chip.Data())
		} else {
			return nil, fmt.Errorf("unkwnown chip size")
		}
	}
	return raw, nil
}

// Detach finalizes and detaches the EasyFlash cartridge, ensuring any pending writes are flushed before shutdown.
func (c *CartridgeEasyFlash) Detach() error {
	if c.writeEnabled {
		if err := c.flushImage(); err != nil {
			return err
		}
	}
	c.stateLow.Shutdown()
	c.stateHigh.Shutdown()
	c.filename = ""
	return nil
}

// flushImage saves the current cartridge data to the specified file based on its type or returns an error if unsuccessful.
func (c *CartridgeEasyFlash) flushImage() error {
	if len(c.filename) == 0 {
		return nil
	}
	if c.filetype == loader.TypeBin {
		return c.binSave(c.filename)
	} else if c.filetype == loader.TypeCrt {
		return c.crtSave(c.filename)
	}
	return fmt.Errorf("unknown cartridget type")
}

// binSave saves the cartridge data to a binary file with the specified filename.
// It returns an error if the filename is invalid or if file creation or data writing fails.
func (c *CartridgeEasyFlash) binSave(filename string) error {
	const size = 0x2000
	if len(filename) == 0 {
		return fmt.Errorf("invalid file name")
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	low := 0
	high := 0
	var lowData []uint8
	var highData []uint8
	for i := 0; i < NBanks; i++ {
		if lowData, err = c.stateLow.ReadInterval(uint(low), uint(low)+size); err != nil {
			return nil
		}
		_, err = fd.Write(lowData)
		if err != nil {
			return err
		}
		if highData, err = c.stateLow.ReadInterval(uint(high), uint(high)+size); err != nil {
			return nil
		}
		_, err = fd.Write(highData)
		if err != nil {
			return err
		}
		low += size
		high += size
	}
	return nil
}

// crtSave saves the current cartridge data to a specified file in the CRT format.
// Returns an error if the save operation fails or is unimplemented.
func (c *CartridgeEasyFlash) crtSave(_ string) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}

// snapshotWriteModule serializes the module state of the CartridgeEasyFlash into the provided snapshot.
func (c *CartridgeEasyFlash) snapshotWriteModule(s *snapshot.Snapshot) error {
	//m := s.NewModule(snapModuleName, SnapMajor, SnapMinor)
	//m.Add("jumper", c.jumper)
	//m.Add("register00", c.register00)
	//m.Add("register00", c.register00)
	//m.Add("ram", c.ram)
	//TODO
	//m.Add("romLBanks", c.romLBanks)
	//m.Add("romHBanks", c.romHBanks)
	//if err := c.stateLow.SnapshotWriteModule(s, flashSnapModuleName); err != nil {
	//	return err
	//}
	//if err := c.stateHigh.SnapshotWriteModule(s, flashSnapModuleName); err != nil {
	//	return err
	//}
	return nil
}
