package board

import (
	"github.com/markel1974/symphony/src/references"
)

// SIDSocket represents a socket connected to a Board for managing or interacting with its state or functionality.
type SIDSocket struct {
	references.IMos6581
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
		IMos6581: nil,
		parent:   parent,
		label:    label,
		fragFreq: fragFreq,
		rasters:  rasters,
	}
	c.hwId = references.IdIMos6581(c.IMos6581, c.label, 0)
	return c
}

func (w *SIDSocket) HardwareId() string {
	return w.hwId
}

// Wire initializes the SIDSocket with the provided Board instance, assigning it to the internal board field.
func (w *SIDSocket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6581, err = references.ComponentToIMos6581(w.component); err != nil {
		return err
	}
	if err = w.IMos6581.Bind(w, w.fragFreq, w.rasters); err != nil {
		return err
	}
	return nil
}
