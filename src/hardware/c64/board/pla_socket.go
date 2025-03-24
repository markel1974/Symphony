package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a socket encapsulating a PLA and its associated board.
type PLASocket struct {
	references.IPlaC64
	vic     references.IVIC
	sid     references.ISID
	cia1    references.ICIA
	cia2    references.ICIA
	cartMan references.ICartridgeManagerC64
	roms    references.IROMLoaderC64
	cfg     *config.Config
}

// NewPLASocket creates a new instance of PLASocket with default uninitialized board and PLA components.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPlaC64: nil,
	}
	return c
}

// Setup initializes the PLASocket with the provided board, PLA, and associated components like VIC, SID, CIAs, and cartridge manager.
func (w *PLASocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	w.cfg = cfg
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
	return nil
}

func (w *PLASocket) Connect() error {
	if err := w.IPlaC64.Setup(w.vic, w.sid, w.cia1, w.cia2, w.cartMan, w.roms, w.cfg); err != nil {
		return err
	}
	return nil
}
