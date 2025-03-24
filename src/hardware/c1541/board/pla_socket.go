package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PLASocket represents a hardware abstraction for the Programmable Logic Array (PLA) in a 1541 drive emulation.
// It integrates VIA components, ROM loading, and configuration management for disk drive emulation functionality.
type PLASocket struct {
	references.IPLAc1541
	via1 references.IVIA
	via2 references.IVIA
	roms references.IROMLoaderC1541
	cfg  *config.Config
}

// NewPLASocket creates and returns a new instance of PLASocket with initial fields set to nil.
func NewPLASocket() *PLASocket {
	c := &PLASocket{
		IPLAc1541: nil,
	}
	return c
}

// Setup initializes the PLASocket by resolving its components from the given map and applying configuration settings.
func (w *PLASocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	w.IPLAc1541, err = references.ComponentsToIPLAc1541(c, 0)
	if err != nil {
		return err
	}
	w.via1, err = references.ComponentsToIVIA(c, 0)
	if err != nil {
		return err
	}
	w.via2, err = references.ComponentsToIVIA(c, 1)
	if err != nil {
		return err
	}
	w.roms, err = references.ComponentsToIROMLoaderC1541(c, 0)
	if err != nil {
		return err
	}
	if err = w.IPLAc1541.Setup(w, w.cfg); err != nil {
		return err
	}
	return nil
}

// Connect establishes the necessary connections between the PLA, VIA components, and ROM loader for the PLASocket.
func (w *PLASocket) Connect() error {
	if err := w.IPLAc1541.Connect(w.via1, w.via2, w.roms); err != nil {
		return err
	}
	return nil
}
