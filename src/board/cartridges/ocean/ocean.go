package ocean

import (
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

type CartridgeOcean struct {
	id        string
	intervals icartridge.RomInterval
	lastData  uint8
	banks     [][]byte
	ioMask    uint8
	currBank  uint8
	game      uint8
	exRom     uint8
	board     icartridge.IExpansion
}

func New() *CartridgeOcean {
	v := icartridge.GetCartridgeSpec(icartridge.CartridgeMode16K)
	return &CartridgeOcean{
		game:      v.Game,
		exRom:     v.ExRom,
		intervals: v.IntervalLow | v.IntervalHigh,
		lastData:  0,
	}
}

func (c *CartridgeOcean) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetMode() == loader.FiletypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

func (c *CartridgeOcean) GetId() string {
	return c.id
}

func (c *CartridgeOcean) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.currBank, addr, data)
		return true
	}
	return false
}

func (c *CartridgeOcean) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
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
	if err := loader.ValidateCartridge(data); err != nil {
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
