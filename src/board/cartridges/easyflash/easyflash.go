package easyflash

//https://skoe.de/easyflash/files/devdocs/EasyFlash-ProgRef.pdf

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/board/cartridges/easyflash/flash"
	"github.com/markel1974/c64emu/src/board/cartridges/icartridge"
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
	"github.com/markel1974/c64emu/src/board/iboard"
	"github.com/markel1974/c64emu/src/board/ram"
	"github.com/markel1974/c64emu/src/board/snapshot"
	"io"
	"os"
)

/* decoding table of the modes */
/* bit3 = jumper, bit2 = mode, bit1 = !exrom, bit0 = game */
var easyflash_memconfig = []uint8{
	/* jumper off, mode 0, trough 00,01,10,11 in game/exrom bits */
	3, /* exrom high, game low, jumper off */
	3, /* Reserved, don't use this */
	1, /* exrom low, game low, jumper off */
	1, /* Reserved, don't use this */
	/* jumper off, mode 1, trough 00,01,10,11 in game/exrom bits */
	2, 3, 0, 1,
	/* jumper on, mode 0, trough 00,01,10,11 in game/exrom bits */
	2, /* exrom high, game low, jumper on */
	3, /* Reserved, don't use this */
	0, /* exrom low, game low, jumper on */
	1, /* Reserved, don't use this */
	/* jumper on, mode 1, trough 00,01,10,11 in game/exrom bits */
	2, 3, 0, 1,
}

type memConfig struct {
	jumper uint8
	mode   uint8
	exrom  uint8
	game   uint8
}

var easyflashMemconfig2 = []*memConfig{
	/* 0 */ {jumper: 0, mode: 0, exrom: 1, game: 0},
	/* 1 */ {jumper: 0, mode: 0, exrom: 1, game: 0},
	/* 2 */ {jumper: 0, mode: 0, exrom: 0, game: 0},
	/* 3 */ {jumper: 0, mode: 0, exrom: 0, game: 0},
	/* 4 */ {jumper: 0, mode: 1, exrom: 1, game: 1},
	/* 5 */ {jumper: 0, mode: 1, exrom: 1, game: 0},
	/* 6 */ {jumper: 0, mode: 1, exrom: 0, game: 1},
	/* 7 */ {jumper: 0, mode: 1, exrom: 0, game: 0},
	/* 8 */ {jumper: 1, mode: 0, exrom: 1, game: 1},
	/* 9 */ {jumper: 1, mode: 0, exrom: 1, game: 1},
	/* A */ {jumper: 1, mode: 0, exrom: 0, game: 1},
	/* B */ {jumper: 1, mode: 0, exrom: 0, game: 1},
	/* C */ {jumper: 1, mode: 1, exrom: 1, game: 1},
	/* D */ {jumper: 1, mode: 1, exrom: 1, game: 0},
	/* E */ {jumper: 1, mode: 1, exrom: 0, game: 1},
	/* F */ {jumper: 1, mode: 1, exrom: 0, game: 0},
}

const (
	CART_RAM_SIZE = 256
)

const (
	CARTRIDGE_NAME_EASYFLASH = "EasyFlash" /* see http://skoe.de/easyflash/ */
	STRING_EASYFLASH         = CARTRIDGE_NAME_EASYFLASH
)

var _eapiam29f040 = []byte{
	0x65, 0x61, 0x70, 0x69, 0xc1, 0x4d,
	0x2f, 0xcd, 0x32, 0x39, 0xc6, 0x30, 0x34, 0x30,
	0x20, 0xd6, 0x31, 0x2e, 0x34, 0x00, 0x08, 0x78,
	0xa5, 0x4b, 0x48, 0xa5, 0x4c, 0x48, 0xa9, 0x60,
	0x85, 0x4b, 0x20, 0x4b, 0x00, 0xba, 0xbd, 0x00,
	0x01, 0x85, 0x4c, 0xca, 0xbd, 0x00, 0x01, 0x85,
	0x4b, 0x18, 0x90, 0x70, 0x4c, 0x67, 0x01, 0x4c,
	0xa4, 0x01, 0x4c, 0x39, 0x02, 0x4c, 0x40, 0x02,
	0x4c, 0x44, 0x02, 0x4c, 0x4e, 0x02, 0x4c, 0x58,
	0x02, 0x4c, 0x8e, 0x02, 0x4c, 0xd9, 0x02, 0x4c,
	0xd9, 0x02, 0x8d, 0x02, 0xde, 0xa9, 0xaa, 0x8d,
	0x55, 0x85, 0xa9, 0x55, 0x8d, 0xaa, 0x82, 0xa9,
	0xa0, 0x8d, 0x55, 0x85, 0xad, 0xf2, 0xdf, 0x8d,
	0x00, 0xde, 0xa9, 0x00, 0x8d, 0xff, 0xff, 0xa2,
	0x07, 0x8e, 0x02, 0xde, 0x60, 0x8d, 0x02, 0xde,
	0xa9, 0xaa, 0x8d, 0x55, 0xe5, 0xa9, 0x55, 0x8d,
	0xaa, 0xe2, 0xa9, 0xa0, 0x8d, 0x55, 0xe5, 0xd0,
	0xdb, 0xa2, 0x55, 0x8e, 0xe3, 0xdf, 0x8c, 0xe4,
	0xdf, 0xa2, 0x85, 0x8e, 0x02, 0xde, 0x8d, 0xff,
	0xff, 0x4c, 0xbb, 0xdf, 0xad, 0xff, 0xff, 0x60,
	0xcd, 0xff, 0xff, 0x60, 0xa2, 0x6f, 0xa0, 0x7f,
	0xb1, 0x4b, 0x9d, 0x80, 0xdf, 0xdd, 0x80, 0xdf,
	0xd0, 0x21, 0x88, 0xca, 0x10, 0xf2, 0xa2, 0x00,
	0xe8, 0x18, 0xbd, 0x80, 0xdf, 0x65, 0x4b, 0x9d,
	0x80, 0xdf, 0xe8, 0xbd, 0x80, 0xdf, 0x65, 0x4c,
	0x9d, 0x80, 0xdf, 0xe8, 0xe0, 0x1e, 0xd0, 0xe8,
	0x18, 0x90, 0x06, 0xa9, 0x01, 0x8d, 0xb9, 0xdf,
	0x38, 0x68, 0x85, 0x4c, 0x68, 0x85, 0x4b, 0xb0,
	0x48, 0xa9, 0xaa, 0xa0, 0xe5, 0x20, 0xd5, 0xdf,
	0xa0, 0x85, 0x20, 0xd5, 0xdf, 0xa9, 0x55, 0xa2,
	0xaa, 0xa0, 0xe2, 0x20, 0xd7, 0xdf, 0xa2, 0xaa,
	0xa0, 0x82, 0x20, 0xd7, 0xdf, 0xa9, 0x90, 0xa0,
	0xe5, 0x20, 0xd5, 0xdf, 0xa0, 0x85, 0x20, 0xd5,
	0xdf, 0xad, 0x00, 0xa0, 0x8d, 0xf1, 0xdf, 0xae,
	0x01, 0xa0, 0x8e, 0xb9, 0xdf, 0xc9, 0x01, 0xd0,
	0x06, 0xe0, 0xa4, 0xd0, 0x02, 0xf0, 0x0c, 0xc9,
	0x20, 0xd0, 0x39, 0xe0, 0xe2, 0xd0, 0x35, 0xf0,
	0x02, 0xb0, 0x50, 0xad, 0x00, 0x80, 0xae, 0x01,
	0x80, 0xc9, 0x01, 0xd0, 0x06, 0xe0, 0xa4, 0xd0,
	0x02, 0xf0, 0x08, 0xc9, 0x20, 0xd0, 0x19, 0xe0,
	0xe2, 0xd0, 0x15, 0xa0, 0x3f, 0x8c, 0x00, 0xde,
	0xae, 0x02, 0x80, 0xd0, 0x13, 0xae, 0x02, 0xa0,
	0xd0, 0x12, 0x88, 0x10, 0xf0, 0x18, 0x90, 0x12,
	0xa9, 0x02, 0xd0, 0x0a, 0xa9, 0x03, 0xd0, 0x06,
	0xa9, 0x04, 0xd0, 0x02, 0xa9, 0x05, 0x8d, 0xb9,
	0xdf, 0x38, 0xa9, 0x00, 0x8d, 0x00, 0xde, 0xa0,
	0xe0, 0xa9, 0xf0, 0x20, 0xd7, 0xdf, 0xa0, 0x80,
	0x20, 0xd7, 0xdf, 0xad, 0xb9, 0xdf, 0xb0, 0x08,
	0xae, 0xf1, 0xdf, 0xa0, 0x40, 0x28, 0x18, 0x60,
	0x28, 0x38, 0x60, 0x8d, 0xb7, 0xdf, 0x8e, 0xb9,
	0xdf, 0x8e, 0xed, 0xdf, 0x8c, 0xba, 0xdf, 0x08,
	0x78, 0x98, 0x29, 0xbf, 0x8d, 0xee, 0xdf, 0xa9,
	0x00, 0x8d, 0x00, 0xde, 0xa9, 0x85, 0xc0, 0xe0,
	0x90, 0x05, 0x20, 0xc1, 0xdf, 0xb0, 0x03, 0x20,
	0x9e, 0xdf, 0xa2, 0x14, 0x20, 0xec, 0xdf, 0xf0,
	0x06, 0xca, 0xd0, 0xf8, 0x18, 0x90, 0x63, 0xad,
	0xf2, 0xdf, 0x8d, 0x00, 0xde, 0x18, 0x90, 0x72,
	0x8d, 0xb7, 0xdf, 0x8e, 0xb9, 0xdf, 0x8c, 0xba,
	0xdf, 0x08, 0x78, 0x98, 0xc0, 0x80, 0xf0, 0x04,
	0xa0, 0xe0, 0xa9, 0xa0, 0x8d, 0xee, 0xdf, 0xc8,
	0xc8, 0xc8, 0xc8, 0xc8, 0xa9, 0xaa, 0x20, 0xd5,
	0xdf, 0xa9, 0x55, 0xa2, 0xaa, 0x88, 0x88, 0x88,
	0x20, 0xd7, 0xdf, 0xa9, 0x80, 0xc8, 0xc8, 0xc8,
	0x20, 0xd5, 0xdf, 0xa9, 0xaa, 0x20, 0xd5, 0xdf,
	0xa9, 0x55, 0xa2, 0xaa, 0x88, 0x88, 0x88, 0x20,
	0xd7, 0xdf, 0xad, 0xb7, 0xdf, 0x8d, 0x00, 0xde,
	0xa2, 0x00, 0x8e, 0xed, 0xdf, 0x88, 0x88, 0xa9,
	0x30, 0x20, 0xd7, 0xdf, 0xa9, 0xff, 0xaa, 0xa8,
	0xd0, 0x24, 0xad, 0xf2, 0xdf, 0x8d, 0x00, 0xde,
	0xa0, 0x80, 0xa9, 0xf0, 0x20, 0xd7, 0xdf, 0xa0,
	0xe0, 0xa9, 0xf0, 0x20, 0xd7, 0xdf, 0x28, 0x38,
	0xb0, 0x02, 0x28, 0x18, 0xac, 0xba, 0xdf, 0xae,
	0xb9, 0xdf, 0xad, 0xb7, 0xdf, 0x60, 0x20, 0xec,
	0xdf, 0xf0, 0x09, 0xca, 0xd0, 0xf8, 0x88, 0xd0,
	0xf5, 0x18, 0x90, 0xce, 0xad, 0xf2, 0xdf, 0x8d,
	0x00, 0xde, 0x18, 0x90, 0xdd, 0x8d, 0xf2, 0xdf,
	0x8d, 0x00, 0xde, 0x60, 0xad, 0xf2, 0xdf, 0x60,
	0x8d, 0xf3, 0xdf, 0x8e, 0xe9, 0xdf, 0x8c, 0xea,
	0xdf, 0x60, 0x8e, 0xf4, 0xdf, 0x8c, 0xf5, 0xdf,
	0x8d, 0xf6, 0xdf, 0x60, 0xad, 0xf2, 0xdf, 0x8d,
	0x00, 0xde, 0x20, 0xe8, 0xdf, 0x8d, 0xb7, 0xdf,
	0x8e, 0xf0, 0xdf, 0x8c, 0xf1, 0xdf, 0xa9, 0x00,
	0x8d, 0xba, 0xdf, 0xf0, 0x3b, 0xad, 0xf4, 0xdf,
	0xd0, 0x10, 0xad, 0xf5, 0xdf, 0xd0, 0x08, 0xad,
	0xf6, 0xdf, 0xf0, 0x0b, 0xce, 0xf6, 0xdf, 0xce,
	0xf5, 0xdf, 0xce, 0xf4, 0xdf, 0x90, 0x45, 0x38,
	0xb0, 0x42, 0x8d, 0xb7, 0xdf, 0x8e, 0xf0, 0xdf,
	0x8c, 0xf1, 0xdf, 0xae, 0xe9, 0xdf, 0xad, 0xea,
	0xdf, 0xc9, 0xa0, 0x90, 0x02, 0x09, 0x40, 0xa8,
	0xad, 0xb7, 0xdf, 0x20, 0x80, 0xdf, 0xb0, 0x24,
	0xee, 0xe9, 0xdf, 0xd0, 0x19, 0xee, 0xea, 0xdf,
	0xad, 0xf3, 0xdf, 0x29, 0xe0, 0xcd, 0xea, 0xdf,
	0xd0, 0x0c, 0xad, 0xf3, 0xdf, 0x0a, 0x0a, 0x0a,
	0x8d, 0xea, 0xdf, 0xee, 0xf2, 0xdf, 0x18, 0xad,
	0xba, 0xdf, 0xf0, 0xa1, 0xac, 0xf1, 0xdf, 0xae,
	0xf0, 0xdf, 0xad, 0xb7, 0xdf, 0x60, 0xff, 0xff,
	0xff, 0xff,
}

const (
	EASYFLASH_N_BANK_BITS = 6
	EASYFLASH_N_BANKS     = 1 << (EASYFLASH_N_BANK_BITS)
	EASYFLASH_BANK_MASK   = (EASYFLASH_N_BANKS) - 1
)

const snap_module_name = "CARTEF"
const flash_snap_module_name = "FLASH040EF"

const (
	SNAP_MAJOR = 0
	SNAP_MINOR = 0
)

/*
static io_source_t easyflash_io1_device = {
	CARTRIDGE_NAME_EASYFLASH, // name of the device
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
	CARTRIDGE_NAME_EASYFLASH, // name of the device
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
	board       iboard.IBoard
	intervals   icartridge.Interval
	game        uint8
	exRom       uint8
	stateLow    *flash.Flash040Context /* the 29F040B statemachine */
	stateHigh   *flash.Flash040Context /* the 29F040B statemachine */
	jumper      int                    /* the jumper */
	crtWrite    int                    /* writing back to crt enabled */
	crtOptimize int                    /* optimizing crt enabled */
	register00  uint8                  /* backup of the registers */
	register02  uint8                  /* backup of the registers */
	ram         []uint8                /* extra RAM */
	filename    string                 /* filename when attached */
	filetype    int
	led         bool
}

func New(game uint8, exRom uint8, lo icartridge.Interval, hi icartridge.Interval) *CartridgeEasyFlash {
	return &CartridgeEasyFlash{
		game:      game,
		exRom:     exRom,
		intervals: lo | hi,
		stateLow:  nil,
		stateHigh: nil,
		filetype:  0,
		jumper:    0,
		led:       false,
		ram:       make([]byte, CART_RAM_SIZE),
	}
}

func (c *CartridgeEasyFlash) Setup(board iboard.IBoard, ldr *loader.CRTLoader) error {
	var rawCart []byte
	c.board = board
	rp := ram.NewInitiator(255, 2, 1, 0x100, 255, 0, 0, 0)
	rp.InitWithPattern(c.ram, CART_RAM_SIZE)
	c.filename = ldr.Name
	c.filetype = 0
	var err error
	if ldr.GetMode() == loader.ModeCrt {
		c.filetype = loader.CARTRIDGE_FILETYPE_CRT
		if rawCart, err = c.crtAttach(ldr); err != nil {
			return err
		}
	} else {
		c.filetype = loader.CARTRIDGE_FILETYPE_BIN
		if rawCart, err = c.binAttach(ldr); err != nil {
			return err
		}
	}
	c.register00 = 0
	//c.jumper = 0
	//for x := uint8(0); x < 8; x++ {
	//	c.controlUpdate(x)
	//}
	//c.jumper = 1
	//for x := uint8(0); x < 8; x++ {
	//	c.controlUpdate(x)
	//}
	//TODO PATCH MOMENTANEA!!!!!
	c.controlUpdate(7)
	//c.controlUpdate(0)
	c.initializeFlash(rawCart)
	return nil
}

func (c *CartridgeEasyFlash) initializeFlash(rawCart []byte) {
	low := make([]byte, 0x80000)
	high := make([]byte, 0x80000)
	// split interleaved low and high banks
	for i := uint(0); i < EASYFLASH_N_BANKS; i++ {
		const size = 0x2000
		start := i * size
		p1 := i * (size * 2)
		p2 := p1 + size
		copy(low[start:start+size], rawCart[p1:p1+size])
		copy(high[start:start+size], rawCart[p2:p2+size])
	}
	c.stateLow = flash.NewFlash040Context(c.board, flash.FLASH040_TYPE_B, low)
	c.stateHigh = flash.NewFlash040Context(c.board, flash.FLASH040_TYPE_B, high)
	eApiCheck := high[0x1800 : 0x1800+4]
	if bytes.Compare(eApiCheck, []byte("eapi")) == 0 {
		eApi := make([]byte, 17)
		for k := 0; k < 16; k++ {
			eApi[k] = c.stateHigh.Peek(uint32(0x1804 + k))
		}
		fmt.Printf("EF: EAPI found (%s)", string(eApi))
		_ = c.stateHigh.StoreInterval(0x1800, 0x1800+768, _eapiam29f040)
	}
}

func (c *CartridgeEasyFlash) controlUpdate(value uint8) {
	c.register02 = value & 0x87 // we only remember led, mode, exrom, game [led 0x80, other 0x07]
	c.led = c.register02&0x80 == 0x80
	jumper := c.jumper << 3
	mxg := int(c.register02) & 0x07
	cfg := jumper | mxg
	memMode := easyflashMemconfig2[cfg]
	c.exRom = memMode.exrom
	c.game = memMode.game
	fmt.Println("register002:", value, "mode:", memMode.mode, "exrom:", memMode.exrom, "game:", memMode.game, "led:", c.led)
}

func (c *CartridgeEasyFlash) Write(i icartridge.Interval, addr uint16, data uint8) bool {
	if i&c.intervals != 0 {
		fmt.Printf("TEST\n")
	}
	//TODO IMPLEMENT
	//if i == ROM_LO {
	//	fmt.Printf("CartridgeEasyFlash can't be write %x => %d\n", addr, data)
	//	return true
	//}
	return false
}

func (c *CartridgeEasyFlash) Read(i icartridge.Interval, addr uint16) (uint8, bool) {
	if i&c.intervals != 0 {
		if i == icartridge.ROM_LO {
			return c.roml_read(addr), true
		} else if i == icartridge.ROM_HI_1 || i == icartridge.ROM_HI_2 {
			return c.romh_read(addr), true
		}
		return 0, false
	}
	return 0, false
}

func (c *CartridgeEasyFlash) IORead(addr uint16) (uint8, bool) {
	ref := addr & 0xff00
	if ref == 0xdf00 {
		v := c.io2_read(addr)
		return v, true
	}
	return 0, false
}

func (c *CartridgeEasyFlash) IOWrite(addr uint16, data uint8) bool {
	ref := addr & 0xff00
	if ref == 0xde00 {
		if addr == 0xde00 {
			c.register00 = data & EASYFLASH_BANK_MASK
			return true
		}
		if addr == 0xde02 {
			c.controlUpdate(data)
			return true
		}
	} else if ref == 0xdf00 {
		c.io2_store(addr, data)
		return true
	}
	return false
}

func (c *CartridgeEasyFlash) GetExRom() uint8 {
	return c.exRom
}

func (c *CartridgeEasyFlash) GetGame() uint8 {
	return c.game
}

func (c *CartridgeEasyFlash) roml_read(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateLow.Read(v)
}

func (c *CartridgeEasyFlash) romh_read(addr uint16) uint8 {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	return c.stateHigh.Read(v)
}

func (c *CartridgeEasyFlash) roml_write(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateLow.Store(v, value)
}

func (c *CartridgeEasyFlash) romh_write(addr uint16, value uint8) {
	v := (uint(c.register00) * 0x2000) + (uint(addr) & 0x1fff)
	c.stateHigh.Store(v, value)
}

// io1_peek - used by monitor [TODO]
func (c *CartridgeEasyFlash) io1_peek(addr uint16) uint8 {
	if addr&2 != 0 {
		return c.register02
	}
	return c.register00
}

func (c *CartridgeEasyFlash) io1_dump() int {
	mode := easyflash_memconfig[(c.jumper<<3)|(int(c.register02)&0x07)]
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

func (c *CartridgeEasyFlash) io2_read(addr uint16) uint8 {
	return c.ram[addr&0xff]
}

func (c *CartridgeEasyFlash) io2_store(addr uint16, value uint8) {
	c.ram[addr&0xff] = value
}

func (c *CartridgeEasyFlash) set_easyflash_jumper(val int) error {
	if val != 0 {
		c.jumper = 1
	} else {
		c.jumper = 0
	}
	return nil
}

func (c *CartridgeEasyFlash) set_easyflash_crt_write(val int) error {
	if val != 0 {
		c.crtWrite = 1
	} else {
		c.crtWrite = 0
	}
	return nil
}

func (c *CartridgeEasyFlash) set_easyflash_crt_optimize(val int) error {
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
	//if err := util_file_load(filename, rawCart, 0x4000*EASYFLASH_N_BANKS, UTIL_FILE_LOAD_SKIP_ADDRESS); err != nil {
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
		if chip.Size == 0x2000 {
			if chip.Bank >= EASYFLASH_N_BANKS || !(chip.Start == 0x8000 || chip.Start == 0xa000 || chip.Start == 0xe000) {
				return nil, fmt.Errorf("invalid start")
			}
			target := (uint(chip.Bank) << 14) | (uint(chip.Start) & uint(chip.Size))
			copy(raw[target:], chip.Data)
		} else if chip.Size == 0x4000 {
			if chip.Bank >= EASYFLASH_N_BANKS || chip.Start != 0x8000 {
				return nil, fmt.Errorf("invalid start")
			}
			target := uint(chip.Bank) << 14
			copy(raw[target:], chip.Data)
		} else {
			return nil, fmt.Errorf("unkwnown chip size")
		}
	}
	return raw, nil
}

func (c *CartridgeEasyFlash) Detach() error {
	if c.crtWrite != 0 {
		if err := c.flush_image(); err != nil {
			return err
		}
	}
	c.stateLow.Shutdown()
	c.stateHigh.Shutdown()
	c.filename = ""
	return nil
}

func (c *CartridgeEasyFlash) flush_image() error {
	if len(c.filename) == 0 {
		return nil
	}
	if c.filetype == loader.CARTRIDGE_FILETYPE_BIN {
		return c.binSave(c.filename)
	} else if c.filetype == loader.CARTRIDGE_FILETYPE_CRT {
		return c.crt_save(c.filename)
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
	for i := 0; i < EASYFLASH_N_BANKS; i++ {
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

func (c *CartridgeEasyFlash) crt_save(filename string) error {
	//TODO IMPLEMENT
	return fmt.Errorf("unimplemented")
}

func (c *CartridgeEasyFlash) snapshotWriteModule(s *snapshot.Snapshot) error {
	m := s.NewModule(snap_module_name, SNAP_MAJOR, SNAP_MINOR)
	m.Add("jumper", uint8(c.jumper))
	m.Add("register00", c.register00)
	m.Add("register00", c.register00)
	m.Add("ram", c.ram)
	//TODO
	//m.Add("romLBanks", c.romLBanks)
	//m.Add("romHBanks", c.romHBanks)
	if err := c.stateLow.SnapshotWriteModule(s, flash_snap_module_name); err != nil {
		return err
	}
	if err := c.stateHigh.SnapshotWriteModule(s, flash_snap_module_name); err != nil {
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
			if c.stateHigh.GetFlashState() == flash.FLASH040_STATE_READ {
				offset := (int(c.register00) * 0x2000) - 0xe000
				base := c.stateHigh.flashData[offset:]
				start := 0xe000
				limit := 0xfffd
				return base, start, limit
			}
			break
		case 0xa000:
			if c.stateHigh.GetFlashState() == flash.FLASH040_STATE_READ {
				offset := (int(c.register00) * 0x2000) - 0xa000
				base := c.stateHigh.flashData[offset:]
				start := 0xa000
				limit := 0xbffd
				return base, start, limit
			}
		case 0x8000:
			if c.stateLow.GetFlashState() == flash.FLASH040_STATE_READ {
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
