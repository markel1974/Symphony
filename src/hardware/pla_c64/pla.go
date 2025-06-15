package pla_c64

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// ReadFn represents a function type that takes a 16-bit unsigned integer as input and returns an 8-bit unsigned integer.
type ReadFn func(uint16) uint8

// WriteFn represents a function that processes a 16-bit address and an 8-bit value.
type WriteFn func(uint16, uint8)

// PLA represents the Programmable Logic Array implementation associated with memory and I/O operations in the C64 emulator.
type PLA struct {
	*component.BaseComponent

	triggerSize int

	ramRead  ReadFn
	ramWrite WriteFn

	cartManRead    func(references.RomInterval, uint16) (uint8, bool)
	cartManWrite   func(references.RomInterval, uint16, uint8) bool
	cartManIORead  func(uint16) (uint8, bool)
	cartManIOWrite func(uint16, uint8) bool
	cartManConfig  func() (uint8, uint8, bool)

	vicLastByte func() uint8

	bankWrite       []WriteFn
	bankRead        []ReadFn
	portWrite       []WriteFn
	portRead        []ReadFn
	memoryMap       *MemoryMap
	memoryConfigIdx int
	memoryConfig    []uint8
	ports           *Ports
	emulatorId      *EmulatorId
	basicRead       ReadFn
	kernalRead      ReadFn
	charRead        ReadFn
	colorRead       ReadFn
	colorWrite      WriteFn
	cfg             *config.Config
	wTriggers       *WriteTriggers
	label           string
}

const bankMask = 0x0f
const bankSize = bankMask + 1

// NewPLA initializes and returns a new instance of the PLA component with its memory map and configurations set up.
func NewPLA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *PLA {
	mm := NewMemoryMap()
	b := &PLA{
		BaseComponent:   component.NewBaseComponent(),
		vicLastByte:     nil,
		triggerSize:     0,
		ramRead:         nil,
		ramWrite:        nil,
		cartManRead:     nil,
		cartManWrite:    nil,
		cartManIORead:   nil,
		cartManIOWrite:  nil,
		bankWrite:       make([]WriteFn, bankSize),
		bankRead:        make([]ReadFn, bankSize),
		portWrite:       make([]WriteFn, bankSize),
		portRead:        make([]ReadFn, bankSize),
		memoryMap:       mm,
		memoryConfig:    mm.Get(0),
		ports:           nil,
		emulatorId:      NewEmulatorId(),
		basicRead:       nil,
		kernalRead:      nil,
		charRead:        nil,
		colorRead:       nil,
		colorWrite:      nil,
		cfg:             nil,
		memoryConfigIdx: -1,
		wTriggers:       nil,
		label:           label,
	}
	b.BaseComponent.Register(factory, parent, Identifier(), b, references.IdIPlaC64(b, label, instance))
	return b
}

// Setup initializes the PLA instance with the provided socket and configuration and returns an error if any issue occurs.
func (b *PLA) Setup() error {
	b.cfg = b.GetFactory().GetConfig()
	b.ports = NewPorts(b.GetFactory(), b, b.label, 0)
	return nil
}

// Bind initializes and connects various components to the PLA, including VIC, SID, CIA1, CIA2, cartridge manager, and ROM loader.
func (b *PLA) Bind(_ references.IPlaC64Socket, vic references.IVIC, sid references.ISID, cia1 references.ICIA, cia2 references.ICIA, cartMan references.ICartridgeManagerC64, ram references.IRamC64, roms references.IRomsC64) error {
	b.vicLastByte = vic.GetLastByte

	b.cartManWrite = cartMan.Write
	b.cartManRead = cartMan.Read
	b.cartManIORead = cartMan.IORead
	b.cartManIOWrite = cartMan.IOWrite
	b.cartManConfig = cartMan.Config

	b.ramRead = ram.Read
	b.ramWrite = ram.Write
	b.triggerSize = ram.Size()
	b.colorRead = ram.ReadColor
	b.colorWrite = ram.WriteColor

	for idx := range b.bankWrite {
		b.bankWrite[idx] = b.ramWrite
	}
	for idx := range b.bankRead {
		b.bankRead[idx] = b.ramRead
	}
	//pla write mapping
	b.bankWrite[0x0] = b.ramWrite0x0000
	b.bankWrite[0xd] = b.ramWrite0xD000

	//pla read mapping
	b.bankRead[0x0] = b.ramRead0x0000
	b.bankRead[0x8] = b.ramRead0x8000
	b.bankRead[0x9] = b.ramRead0x9000
	b.bankRead[0xa] = b.ramRead0xA000
	b.bankRead[0xb] = b.ramRead0xB000
	b.bankRead[0xd] = b.ramRead0xD000
	b.bankRead[0xe] = b.ramRead0xE000
	b.bankRead[0xf] = b.ramRead0xF000

	b.portWrite[0x0] = vic.WriteRegister
	b.portWrite[0x1] = vic.WriteRegister
	b.portWrite[0x2] = vic.WriteRegister
	b.portWrite[0x3] = vic.WriteRegister
	b.portWrite[0x4] = sid.WriteRegister
	b.portWrite[0x5] = sid.WriteRegister
	b.portWrite[0x6] = sid.WriteRegister
	b.portWrite[0x7] = sid.WriteRegister
	b.portWrite[0x8] = b.portWriteColor
	b.portWrite[0x9] = b.portWriteColor
	b.portWrite[0xa] = b.portWriteColor
	b.portWrite[0xb] = b.portWriteColor
	b.portWrite[0xc] = cia1.WriteRegister
	b.portWrite[0xd] = cia2.WriteRegister
	b.portWrite[0xe] = b.portWriteIO
	b.portWrite[0xf] = b.portWriteIO

	b.portRead[0x0] = vic.ReadRegister
	b.portRead[0x1] = vic.ReadRegister
	b.portRead[0x2] = vic.ReadRegister
	b.portRead[0x3] = vic.ReadRegister
	b.portRead[0x4] = sid.ReadRegister
	b.portRead[0x5] = sid.ReadRegister
	b.portRead[0x6] = sid.ReadRegister
	b.portRead[0x7] = sid.ReadRegister
	b.portRead[0x8] = b.portReadColor
	b.portRead[0x9] = b.portReadColor
	b.portRead[0xa] = b.portReadColor
	b.portRead[0xb] = b.portReadColor
	b.portRead[0xc] = cia1.ReadRegister
	b.portRead[0xd] = cia2.ReadRegister
	b.portRead[0xe] = b.portReadIO
	b.portRead[0xf] = b.portReadIO
	b.kernalRead = roms.KernalRead
	b.basicRead = roms.BasicRead
	b.charRead = roms.CharRead
	return nil
}

// Connect establishes a connection using the PLA object and returns an error if the connection fails.
func (b *PLA) Connect() error {
	return nil
}

// Internal checks internal state or behavior specific to the PLA instance and returns a boolean value.
func (b *PLA) Internal() bool {
	return false
}

// Reset reinitializes the internal state of the PLA, including resetting its ports and updating related configurations.
func (b *PLA) Reset() {
	b.ports.Reset()
	b.update()
}

// Emulate runs the emulation process for the PLA, executing its defined logic and behavior within the system's context.
func (b *PLA) Emulate() {
	//
}

// EmulationRequired checks if emulation is needed for the current PLA instance and returns false by default.
func (b *PLA) EmulationRequired() bool {
	return false
}

// update adjusts the PLA's state by synchronizing port and memory configurations. It ensures the system's consistency.
func (b *PLA) update() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	//b.ports.SetTape(tape_sense, tape_write_in, tape_motor_in)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

// RebuildMemoryConfig updates the PLA's current memory configuration based on the cartridge and port settings.
func (b *PLA) RebuildMemoryConfig() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	game := uint8(1)
	exRom := uint8(1)
	if g, x, ok := b.cartManConfig(); ok {
		game = g
		exRom = x
	}
	mcIdx := b.ports.GetMemoryConfig(exRom, game)
	if int(mcIdx) != b.memoryConfigIdx {
		b.memoryConfigIdx = int(mcIdx)
		b.memoryConfig = b.memoryMap.Get(mcIdx)
		//fmt.Printf("SYSTEM MEMORY CONFIG CHANGED [exRom: %d - game: %d] %d -> %v\n", exRom, game, mcIdx, b.memoryConfig)
	}
}

// GetMemoryConfig returns the memory configuration of the PLA as a slice of uint8.
func (b *PLA) GetMemoryConfig() []uint8 {
	return b.memoryConfig
}

// SetMemoryConfig sets the memory configuration for the PLA instance using the provided configuration slice.
func (b *PLA) SetMemoryConfig(cfg []uint8) {
	b.memoryConfig = cfg
}

// SetMemoryEntry sets the memory configuration for the PLA using the provided configuration value.
func (b *PLA) SetMemoryEntry(memConfig uint8) {
	b.memoryConfig = b.memoryMap.Get(memConfig)
}

// ReadCharRom reads a byte from the character ROM at the specified address and returns it.
func (b *PLA) ReadCharRom(addr uint16) uint8 {
	return b.charRead(addr)
}

// ReadColor retrieves the color value from the specified memory address. Returns an 8-bit unsigned integer value.
func (b *PLA) ReadColor(addr uint16) uint8 {
	return b.colorRead(addr)
}

// Read retrieves the value at the specified memory address using the current memory bank configuration.
func (b *PLA) Read(addr uint16) uint8 {
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	return b.bankRead[bank](addr)
}

// ReadDirect accesses and returns the value stored at the specified memory address without any additional logic.
func (b *PLA) ReadDirect(addr uint16) uint8 {
	return b.ramRead(addr)
}

// WriteDirect writes the provided data to the specified address in RAM and executes any write triggers if set.
func (b *PLA) WriteDirect(addr uint16, data uint8) {
	b.ramWrite(addr, data)
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Exec(addr, data)
}

// Write writes a single byte of data to a specified memory address and triggers any assigned write handlers.
func (b *PLA) Write(addr uint16, data uint8) {
	//sta.c64.org/cbm64mem.html
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	b.bankWrite[bank](addr, data)
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Exec(addr, data)
}

// SetWriteTrigger sets a write trigger function for the specified address and returns the trigger ID.
func (b *PLA) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	if b.wTriggers == nil {
		b.wTriggers = NewWriteTriggers(b.triggerSize)
	}
	return b.wTriggers.Add(addr, fn)
}

// RemoveRamTrigger removes a specified RAM trigger identified by the address and ID provided.
func (b *PLA) RemoveRamTrigger(addr uint16, id int) {
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Remove(addr, id)
}

// ramWrite0x0000 writes a byte to the specified RAM address, handling special cases for addresses 0x0000 and 0x0001.
func (b *PLA) ramWrite0x0000(addr uint16, data uint8) {
	if addr == 0 {
		b.ports.SetDir(data)
		b.ramWrite(0, b.vicLastByte())
		b.update()
		return
	} else if addr == 1 {
		b.ports.SetData(data)
		b.ramWrite(1, b.vicLastByte())
		b.update()
		return
	}
	b.ramWrite(addr, data)
}

// ramWrite0xD000 writes a byte of data to the specified address in the 0xD000 memory range based on the memory configuration.
func (b *PLA) ramWrite0xD000(addr uint16, data uint8) {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		b.portWrite[p](addr, data)
		return
	}
	b.ramWrite(addr, data)
}

// ramRead0x0000 reads a byte from RAM or retrieves data from port registers based on the provided address.
func (b *PLA) ramRead0x0000(addr uint16) uint8 {
	if addr == 0 {
		return b.ports.GetDirection()
	} else if addr == 1 {
		return b.ports.GetDataRead()
	}
	return b.ramRead(addr)
}

// ramRead0x8000 reads a byte from RAM at the specified address or from a cartridge if configured in ROM_LO mode.
func (b *PLA) ramRead0x8000(addr uint16) uint8 {
	const bank = 0x8
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartManRead(references.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ramRead(addr)
}

// ramRead0x9000 reads a byte of data from the RAM or cartridge memory at a given address within the 0x9000 bank range.
// It checks the memory configuration of the specified bank and accesses data accordingly.
func (b *PLA) ramRead0x9000(addr uint16) uint8 {
	const bank = 0x9
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartManRead(references.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ramRead(addr)
}

// ramRead0xA000 reads data from memory mapped to the 0xA000 address range based on the current memory configuration.
// For ROH configuration, it attempts to read from cartridge memory.
// For BAS configuration, it returns data from the BASIC ROM bank.
// Defaults to reading from RAM if no other condition is met.
func (b *PLA) ramRead0xA000(addr uint16) uint8 {
	const bank = 0xa
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartManRead(references.ROM_HI_1, addr); ok {
			return v
		}

	} else if b.memoryConfig[bank] == BAS {
		return b.basicRead(addr)
	}
	return b.ramRead(addr)
}

// ramRead0xB000 handles reading from the 0xB000 bank based on memory configuration and address.
// It prioritizes cartridge ROM, BASIC ROM, or falls back to RAM as per the active memory configuration.
// addr represents the memory address to be read within the 0xB000 range.
// Returns the 8-bit value read from the specified memory address.
func (b *PLA) ramRead0xB000(addr uint16) uint8 {
	const bank = 0xb
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartManRead(references.ROM_HI_1, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == BAS {
		return b.basicRead(addr)
	}
	return b.ramRead(addr)
}

// ramRead0xD000 reads data from memory address 0xD000 and handles different configurations like I/O, character, or RAM access.
// It uses the memory configuration to determine the appropriate data source and returns the byte at the specified address.
// If the memory is set to I/O, it delegates the read operation to the appropriate port handler function.
// For character memory configuration, it accesses the char array using the offset masked from the given address.
// Defaults to reading directly from RAM if no special memory configuration is in place for the 0xD000 bank.
func (b *PLA) ramRead0xD000(addr uint16) uint8 {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		return b.portRead[p](addr)
	} else if b.memoryConfig[bank] == CHA {
		return b.charRead(addr)
	}
	return b.ramRead(addr)
}

// ramRead0xE000 handles reading from the 0xE000-0xEFFF memory range based on the current memory configuration.
// It prioritizes cartMan, kernal ROM, or RAM depending on the memory bank setup for the 0xE block.
func (b *PLA) ramRead0xE000(addr uint16) uint8 {
	const bank = 0xe
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartManRead(references.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernalRead(addr)
	}
	return b.ramRead(addr)
}

// ramRead0xF000 reads a byte from the 0xF000 memory range based on the current memory configuration and provided address.
// Depending on configuration, it accesses cartridge ROM, kernal ROM, or internal RAM to retrieve the value.
// The addr parameter specifies the address within the 0xF000 range to read.
func (b *PLA) ramRead0xF000(addr uint16) uint8 {
	const bank = 0xf
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartManRead(references.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernalRead(addr)
	}
	return b.ramRead(addr)
}

// portWriteColor writes a color value to the color buffer at the given address, using only the lower 4 bits of the data.
func (b *PLA) portWriteColor(addr uint16, data uint8) {
	b.colorWrite(addr, data)
}

// portReadColor reads a color value from the specified address in the color memory, merging specific bits from VIC data.
func (b *PLA) portReadColor(addr uint16) uint8 {
	return (b.colorRead(addr) & 0x0f) | (b.vicLastByte() & 0xf0)
}

// portReadIO reads a byte from a specified I/O port address using the provided memory and I/O mappings.
func (b *PLA) portReadIO(addr uint16) uint8 {
	if v, ok := b.cartManIORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vicLastByte()
	}
	return b.emulatorId.Read(addr)
}

// portWriteIO handles writing a byte of data to the specified IO port address.
// If the write operation is handled by the cartMan, it exits early.
func (b *PLA) portWriteIO(addr uint16, data uint8) {
	if ok := b.cartManIOWrite(addr, data); ok {
		return
	}
}
