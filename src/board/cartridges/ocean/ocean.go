package ocean

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/iboard"
)

type CartridgeOcean struct {
	//b0Interval Interval
	//b1Interval Interval
	intervals icartridge.Interval
	lastData  uint8
	banks     [][]byte
	ioMask    uint8
	currBank  uint8
	game      uint8
	exRom     uint8
	board     iboard.IBoard
}

func New(game uint8, exRom uint8, b0 icartridge.Interval, b1 icartridge.Interval) *CartridgeOcean {
	return &CartridgeOcean{
		game:  game,
		exRom: exRom,
		//b0Interval: b0,
		//b1Interval: b1,
		intervals: b0 | b1,
		lastData:  0,
	}
}

func (c *CartridgeOcean) Setup(board iboard.IBoard, ldr *loader.CRTLoader) error {
	c.board = board
	if ldr.GetMode() == loader.ModeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

func (c *CartridgeOcean) Write(i icartridge.Interval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.currBank, addr, data)
		return true
	}
	return false
}

func (c *CartridgeOcean) Read(i icartridge.Interval, addr uint16) (uint8, bool) {
	if i&c.intervals != 0 {
		//if c.b0Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		//if c.b1Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		return c.banks[c.currBank][addr&0x1fff], true
	}
	return 0, false
}

func (c *CartridgeOcean) IORead(addr uint16) (uint8, bool) {
	if (addr & 0xfff0) == 0xde00 {
		return c.lastData, true
	}
	return 0, false
}

func (c *CartridgeOcean) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xde00 {
		//exRomDisabled := (data & 0x80) != 0
		//currBank := data & 0x7f
		currBank := data & c.ioMask & 0x3f
		c.currBank = currBank
		c.lastData = data
		//TODO board.updateMemoryConfig()
		fmt.Printf("BANK SWITCHING %x => %d, %d\n", addr, data, c.currBank)
	}
	return false
}

func (c *CartridgeOcean) GetExRom() uint8 {
	return c.exRom
}

func (c *CartridgeOcean) GetGame() uint8 {
	return c.game
}

func (c *CartridgeOcean) Detach() error {
	//TODO
	return nil
}

func (c *CartridgeOcean) initBin(data []byte) error {
	if err := loader.Validate(data); err != nil {
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

func (c *CartridgeOcean) initCrt(loader *loader.CRTLoader) error {
	c.banks = [][]byte{}
	romSize := 0
	for {
		chip, err := loader.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if (chip.Bank > 63) || ((chip.Start != 0x8000) && (chip.Start != 0xa000)) || (chip.Size != 0x2000) {
			return fmt.Errorf("invalid chip bank")
		}
		c.banks = append(c.banks, chip.Data)
		romSize += int(chip.Size)
	}
	c.ioMask = uint8((romSize >> 13) - 1)
	c.lastData = 0
	c.currBank = 0
	return nil
}
