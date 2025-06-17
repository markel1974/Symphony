package board

import (
	"github.com/markel1974/c64emu/src/references"
)

type ICPUSocketConnections interface {
	ByteReady() bool
}

// CPUSocket represents a CPU socket managing the integration of a 6510 CPU, programmable interrupt controller, memory banks, and VIA.
type CPUSocket struct {
	references.IMos6510
	label       string
	connections ICPUSocketConnections
	parent      references.IComponent
	component   references.IComponent
	pic         references.IMos6510Pic
	pla         references.IC1541Pla
	via2        references.IMos6522
	hwId        string
}

// NewCPUSocket creates and initializes a new CPUSocket instance with default nil values for its fields.
func NewCPUSocket(parent references.IComponent, label string, connections ICPUSocketConnections) *CPUSocket {
	c := &CPUSocket{
		IMos6510:    nil,
		parent:      parent,
		label:       label,
		connections: connections,
		pic:         nil,
		pla:         nil,
	}
	c.hwId = references.IdIMos6510(c.IMos6510, c.label, 0)
	return c
}

func (w *CPUSocket) HardwareId() string {
	return w.hwId
}

// Wire initializes the CPUSocket by linking required components and configuring dependencies using the provided map and config.
// Returns an error if any component setup or binding fails.
func (w *CPUSocket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IMos6510, err = references.ComponentToIMos6510(w.component); err != nil {
		return err
	}
	idPIC := references.IdIMos6510Pic(w.pic, w.label, 0)
	if w.pic, err = references.ComponentToIMos6510Pic(w.parent.GetChildByHardwareId(idPIC)); err != nil {
		return err
	}
	idPLA := references.IdIC1541Pla(w.pla, w.label, 0)
	if w.pla, err = references.ComponentToIC1541Pla(w.parent.GetChildByHardwareId(idPLA)); err != nil {
		return err
	}
	idVIA2 := references.IdIMos6522(w.via2, w.label, 1)
	if w.via2, err = references.ComponentToIMos6522(w.parent.GetChildByHardwareId(idVIA2)); err != nil {
		return err
	}
	if err = w.IMos6510.Bind(w, w.pic, w.pla); err != nil {
		return err
	}
	//w.IMos6510.SetOverflowBranch(w.via2.ByteReady)
	w.IMos6510.SetOverflowBranch(w.connections.ByteReady)
	return nil
}
