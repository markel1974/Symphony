package loader

// Constants representing various cartridge types used in Commodore 64 systems.
// The values include known cartridge types and those that match the CRT IDs.
// Some cartridge types are specific to particular hardware or software implementations.
// Negative values denote cartridges without a ROM image, while non-negative values align with CRT IDs.
// Each cartridge type is mapped to its associated source file for reference.
const (
	// Without a rom image
	CartridgeDIGIMAX            = -100
	CartridgeDQBB               = -101
	CartridgeGEORAM             = -102
	CartridgeISEPIC             = -103
	CartridgeRAMCART            = -104
	CartridgeREU                = -105
	CartridgeSFX_SOUND_EXPANDER = -106
	CartridgeSFX_SOUND_SAMPLER  = -107
	CartridgeMIDI_PASSPORT      = -108
	CartridgeMIDI_DATEL         = -109
	CartridgeMIDI_SEQUENTIAL    = -110
	CartridgeMIDI_NAMESOFT      = -111
	CartridgeMIDI_MAPLIN        = -112
	CartridgeDS12C887RTC        = -113
	CartridgeTFE                = -116
	CartridgeTURBO232           = -117
	CartridgeSWIFTLINK          = -118
	CartridgeACIA               = -119
	CartridgePLUS60K            = -120
	CartridgePLUS256K           = -121
	CartridgeC64_256K           = -122
	CartridgeCPM                = -123
	CartridgeDEBUGCART          = -124
)

// Known cartridge types
const (
	CartridgeULTIMAX      = -6
	CartridgeGENERIC_8KB  = -3
	CartridgeGENERIC_16KB = -2
	CartridgeNONE         = -1
	CartridgeCRT          = 0
)

// the following must match the CRT IDs
const (
	CartridgeActionReplay       = 1
	CartridgeKCSPower           = 2
	CartridgeFinalIII           = 3
	CartridgeSimonsBasic        = 4
	CartridgeOcean              = 5
	CartridgeEXPERT             = 6
	CartridgeFUNPLAY            = 7
	CartridgeSUPER_GAMES        = 8
	CartridgeATOMIC_POWER       = 9
	CartridgeEPYX_FASTLOAD      = 10
	CartridgeWESTERMANN         = 11
	CartridgeREX                = 12
	CartridgeFINAL_I            = 13
	CartridgeMAGIC_FORMEL       = 14
	CartridgeGS                 = 15
	CartridgeWARPSPEED          = 16
	CartridgeDINAMIC            = 17
	CartridgeZAXXON             = 18
	CartridgeMAGIC_DESK         = 19
	CartridgeSUPER_SNAPSHOT_V5  = 20
	CartridgeCOMAL80            = 21
	CartridgeSTRUCTURED_BASIC   = 22
	CartridgeROSS               = 23
	CartridgeDELA_EP64          = 24
	CartridgeDELA_EP7x8         = 25
	CartridgeDELA_EP256         = 26
	CartridgeREX_EP256          = 27
	CartridgeMIKRO_ASSEMBLER    = 28
	CartridgeFINAL_PLUS         = 29
	CartridgeACTION_REPLAY4     = 30
	CartridgeSTARDOS            = 31
	CartridgeEASYFLASH          = 32
	CartridgeEASYFLASH_XBANK    = 33
	CartridgeCAPTURE            = 34
	CartridgeACTION_REPLAY3     = 35
	CartridgeRETRO_REPLAY       = 36
	CartridgeMMC64              = 37
	CartridgeMMC_REPLAY         = 38
	CartridgeIDE64              = 39
	CartridgeSUPER_SNAPSHOT     = 40
	CartridgeIEEE488            = 41
	CartridgeGAME_KILLER        = 42
	CartridgeP64                = 43
	CartridgeEXOS               = 44
	CartridgeFREEZE_FRAME       = 45
	CartridgeFREEZE_MACHINE     = 46
	CartridgeSNAPSHOT64         = 47
	CartridgeSUPER_EXPLODE_V5   = 48
	CartridgeMAGIC_VOICE        = 49
	CartridgeACTION_REPLAY2     = 50
	CartridgeMACH5              = 51
	CartridgeDIASHOW_MAKER      = 52
	CartridgePAGEFOX            = 53
	CartridgeKINGSOFT           = 54
	CartridgeSILVERROCK_128     = 55
	CartridgeFORMEL64           = 56
	CartridgeRGCD               = 57
	CartridgeRRNETMK3           = 58
	CartridgeEASYCALC           = 59
	CartridgeGMOD2              = 60
	CartridgeMAX_BASIC          = 61
	CartridgeGMOD3              = 62
	CartridgeZIPPCODE48         = 63
	CartridgeBLACKBOX8          = 64
	CartridgeBLACKBOX3          = 65
	CartridgeBLACKBOX4          = 66
	CartridgeREX_RAMFLOPPY      = 67
	CartridgeBISPLUS            = 68
	CartridgeSDBOX              = 69
	CartridgeMULTIMAX           = 70
	CartridgeBLACKBOX9          = 71
	CartridgeLT_KERNAL          = 72
	CartridgeRAMLINK            = 73
	CartridgeDREAN              = 74
	CartridgeIEEEFLASH64        = 75
	CartridgeTURTLE_GRAPHICS_II = 76
	CartridgeFREEZE_FRAME_MK2   = 77
	CartridgePARTNER64          = 78
	CartridgeHYPERBASIC         = 79
	CartridgeLAST               = 79
)
