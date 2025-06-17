package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents an implementation of the IC64Pla interface with connections to various C64 components.
// It integrates major C64 subsystems including the VIC, SID, CIAs, cartridge manager, and ROM loader.
type PLASocket struct {
	references.IC64Pla
	label     string
	parent    references.IComponent
	component references.IComponent
	vic       references.IMos6569
	sid       references.IMos6581
	cia1      references.IMos6526
	cia2      references.IMos6526
	cartMan   references.IC64CartridgeManager
	ram       references.IC64Ram
	colorRam  references.IC64ColorRam
	roms      references.IC64Roms
	hwId      string
}

// NewPLASocket initializes and returns a new instance of PLASocket with default nil values for its components.
func NewPLASocket(parent references.IComponent, label string) *PLASocket {
	c := &PLASocket{
		IC64Pla: nil,
		parent:  parent,
		label:   label,
	}
	c.hwId = references.IdIC64Pla(c.IC64Pla, c.label, 0)
	return c
}

func (w *PLASocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the PLASocket by resolving dependencies and setting up its components. Returns an error if any failure occurs.
func (w *PLASocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IC64Pla, err = references.ComponentToIC64Pla(w.component); err != nil {
		return err
	}
	idVIC := references.IdIMos6569(w.vic, w.label, 0)
	if w.vic, err = references.ComponentToIMos6569(w.parent.GetChildByHardwareId(idVIC)); err != nil {
		return err
	}
	idSID := references.IdIMos6581(w.sid, w.label, 0)
	if w.sid, err = references.ComponentToIMos6581(w.parent.GetChildByHardwareId(idSID)); err != nil {
		return err
	}
	idCIA1 := references.IdIMos6526(w.cia1, w.label, 0)
	if w.cia1, err = references.ComponentToIMos6526(w.parent.GetChildByHardwareId(idCIA1)); err != nil {
		return err
	}
	idCIA2 := references.IdIMos6526(w.cia2, w.label, 1)
	if w.cia2, err = references.ComponentToIMos6526(w.parent.GetChildByHardwareId(idCIA2)); err != nil {
		return err
	}
	idCartridgeManager := references.IdIC64CartridgeManager(w.cartMan, w.label, 0)
	if w.cartMan, err = references.ComponentToIC64CartridgeManager(w.parent.GetChildByHardwareId(idCartridgeManager)); err != nil {
		return err
	}
	idRam := references.IdIC64Ram(w.ram, w.label, 0)
	if w.ram, err = references.ComponentToIC64Ram(w.parent.GetChildByHardwareId(idRam)); err != nil {
		return err
	}
	idColorRam := references.IdIC64ColorRam(w.colorRam, w.label, 0)
	if w.colorRam, err = references.ComponentToIC64ColorRam(w.parent.GetChildByHardwareId(idColorRam)); err != nil {
		return err
	}
	idRoms := references.IdIC64Roms(w.roms, w.label, 0)
	if w.roms, err = references.ComponentToIC64Roms(w.parent.GetChildByHardwareId(idRoms)); err != nil {
		return err
	}
	if err = w.IC64Pla.Bind(w, w.vic, w.cartMan, w.ram, w.roms, w.vic, w.sid, w.cia1, w.cia2, w.colorRam); err != nil {
		return err
	}
	return nil
}
