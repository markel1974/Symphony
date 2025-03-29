package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket is a struct that provides integration between a programmable interrupt controller (PIC) and a quartz clock.
type PICSocket struct {
	references.IPIC6510
	label     string
	parent    references.IComponent
	component references.IComponent
	quartz    references.IQuartz
	hwId      string
}

// NewPICSocket creates and returns a new instance of PICSocket with uninitialized dependencies.
func NewPICSocket(parent references.IComponent, label string) *PICSocket {
	s := &PICSocket{
		IPIC6510: nil,
		parent:   parent,
		label:    label,
		quartz:   nil,
	}
	s.hwId = references.IdIPIC6510(s.IPIC6510, s.label, 0)
	return s
}

func (s *PICSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the PICSocket by configuring its dependencies and setting up the IPIC6510 component with the provided config.
func (s *PICSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IPIC6510, err = references.ComponentToIPIC6510(s.component); err != nil {
		return err
	}
	idQuartz := references.IdIQuartz(s.quartz, s.label, 0)
	if s.quartz, err = references.ComponentToIQuartz(s.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	if err = s.IPIC6510.Bind(s, s.quartz); err != nil {
		return err
	}
	return nil
}
