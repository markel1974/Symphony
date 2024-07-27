package reu

import (
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

const (
	REU_NONE = iota
	REU_128K
	REU_256K
	REU_512K
)

type REU struct {
	id        string
	_ex_ram   []uint8 // REU expansion RAM
	_ram_size uint32  // Size of expansion RAM
	_ram_mask uint32  // Expansion RAM address bit mask
	_regs     []uint8 // REU registers
	_old_size int
	_size     int
	expansion icartridge.IExpansion
}

func NewREU() icartridge.ICartridge {
	r := &REU{
		id:        "",
		_ex_ram:   nil,
		_ram_size: 0,
		_ram_mask: 0,
		_regs:     make([]uint8, 16),
		_old_size: 0,
		_size:     0,
		expansion: nil,
	}
	r._regs[0] = 0x40
	for i := 1; i < 11; i++ {
		r._regs[i] = 0
	}
	for i := 11; i < 16; i++ {
		r._regs[i] = 0xff
	}
	return r
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
		reu._regs[i] = 0
	}
	for i := 11; i < 16; i++ {
		reu._regs[i] = 0xff
	}
	if reu._ram_size > 0x20000 {
		reu._regs[0] = 0x50
	} else {
		reu._regs[0] = 0x40
	}
}

func (reu *REU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	//TODO from Setup
	const reuSize = REU_512K
	reu.expansion = board
	reu.id = ldr.GetId()
	if reu._old_size == reuSize {
		return nil
	}
	if reu._old_size != REU_NONE {
		reu._ex_ram = nil
	}
	reu._size = reuSize
	if reu._size != REU_NONE {
		switch reu._size {
		case REU_128K:
			reu._ram_size = 0x20000
		case REU_256K:
			reu._ram_size = 0x40000
		case REU_512K:
			reu._ram_size = 0x80000
		}
		reu._ram_mask = reu._ram_size - 1
		reu._ex_ram = make([]uint8, reu._ram_size)

		// Set size bit in status register
		if reu._ram_size > 0x20000 {
			reu._regs[0] |= 0x10
		} else {
			reu._regs[0] &= 0xef
		}
	}
	reu._old_size = reu._size
	reu.expansion.RamSetWriteTrigger(0xff00, reu.TriggerDMA)
	return nil
}

func (reu *REU) GetExRom() uint8 {
	return 1
}

func (reu *REU) GetGame() uint8 {
	return 1
}

func (reu *REU) IORead(addr uint16) (uint8, bool) {
	if reu._ex_ram == nil {
		return 0, false
	}
	if (addr & 0xfff0) == 0xdf00 {
		reg := addr & 0x0f
		switch reg {
		case 0:
			ret := reu._regs[0]
			reu._regs[0] &= 0x1f
			return ret, true
		case 6:
			return reu._regs[6] | 0xf8, true
		case 9:
			return reu._regs[9] | 0x1f, true
		case 10:
			return reu._regs[10] | 0x3f, true
		default:
			return reu._regs[reg], true
		}
	}
	return 0, false
}

func (reu *REU) IOWrite(addr uint16, data uint8) bool {
	if reu._ex_ram == nil {
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
			reu._regs[1] = data
			if (data & 0x90) == 0x90 {
				reu.executeDma()
			}
		default:
			reu._regs[reg] = data
		}
	}
	return false
}

func (reu *REU) TriggerDMA(_ uint16, _ uint8) {
	// CPU triggered REU by writing to $ff00
	if reu._ex_ram == nil {
		return
	}
	if (reu._regs[1] & 0x90) == 0x80 {
		reu.executeDma()
	}
}

func (reu *REU) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if i == icartridge.ROM_HI_2 && addr == 0xff00 {
		reu.ff00Trigger()
	}
	return false
}

func (reu *REU) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	return 0, false
}

func (reu *REU) Detach() error {
	//TODO IMPLEMENT
	return nil
}

func (reu *REU) ff00Trigger() {
	//TODO IMPLEMENT
}

func (reu *REU) executeDma() {
	reu.expansion.DMALow(true)
	// Get C64 and REU transfer base addresses
	c64Addr := uint16(reu._regs[2]) | (uint16(reu._regs[3]) << 8)
	reuAddr := uint32(reu._regs[4]) | (uint32(reu._regs[5]) << 8) | (uint32(reu._regs[6]) << 16)

	// Calculate transfer length
	length := uint32(reu._regs[7]) | (uint32(reu._regs[8]) << 8)
	if length == 0 {
		length = 0x10000
	}

	// Calculate address increments
	c64Inc := uint32(1)
	reuInc := uint32(1)
	if (uint32(reu._regs[10]) & 0x80) != 0 {
		c64Inc = 0
	}
	if (uint32(reu._regs[10]) & 0x40) != 0 {
		reuInc = 0
	}

	// Do transfer
	switch reu._regs[1] & 3 {
	case 0:
		// C64 -> REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			reu._ex_ram[reuAddr&reu._ram_mask] = reu.expansion.RamRead(c64Addr)
		}

	case 1:
		// C64 <- REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			reu.expansion.RamWrite(c64Addr, reu._ex_ram[reuAddr&reu._ram_mask])
		}

	case 2:
		// C64 <-> REU
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			tmp := reu.expansion.RamRead(c64Addr)
			reu.expansion.RamWrite(c64Addr, reu._ex_ram[reuAddr&reu._ram_mask])
			reu._ex_ram[reuAddr&reu._ram_mask] = tmp
		}

	case 3:
		// Compare
		for ; length > 0; c64Addr, reuAddr, length = (uint16)(c64Inc+c64Inc), reuInc+reuInc, length-1 {
			if reu._ex_ram[reuAddr&reu._ram_mask] != reu.expansion.RamRead(c64Addr) {
				reu._regs[0] |= 0x20
				break
			}
		}
	}

	// Update address and length registers if autoload is off
	if (reu._regs[1] & 0x20) == 0 {
		reu._regs[2] = uint8(c64Addr)
		reu._regs[3] = uint8(c64Addr >> 8)
		reu._regs[4] = uint8(reuAddr)
		reu._regs[5] = uint8(reuAddr >> 8)
		reu._regs[6] = uint8(reuAddr >> 16)
		reu._regs[7] = uint8(length + 1)
		reu._regs[8] = uint8((length + 1) >> 8)
	}

	// Set complete bit in status register
	reu._regs[0] |= 0x40

	// Clear execute bit in command register
	reu._regs[1] &= 0x7f
	reu.expansion.DMALow(false)
}
