package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	references.ISID
	label     string
	parent    references.IComponent
	component references.IComponent
	fragFreq  int
	rasters   int
	hwId      string
}

// NewSIDSocket creates and returns a new instance of SIDSocket with default initialization.
func NewSIDSocket(parent references.IComponent, label string, fragFreq int, rasters int) *SIDSocket {
	c := &SIDSocket{
		ISID:     nil,
		parent:   parent,
		label:    label,
		fragFreq: fragFreq,
		rasters:  rasters,
	}
	c.hwId = references.IdISID(c.ISID, c.label, 0)
	return c
}

func (w *SIDSocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.ISID, err = references.ComponentToISID(w.component); err != nil {
		return err
	}
	if err = w.ISID.Bind(w, w.fragFreq, w.rasters); err != nil {
		return err
	}
	return nil
}
