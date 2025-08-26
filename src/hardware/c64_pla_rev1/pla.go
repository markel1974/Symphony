package c64_pla_rev1

import (
	"log"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

// ReadFn defines a function type that takes a 16-bit unsigned integer as input and returns an 8-bit unsigned integer.
type ReadFn func(uint16) uint8

// WriteFn defines a function type that takes a 16-bit address and an 8-bit data value for write operations.
type WriteFn func(uint16, uint8)

// PLA represents the Programmable Logic Array (PLA) component responsible for memory and I/O address mapping logic.
type PLA struct {
	*component.BaseComponent
	reflect          *PLAReflect
	ramRead          ReadFn
	ramWrite         WriteFn
	cartManRead      ReadFn
	cartManIORead    func(uint16) (uint8, bool)
	cartManIOWrite   func(uint16, uint8) bool
	cartManConfig    func() (uint8, uint8, bool)
	cartManIntervals uint8
	vaSignals        func() uint8
	memoryConfig     []uint8
	bankWrite        []WriteFn
	bankRead         []ReadFn
	u15Write         []WriteFn // represent U15, 74LS138 (mux)
	u15Read          []ReadFn  // represent U15, 74LS138 (mux)
	bankSwitcher     *BankSwitcher
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
	//mm := NewBankSwitcher()
	b := &PLA{
		BaseComponent:    component.NewBaseComponent(),
		vaSignals:        nil,
		ramRead:          nil,
		ramWrite:         nil,
		cartManRead:      nil,
		cartManIORead:    nil,
		cartManIOWrite:   nil,
		bankWrite:        make([]WriteFn, bankSize),
		bankRead:         make([]ReadFn, bankSize),
		u15Write:         make([]WriteFn, bankSize),
		u15Read:          make([]ReadFn, bankSize),
		bankSwitcher:     nil,
		ports:            nil,
		emulatorId:       NewEmulatorId(),
		basicRead:        nil,
		kernalRead:       nil,
		charRead:         nil,
		colorRead:        nil,
		cfg:              nil,
		wTriggers:        nil, //
		label:            label,
		cartManIntervals: 0,
	}
	//b.BaseComponent.Register(factory, parent, Identifier(), instance, b, references.IdIC64Pla(b, label, instance))
	b.reflect = NewPLAReflect(b, factory, parent, Identifier(), instance, references.IdIC64Pla(b, label, instance))
	return b
}

// Setup initializes the PLA instance by configuring its settings and creating the port object.
func (b *PLA) Setup() error {
	b.cfg = b.GetFactory().GetConfig()
	b.bankSwitcher = NewBankSwitcher(b.GetFactory(), b, b.label, 0)
	b.ports = NewPorts(b.GetFactory(), b, b.label, 0)
	return nil
}

//U15 74LS138
//Pin 12	[cs12 $D000–$D3FF] VIC-II
//Pin 04	[ cs4 $D400–$D7FF] SID
//Pin 10	[cs10 $D800–$DBFF] Color RAM
//Pin 14	[cs14 $DC00–$DCFF] CIA1
//Pin 15	[cs15 $DD00–$DDFF] CIA2
//Pin 09	[ cs9 $DF00–$DFFF] I/O 2 (cartridge)
//Pin 11	[cs11 $DE00–$DEFF] I/O 1 (cartridge)

// Bind initializes and binds PLA components such as RAM, ROM, cartridge manager, chip-select elements, and signal handlers.
func (b *PLA) Bind(_ references.IC64PlaSocket, vaSignals references.IC64PlaVASignals, cartMan references.IC64CartridgeManager, ram references.IC64Ram, roms references.IC64Roms, cs12 references.IC64PlaChipSelect, cs4 references.IC64PlaChipSelect, cs14 references.IC64PlaChipSelect, cs15 references.IC64PlaChipSelect, cs10 references.IC64PlaChipSelect) error {
	b.vaSignals = vaSignals.GetVASignal

	b.wTriggers = NewWriteTriggers(ram.Size())
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

	b.bankWrite[0x0] = b.WriteRam0x0000
	b.bankRead[0x0] = b.ReadRam0x0000

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
	b.u15Write[0xe] = b.WritePortIO
	b.u15Write[0xf] = b.WritePortIO

	b.u15Read[0x0] = cs12.ReadRegister
	b.u15Read[0x1] = cs12.ReadRegister
	b.u15Read[0x2] = cs12.ReadRegister
	b.u15Read[0x3] = cs12.ReadRegister
	b.u15Read[0x4] = cs4.ReadRegister
	b.u15Read[0x5] = cs4.ReadRegister
	b.u15Read[0x6] = cs4.ReadRegister
	b.u15Read[0x7] = cs4.ReadRegister
	b.u15Read[0x8] = b.ReadPortColor
	b.u15Read[0x9] = b.ReadPortColor
	b.u15Read[0xa] = b.ReadPortColor
	b.u15Read[0xb] = b.ReadPortColor
	b.u15Read[0xc] = cs14.ReadRegister
	b.u15Read[0xd] = cs15.ReadRegister
	b.u15Read[0xe] = b.ReadPortIO
	b.u15Read[0xf] = b.ReadPortIO

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

// Reset resets the internal state of the PLA, reinitializing ports, and triggering an Update of the memory configuration.
func (b *PLA) Reset() {
	b.ports.Reset()
	b.Update()
}

// Emulate performs the main emulation step for the PLA, updating its state and processing associated logic.
func (b *PLA) Emulate() {
}

// EmulationRequired determines whether emulation is required for the PLA component. Returns false by default.
func (b *PLA) EmulationRequired() bool {
	return false
}

// RebuildMemoryConfig recalculates and applies the current memory configuration based on cartridge and port settings.
func (b *PLA) RebuildMemoryConfig() {
	spec := references.C64CartridgeSpecOff
	if game, exRom, ok := b.cartManConfig(); ok {
		spec = references.GetCartridgeSpec(game, exRom)
	}
	b.cartManIntervals = spec.Intervals
	dir, data := b.ports.Config()
	mcIdx := ((^dir | data) & 0x7) | (spec.ExRom << 3) | (spec.Game << 4)
	b.BankSwitch(int(mcIdx))
}

// Read retrieves the value from the memory mapped by the given address.
func (b *PLA) Read(addr uint16) uint8 {
	bank := addr >> 12
	return b.bankRead[bank](addr)
}

// Write writes a byte of data to a specific memory address based on the bank configuration and triggers write callbacks.
func (b *PLA) Write(addr uint16, data uint8) {
	bank := addr >> 12
	b.bankWrite[bank](addr, data)
}

// WriteExt writes a byte to the specified address using the provided memory configuration, temporarily switching configurations.
func (b *PLA) WriteExt(memConfig int, addr uint16, data uint8) {
	var prevMemConfig = -1
	if memConfig >= 0 {
		prevMemConfig = b.bankSwitcher.ReadConfig()
		b.BankSwitch(memConfig)
	}
	b.Write(addr, data)
	if prevMemConfig >= 0 {
		b.BankSwitch(prevMemConfig)
	}
}

// ReadExt allows reading a specific memory address using a temporary memory configuration.
func (b *PLA) ReadExt(memConfig int, addr uint16) uint8 {
	var prevMemConfig = -1
	if memConfig >= 0 {
		prevMemConfig = b.bankSwitcher.ReadConfig()
		b.BankSwitch(memConfig)
	}
	data := b.Read(addr)
	if prevMemConfig >= 0 {
		b.BankSwitch(prevMemConfig)
	}
	return data
}

// SetWriteTrigger registers a write-trigger at a specified address with a callback and returns the assigned trigger Id.
func (b *PLA) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	rebuild := b.wTriggers.Len() == 0
	w := b.wTriggers.Add(addr, fn)
	if rebuild {
		b.applyMemoryConfig(b.memoryConfig)
	}
	return w
}

// RemoveRamTrigger removes a write-trigger by its unique Id for the specified memory address in the WriteTriggers collection.
func (b *PLA) RemoveRamTrigger(addr uint16, id int) {
	b.wTriggers.Remove(addr, id)
	rebuild := b.wTriggers.Len() == 0
	if rebuild {
		b.applyMemoryConfig(b.memoryConfig)
	}
}

// Update updates the state of the PLA by rebuilding the memory configuration and updating the port settings.
func (b *PLA) Update() {
	//b.ports.WriteTape(tapeSense, tapeWriteIn, tapeMotorIn)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

// WriteRam0x0000 handles writing data to memory address 0x0000, updating port direction and data where applicable.
func (b *PLA) WriteRam0x0000(addr uint16, data uint8) {
	if addr == 0 {
		b.ports.WriteDir(data)
		//b.ramWrite(0, b.vaSignals)
		b.Update()
		return
	} else if addr == 1 {
		b.ports.WriteData(data)
		//b.ramWrite(1, b.vaSignals())
		b.Update()
		return
	}
	b.ramWrite(addr, data)
}

// WriteRamIO writes a byte of data to the specific address by determining the target U15 write function based on the address.
func (b *PLA) WriteRamIO(addr uint16, data uint8) {
	p := (addr >> 8) & 0x0f
	b.u15Write[p](addr, data)
	return
}

// ReadRamIO returns the value of the RAM read at the specified address, using the appropriate U15 read function.
func (b *PLA) ReadRamIO(addr uint16) uint8 {
	p := (addr >> 8) & 0x0f
	return b.u15Read[p](addr)
}

// ReadRam0x0000 reads data from address 0x0000 and processes special cases for addresses 0 and 1.
func (b *PLA) ReadRam0x0000(addr uint16) uint8 {
	if addr == 0 {
		return b.ports.ReadDirection()
	} else if addr == 1 {
		return b.ports.ReadDataRead()
	}
	return b.ramRead(addr)
}

// ReadPortColor reads color data from the Color RAM and combines it with data from the VIC latch for a full byte response.
func (b *PLA) ReadPortColor(addr uint16) uint8 {
	p1 := b.colorRead(addr) & 0x0f // enables the physical Color RAM chip. This chip receives the address and puts the 4 color bits it has stored on data bus lines D0-D3 (lower half).
	p2 := b.vaSignals() & 0xf0     // signals to VIC that Color RAM is being read. VIC responds by putting the last 4 bits of its internal latch (lastByte) on data bus lines D4-D7 (upper half)
	return p1 | p2
}

// ReadPortIO reads a byte from the specified IO port address.
func (b *PLA) ReadPortIO(addr uint16) uint8 {
	if v, ok := b.cartManIORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vaSignals()
	}
	return b.emulatorId.Read(addr)
}

// WritePortIO writes data to the specified IO port address using the cartridge manager's IO write functionality.
func (b *PLA) WritePortIO(addr uint16, data uint8) {
	_ = b.cartManIOWrite(addr, data)
}

// WriteOpenBus writes to an "open bus" state where no specific memory or component is targeted.
func (b *PLA) WriteOpenBus(_ uint16, _ uint8) {
}

// ReadOpenBus reads the current value of the open bus at the specified address and returns its signals as uint8.
func (b *PLA) ReadOpenBus(_ uint16) uint8 {
	return b.vaSignals()
}

// BankSwitch switches the memory bank to the specified bankIndex and applies the corresponding memory configuration.
// It returns true if the switch is successful, otherwise false.
func (b *PLA) BankSwitch(bankIndex int) bool {
	memoryConfig, ok := b.bankSwitcher.WriteConfig(bankIndex)
	if !ok {
		return false
	}
	b.memoryConfig = memoryConfig
	b.applyMemoryConfig(b.memoryConfig)
	//fmt.Printf("SYSTEM MEMORY CONFIG CHANGED  %d -> %v\n", mcIdx, b.memoryConfig)
	return true
}

// applyMemoryConfig configures memory banks based on the provided memory configuration, setting read/write functions accordingly.
func (b *PLA) applyMemoryConfig(memoryConfig []uint8) {
	ramWrite := b.ramWrite
	ramWrite0x0000 := b.WriteRam0x0000
	if b.wTriggers.Len() > 0 {
		ramWrite = func(addr uint16, data uint8) {
			b.ramWrite(addr, data)
			b.wTriggers.Exec(addr, data)
		}
		ramWrite0x0000 = func(addr uint16, data uint8) {
			b.WriteRam0x0000(addr, data)
			b.wTriggers.Exec(addr, data)
		}
	}

	for idx, v := range memoryConfig {
		if idx == 0 {
			b.bankRead[0] = b.ReadRam0x0000
			b.bankWrite[0] = ramWrite0x0000
		} else {
			b.bankRead[idx] = b.ramRead
			b.bankWrite[idx] = ramWrite
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
			b.bankRead[idx] = b.ReadRamIO
			b.bankWrite[idx] = b.WriteRamIO
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
			b.bankRead[idx] = b.ReadOpenBus
			b.bankWrite[idx] = b.WriteOpenBus
		default:
			log.Fatalf("wrong memory config for bank %X: %d", idx, v)
		}
	}
}
