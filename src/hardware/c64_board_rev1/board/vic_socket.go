package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// The VIC data path never passes through the PLA
// The VIC-II accesses video RAM, Color RAM, and Char ROM directly through its own address and data bus,
// without data ever needing to transit through the PLA.

// IVICSocketConnection represents an interface for managing triggers associated with a VIC chip in an emulation environment.
// LastCycleTrigger signals operations to perform at the last cycle of an emulation frame.
// VBlankTrigger triggers the start of a vertical blanking interval.
// RDYLowTrigger sets the RDY (Ready) signal line low or high based on the parameter.
// AECLowTrigger sets the AEC (Address Enable) signal line low or high based on the parameter.
// IRQTrigger generates an interrupt request and passes its designated value.
// IRQClearTrigger clears an interrupt request with the specified value.
type IVICSocketConnection interface {
	LastCycleTrigger()

	VBlankTrigger()

	RDYLowTrigger(v bool)

	AECLowTrigger(v bool)

	IRQTrigger(d uint32)

	IRQClearTrigger(d uint32)
}

// VICSocket encapsulates the connections and components required for VIC chip socket emulation in the C64 system.
// It integrates with display, quartz clocking, PLA memory, and programmable interrupt controller components.
// Provides methods and triggers for managing display rendering and interaction with the C64 hardware components.
type VICSocket struct {
	references.IMos6569
	label        string
	parent       references.IComponent
	component    references.IComponent
	connections  IVICSocketConnection
	db           references.IDisplayBuffer
	ram          references.IC64Ram
	colorRam     references.IC64ColorRam
	rom          references.IC64Roms
	ramRead      func(addr uint16) uint8
	ramReadColor func(addr uint16) uint8
	romCharRead  func(addr uint16) uint8
	quartz       references.IQuartz
	intrId       uint32
	hwId         string
}

// NewVICSocket creates and initializes a new VICSocket instance, setting up necessary connections for video interface control.
func NewVICSocket(parent references.IComponent, label string, connections IVICSocketConnection) *VICSocket {
	v := &VICSocket{
		IMos6569:    nil,
		parent:      parent,
		label:       label,
		connections: connections,
		db:          nil,
		quartz:      nil,
		intrId:      intrIrqVicBit,
	}
	v.hwId = references.IdIMos6569(v, label, 0)
	return v
}

func (v *VICSocket) HardwareId() string {
	return v.hwId
}

// Wire initializes the VICSocket by resolving its dependencies and calling Setup on the IMos6569 component.
func (v *VICSocket) Wire() error {
	var err error
	v.component = v.parent.GetChildByHardwareId(v.HardwareId())
	if v.IMos6569, err = references.ComponentToIMos6569(v.component); err != nil {
		return err
	}
	idRam := references.IdIC64Ram(v.ram, v.label, 0)
	if v.ram, err = references.ComponentToIC64Ram(v.parent.GetChildByHardwareId(idRam)); err != nil {
		return err
	}
	idColorRam := references.IdIC64ColorRam(v.colorRam, v.label, 0)
	if v.colorRam, err = references.ComponentToIC64ColorRam(v.parent.GetChildByHardwareId(idColorRam)); err != nil {
		return err
	}
	idRom := references.IdIC64Roms(v.rom, v.label, 0)
	if v.rom, err = references.ComponentToIC64Roms(v.parent.GetChildByHardwareId(idRom)); err != nil {
		return err
	}
	idQuartz := references.IdIQuartz(v.quartz, v.label, 0)
	if v.quartz, err = references.ComponentToIQuartz(v.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	if err = v.IMos6569.Bind(v); err != nil {
		return err
	}
	v.ramRead = v.ram.Read
	v.ramReadColor = v.colorRam.Read
	v.romCharRead = v.rom.CharRead
	return nil
}

// Cycle retrieves the current clock cycle count from the associated quartz instance.
func (v *VICSocket) Cycle() uint64 {
	return v.quartz.Cycle()
}

func (v *VICSocket) ReadRam(addr uint16) uint8 {
	return v.ramRead(addr)
}

func (v *VICSocket) ReadColorRam(addr uint16) uint8 {
	return v.ramReadColor(addr)
}

func (v *VICSocket) ReadCharRom(addr uint16) uint8 {
	return v.romCharRead(addr)
}

// IRQTrigger triggers an interrupt request by invoking the PIC's TriggerIRQ method with the stored interrupt ID.
func (v *VICSocket) IRQTrigger() {
	v.connections.IRQTrigger(v.intrId)
}

// IRQClearTrigger clears any pending interrupt request associated with the VIC by invoking ClearIRQ on the programmable interrupt controller.
func (v *VICSocket) IRQClearTrigger() {
	v.connections.IRQClearTrigger(v.intrId)
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
