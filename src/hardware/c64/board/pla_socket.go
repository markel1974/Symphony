package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents an implementation of the IPlaC64 interface with connections to various C64 components.
// It integrates major C64 subsystems including the VIC, SID, CIAs, cartridge manager, and ROM loader.
type PLASocket struct {
	references.IPlaC64
	vic     references.IVIC
	sid     references.ISID
	cia1    references.ICIA
	cia2    references.ICIA
	cartMan references.ICartridgeManagerC64
	roms    references.IROMLoaderC64
}

// NewPLASocket initializes and returns a new instance of PLASocket with default nil values for its components.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPlaC64: nil,
	}
	return c
}

// Mount initializes the PLASocket by resolving dependencies and setting up its components. Returns an error if any failure occurs.
func (w *PLASocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if w.IPlaC64, err = references.ComponentsToIPLAc64(cc, label, 0); err != nil {
		return err
	}
	if w.vic, err = references.ComponentsToIVIC(cc, label, 0); err != nil {
		return err
	}
	if w.sid, err = references.ComponentsToISID(cc, label, 0); err != nil {
		return err
	}
	if w.cia1, err = references.ComponentsToICIA(cc, label, 0); err != nil {
		return err
	}
	if w.cia2, err = references.ComponentsToICIA(cc, label, 1); err != nil {
		return err
	}
	if w.cartMan, err = references.ComponentsToICartridgeManagerC64(cc, label, 0); err != nil {
		return err
	}
	if w.roms, err = references.ComponentsToIROMLoaderC64(cc, label, 0); err != nil {
		return err
	}
	if err = w.IPlaC64.Bind(w, w.vic, w.sid, w.cia1, w.cia2, w.cartMan, w.roms); err != nil {
		return err
	}
	return nil
}
