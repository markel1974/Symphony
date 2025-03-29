package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a connection hub for the 6510 CPU, PIC, and PLA components to integrate and interact cohesively.
type CPUSocket struct {
	references.I6510
	label     string
	parent    references.IComponent
	component references.IComponent
	pic       references.IPIC6510
	banks     references.IPlaC64
}

// NewCPUSocket creates and returns a new instance of CPUSocket with its internal references uninitialized.
func NewCPUSocket(parent references.IComponent, label string) *CPUSocket {
	c := &CPUSocket{
		I6510:  nil,
		parent: parent,
		label:  label,
	}
	return c
}

// Mount initializes the CPUSocket with the provided CPU, PIC, and PLA, and sets up the CPU for interaction.
func (w *CPUSocket) Mount() error {
	var err error
	cpuId := references.IdI6510(w.I6510, w.label, 0)
	if w.component = w.parent.GetChildByHardwareId(cpuId); w.component == nil {
		return fmt.Errorf("cpu not found")
	}
	if w.I6510, err = references.ComponentToI6510(w.component); err != nil {
		return err
	}
	picId := references.IdIPIC6510(w.pic, w.label, 0)
	if w.pic, err = references.ComponentToIPIC6510(w.parent.GetChildByHardwareId(picId)); err != nil {
		return err
	}
	plaId := references.IdIPlaC64(w.banks, w.label, 0)
	if w.banks, err = references.ComponentToIPLAc64(w.parent.GetChildByHardwareId(plaId)); err != nil {
		return err
	}
	if err = w.I6510.Bind(w, w.pic, w.banks); err != nil {
		return err
	}
	return nil
}
