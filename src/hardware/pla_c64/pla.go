package pla_c64

import (
	"github.com/markel1974/c64emu/src/common/filler"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type ReadFn func(uint16) uint8

// WriteFn is a function type that represents a write operation with a 16-bit address and an 8-bit data payload.
type WriteFn func(uint16, uint8)

// PLA represents a structure managing memory configurations, ports, and sockets for emulation purposes.
type PLA struct {
	*component.BaseComponent
	factory         references.IComponentFactory
	vic             references.IVIC
	sid             references.ISID
	cia1            references.ICIA
	cia2            references.ICIA
	cartMan         references.ICartridgeManagerC64
	roms            references.IROMLoaderC64
	ram             []byte
	bankWrite       []WriteFn
	bankRead        []ReadFn
	portWrite       []WriteFn
	portRead        []ReadFn
	memoryMap       *MemoryMap
	memoryConfigIdx int
	memoryConfig    []uint8
	ports           *Ports
	emulatorId      *EmulatorId
	basic           []byte
	kernal          []byte
	char            []byte
	color           []byte
	cfg             *config.Config
	wTriggers       *WriteTriggers
}

func NewPLAComponent(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewPLA(parent, factory, suffix)
}

// NewPLA initializes and returns a pointer to a new instance of PLA with default memory and configurations set.
func NewPLA(parent references.IComponent, factory references.IComponentFactory, suffix string) *PLA {
	mm := NewMemoryMap()
	b := &PLA{
		BaseComponent:   component.NewBaseComponent(componentId, suffix),
		factory:         factory,
		vic:             nil,
		sid:             nil,
		cia1:            nil,
		cia2:            nil,
		cartMan:         nil,
		roms:            nil,
		ram:             make([]byte, 0x10000),
		bankWrite:       make([]WriteFn, 0xf+1),
		bankRead:        make([]ReadFn, 0xf+1),
		portWrite:       make([]WriteFn, 0xf+1),
		portRead:        make([]ReadFn, 0xf+1),
		memoryMap:       mm,
		memoryConfig:    mm.Get(0),
		ports:           nil,
		emulatorId:      NewEmulatorId(),
		basic:           []byte{},
		kernal:          []byte{},
		char:            []byte{},
		color:           make([]byte, 0x0400),
		cfg:             nil,
		memoryConfigIdx: -1,
		wTriggers:       nil,
	}
	component.Register(parent, b)
	b.ports = NewPorts(b, "")
	return b
}

func (b *PLA) Setup(vic references.IVIC, sid references.ISID, cia1 references.ICIA, cia2 references.ICIA, cartMan references.ICartridgeManagerC64, roms references.IROMLoaderC64, cfg *config.Config) {
	b.vic = vic
	b.sid = sid
	b.cia1 = cia1
	b.cia2 = cia2
	b.cartMan = cartMan
	b.roms = roms

	b.cfg = cfg
	b.bankWrite[0x0] = b.ramWrite0x0000
	b.bankWrite[0x1] = b.ramWrite0x1000
	b.bankWrite[0x2] = b.ramWrite0x2000
	b.bankWrite[0x3] = b.ramWrite0x3000
	b.bankWrite[0x4] = b.ramWrite0x4000
	b.bankWrite[0x5] = b.ramWrite0x5000
	b.bankWrite[0x6] = b.ramWrite0x6000
	b.bankWrite[0x7] = b.ramWrite0x7000
	b.bankWrite[0x8] = b.ramWrite0x8000
	b.bankWrite[0x9] = b.ramWrite0x9000
	b.bankWrite[0xa] = b.ramWrite0xA000
	b.bankWrite[0xb] = b.ramWrite0xB000
	b.bankWrite[0xc] = b.ramWrite0xC000
	b.bankWrite[0xd] = b.ramWrite0xD000
	b.bankWrite[0xe] = b.ramWrite0xE000
	b.bankWrite[0xf] = b.ramWrite0xF000

	b.bankRead[0x0] = b.ramRead0x0000
	b.bankRead[0x1] = b.ramRead0x1000
	b.bankRead[0x2] = b.ramRead0x2000
	b.bankRead[0x3] = b.ramRead0x3000
	b.bankRead[0x4] = b.ramRead0x4000
	b.bankRead[0x5] = b.ramRead0x5000
	b.bankRead[0x6] = b.ramRead0x6000
	b.bankRead[0x7] = b.ramRead0x7000
	b.bankRead[0x8] = b.ramRead0x8000
	b.bankRead[0x9] = b.ramRead0x9000
	b.bankRead[0xa] = b.ramRead0xA000
	b.bankRead[0xb] = b.ramRead0xB000
	b.bankRead[0xc] = b.ramRead0xC000
	b.bankRead[0xd] = b.ramRead0xD000
	b.bankRead[0xe] = b.ramRead0xE000
	b.bankRead[0xf] = b.ramRead0xF000

	b.portWrite[0x0] = b.vic.WriteRegister
	b.portWrite[0x1] = b.vic.WriteRegister
	b.portWrite[0x2] = b.vic.WriteRegister
	b.portWrite[0x3] = b.vic.WriteRegister
	b.portWrite[0x4] = b.sid.WriteRegister
	b.portWrite[0x5] = b.sid.WriteRegister
	b.portWrite[0x6] = b.sid.WriteRegister
	b.portWrite[0x7] = b.sid.WriteRegister
	b.portWrite[0x8] = b.portWriteColor
	b.portWrite[0x9] = b.portWriteColor
	b.portWrite[0xa] = b.portWriteColor
	b.portWrite[0xb] = b.portWriteColor
	b.portWrite[0xc] = b.cia1.WriteRegister
	b.portWrite[0xd] = b.cia2.WriteRegister
	b.portWrite[0xe] = b.portWriteIO
	b.portWrite[0xf] = b.portWriteIO

	b.portRead[0x0] = b.vic.ReadRegister
	b.portRead[0x1] = b.vic.ReadRegister
	b.portRead[0x2] = b.vic.ReadRegister
	b.portRead[0x3] = b.vic.ReadRegister
	b.portRead[0x4] = b.sid.ReadRegister
	b.portRead[0x5] = b.sid.ReadRegister
	b.portRead[0x6] = b.sid.ReadRegister
	b.portRead[0x7] = b.sid.ReadRegister
	b.portRead[0x8] = b.portReadColor
	b.portRead[0x9] = b.portReadColor
	b.portRead[0xa] = b.portReadColor
	b.portRead[0xb] = b.portReadColor
	b.portRead[0xc] = b.cia1.ReadRegister
	b.portRead[0xd] = b.cia2.ReadRegister
	b.portRead[0xe] = b.portReadIO
	b.portRead[0xf] = b.portReadIO

	ri := filler.New(255, 128, 0, 0, 0, 0, 0, 0)
	ri.InitWithPattern(b.ram, uint(len(b.ram)))

	rc := filler.New(255, 128, 0, 0, 0, 0, 0, filler.InitRandomChanceHalf)
	rc.InitWithPattern(b.color, uint(len(b.color)))
	b.initRom()
}

// Reset reinitializes the state of PLA by resetting its ports and updating internal references or state.
func (b *PLA) Reset() {
	b.ports.Reset()
	b.update()
}

func (b *PLA) Emulate() {
	//
}

//func (b *PLA) AsyncReset() {
//	b.initRom()
//}

func (b *PLA) initRom() {
	b.kernal = b.roms.LoadKernal()
	b.basic = b.roms.LoadBasic()
	b.char = b.roms.LoadChar()
}

// update updates the state of the PLA object by updating ports and rebuilding the memory configuration.
func (b *PLA) update() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	//b.ports.SetTape(tape_sense, tape_write_in, tape_motor_in)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

// RebuildMemoryConfig updates the memory configuration based on the current cartridge and port settings.
func (b *PLA) RebuildMemoryConfig() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	game := uint8(1)
	exRom := uint8(1)
	if g, x, ok := b.cartMan.Config(); ok {
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

// SetMemoryConfig updates the memory configuration of the bank with the provided byte slice configuration.
func (b *PLA) SetMemoryConfig(cfg []uint8) {
	b.memoryConfig = cfg
}

// SetMemoryEntry sets the memory configuration for the PLA instance using the provided memory configuration identifier.
func (b *PLA) SetMemoryEntry(memConfig uint8) {
	b.memoryConfig = b.memoryMap.Get(memConfig)
}

// ReadBasicRom reads a byte from the BASIC ROM at the specified memory address.
func (b *PLA) ReadBasicRom(addr uint16) uint8 {
	return b.basic[addr]
}

// ReadCharRom reads a byte from the character ROM at the given address. It returns the byte located at the specified address.
func (b *PLA) ReadCharRom(addr uint16) uint8 {
	return b.char[addr]
}

// ReadKernalRom reads a byte from the Kernal ROM at the specified memory address.
func (b *PLA) ReadKernalRom(addr uint16) uint8 {
	return b.kernal[addr]
}

// ReadColor reads and returns the color value from the given address in the color bank.
func (b *PLA) ReadColor(addr uint16) uint8 {
	return b.color[addr]
}

// WriteColor writes the given data byte to the specified address in the color memory of the PLA struct.
func (b *PLA) WriteColor(addr uint16, data uint8) {
	b.color[addr] = data
}

// Read retrieves an 8-bit value from the specified memory address using the appropriate bank read function.
func (b *PLA) Read(addr uint16) uint8 {
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	return b.bankRead[bank](addr)
}

// ReadDirect retrieves a byte of data directly from the RAM at the specified address without any additional processing.
func (b *PLA) ReadDirect(addr uint16) uint8 {
	return b.ram[addr]
}

// WriteDirect writes the specified `data` byte to the given `addr` in RAM and triggers any associated write hooks.
func (b *PLA) WriteDirect(addr uint16, data uint8) {
	b.ram[addr] = data
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Exec(addr, data)
}

// Write updates the memory at the specified address with the provided data using the current memory bank configuration.
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

// ramWrite0x0000 writes data to the specified address in RAM, handling special cases for addresses 0 and 1.
func (b *PLA) ramWrite0x0000(addr uint16, data uint8) {
	if addr == 0 {
		b.ports.SetDir(data)
		b.ram[0] = b.vic.GetLastByte()
		b.update()
		return
	} else if addr == 1 {
		b.ports.SetData(data)
		b.ram[1] = b.vic.GetLastByte()
		b.update()
		return
	}
	b.ram[addr] = data
}

// SetWriteTrigger sets a write trigger at the specified address with the given callback function and returns its trigger ID.
func (b *PLA) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	if b.wTriggers == nil {
		b.wTriggers = NewWriteTriggers(len(b.ram))
	}
	return b.wTriggers.Add(addr, fn)
}

// RemoveRamTrigger removes a write trigger associated with the specified memory address and identifier from the bank.
func (b *PLA) RemoveRamTrigger(addr uint16, id int) {
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Remove(addr, id)
}

// ramWrite0x1000 writes a byte of data to the specified address in the RAM within the 0x1000 range.
func (b *PLA) ramWrite0x1000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x2000 writes a single byte of data to the RAM at the specified 16-bit address.
func (b *PLA) ramWrite0x2000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x3000 writes the given data byte to the specified address within the 0x3000 range of the RAM.
func (b *PLA) ramWrite0x3000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x4000 writes a byte of data to the specified RAM address within the range starting at 0x4000.
func (b *PLA) ramWrite0x4000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x5000 writes a single byte of data to the specified address in the RAM within the 0x5000 range.
func (b *PLA) ramWrite0x5000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x6000 writes a byte of data to the specified address in the RAM, starting at the 0x6000 range.
func (b *PLA) ramWrite0x6000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x7000 writes an 8-bit data value to the specified 16-bit address in the RAM within the addressable 0x7000 range.
func (b *PLA) ramWrite0x7000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x8000 writes a single byte of data to the specified address in the RAM starting at 0x8000.
func (b *PLA) ramWrite0x8000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0x9000 writes a byte of data to the RAM at the specified address in the 0x9000 range.
func (b *PLA) ramWrite0x9000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0xA000 writes a byte of data to the RAM at the specified address in the 0xA000 range.
func (b *PLA) ramWrite0xA000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0xB000 writes a byte of data to the specified address in the RAM located at 0xB000.
func (b *PLA) ramWrite0xB000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0xC000 writes a byte of data to the specified address in the RAM starting at 0xC000.
func (b *PLA) ramWrite0xC000(addr uint16, data uint8) {
	b.ram[addr] = data
}

// ramWrite0xD000 writes a byte of data to RAM or an I/O port based on the memory configuration for bank 0xD.
// If the memory is configured as I/O, it determines the port and invokes the relevant port write handler.
// Otherwise, it directly writes the data to the specified RAM address.
func (b *PLA) ramWrite0xD000(addr uint16, data uint8) {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		b.portWrite[p](addr, data)
		return
	}
	b.ram[addr] = data
}

// ramWrite0xE000 writes a byte of data to the specified RAM address within the 0xE000 range in the PLA structure.
func (b *PLA) ramWrite0xE000(addr uint16, data uint8) {
	b.ram[addr] = data
	return
}

// ramWrite0xF000 writes a byte of data to the RAM at the specified 0xF000-based address.
// It directly manipulates the memory location in the RAM slice of the PLA object.
func (b *PLA) ramWrite0xF000(addr uint16, data uint8) {
	b.ram[addr] = data
	return
}

// ramRead0x0000 reads a byte of data from the specified address within the RAM or ports based on the given address input.
func (b *PLA) ramRead0x0000(addr uint16) uint8 {
	if addr == 0 {
		return b.ports.GetDirection()
	} else if addr == 1 {
		return b.ports.GetDataRead()
	}
	return b.ram[addr]
}

// ramRead0x1000 reads a byte from the RAM at the given address within the 0x1000 address space.
func (b *PLA) ramRead0x1000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x2000 reads a byte from the RAM at the specified address within the range 0x2000.
func (b *PLA) ramRead0x2000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x3000 reads a byte from the RAM at the specified address within the 0x3000 memory range.
func (b *PLA) ramRead0x3000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x4000 reads a byte from the RAM at the specified 16-bit address within the 0x4000 memory range.
func (b *PLA) ramRead0x4000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x5000 reads an 8-bit value from the bank's RAM at the specified 16-bit address within range 0x5000.
func (b *PLA) ramRead0x5000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x6000 reads and returns a byte from the RAM at the specified address within the 0x6000 range.
func (b *PLA) ramRead0x6000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x7000 reads a byte from the RAM at the specified address within the 0x7000 memory range.
func (b *PLA) ramRead0x7000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0x8000 reads a byte from RAM or cartridge memory based on the memory configuration at address 0x8000.
func (b *PLA) ramRead0x8000(addr uint16) uint8 {
	const bank = 0x8
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartMan.Read(references.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ram[addr]
}

// ramRead0x9000 reads a byte from the memory bank at address 0x9000 based on the current memory configuration and cartridge mode.
func (b *PLA) ramRead0x9000(addr uint16) uint8 {
	const bank = 0x9
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartMan.Read(references.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ram[addr]
}

// ramRead0xA000 reads a byte from the specified address in the 0xA000 memory bank based on the active memory configuration.
// If the bank is set to "ROH", it attempts to read from the cartridge's high ROM segment.
// If the bank is set to "BAS", it retrieves data from the BASIC ROM.
// Otherwise, it defaults to reading from RAM at the specified address.
func (b *PLA) ramRead0xA000(addr uint16) uint8 {
	const bank = 0xa
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(references.ROM_HI_1, addr); ok {
			return v
		}

	} else if b.memoryConfig[bank] == BAS {
		return b.basic[addr&0x1fff]
	}
	return b.ram[addr]
}

// ramRead0xB000 reads a byte from memory at address 0xB000 based on the current memory configuration.
// It prioritizes cartridge ROM, BASIC memory, or RAM depending on the configuration of the 0xB bank.
// addr is the 16-bit memory address to be read.
// Returns the byte value read from the appropriate memory source.
func (b *PLA) ramRead0xB000(addr uint16) uint8 {
	const bank = 0xb
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(references.ROM_HI_1, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == BAS {
		return b.basic[addr&0x1fff]
	}
	return b.ram[addr]
}

func (b *PLA) ramRead0xC000(addr uint16) uint8 {
	return b.ram[addr]
}

// ramRead0xD000 reads a byte from the memory or I/O port based on the address and current memory configuration for bank 0xD.
// If the configuration is I_O, it reads from an I/O port determined by bits 8-11 of the address.
// If the configuration is CHA, it reads from the character memory at the lower 12 bits of the address.
// Otherwise, it reads directly from the RAM at the specified address.
func (b *PLA) ramRead0xD000(addr uint16) uint8 {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		return b.portRead[p](addr)
	} else if b.memoryConfig[bank] == CHA {
		return b.char[addr&0x0fff]
	}
	return b.ram[addr]
}

// ramRead0xE000 reads a byte from the RAM or ROM mapped to the address 0xE000 based on the current memory configuration.
func (b *PLA) ramRead0xE000(addr uint16) uint8 {
	const bank = 0xe
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(references.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernal[addr&0x1fff]
	}
	return b.ram[addr]
}

// ramRead0xF000 reads a byte from address 0xF000 based on the current memory configuration, supporting ROM, Kernal, or RAM.
func (b *PLA) ramRead0xF000(addr uint16) uint8 {
	const bank = 0xf
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(references.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernal[addr&0x1fff]
	}
	return b.ram[addr]
}

// portWriteColor updates the color buffer at the specified address with the given 4-bit color data.
func (b *PLA) portWriteColor(addr uint16, data uint8) {
	b.color[addr&0x03ff] = data & 0x0f
}

// portReadColor reads the color data at the specified address, combining it with the high nibble of the last VIC byte.
func (b *PLA) portReadColor(addr uint16) uint8 {
	return (b.color[addr&0x03ff] & 0x0f) | (b.vic.GetLastByte() & 0xf0)
}

// portReadIO handles IO read operations for the given address and returns the corresponding byte value.
func (b *PLA) portReadIO(addr uint16) uint8 {
	if v, ok := b.cartMan.IORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vic.GetLastByte()
	}
	return b.emulatorId.Read(addr)
}

// portWriteIO handles writing a byte of data to a specified IO port address, delegating to the cartMan IOWrite method.
func (b *PLA) portWriteIO(addr uint16, data uint8) {
	if ok := b.cartMan.IOWrite(addr, data); ok {
		return
	}
}
