package banks

import (
	"github.com/markel1974/c64emu/src/board/banks/roms"
	"github.com/markel1974/c64emu/src/board/cartridges"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cia"
	"github.com/markel1974/c64emu/src/board/sid"
	"github.com/markel1974/c64emu/src/board/vic"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/filler"
)

const (
	BasicRomFile  = "Basic.rom"
	KernalRomFile = "Kernal.rom"
	CharRomFile   = "Char.rom"
)

type ReadFn func(uint16) uint8
type WriteFn func(uint16, uint8)

type Banks struct {
	vic             *vic.MOS6569
	sid             *sid.MOS6581
	cia1            *cia.MOS6526_1
	cia2            *cia.MOS6526_2
	cartMan         *cartridges.Manager
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
	prefs           *config.Config
	wTriggers       *WriteTriggers
}

func NewBanks() *Banks {
	mm := NewMemoryMap()
	b := &Banks{
		vic:             nil,
		sid:             nil,
		cia1:            nil,
		cia2:            nil,
		cartMan:         nil,
		ram:             make([]byte, 0x10000),
		bankWrite:       make([]WriteFn, 0xf+1),
		bankRead:        make([]ReadFn, 0xf+1),
		portWrite:       make([]WriteFn, 0xf+1),
		portRead:        make([]ReadFn, 0xf+1),
		memoryMap:       mm,
		memoryConfig:    mm.Get(0),
		ports:           NewPorts(),
		emulatorId:      NewEmulatorId(),
		basic:           make([]byte, roms.BASIC_ROM_SIZE),
		kernal:          make([]byte, roms.KERNAL_ROM_SIZE),
		char:            make([]byte, roms.CHAR_ROM_SIZE),
		color:           make([]byte, 0x0400),
		prefs:           nil,
		memoryConfigIdx: -1,
		wTriggers:       nil,
	}
	return b
}

func (b *Banks) Setup(vic *vic.MOS6569, sid *sid.MOS6581, cia1 *cia.MOS6526_1, cia2 *cia.MOS6526_2, cartMan *cartridges.Manager, prefs *config.Config) {
	b.vic = vic
	b.sid = sid
	b.cia1 = cia1
	b.cia2 = cia2
	b.cartMan = cartMan

	b.prefs = prefs
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

func (b *Banks) Reset() {
	b.ports.Reset()
	b.update()
}

//func (b *Banks) AsyncReset() {
//	b.initRom()
//}

func (b *Banks) initRom() {
	romLoader := NewRomLoader()

	//b.kernal = romLoader.Load(roms.BuiltinKernalRom, KernalRomFile)
	//roms.PatchKernalRom(&b.kernal)

	b.kernal = romLoader.Load(roms.BuiltinKernalJiffyRom, KernalRomFile)
	//b.kernal = b.romLoader.Load(builtin_kernal_fast_rom, KernalRomFile)
	//b.kernal = b.romLoader.Load(builtin_kernal_rom, KernalRomFile)
	b.basic = romLoader.Load(roms.BuiltinBasicRom, BasicRomFile)
	b.char = romLoader.Load(roms.BuiltinCharRom, CharRomFile)
}

func (b *Banks) update() {
	//https://sta.c64.org/cbm64mem.html
	//https://codebase64.org/doku.php?id=base:memory_management
	//b.ports.SetTape(tape_sense, tape_write_in, tape_motor_in)
	b.ports.Update()
	b.RebuildMemoryConfig()
}

func (b *Banks) RebuildMemoryConfig() {
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

func (b *Banks) GetMemoryConfig() []uint8 {
	return b.memoryConfig
}

func (b *Banks) SetMemoryConfig(cfg []uint8) {
	b.memoryConfig = cfg
}

func (b *Banks) SetMemoryEntry(memConfig uint8) {
	b.memoryConfig = b.memoryMap.Get(memConfig)
}

func (b *Banks) ReadBasicRom(addr uint16) uint8 {
	return b.basic[addr]
}

func (b *Banks) ReadCharRom(addr uint16) uint8 {
	return b.char[addr]
}

func (b *Banks) ReadKernalRom(addr uint16) uint8 {
	return b.kernal[addr]
}

func (b *Banks) ReadColor(addr uint16) uint8 {
	return b.color[addr]
}

func (b *Banks) WriteColor(addr uint16, data uint8) {
	b.color[addr] = data
}

func (b *Banks) Read(addr uint16) uint8 {
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	return b.bankRead[bank](addr)
}

func (b *Banks) ReadDirect(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) WriteDirect(addr uint16, data uint8) {
	b.ram[addr] = data
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Exec(addr, data)
}

func (b *Banks) Write(addr uint16, data uint8) {
	//sta.c64.org/cbm64mem.html
	//https://www.c64-wiki.com/wiki/Memory_Map#Configurations
	bank := addr >> 12
	b.bankWrite[bank](addr, data)
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Exec(addr, data)
}

func (b *Banks) ramWrite0x0000(addr uint16, data uint8) {
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

func (b *Banks) SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int {
	if b.wTriggers == nil {
		b.wTriggers = NewWriteTriggers(len(b.ram))
	}
	return b.wTriggers.Add(addr, fn)
}

func (b *Banks) RemoveRamTrigger(addr uint16, id int) {
	if b.wTriggers == nil {
		return
	}
	b.wTriggers.Remove(addr, id)
}

func (b *Banks) ramWrite0x1000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x2000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x3000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x4000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x5000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x6000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x7000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x8000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0x9000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0xA000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0xB000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0xC000(addr uint16, data uint8) {
	b.ram[addr] = data
}

func (b *Banks) ramWrite0xD000(addr uint16, data uint8) {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		b.portWrite[p](addr, data)
		return
	}
	b.ram[addr] = data
}

func (b *Banks) ramWrite0xE000(addr uint16, data uint8) {
	b.ram[addr] = data
	return
}

func (b *Banks) ramWrite0xF000(addr uint16, data uint8) {
	b.ram[addr] = data
	return
}

func (b *Banks) ramRead0x0000(addr uint16) uint8 {
	if addr == 0 {
		return b.ports.GetDirection()
	} else if addr == 1 {
		return b.ports.GetDataRead()
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0x1000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x2000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x3000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x4000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x5000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x6000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x7000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0x8000(addr uint16) uint8 {
	const bank = 0x8
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartMan.Read(icartridge.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0x9000(addr uint16) uint8 {
	const bank = 0x9
	if b.memoryConfig[bank] == ROL {
		if v, ok := b.cartMan.Read(icartridge.ROM_LO, addr); ok {
			return v
		}
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0xA000(addr uint16) uint8 {
	const bank = 0xa
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(icartridge.ROM_HI_1, addr); ok {
			return v
		}

	} else if b.memoryConfig[bank] == BAS {
		return b.basic[addr&0x1fff]
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0xB000(addr uint16) uint8 {
	const bank = 0xb
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(icartridge.ROM_HI_1, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == BAS {
		return b.basic[addr&0x1fff]
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0xC000(addr uint16) uint8 {
	return b.ram[addr]
}

func (b *Banks) ramRead0xD000(addr uint16) uint8 {
	const bank = 0xd
	if b.memoryConfig[bank] == I_O {
		p := (addr >> 8) & 0x0f
		return b.portRead[p](addr)
	} else if b.memoryConfig[bank] == CHA {
		return b.char[addr&0x0fff]
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0xE000(addr uint16) uint8 {
	const bank = 0xe
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(icartridge.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernal[addr&0x1fff]
	}
	return b.ram[addr]
}

func (b *Banks) ramRead0xF000(addr uint16) uint8 {
	const bank = 0xf
	if b.memoryConfig[bank] == ROH {
		if v, ok := b.cartMan.Read(icartridge.ROM_HI_2, addr); ok {
			return v
		}
	} else if b.memoryConfig[bank] == KER {
		return b.kernal[addr&0x1fff]
	}
	return b.ram[addr]
}

func (b *Banks) portWriteColor(addr uint16, data uint8) {
	b.color[addr&0x03ff] = data & 0x0f
}

func (b *Banks) portReadColor(addr uint16) uint8 {
	return (b.color[addr&0x03ff] & 0x0f) | (b.vic.GetLastByte() & 0xf0)
}

func (b *Banks) portReadIO(addr uint16) uint8 {
	if v, ok := b.cartMan.IORead(addr); ok {
		return v
	}
	if addr < 0xdfa0 {
		return b.vic.GetLastByte()
	}
	return b.emulatorId.Read(addr)
}

func (b *Banks) portWriteIO(addr uint16, data uint8) {
	if ok := b.cartMan.IOWrite(addr, data); ok {
		return
	}
}

func (b *Banks) getBasicText() (uint16, uint16) {
	start := uint16(b.ram[0x2b]) | (uint16(b.ram[0x2c]) << 8)
	end := uint16(b.ram[0x2d]) | (uint16(b.ram[0x2e]) << 8)
	return start, end
}

func (b *Banks) setBasicText(start uint16, end uint16) {
	s1 := uint8(start) & 0xff
	b.ram[0xac] = s1
	b.ram[0x2b] = s1

	s2 := uint8(start >> 8)
	b.ram[0xad] = s2
	b.ram[0x2c] = s2

	e1 := uint8(end & 0xff)
	b.ram[0xae] = e1
	b.ram[0x31] = e1
	b.ram[0x2f] = e1
	b.ram[0x2d] = e1

	e2 := uint8(end >> 8)
	b.ram[0xaf] = e2
	b.ram[0x32] = e2
	b.ram[0x30] = e2
	b.ram[0x2e] = e2
}

func (b *Banks) Inject(startAddr uint16, data []byte) error {
	size := uint16(len(data))
	start, _ := b.getBasicText()
	const autostartBasicLoad = false
	// load to basic start if requested
	if autostartBasicLoad {
		startAddr = start
	}
	// store data in emu memory
	for i := uint16(0); i < size; i++ {
		b.ram[startAddr+i] = data[i]
	}
	// simulate a basic load
	end := startAddr + size
	b.setBasicText(start, end)
	return nil
}

/*
// used by autostart to locate and "read" kernal output on the current screen
// this function should return whatever the kernal currently uses, regardless
// what is currently visible/active in the UI
// static CHECKYESNO check2(

func (b *Banks) GetCursorParameter(lineOffset int) {
	//uint16_t *screen_addr, uint8_t *cursor_column, uint8_t *line_length, int *blinking
	// CAUTION: this function can be called at any time when the emulation (KERNAL)
	// is in the middle of a screen update. we must make sure that all
	// values are being looked up in an "atomic" way so we don't use a low-
	// and high- byte from before and after an update, leading to invalid values

	// Physical Screen Line Length
	const lineLength = 40
	screenBase := (int(b.ram[0xd1]) + (int(b.ram[0xd2]) * 256)) & ^0x3ff // the upper bits will not change
	//blinking := b.ram[0xcc] == 0 //? 0 : 1;
	// Current Screen Line Address
	screenAddr := screenBase + (int(b.ram[0xd6]) * lineLength)
	// Cursor Column on Current Line
	cursorColumn := int(b.ram[0xd3])
	for cursorColumn >= lineLength {
		cursorColumn -= lineLength
		screenAddr += lineLength
	}
	addr := screenAddr
	addr += lineLength * lineOffset
	for x := 0; x < lineLength; x++ {
		a := uint16(addr+x) & 0xffff
		v := b.ram[a]
		if v != 32 && v != 0 && v != 255 && v != 160 {
			fmt.Println(v)
		}
	}
	//fmt.Println(blinking, screenAddr)
}


*/
