package cartridges

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/iboard"
)

type Cartridge8K struct {
	bank0     []uint8
	intervals Interval
	game      uint8
	exRom     uint8
	board     iboard.IBoard
}

func NewCartridge8K(game uint8, exRom uint8, i Interval) *Cartridge8K {
	return &Cartridge8K{
		game:      game,
		exRom:     exRom,
		intervals: i,
	}
}

func (c *Cartridge8K) Setup(board iboard.IBoard, ldr *loader.CRTLoader) error {
	c.board = board
	if ldr.GetMode() == loader.ModeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

func (c *Cartridge8K) initCrt(ldr *loader.CRTLoader) error {
	//TODO IMPLEMENT
	return fmt.Errorf("uninplemented")
}

func (c *Cartridge8K) initBin(data []byte) error {
	const cSize = 0x2000
	if len(data) != cSize {
		return fmt.Errorf("invalid size")
	}
	if err := validate(data); err != nil {
		return err
	}
	c.bank0 = data
	return nil
}

func (c *Cartridge8K) Write(i Interval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("Cartridge8K can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

func (c *Cartridge8K) Read(i Interval, addr uint16) (uint8, bool) {
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

func (c *Cartridge8K) Detach() error {
	//TODO
	return nil
}
