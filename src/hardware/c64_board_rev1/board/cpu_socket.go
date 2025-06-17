package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a connection hub for the 6510 CPU, PIC, and PLA components to integrate and interact cohesively.
type CPUSocket struct {
	references.IMos6510
	label     string
	parent    references.IComponent
	component references.IComponent
	pic       references.IMos6510Pic
	banks     references.IC64Pla
	hwId      string
}

// NewCPUSocket creates and returns a new instance of CPUSocket with its internal references uninitialized.
func NewCPUSocket(parent references.IComponent, label string) *CPUSocket {
	c := &CPUSocket{
		IMos6510: nil,
		parent:   parent,
		label:    label,
	}
	c.hwId = references.IdIMos6510(c.IMos6510, c.label, 0)
	return c
}

func (w *CPUSocket) HardwareId() string {
	return w.hwId
}

// Wire initializes the CPUSocket with the provided CPU, PIC, and PLA, and sets up the CPU for interaction.
func (w *CPUSocket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6510, err = references.ComponentToIMos6510(w.component); err != nil {
		return err
	}
	picId := references.IdIMos6510Pic(w.pic, w.label, 0)
	if w.pic, err = references.ComponentToIMos6510Pic(w.parent.GetChildByHardwareId(picId)); err != nil {
		return err
	}
	plaId := references.IdIC64Pla(w.banks, w.label, 0)
	if w.banks, err = references.ComponentToIC64Pla(w.parent.GetChildByHardwareId(plaId)); err != nil {
		return err
	}
	if err = w.IMos6510.Bind(w, w.pic, w.banks); err != nil {
		return err
	}
	return nil
}
