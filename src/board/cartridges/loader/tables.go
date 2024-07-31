package loader

type Type int

const (
	TypeBin = Type(iota)
	TypeCrt
)

type MachineType int

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

var _machineContainer = map[string]map[MachineType]bool{
	"C64 CARTRIDGE   ": {MachineC64: true, MachineC64SC: true, MachineC128: true, MachineSCpu64: true},
	"C128 CARTRIDGE  ": {MachineC128: true},
	"CBM2 CARTRIDGE  ": {MachineCbm6x0: true, MachineCbm5x0: true},
	"VIC20 CARTRIDGE ": {MachineVic20: true},
	"PLUS4 CARTRIDGE ": {MachinePlus4: true},
}
