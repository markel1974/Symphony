package loader

// Constants representing various cartridge types used in Commodore 64 systems.
// The values include known cartridge types and those that match the CRT IDs.
// Some cartridge types are specific to particular hardware or software implementations.
// Negative values denote cartridges without a ROM image, while non-negative values align with CRT IDs.
// Each cartridge type is mapped to its associated source file for reference.
const (
	// Without a rom image
	CARTRIDGE_DIGIMAX            = -100 /* digimax.c */
	CARTRIDGE_DQBB               = -101 /* dqbb.c */
	CARTRIDGE_GEORAM             = -102 /* georam.c */
	CARTRIDGE_ISEPIC             = -103 /* isepic.c */
	CARTRIDGE_RAMCART            = -104 /* ramcart.c */
	CARTRIDGE_REU                = -105 /* reu.c */
	CARTRIDGE_SFX_SOUND_EXPANDER = -106 /* sfx_soundexpander.c, fmopl.c */
	CARTRIDGE_SFX_SOUND_SAMPLER  = -107 /* sfx_soundsampler.c */
	CARTRIDGE_MIDI_PASSPORT      = -108 /* c64-midi.c */
	CARTRIDGE_MIDI_DATEL         = -109 /* c64-midi.c */
	CARTRIDGE_MIDI_SEQUENTIAL    = -110 /* c64-midi.c */
	CARTRIDGE_MIDI_NAMESOFT      = -111 /* c64-midi.c */
	CARTRIDGE_MIDI_MAPLIN        = -112 /* c64-midi.c */
	CARTRIDGE_DS12C887RTC        = -113 /* ds12c887rtc.c */
	CARTRIDGE_TFE                = -116 /* ethernetcart.c */
	CARTRIDGE_TURBO232           = -117 /* c64acia1.c */
	CARTRIDGE_SWIFTLINK          = -118 /* c64acia1.c */
	CARTRIDGE_ACIA               = -119 /* c64acia1.c */
	CARTRIDGE_PLUS60K            = -120 /* plus60k.c */
	CARTRIDGE_PLUS256K           = -121 /* plus256k.c */
	CARTRIDGE_C64_256K           = -122 /* c64_256k.c */
	CARTRIDGE_CPM                = -123 /* cpmcart.c */
	CARTRIDGE_DEBUGCART          = -124 /* debugcart.c */

	// Known cartridge types
	CARTRIDGE_ULTIMAX      = -6 /* generic.c */ /* TODO: cartconv (4k and 12k binaries) */
	CARTRIDGE_GENERIC_8KB  = -3 /* generic.c */
	CARTRIDGE_GENERIC_16KB = -2 /* generic.c */
	CARTRIDGE_NONE         = -1
	CARTRIDGE_CRT          = 0

	// the following must match the CRT IDs
	CARTRIDGE_ACTION_REPLAY      = 1  /* actionreplay.c */
	CARTRIDGE_KCS_POWER          = 2  /* kcs.c */
	CARTRIDGE_FINAL_III          = 3  /* final3.c */
	CARTRIDGE_SIMONS_BASIC       = 4  /* simonsbasic.c */
	CARTRIDGE_OCEAN              = 5  /* ocean.c */
	CARTRIDGE_EXPERT             = 6  /* expert.c */
	CARTRIDGE_FUNPLAY            = 7  /* funplay.c */
	CARTRIDGE_SUPER_GAMES        = 8  /* supergames.c */
	CARTRIDGE_ATOMIC_POWER       = 9  /* atomicpower.c */
	CARTRIDGE_EPYX_FASTLOAD      = 10 /* epyxfastload.c */
	CARTRIDGE_WESTERMANN         = 11 /* westermann.c */
	CARTRIDGE_REX                = 12 /* rexutility.c */
	CARTRIDGE_FINAL_I            = 13 /* final.c */
	CARTRIDGE_MAGIC_FORMEL       = 14 /* magicformel.c */
	CARTRIDGE_GS                 = 15 /* gs.c */
	CARTRIDGE_WARPSPEED          = 16 /* warpspeed.c */
	CARTRIDGE_DINAMIC            = 17 /* dinamic.c */
	CARTRIDGE_ZAXXON             = 18 /* zaxxon.c */
	CARTRIDGE_MAGIC_DESK         = 19 /* magicdesk.c */
	CARTRIDGE_SUPER_SNAPSHOT_V5  = 20 /* supersnapshot.c */
	CARTRIDGE_COMAL80            = 21 /* comal80.c */
	CARTRIDGE_STRUCTURED_BASIC   = 22 /* stb.c */
	CARTRIDGE_ROSS               = 23 /* ross.c */
	CARTRIDGE_DELA_EP64          = 24 /* delaep64.c */
	CARTRIDGE_DELA_EP7x8         = 25 /* delaep7x8.c */
	CARTRIDGE_DELA_EP256         = 26 /* delaep256.c */
	CARTRIDGE_REX_EP256          = 27 /* rexep256.c */ /* TODO: cartconv */
	CARTRIDGE_MIKRO_ASSEMBLER    = 28 /* mikroass.c */
	CARTRIDGE_FINAL_PLUS         = 29 /* finalplus.c */ /* TODO: cartconv (24k binary) */
	CARTRIDGE_ACTION_REPLAY4     = 30 /* actionreplay4.c */
	CARTRIDGE_STARDOS            = 31 /* stardos.c */
	CARTRIDGE_EASYFLASH          = 32 /* easyflash.c */
	CARTRIDGE_EASYFLASH_XBANK    = 33 /* easyflash.c */ /* TODO: cartconv (no cart exists?) */
	CARTRIDGE_CAPTURE            = 34 /* capture.c */
	CARTRIDGE_ACTION_REPLAY3     = 35 /* actionreplay3.c */
	CARTRIDGE_RETRO_REPLAY       = 36 /* retroreplay.c */
	CARTRIDGE_MMC64              = 37 /* mmc64.c, spi-sdcard.c */
	CARTRIDGE_MMC_REPLAY         = 38 /* mmcreplay.c, ser-eeprom.c, spi-sdcard.c */
	CARTRIDGE_IDE64              = 39 /* ide64.c */
	CARTRIDGE_SUPER_SNAPSHOT     = 40 /* supersnapshot4.c */
	CARTRIDGE_IEEE488            = 41 /* c64tpi.c */
	CARTRIDGE_GAME_KILLER        = 42 /* gamekiller.c */
	CARTRIDGE_P64                = 43 /* prophet64.c */
	CARTRIDGE_EXOS               = 44 /* exos.c */
	CARTRIDGE_FREEZE_FRAME       = 45 /* freezeframe.c */
	CARTRIDGE_FREEZE_MACHINE     = 46 /* freezemachine.c */
	CARTRIDGE_SNAPSHOT64         = 47 /* snapshot64.c */
	CARTRIDGE_SUPER_EXPLODE_V5   = 48 /* superexplode5.c */
	CARTRIDGE_MAGIC_VOICE        = 49 /* magicvoice.c, tpicore.c, t6721.c */
	CARTRIDGE_ACTION_REPLAY2     = 50 /* actionreplay2.c */
	CARTRIDGE_MACH5              = 51 /* mach5.c */
	CARTRIDGE_DIASHOW_MAKER      = 52 /* diashowmaker.c */
	CARTRIDGE_PAGEFOX            = 53 /* pagefox.c */
	CARTRIDGE_KINGSOFT           = 54 /* kingsoft.c */
	CARTRIDGE_SILVERROCK_128     = 55 /* silverrock128.c */
	CARTRIDGE_FORMEL64           = 56 /* formel64.c */
	CARTRIDGE_RGCD               = 57 /* rgcd.c */
	CARTRIDGE_RRNETMK3           = 58 /* rrnetmk3.c */
	CARTRIDGE_EASYCALC           = 59 /* easycalc.c */
	CARTRIDGE_GMOD2              = 60 /* gmod2.c */
	CARTRIDGE_MAX_BASIC          = 61 /* maxbasic.c */
	CARTRIDGE_GMOD3              = 62 /* gmod3.c */
	CARTRIDGE_ZIPPCODE48         = 63 /* zippcode48.c */
	CARTRIDGE_BLACKBOX8          = 64 /* blackbox8.c */
	CARTRIDGE_BLACKBOX3          = 65 /* blackbox3.c */
	CARTRIDGE_BLACKBOX4          = 66 /* blackbox4.c */
	CARTRIDGE_REX_RAMFLOPPY      = 67 /* rexramfloppy.c */
	CARTRIDGE_BISPLUS            = 68 /* bisplus.c */
	CARTRIDGE_SDBOX              = 69 /* sdbox.c */
	CARTRIDGE_MULTIMAX           = 70 /* multimax.c */
	CARTRIDGE_BLACKBOX9          = 71 /* blackbox9.c */
	CARTRIDGE_LT_KERNAL          = 72 /* ltkernal.c */
	CARTRIDGE_RAMLINK            = 73 /* ramlink.c */
	CARTRIDGE_DREAN              = 74 /* drean.c */
	CARTRIDGE_IEEEFLASH64        = 75 /* ieeeflash64.c */
	CARTRIDGE_TURTLE_GRAPHICS_II = 76 /* turtlegraphics.c */
	CARTRIDGE_FREEZE_FRAME_MK2   = 77 /* freezeframe2.c */
	CARTRIDGE_PARTNER64          = 78 /* partner64.c */
	CARTRIDGE_HYPERBASIC         = 79 /* hyperbasic.c */
	CARTRIDGE_LAST               = 79 /* cartconv: last cartridge in list */
)
