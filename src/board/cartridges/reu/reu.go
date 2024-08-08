package reu

import (
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

const (
	size128K = 0x20000
	size256K = 0x40000
	size512K = 0x80000
	size1M   = 0x100000
	size2M   = 0x200000
	size4M   = 0x400000
	size8M   = 0x800000
	size16M  = 0x1000000
)

const (
	Id128K = "REU128K"
	Id256K = "REU256K"
	Id512K = "REU512K"
	Id1M   = "REU1M"
	Id2M   = "REU2M"
	Id4M   = "REU4M"
	Id8M   = "REU8M"
	Id16M  = "REU16M"
)

type REU struct {
	id        string
	ram       []uint8 // REU RAM
	mask      uint32  // REU RAM address bit mask
	regs      []uint8 // REU registers
	size      int
	expansion icartridge.IExpansion
}

func new(size int) icartridge.ICartridge {
	r := &REU{
		id:        "",
		ram:       nil,
		mask:      0,
		regs:      make([]uint8, 16),
		size:      size,
		expansion: nil,
	}
	r.regs[0] = 0x40
	for i := 1; i < 11; i++ {
		r.regs[i] = 0
	}
	for i := 11; i < 16; i++ {
		r.regs[i] = 0xff
	}
	return r
}

func New128K() icartridge.ICartridge {
	return new(size128K)
}
func New256K() icartridge.ICartridge {
	return new(size256K)
}
func New512K() icartridge.ICartridge {
	return new(size512K)
}
func New1M() icartridge.ICartridge {
	return new(size1M)
}
func New2M() icartridge.ICartridge {
	return new(size2M)
}
func New4M() icartridge.ICartridge {
	return new(size4M)
}
func New8M() icartridge.ICartridge {
	return new(size8M)
}
func New16M() icartridge.ICartridge {
	return new(size16M)
}

func (reu *REU) GetId() string {
	return reu.id
}

func (reu *REU) EmulationRequired() bool {
	return false
}

func (reu *REU) Emulate() {
}

func (reu *REU) Reset() {
	for i := 1; i < 11; i++ {
		reu.regs[i] = 0
	}
	for i := 11; i < 16; i++ {
		reu.regs[i] = 0xff
	}
	if len(reu.ram)-1 > 0x20000 {
		reu.regs[0] = 0x50
	} else {
		reu.regs[0] = 0x40
	}
}

func (reu *REU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	//TODO from Setup
	reu.expansion = board
	reu.id = ldr.GetId()
	reu.ram = nil
	reu.mask = 0
	size := reu.size
	reu.mask = uint32(size) - 1
	reu.ram = make([]uint8, size)
	// Set kind bit in status register
	if size-1 > 0x20000 {
		reu.regs[0] |= 0x10
	} else {
		reu.regs[0] &= 0xef
	}
	reu.expansion.RamSetWriteTrigger(0xff00, reu.triggerDMA)
	return nil
}

func (reu *REU) GetExRom() uint8 {
	return 1
}

func (reu *REU) GetGame() uint8 {
	return 1
}

func (reu *REU) IORead(addr uint16) (uint8, bool) {
	//$DF00-$DF0A, $DF20-$DF2A, $DF40-$DF4A and so on, up to $DFE0-$DFEA
	if reu.ram == nil {
		return 0, false
	}
	if (addr & 0xfff0) == 0xdf00 {
		reg := addr & 0x0f
		switch reg {
		case 0:
			ret := reu.regs[0]
			reu.regs[0] &= 0x1f
			return ret, true
		case 6:
			return reu.regs[6] | 0xf8, true
		case 9:
			return reu.regs[9] | 0x1f, true
		case 10:
			return reu.regs[10] | 0x3f, true
		default:
			return reu.regs[reg], true
		}
	}
	return 0, false
}

func (reu *REU) IOWrite(addr uint16, data uint8) bool {
	//$DF00-$DF0A, $DF20-$DF2A, $DF40-$DF4A and so on, up to $DFE0-$DFEA
	if reu.ram == nil {
		return false
	}
	if (addr & 0xfff0) == 0xdf00 {
		reg := addr & 0x0f
		switch reg {
		case 0, 11, 12, 13, 14, 15:
			// Status register is read-only
			// Unconnected registers
		case 1:
			// Command register
			reu.regs[1] = data
			if (data & 0x90) == 0x90 {
				reu.executeDma()
			}
		default:
			reu.regs[reg] = data
		}
		return true
	}
	return false
}

func (reu *REU) triggerDMA(_ uint16, _ uint8) {
	// CPU triggered REU by writing to $ff00
	if reu.ram == nil {
		return
	}
	if (reu.regs[1] & 0x90) == 0x80 {
		reu.executeDma()
	}
}

func (reu *REU) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	return false
}

func (reu *REU) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

func (reu *REU) Detach() error {
	//TODO IMPLEMENT
	return nil
}

func (reu *REU) executeDma() {
	//TODO IN EMULATE CYCLE ONE BYTE AT TIME

	reu.expansion.DMALow(true)
	// Get C64 and REU transfer base addresses
	c64Addr := uint16(reu.regs[2]) | (uint16(reu.regs[3]) << 8)
	reuAddr := uint32(reu.regs[4]) | (uint32(reu.regs[5]) << 8) | (uint32(reu.regs[6]) << 16)

	// Calculate transfer length
	length := uint32(reu.regs[7]) | (uint32(reu.regs[8]) << 8)
	if length == 0 {
		length = 0x10000
	}

	// Calculate address increments
	c64Inc := uint32(1)
	reuInc := uint32(1)
	if (uint32(reu.regs[10]) & 0x80) != 0 {
		c64Inc = 0
	}
	if (uint32(reu.regs[10]) & 0x40) != 0 {
		reuInc = 0
	}

	// Do transfer
	mode := reu.regs[1] & 3
	switch mode {
	case 0:
		// C64 -> REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			reu.ram[reuAddr&reu.mask] = reu.expansion.Read(c64Addr)
		}

	case 1:
		// C64 <- REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			reu.expansion.Write(c64Addr, reu.ram[reuAddr&reu.mask])
		}

	case 2:
		// C64 <-> REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			tmp := reu.expansion.Read(c64Addr)
			reu.expansion.Write(c64Addr, reu.ram[reuAddr&reu.mask])
			reu.ram[reuAddr&reu.mask] = tmp
		}

	case 3:
		// Compare
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			if reu.ram[reuAddr&reu.mask] != reu.expansion.Read(c64Addr) {
				reu.regs[0] |= 0x20
				break
			}
		}
	}

	// Update address and length registers if autoload is off
	if (reu.regs[1] & 0x20) == 0 {
		reu.regs[2] = uint8(c64Addr)
		reu.regs[3] = uint8(c64Addr >> 8)
		reu.regs[4] = uint8(reuAddr)
		reu.regs[5] = uint8(reuAddr >> 8)
		reu.regs[6] = uint8(reuAddr >> 16)
		reu.regs[7] = uint8(length + 1)
		reu.regs[8] = uint8((length + 1) >> 8)
	}

	// Set complete bit in status register
	reu.regs[0] |= 0x40

	// Clear exec bit in command register
	reu.regs[1] &= 0x7f
	reu.expansion.DMALow(false)
}
