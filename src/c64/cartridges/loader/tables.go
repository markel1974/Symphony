package loader

// Type represents an integer-based enumeration for defining specific types or categories in a program.
type Type int

// TypeBin represents a binary type constant derived from the Type enum.
// TypeCrt represents a certificate type constant derived from the Type enum.
const (
	TypeBin = Type(iota)
	TypeCrt
)

// MachineType defines an enumeration used to represent various machine types in a system.
type MachineType int

// MachineC64 represents the Commodore 64 machine type.
// MachineC128 represents the Commodore 128 machine type.
// MachineVic20 represents the Commodore VIC-20 machine type.
// MachinePet represents the Commodore PET machine type.
// MachineCbm5x0 represents the CBM-II 500 series machine type.
// MachineCbm6x0 represents the CBM-II 600 series machine type.
// MachinePlus4 represents the Commodore Plus/4 machine type.
// MachineC64DTV represents the Commodore 64 Direct-to-TV machine type.
// MachineC64SC represents the Commodore 64 SC machine type.
// MachineVSid represents the virtual SID chip machine type.
// MachineSCpu64 represents the SuperCPU 64 machine type.
const (
	MachineC64 = MachineType(iota)
	MachineC128
	MachineVic20
	MachinePet
	MachineCbm5x0
	MachineCbm6x0
	MachinePlus4
	MachineC64DTV
	MachineC64SC
	MachineVSid
	MachineSCpu64
)

// _machineContainer maps cartridge identifiers to supported MachineType configurations for validation purposes.
var _machineContainer = map[string]map[MachineType]bool{
	"C64 CARTRIDGE   ": {MachineC64: true, MachineC64SC: true, MachineC128: true, MachineSCpu64: true},
	"C128 CARTRIDGE  ": {MachineC128: true},
	"CBM2 CARTRIDGE  ": {MachineCbm6x0: true, MachineCbm5x0: true},
	"VIC20 CARTRIDGE ": {MachineVic20: true},
	"PLUS4 CARTRIDGE ": {MachinePlus4: true},
}
