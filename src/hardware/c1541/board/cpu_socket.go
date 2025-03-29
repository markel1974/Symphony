package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CPUSocket represents a CPU socket managing the integration of a 6510 CPU, programmable interrupt controller, memory banks, and VIA.
type CPUSocket struct {
	references.I6510
	label     string
	parent    references.IComponent
	component references.IComponent
	pic       references.IPIC6510
	pla       references.IPLAc1541
	via2      references.IVIA
}

// NewCPUSocket creates and initializes a new CPUSocket instance with default nil values for its fields.
func NewCPUSocket(parent references.IComponent, label string) *CPUSocket {
	c := &CPUSocket{
		I6510:  nil,
		parent: parent,
		label:  label,
		pic:    nil,
		pla:    nil,
	}
	return c
}

// Mount initializes the CPUSocket by linking required components and configuring dependencies using the provided map and config.
// Returns an error if any component setup or binding fails.
func (w *CPUSocket) Mount() error {
	var err error
	idPIC := references.IdIPIC6510(w.pic, w.label, 0)
	if w.pic, err = references.ComponentToIPIC6510(w.parent.GetChildByHardwareId(idPIC)); err != nil {
		return err
	}
	idPLA := references.IdIPLAc1541(w.pla, w.label, 0)
	if w.pla, err = references.ComponentToIPLAc1541(w.parent.GetChildByHardwareId(idPLA)); err != nil {
		return err
	}
	idI6510 := references.IdI6510(w.I6510, w.label, 0)
	if w.I6510, err = references.ComponentToI6510(w.parent.GetChildByHardwareId(idI6510)); err != nil {
		return err
	}
	idVIA2 := references.IdIVIA(w.via2, w.label, 1)
	if w.via2, err = references.ComponentToIVIA(w.parent.GetChildByHardwareId(idVIA2)); err != nil {
		return err
	}
	if err = w.I6510.Bind(w, w.pic, w.pla); err != nil {
		return err
	}
	w.I6510.SetOverflowBranch(w.via2.ByteReady)
	return nil
}
