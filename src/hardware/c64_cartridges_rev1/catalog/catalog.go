package catalog

import (
	"github.com/markel1974/c64emu/src/references"
)

// Constants representing various cartridges without a rom image used in Commodore 64 systems.
// The values include known cartridge types and those that match the CRT IDs.
// Some cartridge types are specific to particular hardware or software implementations.
// Negative values denote cartridges without a ROM image, while non-negative values align with CRT IDs.
// Each cartridge type is mapped to its associated source file for reference.
const (
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

// the following must match the CRT IDs
const (
	CartridgeCRT = iota
	CartridgeActionReplay
	CartridgeKCSPower
	CartridgeFinalIII
	CartridgeSimonsBasic
	CartridgeOcean
	CartridgeExpert
	CartridgeFunPlay
	CartridgeSuperGames
	CartridgeAtomicPower
	CartridgeEpyxFastLoad
	CartridgeWesterMann
	CartridgeREX
	CartridgeFinalI
	CartridgeMagicFormel
	CartridgeGS
	CartridgeWARPSPEED
	CartridgeDINAMIC
	CartridgeZAXXON
	CartridgeMagicDesk
	CartridgeSUPER_SNAPSHOT_V5
	CartridgeCOMAL80
	CartridgeSTRUCTURED_BASIC
	CartridgeROSS
	CartridgeDELA_EP64
	CartridgeDELA_EP7x8
	CartridgeDELA_EP256
	CartridgeREX_EP256
	CartridgeMIKRO_ASSEMBLER
	CartridgeFINAL_PLUS
	CartridgeACTION_REPLAY4
	CartridgeSTARDOS
	CartridgeEasyFlash
	CartridgeEASYFLASH_XBANK
	CartridgeCAPTURE
	CartridgeACTION_REPLAY3
	CartridgeRETRO_REPLAY
	CartridgeMMC64
	CartridgeMMC_REPLAY
	CartridgeIDE64
	CartridgeSUPER_SNAPSHOT
	CartridgeIEEE488
	CartridgeGAME_KILLER
	CartridgeP64
	CartridgeEXOS
	CartridgeFREEZE_FRAME
	CartridgeFREEZE_MACHINE
	CartridgeSNAPSHOT64
	CartridgeSUPER_EXPLODE_V5
	CartridgeMAGIC_VOICE
	CartridgeACTION_REPLAY2
	CartridgeMACH5
	CartridgeDIASHOW_MAKER
	CartridgePAGEFOX
	CartridgeKINGSOFT
	CartridgeSILVERROCK_128
	CartridgeFORMEL64
	CartridgeRGCD
	CartridgeRRNETMK3
	CartridgeEASYCALC
	CartridgeGMOD2
	CartridgeMAX_BASIC
	CartridgeGMOD3
	CartridgeZIPPCODE48
	CartridgeBLACKBOX8
	CartridgeBLACKBOX3
	CartridgeBLACKBOX4
	CartridgeREX_RAMFLOPPY
	CartridgeBISPLUS
	CartridgeSDBOX
	CartridgeMULTIMAX
	CartridgeBLACKBOX9
	CartridgeLT_KERNAL
	CartridgeRAMLINK
	CartridgeDREAN
	CartridgeIEEEFLASH64
	CartridgeTURTLE_GRAPHICS_II
	CartridgeFREEZE_FRAME_MK2
	CartridgePARTNER64
	CartridgeHYPERBASIC
)

// _registerHardware is a map that associates hardware names with factory functions to create IC64Cartridge instances.
var _registerHardware = make(map[string]func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge)

// _registerType is a mapping that associates an integer type identifier with a factory function for creating IC64Cartridge instances.
var _registerType = make(map[int]func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge)

// _registerSize is a map associating sizes to factory functions that generate instances of IC64Cartridge components.
var _registerSize = make(map[int]func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge)

// _registerSizeDefault is a variable that stores a function for registering a default C64-compatible cartridge size.
var _registerSizeDefault func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge

// RegisterHardware registers a hardware component with a given name and factory function for creating C64 cartridges.
func RegisterHardware(name string, factory func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge) {
	_registerHardware[name] = factory
}

// RegisterType registers a cartridge factory function to a specific type for later instantiation.
func RegisterType(kind int, factory func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge) {
	_registerType[kind] = factory
}

// RegisterSize registers a cartridge size with its corresponding factory for creating instances of IC64Cartridge.
// The `size` parameter specifies the cartridge size to register.
// The `factory` parameter is a function to create an IC64Cartridge with component references and configuration details.
func RegisterSize(size int, factory func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge) {
	_registerSize[size] = factory
}

// RegisterSizeDefault sets the default factory function for creating IC64Cartridge instances with specified parameters.
func RegisterSizeDefault(factory func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge) {
	_registerSizeDefault = factory
}

// ByHardware returns a function that binds an IComponent, IComponentFactory, string, and int to produce an IC64Cartridge instance.
// The returned function is retrieved from the _registerHardware map using the provided name as the key.
func ByHardware(name string) func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge {
	return _registerHardware[name]
}

// ByType returns a function to create an IC64Cartridge instance based on the provided kind and parameters.
func ByType(kind int) func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge {
	return _registerType[kind]
}

// BySize returns a function that generates an IC64Cartridge implementation based on the provided size.
func BySize(size int) func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge {
	return _registerSize[size]
}

// BySizeDefault provides a default function for registering a C64 cartridge with specific components and configurations.
func BySizeDefault() func(references.IComponent, references.IComponentFactory, string, int) references.IC64Cartridge {
	return _registerSizeDefault
}
