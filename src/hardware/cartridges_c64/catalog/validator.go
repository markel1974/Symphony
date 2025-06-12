package catalog

import (
	"errors"
	"fmt"
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

// ValidateMachine checks if the given MachineType and id combination is valid based on predefined mappings. Returns error if invalid.
func ValidateMachine(m MachineType, id string) error {
	supported, ok := _machineContainer[id]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	_, ok = supported[m]
	if !ok {
		return fmt.Errorf("invalid crt header")
	}
	return nil
}

// ValidateCartridge checks if the provided cartridge data matches specific byte patterns and returns an error if invalid.
func ValidateCartridge(data []byte) error {
	i4 := data[0x4]
	i5 := data[0x5]
	i6 := data[0x6]
	i7 := data[0x7]
	i8 := data[0x8]
	if i4 == 0xc3 && i5 == 0xc2 && i6 == 0xcd && i7 == 0x38 && i8 == 0x30 {
		return nil
		//valid cartridge
	}
	return errors.New("invalid cartridge")
}
