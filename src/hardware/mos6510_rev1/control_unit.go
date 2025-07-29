package mos6510_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
	"reflect"
)

// ControlUnit represents a central unit for managing operations and their mappings within a system or application.
type ControlUnit struct {
	*component.BaseComponent
	opContainer        map[reflect.Value]string
	opReverseContainer map[string]reflect.Value
}

// NewControlUnit creates and initializes a new instance of ControlUnit with the specified parent, factory, label, and instance number.
// It sets up internal operation mappings for the ControlUnit and prepares it for handling instructions.
func NewControlUnit(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *ControlUnit {
	er := &ControlUnit{
		BaseComponent:      component.NewBaseComponent(),
		opContainer:        make(map[reflect.Value]string),
		opReverseContainer: make(map[string]reflect.Value),
	}

	er.addOpId(er.InstOpINI, "instOpINI")
	er.addOpId(er.InstOpHalt, "instOpHalt")
	er.addOpId(er.InstOpPHA, "instOpPHA")
	er.addOpId(er.InstOpPHA1, "instOpPHA1")
	er.addOpId(er.InstOpPLA, "instOpPLA")
	er.addOpId(er.InstOpPLA1, "instOpPLA1")
	er.addOpId(er.InstOpPLA2, "instOpPLA2")
	er.addOpId(er.InstOpPHP, "instOpPHP")
	er.addOpId(er.InstOpPHP1, "instOpPHP1")
	er.addOpId(er.InstOpPLP, "instOpPLP")
	er.addOpId(er.InstOpPLP1, "instOpPLP1")
	er.addOpId(er.InstOpPLP2, "instOpPLP2")
	er.addOpId(er.InstOpNMI, "instOpNMI")
	er.addOpId(er.InstOpNMI1, "instOpNMI1")
	er.addOpId(er.InstOpNMI2, "instOpNMI2")
	er.addOpId(er.InstOpNMI3, "instOpNMI3")
	er.addOpId(er.InstOpNMI4, "instOpNMI4")
	er.addOpId(er.InstOpNMI5, "instOpNMI5")
	er.addOpId(er.InstOpNMI6, "instOpNMI6")
	er.addOpId(er.InstApZER, "instApZER")
	er.addOpId(er.InstApZERx, "instApZERx")
	er.addOpId(er.InstApZERx1, "instApZERx1")
	er.addOpId(er.InstApZERy, "instApZERy")
	er.addOpId(er.InstApZERy1, "instApZERy1")
	er.addOpId(er.InstApABS, "instApABS")
	er.addOpId(er.InstApABS1, "instApABS1")
	er.addOpId(er.InstApABSx, "instApABSx")
	er.addOpId(er.InstApABSx1, "instApABSx1")
	er.addOpId(er.InstApABSx2, "instApABSx2")
	er.addOpId(er.InstApABSx3, "instApABSx3")
	er.addOpId(er.InstApABSy, "instApABSy")
	er.addOpId(er.InstApABSy1, "instApABSy1")
	er.addOpId(er.InstApABSy2, "instApABSy2")
	er.addOpId(er.InstApABSy3, "instApABSy3")
	er.addOpId(er.InstApINDx, "instApINDx")
	er.addOpId(er.InstApINDx1, "instApINDx1")
	er.addOpId(er.InstApINDx2, "instApINDx2")
	er.addOpId(er.InstApINDx3, "instApINDx3")
	er.addOpId(er.InstApINDy, "instApINDy")
	er.addOpId(er.InstApINDy1, "instApINDy1")
	er.addOpId(er.InstApINDy2, "instApINDy2")
	er.addOpId(er.InstApINDy3, "instApINDy3")
	er.addOpId(er.InstApINDy4, "instApINDy4")
	er.addOpId(er.InstAeABSx, "instAeABSx")
	er.addOpId(er.InstAeABSx1, "instAeABSx1")
	er.addOpId(er.InstAeABSx2, "instAeABSx2")
	er.addOpId(er.InstAeABSy, "instAeABSy")
	er.addOpId(er.InstAeABSy1, "instAeABSy1")
	er.addOpId(er.InstAeABSy2, "instAeABSy2")
	er.addOpId(er.InstAeINDy, "instAeINDy")
	er.addOpId(er.InstAeINDy1, "instAeINDy1")
	er.addOpId(er.InstAeINDy2, "instAeINDy2")
	er.addOpId(er.InstAeINDy3, "instAeINDy3")
	er.addOpId(er.InstMpZER, "instMpZER")
	er.addOpId(er.InstMpZERx, "instMpZERx")
	er.addOpId(er.InstMpZERx1, "instMpZERx1")
	er.addOpId(er.InstMpABS, "instMpABS")
	er.addOpId(er.InstMpABS1, "instMpABS1")
	er.addOpId(er.InstMpABSx, "instMpABSx")
	er.addOpId(er.InstMpABSx1, "instMpABSx1")
	er.addOpId(er.InstMpABSx2, "instMpABSx2")
	er.addOpId(er.InstMpABSx3, "instMpABSx3")
	er.addOpId(er.InstMpABSy, "instMpABSy")
	er.addOpId(er.InstMpABSy1, "instMpABSy1")
	er.addOpId(er.InstMpABSy2, "instMpABSy2")
	er.addOpId(er.InstMpABSy3, "instMpABSy3")
	er.addOpId(er.InstMpINDx, "instMpINDx")
	er.addOpId(er.InstMpINDx1, "instMpINDx1")
	er.addOpId(er.InstMpINDx2, "instMpINDx2")
	er.addOpId(er.InstMpINDx3, "instMpINDx3")
	er.addOpId(er.InstMpINDy, "instMpINDy")
	er.addOpId(er.InstMpINDy1, "instMpINDy1")
	er.addOpId(er.InstMpINDy2, "instMpINDy2")
	er.addOpId(er.InstMpINDy3, "instMpINDy3")
	er.addOpId(er.InstMpINDy4, "instMpINDy4")
	er.addOpId(er.InstOpRMW, "instOpRMW")
	er.addOpId(er.InstOpRMW1, "instOpRMW1")
	er.addOpId(er.InstOpLDA, "instOpLDA")
	er.addOpId(er.InstOiLDA, "instOiLDA")
	er.addOpId(er.InstOpLDX, "instOpLDX")
	er.addOpId(er.InstOiLDX, "instOiLDX")
	er.addOpId(er.InstOpLDY, "instOpLDY")
	er.addOpId(er.InstOiLDY, "instOiLDY")
	er.addOpId(er.InstOpSTA, "instOpSTA")
	er.addOpId(er.InstOpSTX, "instOpSTX")
	er.addOpId(er.InstOpSTY, "instOpSTY")
	er.addOpId(er.InstOpIRQ, "instOpIRQ")
	er.addOpId(er.InstOpIRQ1, "instOpIRQ1")
	er.addOpId(er.InstOpIRQ2, "instOpIRQ2")
	er.addOpId(er.InstOpIRQ3, "instOpIRQ3")
	er.addOpId(er.InstOpIRQ4, "instOpIRQ4")
	er.addOpId(er.InstOpIRQ5, "instOpIRQ5")
	er.addOpId(er.InstOpIRQ6, "instOpIRQ6")
	er.addOpId(er.InstOiNOP, "instOiNOP")
	er.addOpId(er.InstOaNOP, "instOaNOP")
	er.addOpId(er.InstOpLAX, "instOpLAX")
	er.addOpId(er.InstOpSAX, "instOpSAX")
	er.addOpId(er.InstOpSLO, "instOpSLO")
	er.addOpId(er.InstOpRLA, "instOpRLA")
	er.addOpId(er.InstOpSRE, "instOpSRE")
	er.addOpId(er.InstOpRRA, "instOpRRA")
	er.addOpId(er.InstOpDCP, "instOpDCP")
	er.addOpId(er.InstOpISB, "instOpISB")
	er.addOpId(er.InstOiANC, "instOiANC")
	er.addOpId(er.InstOiASR, "instOiASR")
	er.addOpId(er.InstOiARR, "instOiARR")
	er.addOpId(er.InstOiANE, "instOiANE")
	er.addOpId(er.InstOiLXA, "instOiLXA")
	er.addOpId(er.InstOiSBX, "instOiSBX")
	er.addOpId(er.InstOpLAS, "instOpLAS")
	er.addOpId(er.InstOpSHS, "instOpSHS")
	er.addOpId(er.InstOpSHY, "instOpSHY")
	er.addOpId(er.InstOpSHX, "instOpSHX")
	er.addOpId(er.InstOpSHA, "instOpSHA")
	er.addOpId(er.InstOpJAM, "instOpJAM")
	er.addOpId(er.InstOpTAX, "instOpTAX")
	er.addOpId(er.InstOpTXA, "instOpTXA")
	er.addOpId(er.InstOpTAY, "instOpTAY")
	er.addOpId(er.InstOpTYA, "instOpTYA")
	er.addOpId(er.InstOpTSX, "instOpTSX")
	er.addOpId(er.InstOpTXS, "instOpTXS")
	er.addOpId(er.InstOpSEC, "instOpSEC")
	er.addOpId(er.InstOpCLC, "instOpCLC")
	er.addOpId(er.InstOpSED, "instOpSED")
	er.addOpId(er.InstOpCLD, "instOpCLD")
	er.addOpId(er.InstOpSEI, "instOpSEI")
	er.addOpId(er.InstOpCLI, "instOpCLI")
	er.addOpId(er.InstOpCLV, "instOpCLV")
	er.addOpId(er.InstOpNOP, "instOpNOP")
	er.addOpId(er.InstOpADC, "instOpADC")
	er.addOpId(er.InstOiADC, "instOiADC")
	er.addOpId(er.InstOpSBC, "instOpSBC")
	er.addOpId(er.InstOiSBC, "instOiSBC")
	er.addOpId(er.InstOpINX, "instOpINX")
	er.addOpId(er.InstOpDEX, "instOpDEX")
	er.addOpId(er.InstOpINY, "instOpINY")
	er.addOpId(er.InstOpDEY, "instOpDEY")
	er.addOpId(er.InstOpINC, "instOpINC")
	er.addOpId(er.InstOpDEC, "instOpDEC")
	er.addOpId(er.InstOpAND, "instOpAND")
	er.addOpId(er.InstOiAND, "instOiAND")
	er.addOpId(er.InstOpORA, "instOpORA")
	er.addOpId(er.InstOiOPA, "instOiOPA")
	er.addOpId(er.InstOpEOR, "instOpEOR")
	er.addOpId(er.InstOiEOR, "instOiEOR")
	er.addOpId(er.InstOpCMP, "instOpCMP")
	er.addOpId(er.InstOiCMP, "instOiCMP")
	er.addOpId(er.InstOpCPX, "instOpCPX")
	er.addOpId(er.InstOiCPX, "instOiCPX")
	er.addOpId(er.InstOpCPY, "instOpCPY")
	er.addOpId(er.InstOiCPY, "instOiCPY")
	er.addOpId(er.InstOpBIT, "instOpBIT")
	er.addOpId(er.InstOpASL, "instOpASL")
	er.addOpId(er.InstOaASL, "instOaASL")
	er.addOpId(er.InstOpLSR, "instOpLSR")
	er.addOpId(er.InstOaLSR, "instOaLSR")
	er.addOpId(er.InstOpROL, "instOpROL")
	er.addOpId(er.InstOaROL, "instOaROL")
	er.addOpId(er.InstOpROR, "instOpROR")
	er.addOpId(er.InstOaROR, "instOaROR")
	er.addOpId(er.InstOpJMP, "instOpJMP")
	er.addOpId(er.InstOpJMP1, "instOpJMP1")
	er.addOpId(er.InstOiJMP, "instOiJMP")
	er.addOpId(er.InstOiJMP1, "instOiJMP1")
	er.addOpId(er.InstOpJSR, "instOpJSR")
	er.addOpId(er.InstOpJSR1, "instOpJSR1")
	er.addOpId(er.InstOpJSR2, "instOpJSR2")
	er.addOpId(er.InstOpJSR3, "instOpJSR3")
	er.addOpId(er.InstOpJSR4, "instOpJSR4")
	er.addOpId(er.InstOpRTS, "instOpRTS")
	er.addOpId(er.InstOpRTS1, "instOpRTS1")
	er.addOpId(er.InstOpRTS2, "instOpRTS2")
	er.addOpId(er.InstOpRTS3, "instOpRTS3")
	er.addOpId(er.InstOpRTS4, "instOpRTS4")
	er.addOpId(er.InstOpRTI, "instOpRTI")
	er.addOpId(er.InstOpRTI1, "instOpRTI1")
	er.addOpId(er.InstOpRTI2, "instOpRTI2")
	er.addOpId(er.InstOpRTI3, "instOpRTI3")
	er.addOpId(er.InstOpRTI4, "instOpRTI4")
	er.addOpId(er.InstOpBRK, "instOpBRK")
	er.addOpId(er.InstOpBRK1, "instOpBRK1")
	er.addOpId(er.InstOpBRK2, "instOpBRK2")
	er.addOpId(er.InstOpBRK3, "instOpBRK3")
	er.addOpId(er.InstOpBRK4, "instOpBRK4")
	er.addOpId(er.InstOpBRK5, "instOpBRK5")
	er.addOpId(er.InstOpBCS, "instOpBCS")
	er.addOpId(er.InstOpBCC, "instOpBCC")
	er.addOpId(er.InstOpBEQ, "instOpBEQ")
	er.addOpId(er.InstOpBNE, "instOpBNE")
	er.addOpId(er.InstOpBVS, "instOpBVS")
	er.addOpId(er.InstOpBVC, "instOpBVC")
	er.addOpId(er.InstOpBMI, "instOpBMI")
	er.addOpId(er.InstOpBPL, "instOpBPL")
	er.addOpId(er.InstOpBRAnp, "instOpBRAnp")
	er.addOpId(er.InstOpBRAbp, "instOpBRAbp")
	er.addOpId(er.InstOpBRAbp1, "instOpBRAbp1")
	er.addOpId(er.InstOpBRAfp, "instOpBRAfp")
	er.addOpId(er.InstOpBRAfp1, "instOpBRAfp1")

	er.BaseComponent.Register(factory, parent, "control", instance, er, references.IdInternalComponent(label, instance, "ControlUnit"))

	return er
}

// Setup initializes the ControlUnit, preparing it for operation and ensuring all necessary configurations are applied.
func (er *ControlUnit) Setup() error {
	return nil
}

// Connect establishes a connection to the control unit and returns an error if the operation fails.
func (er *ControlUnit) Connect() error {
	return nil
}

// EmulationRequired checks if the control unit requires emulation, returning true if emulation is needed, otherwise false.
func (er *ControlUnit) EmulationRequired() bool {
	return false
}

// Emulate performs the core emulation process for the ControlUnit, simulating its behavior based on defined parameters.
func (er *ControlUnit) Emulate() {
}

// Internal returns a boolean indicating an internal state or status within the ControlUnit.
func (er *ControlUnit) Internal() bool {
	return true
}

// Reset reinitializes the ControlUnit, clearing its state and preparing it for a new operation cycle.
func (er *ControlUnit) Reset() {
}

// GetInstOpINI returns a function that represents the InstOpINI operation associated with the ControlUnit instance.
func (er *ControlUnit) GetInstOpINI() func(cpu *CPU) {
	return er.InstOpINI
}

// GetInstOpHalt returns a function that represents the halt operation for the CPU, as defined in the control unit.
func (er *ControlUnit) GetInstOpHalt() func(cpu *CPU) {
	return er.InstOpHalt
}

// GetInstOpNMI returns a function that represents the InstOpNMI operation handler for the CPU.
func (er *ControlUnit) GetInstOpNMI() func(cpu *CPU) {
	return er.InstOpNMI
}

// GetInstOpIRQ returns a function that handles the IRQ operation of the Instruction Unit for the given CPU.
func (er *ControlUnit) GetInstOpIRQ() func(cpu *CPU) {
	return er.InstOpIRQ
}

// GetInstOpBRAbp returns a function pointer to the InstOpBRAbp method of the ControlUnit instance.
func (er *ControlUnit) GetInstOpBRAbp() func(cpu *CPU) {
	return er.InstOpBRAbp
}

// GetInstOpBRAfp returns a function that represents the branch operation for floating-point instructions.
func (er *ControlUnit) GetInstOpBRAfp() func(cpu *CPU) {
	return er.InstOpBRAfp
}

// GetInstOpBRAnp returns a function that handles the instruction operation for BRAnp in the CPU.
func (er *ControlUnit) GetInstOpBRAnp() func(cpu *CPU) {
	return er.InstOpBRAnp
}

// CreateModeTable initializes and returns a lookup table of functions mapped to CPU instructions.
func (er *ControlUnit) CreateModeTable() []func(*CPU) {
	modeTable := []func(*CPU){
		er.InstOpBRK, er.InstApINDx, er.InstOpJAM, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // 00
		er.InstOpPHP, er.InstOiOPA, er.InstOaASL, er.InstOiANC, er.InstApABS, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBPL, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // 10
		er.InstOpCLC, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
		er.InstOpJSR, er.InstApINDx, er.InstOpJAM, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // 20
		er.InstOpPLP, er.InstOiAND, er.InstOaROL, er.InstOiANC, er.InstApABS, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBMI, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // 30
		er.InstOpSEC, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
		er.InstOpRTI, er.InstApINDx, er.InstOpJAM, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // 40
		er.InstOpPHA, er.InstOiEOR, er.InstOaLSR, er.InstOiASR, er.InstOpJMP, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBVC, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // 50
		er.InstOpCLI, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
		er.InstOpRTS, er.InstApINDx, er.InstOpJAM, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // 60
		er.InstOpPLA, er.InstOiADC, er.InstOaROR, er.InstOiARR, er.InstApABS, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBVS, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // 70
		er.InstOpSEI, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
		er.InstOiNOP, er.InstApINDx, er.InstOiNOP, er.InstApINDx, er.InstApZER, er.InstApZER, er.InstApZER, er.InstApZER, // 80
		er.InstOpDEY, er.InstOiNOP, er.InstOpTXA, er.InstOiANE, er.InstApABS, er.InstApABS, er.InstApABS, er.InstApABS,
		er.InstOpBCC, er.InstApINDy, er.InstOpJAM, er.InstApINDy, er.InstApZERx, er.InstApZERx, er.InstApZERy, er.InstApZERy, // 90
		er.InstOpTYA, er.InstApABSy, er.InstOpTXS, er.InstApABSy, er.InstApABSx, er.InstApABSx, er.InstApABSy, er.InstApABSy,
		er.InstOiLDY, er.InstApINDx, er.InstOiLDX, er.InstApINDx, er.InstApZER, er.InstApZER, er.InstApZER, er.InstApZER, // a0
		er.InstOpTAY, er.InstOiLDA, er.InstOpTAX, er.InstOiLXA, er.InstApABS, er.InstApABS, er.InstApABS, er.InstApABS,
		er.InstOpBCS, er.InstAeINDy, er.InstOpJAM, er.InstAeINDy, er.InstApZERx, er.InstApZERx, er.InstApZERy, er.InstApZERy, // b0
		er.InstOpCLV, er.InstAeABSy, er.InstOpTSX, er.InstAeABSy, er.InstAeABSx, er.InstAeABSx, er.InstAeABSy, er.InstAeABSy,
		er.InstOiCPY, er.InstApINDx, er.InstOiNOP, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // c0
		er.InstOpINY, er.InstOiCMP, er.InstOpDEX, er.InstOiSBX, er.InstApABS, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBNE, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // d0
		er.InstOpCLD, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
		er.InstOiCPX, er.InstApINDx, er.InstOiNOP, er.InstMpINDx, er.InstApZER, er.InstApZER, er.InstMpZER, er.InstMpZER, // e0
		er.InstOpINX, er.InstOiSBC, er.InstOpNOP, er.InstOiSBC, er.InstApABS, er.InstApABS, er.InstMpABS, er.InstMpABS,
		er.InstOpBEQ, er.InstAeINDy, er.InstOpJAM, er.InstMpINDy, er.InstApZERx, er.InstApZERx, er.InstMpZERx, er.InstMpZERx, // f0
		er.InstOpSED, er.InstAeABSy, er.InstOpNOP, er.InstMpABSy, er.InstAeABSx, er.InstAeABSx, er.InstMpABSx, er.InstMpABSx,
	}
	return modeTable
}

// CreateOpTable initializes and returns an operation table as a slice of functions for CPU instruction execution.
func (er *ControlUnit) CreateOpTable() []func(cpu *CPU) {
	opTable := []func(*CPU){
		er.InstOpJAM, er.InstOpORA, er.InstOpJAM, er.InstOpSLO, er.InstOaNOP, er.InstOpORA, er.InstOpASL, er.InstOpSLO, // 00
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOaNOP, er.InstOpORA, er.InstOpASL, er.InstOpSLO,
		er.InstOpJAM, er.InstOpORA, er.InstOpJAM, er.InstOpSLO, er.InstOaNOP, er.InstOpORA, er.InstOpASL, er.InstOpSLO, // 10
		er.InstOpJAM, er.InstOpORA, er.InstOpJAM, er.InstOpSLO, er.InstOaNOP, er.InstOpORA, er.InstOpASL, er.InstOpSLO,
		er.InstOpJAM, er.InstOpAND, er.InstOpJAM, er.InstOpRLA, er.InstOpBIT, er.InstOpAND, er.InstOpROL, er.InstOpRLA, // 20
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpBIT, er.InstOpAND, er.InstOpROL, er.InstOpRLA,
		er.InstOpJAM, er.InstOpAND, er.InstOpJAM, er.InstOpRLA, er.InstOaNOP, er.InstOpAND, er.InstOpROL, er.InstOpRLA, // 30
		er.InstOpJAM, er.InstOpAND, er.InstOpJAM, er.InstOpRLA, er.InstOaNOP, er.InstOpAND, er.InstOpROL, er.InstOpRLA,
		er.InstOpJAM, er.InstOpEOR, er.InstOpJAM, er.InstOpSRE, er.InstOaNOP, er.InstOpEOR, er.InstOpLSR, er.InstOpSRE, // 40
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpEOR, er.InstOpLSR, er.InstOpSRE,
		er.InstOpJAM, er.InstOpEOR, er.InstOpJAM, er.InstOpSRE, er.InstOaNOP, er.InstOpEOR, er.InstOpLSR, er.InstOpSRE, // 50
		er.InstOpJAM, er.InstOpEOR, er.InstOpJAM, er.InstOpSRE, er.InstOaNOP, er.InstOpEOR, er.InstOpLSR, er.InstOpSRE,
		er.InstOpJAM, er.InstOpADC, er.InstOpJAM, er.InstOpRRA, er.InstOaNOP, er.InstOpADC, er.InstOpROR, er.InstOpRRA, // 60
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOiJMP, er.InstOpADC, er.InstOpROR, er.InstOpRRA,
		er.InstOpJAM, er.InstOpADC, er.InstOpJAM, er.InstOpRRA, er.InstOaNOP, er.InstOpADC, er.InstOpROR, er.InstOpRRA, // 70
		er.InstOpJAM, er.InstOpADC, er.InstOpJAM, er.InstOpRRA, er.InstOaNOP, er.InstOpADC, er.InstOpROR, er.InstOpRRA,
		er.InstOpJAM, er.InstOpSTA, er.InstOpJAM, er.InstOpSAX, er.InstOpSTY, er.InstOpSTA, er.InstOpSTX, er.InstOpSAX, // 80
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpSTY, er.InstOpSTA, er.InstOpSTX, er.InstOpSAX,
		er.InstOpJAM, er.InstOpSTA, er.InstOpJAM, er.InstOpSHA, er.InstOpSTY, er.InstOpSTA, er.InstOpSTX, er.InstOpSAX, // 90
		er.InstOpJAM, er.InstOpSTA, er.InstOpJAM, er.InstOpSHS, er.InstOpSHY, er.InstOpSTA, er.InstOpSHX, er.InstOpSHA,
		er.InstOpJAM, er.InstOpLDA, er.InstOpJAM, er.InstOpLAX, er.InstOpLDY, er.InstOpLDA, er.InstOpLDX, er.InstOpLAX, // a0
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpLDY, er.InstOpLDA, er.InstOpLDX, er.InstOpLAX,
		er.InstOpJAM, er.InstOpLDA, er.InstOpJAM, er.InstOpLAX, er.InstOpLDY, er.InstOpLDA, er.InstOpLDX, er.InstOpLAX, // b0
		er.InstOpJAM, er.InstOpLDA, er.InstOpJAM, er.InstOpLAS, er.InstOpLDY, er.InstOpLDA, er.InstOpLDX, er.InstOpLAX,
		er.InstOpJAM, er.InstOpCMP, er.InstOpJAM, er.InstOpDCP, er.InstOpCPY, er.InstOpCMP, er.InstOpDEC, er.InstOpDCP, // c0
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpCPY, er.InstOpCMP, er.InstOpDEC, er.InstOpDCP,
		er.InstOpJAM, er.InstOpCMP, er.InstOpJAM, er.InstOpDCP, er.InstOaNOP, er.InstOpCMP, er.InstOpDEC, er.InstOpDCP, // d0
		er.InstOpJAM, er.InstOpCMP, er.InstOpJAM, er.InstOpDCP, er.InstOaNOP, er.InstOpCMP, er.InstOpDEC, er.InstOpDCP,
		er.InstOpJAM, er.InstOpSBC, er.InstOpJAM, er.InstOpISB, er.InstOpCPX, er.InstOpSBC, er.InstOpINC, er.InstOpISB, // e0
		er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpJAM, er.InstOpCPX, er.InstOpSBC, er.InstOpINC, er.InstOpISB,
		er.InstOpJAM, er.InstOpSBC, er.InstOpJAM, er.InstOpISB, er.InstOaNOP, er.InstOpSBC, er.InstOpINC, er.InstOpISB, // f0
		er.InstOpJAM, er.InstOpSBC, er.InstOpJAM, er.InstOpISB, er.InstOaNOP, er.InstOpSBC, er.InstOpINC, er.InstOpISB,
	}
	return opTable
}

// addOpId associates a function of type func(cpu *CPU) with a string identifier for reverse mapping.
func (er *ControlUnit) addOpId(v func(cpu *CPU), id string) {
	r := reflect.ValueOf(v)
	er.opContainer[r] = id
	er.opReverseContainer[id] = r
}

// GetOpId retrieves the operation ID and a boolean indicating its existence for a given CPU operation function.
func (er *ControlUnit) GetOpId(v func(cpu *CPU)) (string, bool) {
	if v == nil {
		return "", false
	}
	r := reflect.ValueOf(v)
	ret, ok := er.opContainer[r]
	if !ok {
		return "", false
	}
	return ret, true
}

// GetOpFn retrieves an operation function associated with the given string key and returns it along with a success flag.
func (er *ControlUnit) GetOpFn(v string) (func(cpu *CPU), bool) {
	r, ok := er.opReverseContainer[v]
	if !ok {
		return nil, false
	}
	ret, ok := r.Interface().(func(cpu *CPU))
	if !ok {
		return nil, false
	}
	return ret, true
}
