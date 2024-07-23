package easyflash

//https://skoe.de/easyflash/files/devdocs/EasyFlash-ProgRef.pdf

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/easyflash/flash"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/filler"

	//"github.com/markel1974/c64emu/src/board/ram"
	"github.com/markel1974/c64emu/src/board/snapshot"
	"io"
	"os"
)

/*
static io_source_t easyflash_io1_device = {
	CartridgeNameEasyFlash, // name of the device
	IO_DETACH_CART,           // use cartridge ID to detach the device when involved in a read-collision
	IO_DETACH_NO_RESOURCE,    // does not use a resource for detach
	0xde00, 0xdeff, 0x03,     // range for the device, regs:$de00-$de03, mirrors:$de04-$deff
	0,                        // read is never valid, regs are write only
	easyflash_io1_store,      // store function
	NULL,                     // NO poke function < write without side effects (used by monitor)
	NULL,                     // NO read function
	easyflash_io1_peek,       // peek function < read without side effects (used by monitor)
	easyflash_io1_peek,       // device state information dump function  !< print detailed state for this i/o device (used by monitor)
	CARTRIDGE_EASYFLASH,      // cartridge ID
	IO_PRIO_NORMAL,           // normal priority, device read needs to be checked for collisions
	0,                        // insertion order, gets filled in by the registration function
	IO_MIRROR_NONE            // NO mirroring
}
*/

/*
The following structure is used to register the I/O address range used by a certain device/chip/cartridge.
 *
 * The address_mask determines if the defined address range is mirrored across a bigger range, the mask tells
 * the I/O read/write handler how many bits of the address being accessed are valid for an I/O read or write.
 *
 * Some examples:
 *
 * start_address | end_address | mask | primary register address | address mirrors
 * -------------------------------------------------------------------------------
 * $de00         | $de0f       | $00  | $de00                    | $de01-$de0f (mirrors of 0xde00)
 * $de00         | $de03       | $01  | $de00-$de01              | $de02-$de03 ($de02 mirrors $de00 and $de03 mirrors $de01)
 * $de00         | $deff       | $0f  | $de00-$de0f              | $de10-$deff (15 blocks of 16 bytes mirroring $de00-$de0f)
 * $de00         | $deff       | $7f  | $de00-$de7f              | $de80-$deff (1 block of 128 bytes mirroring $de00-$de7f)
 * $de80         | $deff       | $7f  | $de80-$deff              | no mirrors
 * $de00         | $de0f       | $0f  | $de00-$de0f              | no mirrors
 * $de00         | $deff       | $ff  | $de00-$deff              | no mirrors

*/
/*
static io_source_t easyflash_io2_device = {
	CartridgeNameEasyFlash, // name of the device
	IO_DETACH_CART,           // use cartridge ID to detach the device when involved in a read-collision
	IO_DETACH_NO_RESOURCE,    // does not use a resource for detach
	0xdf00, 0xdfff, 0xff,     // range for the device, regs:$df00-$dfff
	1,                        // read is always valid
	easyflash_io2_store,      // store function
	NULL,                     // NO poke function
	easyflash_io2_read,       // read function
	easyflash_io2_read,       // peek function, same implementation
	NULL,                     // device state information dump function
	CARTRIDGE_EASYFLASH,      // cartridge ID
	IO_PRIO_NORMAL,           // normal priority, device read needs to be checked for collisions
	0,                        // insertion order, gets filled in by the registration function
	IO_MIRROR_NONE            // NO mirroring
};
*/

/*
0xde00, 0xdeff, 0x03,     // range for the device, regs:$de00-$de03, mirrors:$de04-$deff
0,                        // read is never valid, regs are write only
easyflash_io1_store,      // store function
NULL,                     // NO poke function
NULL,                     // NO read function
easyflash_io1_peek,       // peek function
easyflash_io1_peek,       // device state information dump function

0xdf00, 0xdfff, 0xff,     // range for the device, regs:$df00-$dfff
1,                        // read is always valid
easyflash_io2_store,      // store function
NULL,                     // NO poke function
easyflash_io2_read,       // read function
easyflash_io2_read,       // peek function, same implementation
*/

type CartridgeEasyFlash struct {
	board           icartridge.IExpansion
	intervalLo      icartridge.RomInterval
	intervalHi      icartridge.RomInterval
	memoryConfigIdx int
	id              string
	game            uint8
	exRom           uint8
	stateLow        *flash.Flash040 /* the 29F040B statemachine */
	stateHigh       *flash.Flash040 /* the 29F040B statemachine */
	jumper          int             /* the jumper */
	crtWrite        int             /* writing back to crt enabled */
	crtOptimize     int             /* optimizing crt enabled */
	register00      uint8           /* backup of the registers */
	register02      uint8           /* backup of the registers */
	ram             []uint8         /* extra RAM */
	filename        string          /* filename when attached */
	filetype        loader.Filetype
	led             bool
}

func New(game uint8, exRom uint8, lo icartridge.RomInterval, hi icartridge.RomInterval) *CartridgeEasyFlash {
	return &CartridgeEasyFlash{
		game:            game,
		exRom:           exRom,
		intervalHi:      hi,
		intervalLo:      lo,
		stateLow:        nil,
		stateHigh:       nil,
		filetype:        0,
		jumper:          0,
		led:             false,
		ram:             make([]byte, CartRamSize),
		memoryConfigIdx: -1,
	}
}

func (c *CartridgeEasyFlash) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	var rawCart []byte
	c.board = board
	c.id = ldr.GetId()
	rp := filler.New(255, 2, 1, 0x100, 255, 0, 0, 0)
	rp.InitWithPattern(c.ram, CartRamSize)
	c.filename = ldr.Name
	c.filetype = 0
	var err error
	if ldr.GetMode() == loader.FiletypeCrt {
		c.filetype = loader.FiletypeCrt
		if rawCart, err = c.crtAttach(ldr); err != nil {
			return err
		}
	} else {
		c.filetype = loader.FiletypeBin
		if rawCart, err = c.binAttach(ldr); err != nil {
			return err
		}
	}
	c.register00 = 0
	c.controlUpdate(0, false)
	//c.controlUpdate(0, false)
	c.initializeFlash(rawCart)
	return nil
}

func (c *CartridgeEasyFlash) GetId() string {
	return c.id
}

func (c *CartridgeEasyFlash) initializeFlash(rawCart []byte) {
	low := make([]byte, 0x80000)
	high := make([]byte, 0x80000)
	// split interleaved low and high banks
	for i := uint(0); i < NBanks; i++ {
		const size = 0x2000
		start := i * size
		p1 := i * (size * 2)
		p2 := p1 + size
		copy(low[start:start+size], rawCart[p1:p1+size])
		copy(high[start:start+size], rawCart[p2:p2+size])
	}
	c.stateLow = flash.NewFlash040(c.board, flash.KindB, low)
	c.stateHigh = flash.NewFlash040(c.board, flash.KindB, high)

	eApiCheck := high[0x1800 : 0x1800+4]
	if bytes.Compare(eApiCheck, []byte("eapi")) == 0 {
		eApi := make([]byte, 17)
		for k := 0; k < 16; k++ {
			eApi[k] = c.stateHigh.Peek(uint32(0x1804 + k))
		}
		_ = c.stateHigh.StoreInterval(0x1800, _eApiAM29f040)
	}
}

func (c *CartridgeEasyFlash) controlUpdate(value uint8, update bool) {
	c.register02 = value & 0x87 // we only remember led, mode, exrom, game [led 0x80, other 0x07]
	mcIdx := icartridge.CartridgeModeOff
	//jumper := c.jumper << 3
	mxg := value & 0x07
	switch mxg {
	case 0:
		//GAME from jumper, EXROM high (i.e. Ultimax or Off)
		if c.jumper == 0 {
			mcIdx = icartridge.CartridgeModeUltimax
		} else {
			mcIdx = icartridge.CartridgeModeOff
		}
	case 1:
		//Reserved, don’t use this
	case 2:
		//GAME from jumper, EXROM low (i.e. 16K or 8K)
		if c.jumper == 0 {
			mcIdx = icartridge.CartridgeMode16K
		} else {
			mcIdx = icartridge.CartridgeMode8K
		}
	case 3:
		//Reserved, don’t use this
	case 4:
		//Cartridge ROM off (RAM at $DF00 still available)
		mcIdx = icartridge.CartridgeModeOff
	case 5:
		//Ultimax (Low bank at $8000, high bank at $e000) GAME = 0, EXROM = 1
		mcIdx = icartridge.CartridgeModeUltimax
	case 6:
		// 8k Cartridge (Low bank at $8000) GAME = 1, EXROM = 0
		mcIdx = icartridge.CartridgeMode8K
	case 7:
		//16k cartridge (Low bank at $8000, high bank at $a000)
		mcIdx = icartridge.CartridgeMode16K
	}
	if int(mcIdx) != c.memoryConfigIdx {
		c.memoryConfigIdx = int(mcIdx)
		v := icartridge.GetCartridgeSpec(mcIdx)
		c.game = v.Game
		c.exRom = v.ExRom
		c.intervalLo = v.IntervalLow
		c.intervalHi = v.IntervalHigh
		//fmt.Println("register002:", value, "mode:", memMode.mode, "exrom:", memMode.exrom, "game:", memMode.game, "led:", c.led)
		fmt.Println("EASYFLASH MEMORY CONFIG CHANGED:", mxg, "exrom:", c.exRom, "game:", c.game)
		if update {
			c.board.RebuildMemoryConfig()
		}
	}
	if led := value&0x80 == 0x80; led != c.led {
		c.led = led
		fmt.Println("EASYFLASH LED:", c.led)
	}
}

func (c *CartridgeEasyFlash) Write(i icartridge.RomInterval, addr uint16, data uint8) bool {
	if c.intervalLo == i {
		fmt.Printf("EASYFLASH Write LOW NOT DEFINED\n")
	} else if c.intervalHi == i {
		fmt.Printf("EASYFLASH Write HIGH NOT DEFINED\n")
	}
	//TODO IMPLEMENT
	//if i == ROM_LO {
	//	fmt.Printf("CartridgeEasyFlash can't be write %x => %d\n", addr, data)
	//	return true
	//}
	return false
}

func (c *CartridgeEasyFlash) Read(i icartridge.RomInterval, addr uint16) (uint8, bool) {
	if c.intervalLo == i {
		return c.romLRead(addr), true
	} else if c.intervalHi == i {
		return c.romHRead(addr), true
	}
	return 0, false
}

func (c *CartridgeEasyFlash) IORead(addr uint16) (uint8, bool) {
	ref := addr & 0xff00
	if ref == 0xdf00 {
		v := c.io2Read(addr)
		//fmt.Printf("EASYFLASH RAM READ: %x => %d\n", addr, v)
		return v, true
	}
	fmt.Printf("EASYFLASH RAM READ OUTSIDE: %x\n", addr)
	return 0, false
}

func (c *CartridgeEasyFlash) IOWrite(addr uint16, data uint8) bool {
	ref := addr & 0xff00
	if ref == 0xde00 {
		if addr == 0xde00 {
			c.register00 = data & BankMask
			return true
		}
		if addr == 0xde02 {
			c.controlUpdate(data, true)
			return true
		}
	} else if ref == 0xdf00 {
		//fmt.Printf("EASYFLASH RAM WRITE: %x => %d\n", addr, data)
		c.io2Store(addr, data)
		return true
	}
	fmt.Printf("EASYFLASH RAM WRITE OUSIDE: %x => %d\n", addr, data)
	return false
}

func (c *CartridgeEasyFlash) GetExRom() uint8 {
	return c.exRom
}

func (c *CartridgeEasyFlash) GetGame() uint8 {
	return c.game
}

func (c *CartridgeEasyFlash) romLRead(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateLow.Read(v)
}

func (c *CartridgeEasyFlash) romHRead(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateHigh.Read(v)
}

func (c *CartridgeEasyFlash) romLWrite(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateLow.Store(v, value)
}

func (c *CartridgeEasyFlash) romHWrite(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateHigh.Store(v, value)
}

// io1Peek - used by monitor [TODO]
func (c *CartridgeEasyFlash) io1Peek(addr uint16) uint8 {
	if addr&2 != 0 {
		return c.register02
	}
	return c.register00
}

func (c *CartridgeEasyFlash) io1Dump() int {
	mode := _easyFlashMemConfig[(c.jumper<<3)|(int(c.register02)&0x07)]
	bank := c.register00
	led := "false"
	if c.register02&0x80 != 0 {
		led = "true"
	}
	jumper := "off"
	if c.jumper != 0 {
		jumper = "on"
	}
	fmt.Printf("Mode: %d, Bank: %d, LED %s, jumper %s\n", mode, bank, led, jumper)
	//fmt.Printf("EAPI found: %s\n", mode, bank, led, jumper)
	//mon_out("EAPI found: %s\n", (memcmp(&romHBanks[0x1800], "eapi", 4) == 0) ? "yes" : "no")
	return 0
}

func (c *CartridgeEasyFlash) io2Read(addr uint16) uint8 {
	return c.ram[addr&0xff]
}

func (c *CartridgeEasyFlash) io2Store(addr uint16, value uint8) {
	c.ram[addr&0xff] = value
}

func (c *CartridgeEasyFlash) setEasyFlashJumper(val int) error {
	if val != 0 {
		c.jumper = 1
	} else {
		c.jumper = 0
	}
	return nil
}

func (c *CartridgeEasyFlash) setEasyFlashCrtWrite(val int) error {
	if val != 0 {
		c.crtWrite = 1
	} else {
		c.crtWrite = 0
	}
	return nil
}

func (c *CartridgeEasyFlash) setEasyFlashCrtOptimize(val int) error {
	if val != 0 {
		c.crtOptimize = 1
	} else {
		c.crtOptimize = 0
	}
	return nil
}

func (c *CartridgeEasyFlash) write_chip_if_not_empty(fd io.Writer, chip *loader.CrtChipHeader) error {
	for i := uint16(0); i < chip.Size; i++ {
		if (chip.Data[i] != 0xff) || (c.crtOptimize == 0) {
			if err := chip.Write(fd); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (c *CartridgeEasyFlash) binAttach(ldr *loader.CRTLoader) ([]byte, error) {
	rawCart := make([]uint8, 0x100000)
	for idx := range rawCart {
		rawCart[idx] = 0xff
	}
	copy(rawCart, ldr.GetData())
	return rawCart, nil
	//if err := util_file_load(filename, rawCart, 0x4000*NBanks, UTIL_FILE_LOAD_SKIP_ADDRESS); err != nil {
	//	return err
	//}
	//if len(rawCart) > len(c.rawCart) {
	//	return fmt.Errorf("invalid cartridge size")
	//}
}

func (c *CartridgeEasyFlash) crtAttach(ldr *loader.CRTLoader) ([]byte, error) {
	raw := make([]uint8, 0x100000)
	for idx := range raw {
		raw[idx] = 0xff
	}
	for {
		chip, err := ldr.ReadChipHeader()
		if err != nil {
			return nil, err
		}
		if chip == nil {
			break
		}
		if uint16(len(chip.Data)) != chip.Size {
			return nil, fmt.Errorf("invalid chip size")
		}
		if chip.Size == 0x2000 {
			if chip.Bank >= NBanks || !(chip.Start == 0x8000 || chip.Start == 0xa000 || chip.Start == 0xe000) {
				return nil, fmt.Errorf("invalid start")
			}
			index := uint(chip.Bank) << 14
			offset := uint(chip.Start) & uint(chip.Size)
			target := index | offset
			copy(raw[target:target+uint(chip.Size)], chip.Data)
		} else if chip.Size == 0x4000 {
			if chip.Bank >= NBanks || chip.Start != 0x8000 {
				return nil, fmt.Errorf("invalid start")
			}
			target := uint(chip.Bank) << 14
			copy(raw[target:target+uint(chip.Size)], chip.Data)
		} else {
			return nil, fmt.Errorf("unkwnown chip size")
		}
	}
	return raw, nil
}

func (c *CartridgeEasyFlash) Detach() error {
	if c.crtWrite != 0 {
		if err := c.flushImage(); err != nil {
			return err
		}
	}
	c.stateLow.Shutdown()
	c.stateHigh.Shutdown()
	c.filename = ""
	return nil
}

func (c *CartridgeEasyFlash) flushImage() error {
	if len(c.filename) == 0 {
		return nil
	}
	if c.filetype == loader.FiletypeBin {
		return c.binSave(c.filename)
	} else if c.filetype == loader.FiletypeCrt {
		return c.crtSave(c.filename)
	}
	return fmt.Errorf("unknown cartridget type")
}

func (c *CartridgeEasyFlash) binSave(filename string) error {
	const size = 0x2000
	if len(filename) == 0 {
		return fmt.Errorf("invalid file name")
	}
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	low := 0  //c.stateLow.flash_data
	high := 0 //c.stateHigh.flash_data
	var lowData []uint8
	var highData []uint8
	for i := 0; i < NBanks; i++ {
		if lowData, err = c.stateLow.ReadInterval(uint(low), uint(low)+size); err != nil {
			return nil
		}
		_, err = fd.Write(lowData)
		if err != nil {
			return err
		}
		if highData, err = c.stateLow.ReadInterval(uint(high), uint(high)+size); err != nil {
			return nil
		}
		_, err = fd.Write(highData)
		if err != nil {
			return err
		}
		low += size
		high += size
	}
	return nil
}

func (c *CartridgeEasyFlash) crtSave(filename string) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}

func (c *CartridgeEasyFlash) snapshotWriteModule(s *snapshot.Snapshot) error {
	m := s.NewModule(snapModuleName, SnapMajor, SnapMinor)
	m.Add("jumper", uint8(c.jumper))
	m.Add("register00", c.register00)
	m.Add("register00", c.register00)
	m.Add("ram", c.ram)
	//TODO
	//m.Add("romLBanks", c.romLBanks)
	//m.Add("romHBanks", c.romHBanks)
	if err := c.stateLow.SnapshotWriteModule(s, flashSnapModuleName); err != nil {
		return err
	}
	if err := c.stateHigh.SnapshotWriteModule(s, flashSnapModuleName); err != nil {
		return err
	}
	return nil
}

func (c *CartridgeEasyFlash) printConfigValue(val uint8, jumper uint8) {
	mode := 0
	exrom := uint8(1)
	game := uint8(1)
	if bitM := val&0x04 == 0x04; bitM {
		mode = 1
	}
	if bitX := val&0x02 == 0x02; bitX {
		exrom = 0
	}
	if mode != 0 {
		if bitG := val&0x01 == 0x01; bitG {
			game = 0
		}
	} else {
		game = jumper
	}
	fmt.Printf("/* %d */{jumper: %d, mode: %d, exrom: %d, game: %d},\n", val, jumper, mode, exrom, game)
	//fmt.Println("led:", led, "mode:", mode, "exrom:", exrom, "game:", game)
}

/*
func (c *CartridgeEasyFlash) mmu_translate(addr uint32) ([]uint8, int, int) {
	if c.stateHigh != nil && c.stateLow != nil {
		switch addr & 0xe000 {
		case 0xe000:
			if c.stateHigh.GetFlashState() == flash.StateRead {
				offset := (int(c.register00) * 0x2000) - 0xe000
				base := c.stateHigh.flashData[offset:]
				start := 0xe000
				limit := 0xfffd
				return base, start, limit
			}
			break
		case 0xa000:
			if c.stateHigh.GetFlashState() == flash.StateRead {
				offset := (int(c.register00) * 0x2000) - 0xa000
				base := c.stateHigh.flashData[offset:]
				start := 0xa000
				limit := 0xbffd
				return base, start, limit
			}
		case 0x8000:
			if c.stateLow.GetFlashState() == flash.StateRead {
				offset := (int(c.register00) * 0x2000) - 0x8000
				base := c.stateLow.flashData[offset:]
				start := 0x8000
				limit := 0x9ffd
				return base, start, limit
			}
		}
	}
	return nil, 0, 0
}
*/

/*
static const cmdline_option_t cmdline_options[] =
{
    { "-easyflashjumper", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashJumper", (resource_value_t)1,
      NULL, "Enable EasyFlash jumper" },
    { "+easyflashjumper", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashJumper", (resource_value_t)0,
      NULL, "Disable EasyFlash jumper" },
    { "-easyflashcrtwrite", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashWriteCRT", (resource_value_t)1,
      NULL, "Enable writing to EasyFlash .crt image" },
    { "+easyflashcrtwrite", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashWriteCRT", (resource_value_t)0,
      NULL, "Disable writing to EasyFlash .crt image" },
    { "-easyflashcrtoptimize", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashOptimizeCRT", (resource_value_t)1,
      NULL, "Enable EasyFlash .crt image optimize on write" },
    { "+easyflashcrtoptimize", SET_RESOURCE, CMDLINE_ATTRIB_NONE,
      NULL, NULL, "EasyFlashOptimizeCRT", (resource_value_t)0,
      NULL, "Disable writing to EasyFlash .crt image" },
    CMDLINE_LIST_END
};

int easyflash_cmdline_options_init(void)
{
    return cmdline_register_options(cmdline_options);
}
func (c *CartridgeEasyFlash) easyflash_cmdline_options_init(void) int {
{
return cmdline_register_options(cmdline_options);
}
*/
