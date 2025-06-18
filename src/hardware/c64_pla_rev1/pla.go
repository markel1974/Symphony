package c64_pla_rev1

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

	cartManRead      ReadFn
	cartManIORead    func(uint16) (uint8, bool)
	cartManIOWrite   func(uint16, uint8) bool
	cartManConfig    func() (uint8, uint8, bool)
	cartManIntervals uint8

	vaSignals func() uint8

	bankWrite       []WriteFn
	bankRead        []ReadFn
	u15Write        []WriteFn // represent U15, 74LS138 (mux)
	u15Read         []ReadFn  // represent U15, 74LS138 (mux)
	memoryMap       *MemoryMap
	memoryConfigIdx int
	memoryConfig    []uint8
	ports           *Ports
	emulatorId      *EmulatorId
	basicRead       ReadFn
	kernalRead      ReadFn
	charRead        ReadFn
	colorRead       ReadFn
	//colorWrite      WriteFn
	cfg       *config.Config
	wTriggers *WriteTriggers
	label     string
}

const bankMask = 0x0f
const bankSize = bankMask + 1

// NewPLA initializes and returns a new instance of the PLA component with its memory map and configurations set up.
func NewPLA(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *PLA {
	mm := NewMemoryMap()
	b := &PLA{
		BaseComponent:    component.NewBaseComponent(),
		vaSignals:        nil,
		triggerSize:      0,
		ramRead:          nil,
		ramWrite:         nil,
		cartManRead:      nil,
		cartManIORead:    nil,
		cartManIOWrite:   nil,
		bankWrite:        make([]WriteFn, bankSize),
		bankRead:         make([]ReadFn, bankSize),
		u15Write:         make([]WriteFn, bankSize),
		u15Read:          make([]ReadFn, bankSize),
		memoryMap:        mm,
		memoryConfig:     mm.Get(0),
		ports:            nil,
		emulatorId:       NewEmulatorId(),
		basicRead:        nil,
		kernalRead:       nil,
		charRead:         nil,
		colorRead:        nil,
		cfg:              nil,
		memoryConfigIdx:  -1,
		wTriggers:        nil,
		label:            label,
		cartManIntervals: 0,
	}
	b.BaseComponent.Register(factory, parent, Identifier(), b, references.IdIC64Pla(b, label, instance))
	return b
}

// Setup initializes the PLA instance with the provided socket and configuration and returns an error if any issue occurs.
func (b *PLA) Setup() error {
	b.cfg = b.GetFactory().GetConfig()
	b.ports = NewPorts(b.GetFactory(), b, b.label, 0)
	return nil
}

//U15 74LS138
//Pin 12	VIC-II	$D000–$D3FF
//Pin 04	SID	$D400–$D7FF
//Pin 10	Color RAM	$D800–$DBFF
//Pin 14	CIA1	$DC00–$DCFF
//Pin 15	CIA2	$DD00–$DDFF
//Pin 09	I/O2 (cartucce)	$DF00–$DFFF
//Pin 11	I/O1 (cartucce)	$DE00–$DEFF

// Bind initializes and connects various components to the PLA, including VIC, SID, CIA1, CIA2, cartridge manager, and ROM loader.
func (b *PLA) Bind(_ references.IC64PlaSocket, vaSignals references.IC64PlaVASignals, cartMan references.IC64CartridgeManager, ram references.IC64Ram, roms references.IC64Roms, cs12 references.IC64PlaChipSelect, cs4 references.IC64PlaChipSelect, cs14 references.IC64PlaChipSelect, cs15 references.IC64PlaChipSelect, cs10 references.IC64PlaChipSelect) error {
	b.vaSignals = vaSignals.GetVASignal
	b.triggerSize = ram.Size()

	b.cartManRead = cartMan.Read

	//MUST BE cs9 - cs11
	b.cartManIORead = cartMan.IORead
	b.cartManIOWrite = cartMan.IOWrite

	b.cartManConfig = cartMan.Config

	b.ramRead = ram.Read
	b.ramWrite = ram.Write

	b.colorRead = cs10.ReadRegister

	for idx := range b.bankWrite {
		b.bankWrite[idx] = b.ramWrite
	}
	for idx := range b.bankRead {
		b.bankRead[idx] = b.ramRead
	}
	//pla write mapping
	b.bankWrite[0x0] = b.ramWrite0x0000
	//pla read mapping
	b.bankRead[0x0] = b.ramRead0x0000

	b.u15Write[0x0] = cs12.WriteRegister
	b.u15Write[0x1] = cs12.WriteRegister
	b.u15Write[0x2] = cs12.WriteRegister
	b.u15Write[0x3] = cs12.WriteRegister
	b.u15Write[0x4] = cs4.WriteRegister
	b.u15Write[0x5] = cs4.WriteRegister
	b.u15Write[0x6] = cs4.WriteRegister
	b.u15Write[0x7] = cs4.WriteRegister
	b.u15Write[0x8] = cs10.WriteRegister //b.portWriteColor
	b.u15Write[0x9] = cs10.WriteRegister //b.portWriteColor
	b.u15Write[0xa] = cs10.WriteRegister //b.portWriteColor
	b.u15Write[0xb] = cs10.WriteRegister //b.portWriteColor
	b.u15Write[0xc] = cs14.WriteRegister
	b.u15Write[0xd] = cs15.WriteRegister
	b.u15Write[0xe] = b.portWriteIO
	b.u15Write[0xf] = b.portWriteIO

	b.u15Read[0x0] = cs12.ReadRegister
	b.u15Read[0x1] = cs12.ReadRegister
	b.u15Read[0x2] = cs12.ReadRegister
	b.u15Read[0x3] = cs12.ReadRegister
	b.u15Read[0x4] = cs4.ReadRegister
	b.u15Read[0x5] = cs4.ReadRegister
	b.u15Read[0x6] = cs4.ReadRegister
	b.u15Read[0x7] = cs4.ReadRegister
	b.u15Read[0x8] = b.portReadColor
	b.u15Read[0x9] = b.portReadColor
	b.u15Read[0xa] = b.portReadColor
	b.u15Read[0xb] = b.portReadColor
	b.u15Read[0xc] = cs14.ReadRegister
	b.u15Read[0xd] = cs15.ReadRegister
	b.u15Read[0xe] = b.portReadIO
	b.u15Read[0xf] = b.portReadIO

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

// RebuildMemoryConfig updates the PLA's current memory configuration based on the cartridge and port settings.
func (b *PLA) RebuildMemoryConfig() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	spec := references.C64CartridgeSpecOff
	if game, exRom, ok := b.cartManConfig(); ok {
		spec = references.GetCartridgeSpec(game, exRom)
	}
	b.cartManIntervals = spec.Intervals
	dir, data := b.ports.Config()
	mcIdx := ((^dir | data) & 0x7) | (spec.ExRom << 3) | (spec.Game << 4)
	b.applyMemoryConfig(int(mcIdx))
}

// update adjusts the PLA's state by synchronizing port and memory configurations. It ensures the system's consistency.
func (b *PLA) update() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	//b.ports.SetTape(tape_sense, tape_write_in, tape_motor_in)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

// ExtWrite writes a byte to the specified address with an optional memory configuration adjustment.
func (b *PLA) ExtWrite(memConfig int, addr uint16, data uint8) {
	var prevMemConfig = -1
	if memConfig >= 0 {
		prevMemConfig = b.memoryConfigIdx
		b.applyMemoryConfig(memConfig)
	}
	b.Write(addr, data)
	if prevMemConfig >= 0 {
		b.applyMemoryConfig(prevMemConfig)
	}
}

// ExtRead reads a byte of data from the specified memory address with a temporary memory configuration applied.
func (b *PLA) ExtRead(memConfig int, addr uint16) uint8 {
	var prevMemConfig = -1
	if memConfig >= 0 {
		prevMemConfig = b.memoryConfigIdx
		b.applyMemoryConfig(memConfig)
	}
	data := b.Read(addr)
	if prevMemConfig >= 0 {
		b.applyMemoryConfig(prevMemConfig)
	}
	return data
}

// Read retrieves the value at the specified memory address using the current memory bank configuration.
func (b *PLA) Read(addr uint16) uint8 {
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	return b.bankRead[bank](addr)
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
		//b.ramWrite(0, b.vicLastByte())
		b.update()
		return
	} else if addr == 1 {
		b.ports.SetData(data)
		//b.ramWrite(1, b.vicLastByte())
		b.update()
		return
	}
	b.ramWrite(addr, data)
}

// ramWrite0xD000 writes a byte of data to the specified address in the 0xD000 memory range based on the memory configuration.
func (b *PLA) ramWrite0xD000_I_O(addr uint16, data uint8) {
	p := (addr >> 8) & 0x0f
	b.u15Write[p](addr, data)
	return
}

func (b *PLA) ramRead0xD000_I_O(addr uint16) uint8 {
	p := (addr >> 8) & 0x0f
	return b.u15Read[p](addr)
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

// applyMemoryConfig updates the current memory configuration if the provided index differs from the current configuration.
// It modifies the memoryConfigIdx and memoryConfig fields and returns true if the configuration was updated, false otherwise.
func (b *PLA) applyMemoryConfig(mcIdx int) bool {
	if mcIdx == b.memoryConfigIdx {
		return false
	}
	b.memoryConfigIdx = mcIdx
	b.memoryConfig = b.memoryMap.Get(uint8(mcIdx))

	if b.memoryConfig[0xd] == I_O {
		b.bankWrite[0xd] = b.ramWrite0xD000_I_O
	} else {
		b.bankWrite[0xd] = b.ramWrite
	}

	if mc := b.memoryConfig[0x8]; mc == ROL && b.cartManIntervals&references.ROM_LO == references.ROM_LO {
		b.bankRead[0x8] = b.cartManRead
	} else {
		b.bankRead[0x8] = b.ramRead
	}

	if mc := b.memoryConfig[0x9]; mc == ROL && b.cartManIntervals&references.ROM_LO == references.ROM_LO {
		b.bankRead[0x9] = b.cartManRead
	} else {
		b.bankRead[0x9] = b.ramRead
	}

	if mc := b.memoryConfig[0xa]; mc == BAS {
		b.bankRead[0xa] = b.basicRead
	} else if mc == ROH && b.cartManIntervals&references.ROM_HI_1 == references.ROM_HI_1 {
		b.bankRead[0xa] = b.cartManRead
	} else {
		b.bankRead[0xa] = b.ramRead
	}

	if mc := b.memoryConfig[0xb]; mc == BAS {
		b.bankRead[0xb] = b.basicRead
	} else if mc == ROH && b.cartManIntervals&references.ROM_HI_1 == references.ROM_HI_1 {
		b.bankRead[0xb] = b.cartManRead
	} else {
		b.bankRead[0xb] = b.ramRead
	}

	if mc := b.memoryConfig[0xd]; mc == I_O {
		b.bankRead[0xd] = b.ramRead0xD000_I_O
	} else if mc == CHA {
		b.bankRead[0xd] = b.charRead
	} else {
		b.bankRead[0xd] = b.ramRead
	}

	if mc := b.memoryConfig[0xe]; mc == KER {
		b.bankRead[0xe] = b.kernalRead
	} else if mc == ROH && b.cartManIntervals&references.ROM_HI_2 == references.ROM_HI_2 {
		b.bankRead[0xe] = b.cartManRead
	} else {
		b.bankRead[0xe] = b.ramRead
	}

	if mc := b.memoryConfig[0xf]; mc == KER {
		b.bankRead[0xf] = b.kernalRead
	} else if mc == ROH && b.cartManIntervals&references.ROM_HI_2 == references.ROM_HI_2 {
		b.bankRead[0xf] = b.cartManRead
	} else {
		b.bankRead[0xf] = b.ramRead
	}
	//fmt.Printf("SYSTEM MEMORY CONFIG CHANGED [exRom: %d - game: %d] %d -> %v\n", exRom, game, mcIdx, b.memoryConfig)
	return true
}

// portReadColor reads a color value from the specified address in the color memory, merging specific bits from VIC data.
func (b *PLA) portReadColor(addr uint16) uint8 {
	p1 := b.colorRead(addr) & 0x0f // enables the physical Color RAM chip. This chip receives the address and puts the 4 color bits it has stored on data bus lines D0-D3 (lower half).
	p2 := b.vaSignals() & 0xf0     // signals to VIC that Color RAM is being read. VIC responds by putting the last 4 bits of its internal latch (lastByte) on data bus lines D4-D7 (upper half)
	return p1 | p2
}

// portReadIO reads a byte from a specified I/O port address using the provided memory and I/O mappings.
func (b *PLA) portReadIO(addr uint16) uint8 {
	if v, ok := b.cartManIORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vaSignals()
	}
	return b.emulatorId.Read(addr)
}

// portWriteIO handles writing a byte of data to the specified IO port address.
// If the write operation is handled by the cartMan, it exits early.
func (b *PLA) portWriteIO(addr uint16, data uint8) {
	_ = b.cartManIOWrite(addr, data)
}
