package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a hardware abstraction for the Programmable Logic Array (PLA) in a 1541 drive emulation.
// It integrates VIA components, ROM loading, and configuration management for disk drive emulation functionality.
type PLASocket struct {
	references.IPLAc1541
	via1      references.IVIA
	via2      references.IVIA
	romLoader references.IROMLoaderC1541
	cfg       *config.Config
}

// NewPLASocket creates and returns a new instance of PLASocket with initial fields set to nil.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPLAc1541: nil,
	}
	return c
}

// Mount initializes the PLASocket by resolving its components from the given map and applying configuration settings.
func (w *PLASocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	w.via1, err = references.ComponentsToIVIA(cc, label, 0)
	if err != nil {
		return err
	}
	w.via2, err = references.ComponentsToIVIA(cc, label, 1)
	if err != nil {
		return err
	}
	w.romLoader, err = references.ComponentsToIROMLoaderC1541(cc, label, 0)
	if err != nil {
		return err
	}
	w.IPLAc1541, err = references.ComponentsToIPLAc1541(cc, label, 0)
	if err != nil {
		return err
	}
	if err = w.IPLAc1541.Bind(w, w.via1, w.via2, w.romLoader); err != nil {
		return err
	}
	return nil
}
