package magicdesk

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
)

type CartridgeMagicDesk struct {
	id       string
	spec     *icartridge.CartridgeSpec
	banks    [][]byte
	bankMask uint8
	regVal   uint8
	slot     uint8
	board    icartridge.IExpansion
}

func GetType() int {
	return loader.CARTRIDGE_MAGIC_DESK
}

func New() icartridge.ICartridge {
	return &CartridgeMagicDesk{
		spec:     icartridge.GetCartridgeSpec(icartridge.CartridgeMode8K),
		bankMask: 0x7f,
		regVal:   0,
		slot:     0,
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
	if (i & (c.spec.IntervalLow | c.spec.IntervalHigh)) != 0 {
		fmt.Printf("CartridgeOcean can't be write [bank %d] %x => %d\n", c.slot, addr, data)
		return true
	}
	return false
}

func (c *CartridgeMagicDesk) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if (i & (c.spec.IntervalLow | c.spec.IntervalHigh)) != 0 {
		//if c.b0Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		//if c.b1Interval == i {
		//	return c.banks[c.currBank][addr&0x1fff], true
		//}
		return c.banks[c.slot][addr&0x1fff], true
	}
	return 0, false
}

func (c *CartridgeMagicDesk) IORead(addr uint16) (uint8, bool) {
	if (addr & 0xfff0) == 0xde00 {
		return c.regVal, true
	}
	return 0, false
}

func (c *CartridgeMagicDesk) IOWrite(addr uint16, data uint8) bool {
	if (addr & 0xfff0) == 0xde00 {
		c.regVal = data & (0x80 | c.bankMask)
		c.slot = data & c.bankMask
		fmt.Println("magic desk slot", c.slot)
		var spec *icartridge.CartridgeSpec
		if (data & 0x80) != 0 {
			spec = icartridge.GetCartridgeSpec(icartridge.CartridgeModeOff)
		} else {
			spec = icartridge.GetCartridgeSpec(icartridge.CartridgeMode8K)
		}
		if spec != c.spec {
			fmt.Println("magic desk changing config", c.spec)
			c.spec = spec
			c.board.GameExRomConfigChanged()
		}
	}
	return false
}

func (c *CartridgeMagicDesk) GetExRom() uint8 {
	return c.spec.ExRom
}

func (c *CartridgeMagicDesk) GetGame() uint8 {
	return c.spec.Game
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
	c.banks = [][]byte{}
	c.bankMask = 0x7f
	c.regVal = 0
	c.slot = 0
	switch len(data) {
	case 0x100000:
		c.bankMask = 0x3f
	case 0x80000:
		c.bankMask = 0x1f
	case 0x40000:
		c.bankMask = 0x0f
	case 0x20000:
		c.bankMask = 0x07
	case 0x10000:
		c.bankMask = 0x03
	default:
		return fmt.Errorf("unsupported size")
	}
	start := 0
	for start < len(data) {
		end := start + 0x2000
		c.banks = append(c.banks, data[start:end])
		start += end
	}
	return nil
}

func (c *CartridgeMagicDesk) initCrt(loader *loader.CRTLoader) error {
	c.banks = [][]byte{}
	lastBank := uint16(0)
	c.bankMask = 0x7f
	c.regVal = 0
	c.slot = 0
	for {
		chip, err := loader.ReadChipHeader()
		if chip == nil {
			break
		}
		if err != nil {
			return err
		}
		if (chip.Bank > 128) || ((chip.Start != 0x8000) && (chip.Start != 0xa000)) || (chip.Size != 0x2000) {
			return fmt.Errorf("invalid chip bank")
		}
		c.banks = append(c.banks, chip.Data)
		if chip.Bank > lastBank {
			lastBank = chip.Bank
		}
	}
	if lastBank >= 128 {
		return fmt.Errorf("chip has more than 128 banks")
	}
	if lastBank >= 64 {
		c.bankMask = 0x7f // min 65, max 128 banks
	} else if lastBank >= 32 {
		c.bankMask = 0x3f // min 33, max 64 banks
	} else if lastBank >= 16 {
		c.bankMask = 0x1f // min 17, max 32 banks
	} else if lastBank >= 8 {
		c.bankMask = 0x0f // min 9, max 16 banks
	} else if lastBank >= 4 {
		c.bankMask = 0x07 // min 5, max 8 banks
	} else {
		c.bankMask = 0x03 // max 4 banks
	}
	return nil
}
