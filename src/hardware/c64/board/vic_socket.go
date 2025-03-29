package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// IVICSocketConnection defines an interface for responding to various triggers in the VIC socket connection lifecycle.
// LastCycleTrigger is invoked on the completion of the last cycle.
// VBlankTrigger handles vertical blanking interval initiation.
// RDYLowTrigger sets the RDY (Ready) line to low or high based on the boolean value.
// AECLowTrigger sets the AEC (Address Enable) line to low or high based on the boolean value.
type IVICSocketConnection interface {
	LastCycleTrigger()

	VBlankTrigger()

	RDYLowTrigger(v bool)

	AECLowTrigger(v bool)
}

// VICSocket encapsulates the connections and components required for VIC chip socket emulation in the C64 system.
// It integrates with display, quartz clocking, PLA memory, and programmable interrupt controller components.
// Provides methods and triggers for managing display rendering and interaction with the C64 hardware components.
type VICSocket struct {
	references.IVIC
	label       string
	parent      references.IComponent
	component   references.IComponent
	connections IVICSocketConnection
	db          references.IDisplayBuffer
	pic         references.IPIC6510
	pla         references.IPlaC64
	quartz      references.IQuartz
	intrId      uint32
}

// NewVICSocket creates and initializes a new VICSocket instance, setting up necessary connections for video interface control.
func NewVICSocket(parent references.IComponent, label string, connections IVICSocketConnection) *VICSocket {
	return &VICSocket{
		IVIC:        nil,
		parent:      parent,
		label:       label,
		connections: connections,
		db:          nil,
		pic:         nil,
		pla:         nil,
		quartz:      nil,
		intrId:      intrIrqVicBit,
	}
}

// Mount initializes the VICSocket by resolving its dependencies and calling Setup on the IVIC component.
func (v *VICSocket) Mount() error {
	var err error
	idIVIC := references.IdIVIC(v.IVIC, v.label, 0)
	if v.IVIC, err = references.ComponentToIVIC(v.parent.GetChildByHardwareId(idIVIC)); err != nil {
		return err
	}
	idPIC := references.IdIPIC6510(v.pic, v.label, 0)
	if v.pic, err = references.ComponentToIPIC6510(v.parent.GetChildByHardwareId(idPIC)); err != nil {
		return err
	}
	idPla := references.IdIPlaC64(v.pla, v.label, 0)
	if v.pla, err = references.ComponentToIPLAc64(v.parent.GetChildByHardwareId(idPla)); err != nil {
		return err
	}
	idQuartz := references.IdIQuartz(v.quartz, v.label, 0)
	if v.quartz, err = references.ComponentToIQuartz(v.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	if err = v.IVIC.Bind(v); err != nil {
		return err
	}
	return nil
}

// Cycle retrieves the current clock cycle count from the associated quartz instance.
func (v *VICSocket) Cycle() uint64 {
	return v.quartz.Cycle()
}

// GetBanks returns the IVICBanks interface, allowing access to methods for reading from VIC memory regions.
func (v *VICSocket) GetBanks() references.IVICBanks {
	return v.pla
}

// IRQTrigger triggers an interrupt request by invoking the PIC's TriggerIRQ method with the stored interrupt ID.
func (v *VICSocket) IRQTrigger() {
	v.pic.TriggerIRQ(v.intrId)
}

// IRQClear clears any pending interrupt request associated with the VIC by invoking ClearIRQ on the programmable interrupt controller.
func (v *VICSocket) IRQClear() {
	v.pic.ClearIRQ(v.intrId)
}

// BALow toggles the BA (Bus Available) line state by triggering the associated RDYLow event with the given boolean value.
func (v *VICSocket) BALow(d bool) {
	v.connections.RDYLowTrigger(d)
}

// AECLow sets the AEC (Address Enable) line state to low or high based on the provided boolean parameter.
func (v *VICSocket) AECLow(d bool) {
	v.connections.AECLowTrigger(d)
}

// LastCycle triggers the LastCycle logic on the connected VICSocketConnection implementation.
func (v *VICSocket) LastCycle() {
	v.connections.LastCycleTrigger()
}

// VBlank triggers the vertical blanking interval event by invoking the VBlankTrigger method on the connected interface.
func (v *VICSocket) VBlank() {
	v.connections.VBlankTrigger()
}
