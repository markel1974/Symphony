package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents an implementation of the IPlaC64 interface with connections to various C64 components.
// It integrates major C64 subsystems including the VIC, SID, CIAs, cartridge manager, and ROM loader.
type PLASocket struct {
	references.IPlaC64
	label     string
	parent    references.IComponent
	component references.IComponent
	vic       references.IVIC
	sid       references.ISID
	cia1      references.ICIA
	cia2      references.ICIA
	cartMan   references.ICartridgeManagerC64
	ram       references.IRamC64
	roms      references.IROMLoaderC64
	hwId      string
}

// NewPLASocket initializes and returns a new instance of PLASocket with default nil values for its components.
func NewPLASocket(parent references.IComponent, label string) *PLASocket {
	c := &PLASocket{
		IPlaC64: nil,
		parent:  parent,
		label:   label,
	}
	c.hwId = references.IdIPlaC64(c.IPlaC64, c.label, 0)
	return c
}

func (w *PLASocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the PLASocket by resolving dependencies and setting up its components. Returns an error if any failure occurs.
func (w *PLASocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IPlaC64, err = references.ComponentToIPLAc64(w.component); err != nil {
		return err
	}
	idVIC := references.IdIVIC(w.vic, w.label, 0)
	if w.vic, err = references.ComponentToIVIC(w.parent.GetChildByHardwareId(idVIC)); err != nil {
		return err
	}
	idSID := references.IdISID(w.sid, w.label, 0)
	if w.sid, err = references.ComponentToISID(w.parent.GetChildByHardwareId(idSID)); err != nil {
		return err
	}
	idCIA1 := references.IdICIA(w.cia1, w.label, 0)
	if w.cia1, err = references.ComponentToICIA(w.parent.GetChildByHardwareId(idCIA1)); err != nil {
		return err
	}
	idCIA2 := references.IdICIA(w.cia2, w.label, 1)
	if w.cia2, err = references.ComponentToICIA(w.parent.GetChildByHardwareId(idCIA2)); err != nil {
		return err
	}
	idCartridgeManager := references.IdICartridgeManagerC64(w.cartMan, w.label, 0)
	if w.cartMan, err = references.ComponentToICartridgeManagerC64(w.parent.GetChildByHardwareId(idCartridgeManager)); err != nil {
		return err
	}
	idRam := references.IdIRamC64(w.ram, w.label, 0)
	if w.ram, err = references.ComponentToIRamC64(w.parent.GetChildByHardwareId(idRam)); err != nil {
		return err
	}
	idRomLoader := references.IdIROMLoaderC64(w.roms, w.label, 0)
	if w.roms, err = references.ComponentToIROMLoaderC64(w.parent.GetChildByHardwareId(idRomLoader)); err != nil {
		return err
	}
	if err = w.IPlaC64.Bind(w, w.vic, w.sid, w.cia1, w.cia2, w.cartMan, w.ram, w.roms); err != nil {
		return err
	}
	return nil
}
