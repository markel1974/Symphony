package c64_pla_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// ReadFn defines a function type that takes a 16-bit unsigned integer as input and returns an 8-bit unsigned integer.
type ReadFn func(uint16) uint8

// WriteFn defines a function type that takes a 16-bit address and an 8-bit data value for write operations.
type WriteFn func(uint16, uint8)

// PLA represents the Programmable Logic Array (PLA) component responsible for memory and I/O address mapping logic.
type PLA struct {
	*component.BaseComponent
	triggerSize      int
	ramRead          ReadFn
	ramWrite         WriteFn
	cartManRead      ReadFn
	cartManIORead    func(uint16) (uint8, bool)
	cartManIOWrite   func(uint16, uint8) bool
	cartManConfig    func() (uint8, uint8, bool)
	cartManIntervals uint8
	vaSignals        func() uint8
	bankWrite        []WriteFn
	bankRead         []ReadFn
	u15Write         []WriteFn // represent U15, 74LS138 (mux)
	u15Read          []ReadFn  // represent U15, 74LS138 (mux)
	memoryMap        *MemoryMap
	memoryConfigIdx  int
	memoryConfig     []uint8
	ports            *Ports
	basicRead        ReadFn
	kernalRead       ReadFn
	charRead         ReadFn
	colorRead        ReadFn
	cfg              *config.Config
	wTriggers        *WriteTriggers
	label            string
	emulatorId       *EmulatorId
}

// bankMask is a constant used as a bitmask to extract or manipulate specific bits in a value, often for bank-related operations.
const bankMask = 0x0f

// bankSize represents the total size of the bank, calculated as bankMask + 1.
const bankSize = bankMask + 1

// NewPLA initializes and returns a new PLA (Programmable Logic Array) component with the specified configuration.
// It registers the component with the given parent and factory, assigns a label, and sets the instance number.
// The function also sets up memory mapping, triggers, and other internal properties specific to the PLA.
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

// Setup initializes the PLA instance by configuring its settings and creating the port object.
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

// Bind initializes and binds PLA components such as RAM, ROM, cartridge manager, chip-select elements, and signal handlers.
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

// Connect initializes and sets up the PLA component, enabling necessary configurations for proper operation.
func (b *PLA) Connect() error {
	return nil
}

// Internal checks and returns whether the PLA's internal operation mode is enabled or active.
func (b *PLA) Internal() bool {
	return false
}

// Reset resets the internal state of the PLA, reinitializing ports, and triggering an update of the memory configuration.
func (b *PLA) Reset() {
	b.ports.Reset()
	b.update()
}

// Emulate performs the main emulation step for the PLA, updating its state and processing associated logic.
func (b *PLA) Emulate() {
	//
}

// EmulationRequired determines whether emulation is required for the PLA component. Returns false by default.
func (b *PLA) EmulationRequired() bool {
	return false
}

// RebuildMemoryConfig recalculates and applies the current memory configuration based on cartridge and port settings.
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

// update updates the state of the PLA by rebuilding the memory configuration and updating the port settings.
func (b *PLA) update() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	//b.ports.SetTape(tape_sense, tape_write_in, tape_motor_in)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

// ExtWrite writes a byte to the specified address using the provided memory configuration, temporarily switching configurations.
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

// ExtRead allows reading a specific memory address using a temporary memory configuration.
// The original memory configuration is restored after the read operation.
// memConfig specifies the temporary memory configuration to apply.
// addr is the address to read from.
// Returns the data read from the given memory address.
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

// Read retrieves the value from the memory mapped by the given address.
// The address is used to determine the memory bank by shifting its bits.
// The bank's read function is then invoked with the address.
func (b *PLA) Read(addr uint16) uint8 {
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	return b.bankRead[bank](addr)
}

// Write writes a byte of data to a specific memory address based on the bank configuration and triggers write callbacks.
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

// SetWriteTrigger registers a write-trigger at a specified address with a callback and returns the assigned trigger ID.
func (b *PLA) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	if b.wTriggers == nil {
		b.wTriggers = NewWriteTriggers(b.triggerSize)
	}
	return b.wTriggers.Add(addr, fn)
}

// RemoveRamTrigger removes a write-trigger by its unique ID for the specified memory address in the WriteTriggers collection.
func (b *PLA) RemoveRamTrigger(addr uint16, id int) {
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Remove(addr, id)
}

// ramWrite0x0000 handles writing data to memory address 0x0000, updating port direction and data where applicable.
func (b *PLA) ramWrite0x0000(addr uint16, data uint8) {
	if addr == 0 {
		b.ports.SetDir(data)
		//b.ramWrite(0, b.vaSignals)
		b.update()
		return
	} else if addr == 1 {
		b.ports.SetData(data)
		//b.ramWrite(1, b.vaSignals())
		b.update()
		return
	}
	b.ramWrite(addr, data)
}

// ramWriteIO writes a byte of data to the specific address by determining the target U15 write function based on the address.
func (b *PLA) ramWriteIO(addr uint16, data uint8) {
	p := (addr >> 8) & 0x0f
	b.u15Write[p](addr, data)
	return
}

// ramReadIO returns the value of the RAM read at the specified address, using the appropriate U15 read function.
func (b *PLA) ramReadIO(addr uint16) uint8 {
	p := (addr >> 8) & 0x0f
	return b.u15Read[p](addr)
}

// ramRead0x0000 reads data from address 0x0000 and processes special cases for addresses 0 and 1.
// For address 0, it returns the direction register value from the ports.
// For address 1, it returns the data register value from the ports.
// For other addresses, it delegates to the generic RAM read function.
func (b *PLA) ramRead0x0000(addr uint16) uint8 {
	if addr == 0 {
		return b.ports.GetDirection()
	} else if addr == 1 {
		return b.ports.GetDataRead()
	}
	return b.ramRead(addr)
}

// portReadColor reads color data from the Color RAM and combines it with data from the VIC latch for a full byte response.
func (b *PLA) portReadColor(addr uint16) uint8 {
	p1 := b.colorRead(addr) & 0x0f // enables the physical Color RAM chip. This chip receives the address and puts the 4 color bits it has stored on data bus lines D0-D3 (lower half).
	p2 := b.vaSignals() & 0xf0     // signals to VIC that Color RAM is being read. VIC responds by putting the last 4 bits of its internal latch (lastByte) on data bus lines D4-D7 (upper half)
	return p1 | p2
}

// portReadIO reads a byte from the specified IO port address.
// It prioritizes cartridge IO reads, falls back to VA signal reads, or accesses the emulator ID if applicable.
func (b *PLA) portReadIO(addr uint16) uint8 {
	if v, ok := b.cartManIORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vaSignals()
	}
	return b.emulatorId.Read(addr)
}

// portWriteIO writes data to the specified IO port address using the cartridge manager's IO write functionality.
func (b *PLA) portWriteIO(addr uint16, data uint8) {
	_ = b.cartManIOWrite(addr, data)
}

// writeOpenBus writes to an "open bus" state where no specific memory or component is targeted.
// Used for situations where writes do not directly influence system behavior or state.
func (b *PLA) writeOpenBus(_ uint16, _ uint8) {
	//
}

// readOpenBus reads the current value of the open bus at the specified address and returns its signals as uint8.
func (b *PLA) readOpenBus(_ uint16) uint8 {
	return b.vaSignals()
}

// applyMemoryConfig updates the current memory configuration based on the given index and returns true if a change occurred.
func (b *PLA) applyMemoryConfig(mcIdx int) bool {
	if mcIdx == b.memoryConfigIdx {
		return false
	}
	b.memoryConfigIdx = mcIdx
	b.memoryConfig = b.memoryMap.Get(uint8(mcIdx))

	for idx, v := range b.memoryConfig {
		if idx == 0 {
			b.bankRead[0] = b.ramRead0x0000
			b.bankWrite[0] = b.ramWrite0x0000
		} else {
			b.bankRead[idx] = b.ramRead
			b.bankWrite[idx] = b.ramWrite
		}
		switch v {
		case RAM:
			// Default is RAM, no action needed.
		case KER:
			b.bankRead[idx] = b.kernalRead
		case BAS:
			b.bankRead[idx] = b.basicRead
		case CHA:
			b.bankRead[idx] = b.charRead
		case I_O:
			b.bankRead[idx] = b.ramReadIO
			b.bankWrite[idx] = b.ramWriteIO
		case ROL:
			if b.cartManIntervals&references.ROM_LO == references.ROM_LO {
				b.bankRead[idx] = b.cartManRead
			}
		case ROH:
			if (idx == 0xa || idx == 0xb) && (b.cartManIntervals&references.ROM_HI_1 == references.ROM_HI_1) {
				b.bankRead[idx] = b.cartManRead
			} else if (idx == 0xe || idx == 0xf) && (b.cartManIntervals&references.ROM_HI_2 == references.ROM_HI_2) {
				b.bankRead[idx] = b.cartManRead
			}
		case UND:
			b.bankRead[idx] = b.readOpenBus
			b.bankWrite[idx] = b.writeOpenBus
		default:
			log.Fatalf("wrong memory config for bank %X: %d", idx, v)
		}
	}
	//fmt.Printf("SYSTEM MEMORY CONFIG CHANGED  %d -> %v\n", mcIdx, b.memoryConfig)
	return true
}

/*
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
		b.bankRead[0xd] = b.ramReadIO
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
*/
