package reu

import (
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/components/board"
	"strconv"
)

// Predefined constants representing memory sizes in bytes.
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

// Predefined constants representing different memory unit identifiers.
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

// REU represents a RAM Expansion Unit (REU) type with attributes for RAM, registers, IRQ masking, size, and expansion handling.
type REU struct {
	*board.BaseComponent
	loaderId  string
	ram       []uint8 // REU RAM
	mask      uint32  // REU RAM address bit mask
	regs      []uint8 // REU registers
	size      int
	irqMask   uint8
	expansion icartridge.IExpansion
}

// newReu initializes a new REU instance with a given size and returns an ICartridge implementation.
// It sets up REU registers, memory size, and RAM contents, and performs an initial reset.
func newReu(parentNode *board.Node, suffix string, size int) icartridge.ICartridge {
	r := &REU{
		BaseComponent: board.NewBaseComponent("reu", suffix, nil),
		loaderId:      "reu" + strconv.Itoa(size),
		regs:          make([]uint8, 16),
		size:          size,
		expansion:     nil,
		mask:          uint32(size) - 1,
		ram:           make([]uint8, size),
		irqMask:       0,
	}
	r.SetNode(board.CreateNode(parentNode, r))
	r.Reset()
	return r
}

// New128K creates and returns a new 128K REU cartridge implementing the ICartridge interface.
func New128K(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size128K)
}

// New256K creates a new 256K REU (RAM Expansion Unit) cartridge and returns it as an ICartridge interface.
func New256K(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size256K)
}

// New512K creates and returns a new instance of an ICartridge with a memory size of 512K.
func New512K(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size512K)
}

// New1M creates a new 1M REU (RAM Expansion Unit) cartridge implementing the ICartridge interface.
func New1M(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size1M)
}

// New2M creates and returns a new 2MB REU cartridge instance implementing the ICartridge interface.
func New2M(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size2M)
}

// New4M creates and returns a new 4MB REU cartridge implementing the ICartridge interface.
func New4M(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size4M)
}

// New8M creates and returns a new 8 MB REU (RAM Expansion Unit) cartridge implementing the ICartridge interface.
func New8M(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size8M)
}

// New16M creates a new ICartridge instance with 16MB of memory, utilizing the REU (RAM Expansion Unit) implementation.
func New16M(node *board.Node, suffix string) icartridge.ICartridge {
	return newReu(node, suffix, size16M)
}

// GetLoaderId returns the identifier string of the REU instance.
func (reu *REU) GetLoaderId() string {
	return reu.loaderId
}

// EmulationRequired determines whether emulation is required for the REU instance. Always returns false.
func (reu *REU) EmulationRequired() bool {
	return false
}

// Emulate performs the main emulation logic for the REU device, handling internal operations and DMA processing.
func (reu *REU) Emulate() {
}

// Reset initializes the REU's registers and status based on its RAM size.
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

// Setup initializes the REU instance by configuring its expansion board, identifying the loader, and setting up memory.
func (reu *REU) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	//TODO from Setup
	reu.expansion = board
	reu.loaderId = ldr.GetId()
	// Set kind bit in status register
	if reu.size-1 > 0x20000 {
		reu.regs[0] |= 0x10
	} else {
		reu.regs[0] &= 0xef
	}
	data := ldr.GetData()
	if len(data) > 0 {
		copy(reu.ram, data)
	}
	reu.expansion.RamSetWriteTrigger(0xff00, reu.triggerDMA)
	return nil
}

// GetExRom returns the value of the EXROM line, typically indicating cartridge configuration for memory mapping.
func (reu *REU) GetExRom() uint8 {
	return 1
}

// GetGame returns the game state as a uint8 value.
func (reu *REU) GetGame() uint8 {
	return 1
}

// IORead reads data from the specified REU address. Returns the read value and a boolean indicating success or failure.
// The address range is restricted to specific REU registers: $DF00-$DF0A, $DF20-$DF2A, and so on, up to $DFE0-$DFEA.
// Returns the value of the register or a modified value based on the accessed register type.
func (reu *REU) IORead(addr uint16) (uint8, bool) {
	//$DF00-$DF0A, $DF20-$DF2A, $DF40-$DF4A and so on, up to $DFE0-$DFEA
	if (addr & 0xfff0) != 0xdf00 {
		return 0, false
	}
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

// IOWrite processes a write operation to the REU's I/O registers based on the provided address and data.
// It validates the address range, updates the REU's internal registers, and triggers DMA execution if necessary.
// Returns true if the write operation is valid, otherwise false.
func (reu *REU) IOWrite(addr uint16, data uint8) bool {
	//$DF00-$DF0A, $DF20-$DF2A, $DF40-$DF4A and so on, up to $DFE0-$DFEA
	if (addr & 0xfff0) != 0xdf00 {
		return false
	}
	reg := addr & 0x0f
	switch reg {
	case 0, 11, 12, 13, 14, 15:
		// Status register is read-only
		// Unconnected registers
	case 1:
		// Command register $DF01
		reu.regs[1] = data
		if (data & 0x90) == 0x90 {
			reu.executeDma()
		}
	case 2:
		//c64base  = $DF02
		reu.regs[reg] = data
	case 3:
		reu.regs[reg] = data
	case 4:
		//reuBase  = $DF04
		reu.regs[reg] = data
	case 5:
		reu.regs[reg] = data
	case 6:
		reu.regs[reg] = data
	case 7:
		//transfer len = $DF07
		reu.regs[reg] = data
	case 8:
		reu.regs[reg] = data
	case 9:
		//irqMask  = $DF09
		reu.regs[reg] = data
		reu.irqMask = data
	case 10:
		//control  = $DF0A
		reu.regs[reg] = data
	default:
		reu.regs[reg] = data
	}
	return true
}

// triggerDMA handles DMA operations triggered by the CPU when writing to address $ff00 in REU.
func (reu *REU) triggerDMA(_ uint16, _ uint8) {
	// CPU triggered REU by writing to $ff00
	if (reu.regs[1] & 0x90) == 0x80 {
		reu.executeDma()
	}
}

// Write attempts to write to a given ROM interval, address, and value, returning false as operation is unsupported.
func (reu *REU) Write(_ icartridge.RomInterval, _ uint16, _ uint8) bool {
	return false
}

// Read retrieves data from the specified ROM interval and memory address. Returns the data and a boolean indicating success.
func (reu *REU) Read(_ icartridge.RomInterval, _ uint16) (uint8, bool) {
	return 0, false
}

// Detach releases resources or connections held by the REU instance and prepares it for a safe disconnection.
func (reu *REU) Detach() error {
	//TODO IMPLEMENT
	return nil
}

// executeDma performs a Direct Memory Access (DMA) transfer between the C64 and the REU based on the configured mode.
func (reu *REU) executeDma() {
	//TODO IN EMULATE ONE BYTE PER CYCLE
	//reu.expansion.SetDMALow(true)
	// Get C64 and REU transfer base addresses
	c64Addr := uint16(reu.regs[2]) | (uint16(reu.regs[3]) << 8)
	reuAddr := uint32(reu.regs[4]) | (uint32(reu.regs[5]) << 8) | (uint32(reu.regs[6]) << 16)

	// Calculate transfer length
	length := int(reu.regs[7]) | (int(reu.regs[8]) << 8)
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
		for ; length > 0; length-- {
			reu.ram[reuAddr&reu.mask] = reu.expansion.Read(c64Addr)
			c64Addr += uint16(c64Inc)
			reuAddr += reuInc
		}

	case 1:
		// C64 <- REU
		for ; length > 0; length-- {
			reu.expansion.Write(c64Addr, reu.ram[reuAddr&reu.mask])
			c64Addr += uint16(c64Inc)
			reuAddr += reuInc
		}
	case 2:
		// C64 <-> REU
		for ; length > 0; length-- {
			tmp := reu.expansion.Read(c64Addr)
			reu.expansion.Write(c64Addr, reu.ram[reuAddr&reu.mask])
			reu.ram[reuAddr&reu.mask] = tmp
			c64Addr += uint16(c64Inc)
			reuAddr += reuInc
		}

	case 3:
		// Compare
		for ; length > 0; length-- {
			if reu.ram[reuAddr&reu.mask] != reu.expansion.Read(c64Addr) {
				reu.regs[0] |= 0x20
				break
			}
			c64Addr += uint16(c64Inc)
			reuAddr += reuInc
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
	//reu.expansion.SetDMALow(false)
}
