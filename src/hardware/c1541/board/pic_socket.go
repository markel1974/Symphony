package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket is a struct that provides integration between a programmable interrupt controller (PIC) and a quartz clock.
type PICSocket struct {
	references.IPIC6510
	quartz references.IQuartz
}

// NewPICSocket creates and returns a new instance of PICSocket with uninitialized dependencies.
func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
		quartz:   nil,
	}
}

// Mount initializes the PICSocket by configuring its dependencies and setting up the IPIC6510 component with the provided config.
func (s *PICSocket) Mount(c map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if s.quartz, err = references.ComponentsToIQuartz(c, label, 0); err != nil {
		return err
	}
	if s.IPIC6510, err = references.ComponentsToIPIC6510(c, label, 0); err != nil {
		return err
	}
	if err = s.IPIC6510.Bind(s, s.quartz); err != nil {
		return err
	}
	return nil
}
