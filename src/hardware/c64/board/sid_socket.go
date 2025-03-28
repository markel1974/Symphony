package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	references.ISID
	fragFreq int
	rasters  int
}

// NewSIDSocket creates and returns a new instance of SIDSocket with default initialization.
func NewSIDSocket(fragFreq int, rasters int) *SIDSocket {
	c := &SIDSocket{
		ISID:     nil,
		fragFreq: fragFreq,
		rasters:  rasters,
	}
	return c
}

// Mount initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if w.ISID, err = references.ComponentsToISID(cc, label, 0); err != nil {
		return err
	}
	if err = w.ISID.Bind(w, w.fragFreq, w.rasters); err != nil {
		return err
	}
	return nil
}
