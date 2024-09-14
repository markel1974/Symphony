package magicdesk

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
)

type CartridgeMagicDesk struct {
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

func GetType() int {
	return loader.CARTRIDGE_MAGIC_DESK
}

func New() icartridge.ICartridge {
	v := icartridge.GetCartridgeSpec(icartridge.CartridgeMode16K)
	return &CartridgeMagicDesk{
		game:      v.Game,
		exRom:     v.ExRom,
		intervals: v.IntervalLow | v.IntervalHigh,
		lastData:  0,
	}
}

func (c *CartridgeMagicDesk) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetType() == loader.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initBin(ldr.GetData())
}

func (c *CartridgeMagicDesk) Reset() {

}

func (c *CartridgeMagicDesk) GetId() string {
	return c.id
}

func (c *CartridgeMagicDesk) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.currBank, addr, data)
		return true
	}
	return false
}

func (c *CartridgeMagicDesk) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
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

func (c *CartridgeMagicDesk) IORead(addr uint16) (uint8, bool) {
	if (addr & 0xfff0) == 0xde00 {
		return c.lastData, true
	}
	return 0, false
}

func (c *CartridgeMagicDesk) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xde00 {
		//exRomDisabled := (data & 0x80) != 0
		//currBank := data & 0x7f
		currBank := (data & c.ioMask) & 0x3f
		c.currBank = currBank
		c.lastData = data
		fmt.Printf("[MAGIC DESK] Bank switching %x => %d, %d\n", addr, data, c.currBank)
	}
	return false
}

func (c *CartridgeMagicDesk) GetExRom() uint8 {
	return c.exRom
}

func (c *CartridgeMagicDesk) GetGame() uint8 {
	return c.game
}

func (c *CartridgeMagicDesk) Detach() error {
	//TODO
	return nil
}

func (c *CartridgeMagicDesk) EmulationRequired() bool {
	return false
}

func (c *CartridgeMagicDesk) Emulate() {

}

func (c *CartridgeMagicDesk) initBin(data []byte) error {
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

func (c *CartridgeMagicDesk) initCrt(loader *loader.CRTLoader) error {
	c.banks = [][]byte{}
	//c.exRom = uint8(loader.ExRom)
	//c.game = uint8(loader.Game)

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
