package cartridge16k

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

type Cartridge16K struct {
	id         string
	b0Interval icartridge.RomInterval
	b1Interval icartridge.RomInterval
	bank0      []uint8
	bank1      []uint8
	intervals  icartridge.RomInterval
	game       uint8
	exRom      uint8
	board      icartridge.IExpansion
}

func New(ultimax bool) *Cartridge16K {
	v := icartridge.GetCartridgeSpec(icartridge.CartridgeMode16K)
	if ultimax {
		v = icartridge.GetCartridgeSpec(icartridge.CartridgeModeUltimax)
	}
	return &Cartridge16K{
		game:       v.Game,
		exRom:      v.ExRom,
		b0Interval: v.IntervalLow,
		b1Interval: v.IntervalHigh,
		intervals:  v.IntervalLow | v.IntervalHigh,
	}
}

func (c *Cartridge16K) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetMode() == loader.FiletypeCrt {
		return c.initBin(ldr)
	}
	return c.initRaw(ldr.GetData())
}

func (c *Cartridge16K) GetId() string {
	return c.id
}

func (c *Cartridge16K) initBin(ldr *loader.CRTLoader) error {
	//TODO IMPLEMENT
	return fmt.Errorf("uninplemented")
}

func (c *Cartridge16K) initRaw(data []byte) error {
	const cSize = 0x4000
	if len(data) != cSize {
		return fmt.Errorf("invalid size")
	}
	if err := loader.ValidateCartridge(data); err != nil {
		return err
	}
	c.bank0 = data[:0x2000]
	c.bank1 = data[0x2000:]
	return nil
}

func (c *Cartridge16K) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("Cartridge16K can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

func (c *Cartridge16K) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if i&c.intervals != 0 {
		if c.b0Interval == i {
			return c.bank0[addr&0x1fff], true
		}
		if c.b1Interval == i {
			return c.bank1[addr&0x1fff], true
		}
	}
	return 0, false
}

func (c *Cartridge16K) IORead(addr uint16) (uint8, bool) {
	return 0, false
}

func (c *Cartridge16K) IOWrite(addr uint16, data uint8) bool {
	return false
}

func (c *Cartridge16K) GetExRom() uint8 {
	return c.exRom
}

func (c *Cartridge16K) GetGame() uint8 {
	return c.game
}

func (c *Cartridge16K) Detach() error {
	//TODO
	return nil
}
