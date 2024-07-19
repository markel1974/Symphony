package cartridge16k

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/iboard"
)

type Cartridge16K struct {
	b0Interval icartridge.Interval
	b1Interval icartridge.Interval
	bank0      []uint8
	bank1      []uint8
	intervals  icartridge.Interval
	game       uint8
	exRom      uint8
	board      iboard.IBoard
}

func New(game uint8, exRom uint8, b0 icartridge.Interval, b1 icartridge.Interval) *Cartridge16K {
	return &Cartridge16K{
		game:       game,
		exRom:      exRom,
		b0Interval: b0,
		b1Interval: b1,
		intervals:  b0 | b1,
	}
}

func (c *Cartridge16K) Setup(board iboard.IBoard, ldr *loader.CRTLoader) error {
	c.board = board
	if ldr.GetMode() == loader.ModeCrt {
		return c.initBin(ldr)
	}
	return c.initRaw(ldr.GetData())
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
	if err := loader.Validate(data); err != nil {
		return err
	}
	c.bank0 = data[:0x2000]
	c.bank1 = data[0x2000:]
	return nil
}

func (c *Cartridge16K) Write(i icartridge.Interval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("Cartridge16K can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

func (c *Cartridge16K) Read(i icartridge.Interval, addr uint16) (uint8, bool) {
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
