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

// Setup initializes the PLASocket by resolving dependencies and setting up its components. Returns an error if any failure occurs.
func (w *PLASocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if w.IPlaC64, err = references.ComponentsToIPLAc64(c, 0); err != nil {
		return err
	}
	if w.vic, err = references.ComponentsToIVIC(c, 0); err != nil {
		return err
	}
	if w.sid, err = references.ComponentsToISID(c, 0); err != nil {
		return err
	}
	if w.cia1, err = references.ComponentsToICIA(c, 0); err != nil {
		return err
	}
	if w.cia2, err = references.ComponentsToICIA(c, 1); err != nil {
		return err
	}
	if w.cartMan, err = references.ComponentsToICartridgeManagerC64(c, 0); err != nil {
		return err
	}
	if w.roms, err = references.ComponentsToIROMLoaderC64(c, 0); err != nil {
		return err
	}
	if err = w.IPlaC64.Setup(w, cfg); err != nil {
		return err
	}
	return nil
}

// Connect establishes connections between PLASocket components and initializes the IPlaC64 interface for proper functionality.
func (w *PLASocket) Connect() error {
	if err := w.IPlaC64.Connect(w.vic, w.sid, w.cia1, w.cia2, w.cartMan, w.roms); err != nil {
		return err
	}
	return nil
}
