package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a hardware abstraction for the Programmable Logic Array (PLA) in a 1541 drive emulation.
// It integrates VIA components, ROM loading, and configuration management for disk drive emulation functionality.
type PLASocket struct {
	references.IC1541Pla
	label     string
	parent    references.IComponent
	component references.IComponent
	via1      references.IMos6522
	via2      references.IMos6522
	ram       references.IC1541Ram
	romLoader references.IC1541Roms
	cfg       *config.Config
	hwId      string
}

// NewPLASocket creates and returns a new instance of PLASocket with initial fields set to nil.
func NewPLASocket(parent references.IComponent, label string) *PLASocket {
	c := &PLASocket{
		IC1541Pla: nil,
		parent:    parent,
		label:     label,
	}
	c.hwId = references.IdIC1541Pla(c.IC1541Pla, c.label, 0)
	return c
}

func (w *PLASocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the PLASocket by resolving its components from the given map and applying configuration settings.
func (w *PLASocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IC1541Pla, err = references.ComponentToIC1541Pla(w.component); err != nil {
		return err
	}
	idVIA1 := references.IdIMos6522(w.via1, w.label, 0)
	if w.via1, err = references.ComponentToIMos6522(w.parent.GetChildByHardwareId(idVIA1)); err != nil {
		return err
	}
	idVIA2 := references.IdIMos6522(w.via2, w.label, 1)
	if w.via2, err = references.ComponentToIMos6522(w.parent.GetChildByHardwareId(idVIA2)); err != nil {
		return err
	}
	idRam := references.IdIC1541Ram(w.ram, w.label, 0)
	if w.ram, err = references.ComponentToIC1541Ram(w.parent.GetChildByHardwareId(idRam)); err != nil {
		return err
	}
	idRomLoader := references.IdIC1541Roms(w.romLoader, w.label, 0)
	if w.romLoader, err = references.ComponentToIC1541Roms(w.parent.GetChildByHardwareId(idRomLoader)); err != nil {
		return err
	}
	if err = w.IC1541Pla.Bind(w, w.via1, w.via2, w.ram, w.romLoader); err != nil {
		return err
	}
	return nil
}
