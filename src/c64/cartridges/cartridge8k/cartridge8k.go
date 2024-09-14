package cartridge8k

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
)

const cSize = 0x2000

type Cartridge8K struct {
	id        string
	bank0     []uint8
	intervals icartridge.RomInterval
	game      uint8
	exRom     uint8
	board     icartridge.IExpansion
}

func GetType() int {
	return loader.CARTRIDGE_GENERIC_8KB
}

func New() icartridge.ICartridge {
	v := icartridge.GetCartridgeSpec(icartridge.CartridgeMode8K)
	return &Cartridge8K{
		game:      v.Game,
		exRom:     v.ExRom,
		intervals: v.IntervalLow | v.IntervalHigh,
	}
}

func (c *Cartridge8K) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetType() == loader.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

func (c *Cartridge8K) Reset() {

}

func (c *Cartridge8K) GetId() string {
	return c.id
}

func (c *Cartridge8K) initCrt(ldr *loader.CRTLoader) error {
	chip, err := ldr.ReadChipHeader()
	if chip == nil {
		return fmt.Errorf("nil chip")
	}
	if err != nil {
		return err
	}
	if chip.Start != 0x8000 {
		return fmt.Errorf("invalid chip start")
	}
	if chip.Size != cSize {
		return fmt.Errorf("invalid chip size")
	}
	c.bank0 = chip.Data
	return nil
}

func (c *Cartridge8K) initBin(data []byte) error {
	if len(data) != cSize {
		return fmt.Errorf("invalid size")
	}
	if err := loader.ValidateCartridge(data); err != nil {
		return err
	}
	c.bank0 = data
	return nil
}

func (c *Cartridge8K) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("Cartridge8K can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

func (c *Cartridge8K) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if i&c.intervals != 0 {
		v := c.bank0[addr&0x1fff]
		return v, true
	}
	return 0, false
}

func (c *Cartridge8K) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

func (c *Cartridge8K) IOWrite(addr uint16, data uint8) bool {
	return false
}

func (c *Cartridge8K) GetExRom() uint8 {
	return c.exRom
}

func (c *Cartridge8K) GetGame() uint8 {
	return c.game
}

func (c *Cartridge8K) EmulationRequired() bool {
	return false
}

func (c *Cartridge8K) Emulate() {

}

func (c *Cartridge8K) Detach() error {
	//TODO
	return nil
}
