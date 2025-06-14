package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a hardware abstraction for the Programmable Logic Array (PLA) in a 1541 drive emulation.
// It integrates VIA components, ROM loading, and configuration management for disk drive emulation functionality.
type PLASocket struct {
	references.IPLAc1541
	label     string
	parent    references.IComponent
	component references.IComponent
	via1      references.IVIA
	via2      references.IVIA
	ram       references.IRamC1541
	romLoader references.IROMLoaderC1541
	cfg       *config.Config
	hwId      string
}

// NewPLASocket creates and returns a new instance of PLASocket with initial fields set to nil.
func NewPLASocket(parent references.IComponent, label string) *PLASocket {
	c := &PLASocket{
		IPLAc1541: nil,
		parent:    parent,
		label:     label,
	}
	c.hwId = references.IdIPLAc1541(c.IPLAc1541, c.label, 0)
	return c
}

func (w *PLASocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the PLASocket by resolving its components from the given map and applying configuration settings.
func (w *PLASocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IPLAc1541, err = references.ComponentToIPLAc1541(w.component); err != nil {
		return err
	}
	idVIA1 := references.IdIVIA(w.via1, w.label, 0)
	if w.via1, err = references.ComponentToIVIA(w.parent.GetChildByHardwareId(idVIA1)); err != nil {
		return err
	}
	idVIA2 := references.IdIVIA(w.via2, w.label, 1)
	if w.via2, err = references.ComponentToIVIA(w.parent.GetChildByHardwareId(idVIA2)); err != nil {
		return err
	}
	idRam := references.IdIRamC1541(w.ram, w.label, 0)
	if w.ram, err = references.ComponentToIRamC1541(w.parent.GetChildByHardwareId(idRam)); err != nil {
		return err
	}
	idRomLoader := references.IdIROMLoaderC1541(w.romLoader, w.label, 0)
	if w.romLoader, err = references.ComponentToIROMLoaderC1541(w.parent.GetChildByHardwareId(idRomLoader)); err != nil {
		return err
	}
	if err = w.IPLAc1541.Bind(w, w.via1, w.via2, w.ram, w.romLoader); err != nil {
		return err
	}
	return nil
}
