package cartridge16k

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
)

///Users/tinmr305/Desktop/emu/vice-emu-code-r45201-trunk-vice/src/c64/cart/c64-generic.c

const cSize = 0x4000

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

func GetType() int {
	return loader.CARTRIDGE_CRT
}

func New() icartridge.ICartridge {
	return &Cartridge16K{
		game:       1,
		exRom:      1,
		b0Interval: 0,
		b1Interval: 0,
		intervals:  0,
	}
}

func (c *Cartridge16K) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	v := icartridge.GetCartridgeSpec(icartridge.CartridgeMode16K)
	c.game = v.Game
	c.exRom = v.ExRom
	c.b0Interval = v.IntervalLow
	c.b1Interval = v.IntervalHigh
	c.intervals = v.IntervalLow | v.IntervalHigh
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetType() == loader.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initRaw(ldr.GetData())
}

func (c *Cartridge16K) Reset() {

}

func (c *Cartridge16K) GetId() string {
	return c.id
}

func (c *Cartridge16K) initCrt(ldr *loader.CRTLoader) error {
	const bankSize = cSize / 2
	chip, err := ldr.ReadChipHeader()
	if chip == nil {
		return fmt.Errorf("nil chip")
	}
	if err != nil {
		return err
	}
	if (chip.Start == 0x8000) && (chip.Size == cSize) {
		c.bank0 = make([]uint8, bankSize)
		c.bank1 = make([]uint8, bankSize)
		copy(c.bank0, chip.Data[:bankSize])
		copy(c.bank1, chip.Data[bankSize:])
		return nil
	}

	if (chip.Start == 0x8000) && (chip.Size == bankSize) {
		chip2, err2 := ldr.ReadChipHeader()
		if chip2 == nil {
			return fmt.Errorf("nil chip")
		}
		if err2 != nil {
			return err
		}
		c.bank0 = make([]uint8, bankSize)
		c.bank1 = make([]uint8, bankSize)
		copy(c.bank0, chip.Data)
		copy(c.bank1, chip2.Data)
		return nil
	}

	if (chip.Start == 0xe000) && (chip.Size > 0) && (uint(chip.Size)+uint(chip.Start)) == 0x10000 {
		//TODO ULTIMAX
		v := icartridge.GetCartridgeSpec(icartridge.CartridgeModeUltimax)
		c.game = v.Game
		c.exRom = v.ExRom
	}

	return fmt.Errorf("unsupported crt")

	/*
		if (chip.start >= 0xe000 && chip.size > 0 && (chip.size + chip.start) == 0x10000) {
		        if (crt_read_chip(rawcart, chip.start & 0x3fff, &chip, fd)) {
		            return -1;
		        }
		        if (generic_common_attach(CARTRIDGE_ULTIMAX) < 0) {
		            return -1;
		        }
		        return CARTRIDGE_ULTIMAX;
		    }
	*/
}

func (c *Cartridge16K) initRaw(data []byte) error {
	if len(data) != cSize {
		return fmt.Errorf("invalid size")
	}
	if err := loader.ValidateCartridge(data); err != nil {
		return err
	}
	const bankSize = cSize / 2
	c.bank0 = make([]uint8, bankSize)
	c.bank1 = make([]uint8, bankSize)
	copy(c.bank0, data[:bankSize])
	copy(c.bank1, data[bankSize:])
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

func (c *Cartridge16K) EmulationRequired() bool {
	return false
}

func (c *Cartridge16K) Emulate() {

}

func (c *Cartridge16K) Detach() error {
	//TODO
	return nil
}
