package easyflash

//https://skoe.de/easyflash/files/devdocs/EasyFlash-ProgRef.pdf

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/c64/cartridges/easyflash/flash"
	"github.com/markel1974/c64emu/src/c64/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/c64/cartridges/loader"
	"github.com/markel1974/c64emu/src/filler"

	//"github.com/markel1974/c64emu/src/board/ram"
	"github.com/markel1974/c64emu/src/c64/snapshot"
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
	jumper          bool
	writeEnabled    bool    /* writing back to crt enabled */
	optimize        bool    /* optimizing crt enabled */
	register00      uint8   /* backup of the registers */
	register02      uint8   /* backup of the registers */
	ram             []uint8 /* extra RAM */
	filename        string  /* filename when attached */
	filetype        loader.Type
	led             bool
	updateEApi      bool
}

func GetType() int {
	return loader.CARTRIDGE_EASYFLASH
}

func New() icartridge.ICartridge {
	return &CartridgeEasyFlash{
		game:            1,
		exRom:           1,
		register00:      0,
		register02:      0,
		intervalHi:      icartridge.ROM_HI_1,
		intervalLo:      icartridge.ROM_LO,
		stateLow:        nil,
		stateHigh:       nil,
		filetype:        0,
		jumper:          false,
		led:             false,
		ram:             make([]byte, CartRamSize),
		memoryConfigIdx: -1,
		updateEApi:      true,
	}
}

func (c *CartridgeEasyFlash) Setup(board icartridge.IExpansion, ldr *loader.CRTLoader) error {
	var rawCart []byte
	c.board = board
	c.id = ldr.GetId()
	c.game = uint8(ldr.Game)
	c.exRom = uint8(ldr.ExRom)
	rp := filler.New(255, 2, 1, 0x100, 255, 0, 0, 0)
	rp.InitWithPattern(c.ram, CartRamSize)
	c.filename = ldr.Name
	c.filetype = 0
	var err error
	if ldr.GetType() == loader.TypeCrt {
		c.filetype = loader.TypeCrt
		if rawCart, err = c.crtAttach(ldr); err != nil {
			return err
		}
	} else {
		c.filetype = loader.TypeBin
		if rawCart, err = c.binAttach(ldr); err != nil {
			return err
		}
	}
	c.register00 = 0
	c.controlUpdate(0, false)
	c.initialize(rawCart)
	return nil
}

func (c *CartridgeEasyFlash) Reset() {

}

func (c *CartridgeEasyFlash) GetId() string {
	return c.id
}

func (c *CartridgeEasyFlash) GetExRom() uint8 {
	return c.exRom
}

func (c *CartridgeEasyFlash) GetGame() uint8 {
	return c.game
}

func (c *CartridgeEasyFlash) EmulationRequired() bool {
	return false
}

func (c *CartridgeEasyFlash) Emulate() {
}

func (c *CartridgeEasyFlash) SetJumper(j bool) {
	c.jumper = j
}

func (c *CartridgeEasyFlash) SetWriteEnabled(e bool) {
	c.writeEnabled = e
}

func (c *CartridgeEasyFlash) SetOptimize(o bool) {
	c.optimize = o
}

func (c *CartridgeEasyFlash) SetUpdateEApi(e bool) {
	c.updateEApi = e
}

func (c *CartridgeEasyFlash) initialize(rawCart []byte) {
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

	if c.updateEApi {
		eApiHeader := high[EApiStartAddress : EApiStartAddress+len(EApiHeader)]
		if bytes.Compare(eApiHeader, []byte(EApiHeader)) == 0 {
			//updating eapi
			//eApi := make([]byte, 17)
			//for k := 0; k < 16; k++ {
			//	eApi[k] = c.stateHigh.Peek(uint32(0x1804 + k))
			//}
			_ = c.stateHigh.StoreInterval(0x1800, _eApiAM29f040)
		}
	}
}

func (c *CartridgeEasyFlash) controlUpdate(value uint8, update bool) {
	c.register02 = value & 0x87 // we only remember led, mode, exrom, game [led 0x80, other 0x07]
	mode := icartridge.CartridgeModeOff
	mxg := value & 0x07
	switch mxg {
	case 0:
		//GAME from jumper, EXROM high (i.e. Ultimax or Off)
		if !c.jumper {
			mode = icartridge.CartridgeModeUltimax
		} else {
			mode = icartridge.CartridgeModeOff
		}
	case 1, 3:
		//Reserved, don’t use these
	case 2:
		//GAME from jumper, EXROM low (i.e. 16K or 8K)
		if !c.jumper {
			mode = icartridge.CartridgeMode16K
		} else {
			mode = icartridge.CartridgeMode8K
		}
	case 4:
		//Cartridge ROM off (RAM at $DF00 still available)
		mode = icartridge.CartridgeModeOff
	case 5:
		//Ultimax (Low bank at $8000, high bank at $e000) GAME = 0, EXROM = 1
		mode = icartridge.CartridgeModeUltimax
	case 6:
		//8k Cartridge (Low bank at $8000) GAME = 1, EXROM = 0
		mode = icartridge.CartridgeMode8K
	case 7:
		//16k cartridge (Low bank at $8000, high bank at $a000)
		mode = icartridge.CartridgeMode16K
	}
	if int(mode) != c.memoryConfigIdx {
		c.memoryConfigIdx = int(mode)
		v := icartridge.GetCartridgeSpec(mode)
		c.game = v.Game
		c.exRom = v.ExRom
		c.intervalLo = v.IntervalLow
		c.intervalHi = v.IntervalHigh
		//fmt.Println("EASYFLASH MEMORY CONFIG CHANGED:", mxg, "exrom:", c.exRom, "game:", c.game, "LO", c.intervalLo, "HIGH", c.intervalHi)
		if update {
			c.board.GameExRomConfigChanged()
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
	bank := addr & 0xff00
	if bank == 0xdf00 {
		v := c.io2Read(addr)
		return v, true
	} else if bank == 0xde00 {
		// read is never valid, regs are write only
		fmt.Printf("EASYFLASH WARNING: regs $de00-$deff are write only -> $%x\n", addr)
		return 0, false
	}
	return 0, false
}

func (c *CartridgeEasyFlash) IOWrite(addr uint16, data uint8) bool {
	bank := addr & 0xff00
	if bank == 0xde00 {
		switch addr & 0x2 {
		case 0:
			c.register00 = data & BankMask
			return true
		case 1:
			return true
		case 2:
			c.controlUpdate(data, true)
			return true
		case 3:
			return true
		}
	} else if bank == 0xdf00 {
		c.io2Store(addr, data)
		return true
	}
	return false
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

func (c *CartridgeEasyFlash) io2Read(addr uint16) uint8 {
	return c.ram[addr&0xff]
}

func (c *CartridgeEasyFlash) io2Store(addr uint16, value uint8) {
	c.ram[addr&0xff] = value
}

func (c *CartridgeEasyFlash) writeChipIfNotEmpty(fd io.Writer, chip *loader.CrtChipHeader) error {
	for i := uint16(0); i < chip.Size; i++ {
		if (chip.Data[i] != 0xff) || !c.optimize {
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
	if c.writeEnabled {
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
	if c.filetype == loader.TypeBin {
		return c.binSave(c.filename)
	} else if c.filetype == loader.TypeCrt {
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
	low := 0
	high := 0
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
	m.Add("jumper", c.jumper)
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
