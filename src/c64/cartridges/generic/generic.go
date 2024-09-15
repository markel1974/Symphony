package generic

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
)

///Users/tinmr305/Desktop/emu/vice-emu-code-r45201-trunk-vice/src/c64/cart/c64-generic.c

const cSize16K = 0x4000
const cSize8K = 0x2000

type Generic struct {
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
	return &Generic{
		game:       0,
		exRom:      0,
		b0Interval: 0,
		b1Interval: 0,
		intervals:  0,
	}
}

func (c *Generic) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	c.board = board
	c.id = ldr.GetId()
	if ldr.GetType() == loader.TypeCrt {
		return c.initCrt(ldr)
	}
	return c.initRaw(ldr.GetData())
}

func (c *Generic) Reset() {

}

func (c *Generic) GetId() string {
	return c.id
}

func (c *Generic) initCrt(ldr *loader.CRTLoader) error {
	c.bank0 = make([]uint8, cSize8K)
	c.bank1 = make([]uint8, cSize8K)
	chip1, err := ldr.ReadChipHeader()
	if err != nil {
		return err
	}
	if chip1 == nil {
		return fmt.Errorf("nil chip")
	}
	if chip1.Start == 0x8000 {
		if chip1.Size == cSize8K {
			if chip2, _ := ldr.ReadChipHeader(); chip2 == nil {
				copy(c.bank0, chip1.Data)
				c.applyConfig(icartridge.CartridgeMode8K)
				return nil
			} else if chip2.Size == cSize8K {
				if chip2.Start == 0x8000 {
					copy(c.bank0, chip1.Data)
					copy(c.bank1, chip2.Data)
					c.applyConfig(icartridge.CartridgeMode16K)
					return nil
				} else if chip2.Start == 0xe000 {
					copy(c.bank0, chip1.Data)
					copy(c.bank1, chip2.Data)
					c.applyConfig(icartridge.CartridgeModeUltimax)
					return nil
				}
			}
		} else if chip1.Size == cSize16K {
			copy(c.bank0, chip1.Data[:cSize8K])
			copy(c.bank1, chip1.Data[cSize8K:])
			c.applyConfig(icartridge.CartridgeMode16K)
			return nil
		}
	}
	return fmt.Errorf("unsupported crt")
}

func (c *Generic) initRaw(data []byte) error {
	c.bank0 = make([]uint8, cSize8K)
	c.bank1 = make([]uint8, cSize8K)
	if err := loader.ValidateCartridge(data); err != nil {
		return err
	}
	if len(data) == cSize8K {
		copy(c.bank0, data)
		c.applyConfig(icartridge.CartridgeMode8K)
		return nil
	}
	if len(data) == cSize16K {
		copy(c.bank0, data[:cSize8K])
		copy(c.bank1, data[cSize8K:])
		c.applyConfig(icartridge.CartridgeMode16K)
		return nil
	}
	return fmt.Errorf("invalid size")
}

func (c *Generic) applyConfig(ct icartridge.CartridgeMode) {
	v := icartridge.GetCartridgeSpec(ct)
	c.game = v.Game
	c.exRom = v.ExRom
	c.b0Interval = v.IntervalLow
	c.b1Interval = v.IntervalHigh
	c.intervals = v.IntervalLow | v.IntervalHigh
}

func (c *Generic) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if (i & c.intervals) != 0 {
		fmt.Printf("Generic Cartridge can't be write %x => %d\n", addr, data)
		return true
	}
	return false
}

func (c *Generic) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if (i & c.intervals) != 0 {
		if c.b0Interval == i {
			return c.bank0[(addr & 0x1fff)], true
		}
		if c.b1Interval == i {
			return c.bank1[(addr & 0x1fff)], true
		}
	}
	return 0, false
}

func (c *Generic) IORead(_ uint16) (uint8, bool) {
	return 0, false
}

func (c *Generic) IOWrite(_ uint16, _ uint8) bool {
	return false
}

func (c *Generic) GetExRom() uint8 {
	return c.exRom
}

func (c *Generic) GetGame() uint8 {
	return c.game
}

func (c *Generic) EmulationRequired() bool {
	return false
}

func (c *Generic) Emulate() {
}

func (c *Generic) Detach() error {
	//TODO
	return nil
}
