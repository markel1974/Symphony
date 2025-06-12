package catalog

import (
	"github.com/markel1974/c64emu/src/references"
)

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

// _registerHardware is a map that associates hardware names with factory functions to create ICartridgeC64 instances.
var _registerHardware = make(map[string]func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64)

// _registerType is a mapping that associates an integer type identifier with a factory function for creating ICartridgeC64 instances.
var _registerType = make(map[int]func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64)

// _registerSize is a map associating sizes to factory functions that generate instances of ICartridgeC64 components.
var _registerSize = make(map[int]func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64)

// _registerSizeDefault is a variable that stores a function for registering a default C64-compatible cartridge size.
var _registerSizeDefault func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64

// RegisterHardware registers a hardware component with a given name and factory function for creating C64 cartridges.
func RegisterHardware(name string, factory func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64) {
	_registerHardware[name] = factory
}

// RegisterType registers a cartridge factory function to a specific type for later instantiation.
func RegisterType(kind int, factory func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64) {
	_registerType[kind] = factory
}

// RegisterSize registers a cartridge size with its corresponding factory for creating instances of ICartridgeC64.
// The `size` parameter specifies the cartridge size to register.
// The `factory` parameter is a function to create an ICartridgeC64 with component references and configuration details.
func RegisterSize(size int, factory func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64) {
	_registerSize[size] = factory
}

// RegisterSizeDefault sets the default factory function for creating ICartridgeC64 instances with specified parameters.
func RegisterSizeDefault(factory func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64) {
	_registerSizeDefault = factory
}

// ByHardware returns a function that binds an IComponent, IComponentFactory, string, and int to produce an ICartridgeC64 instance.
// The returned function is retrieved from the _registerHardware map using the provided name as the key.
func ByHardware(name string) func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64 {
	return _registerHardware[name]
}

// ByType returns a function to create an ICartridgeC64 instance based on the provided kind and parameters.
func ByType(kind int) func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64 {
	return _registerType[kind]
}

// BySize returns a function that generates an ICartridgeC64 implementation based on the provided size.
func BySize(size int) func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64 {
	return _registerSize[size]
}

// BySizeDefault provides a default function for registering a C64 cartridge with specific components and configurations.
func BySizeDefault() func(references.IComponent, references.IComponentFactory, string, int) references.ICartridgeC64 {
	return _registerSizeDefault
}
