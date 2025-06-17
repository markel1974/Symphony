package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket provides an implementation that integrates the IMos6510Pic programmable interrupt controller and IQuartz clock system.
type PICSocket struct {
	references.IMos6510Pic
	label     string
	parent    references.IComponent
	component references.IComponent
	quartz    references.IQuartz
	hwId      string
}

// NewPICSocket creates and returns a new instance of PICSocket with uninitialized IMos6510Pic and quartz properties.
func NewPICSocket(parent references.IComponent, label string) *PICSocket {
	s := &PICSocket{
		IMos6510Pic: nil,
		parent:      parent,
		label:       label,
	}
	s.hwId = references.IdIMos6510Pic(s.IMos6510Pic, s.label, 0)
	return s
}

func (s *PICSocket) HardwareId() string {
	return s.hwId
}

// Wire initializes the PICSocket by resolving and configuring referenced components. Returns error if setup fails.
func (s *PICSocket) Wire() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IMos6510Pic, err = references.ComponentToIMos6510Pic(s.component); err != nil {
		return err
	}
	idQuartz := references.IdIQuartz(s.quartz, s.label, 0)
	if s.quartz, err = references.ComponentToIQuartz(s.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	if err = s.IMos6510Pic.Bind(s, s.quartz); err != nil {
		return err
	}
	return nil
}
