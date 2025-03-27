package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// PICSocket provides an implementation that integrates the IPIC6510 programmable interrupt controller and IQuartz clock system.
type PICSocket struct {
	references.IPIC6510
	quartz references.IQuartz
}

// NewPICSocket creates and returns a new instance of PICSocket with uninitialized IPIC6510 and quartz properties.
func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
	}
}

// Setup initializes the PICSocket by resolving and configuring referenced components. Returns error if setup fails.
func (s *PICSocket) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if s.IPIC6510, err = references.ComponentsToIPIC6510(cc, 0); err != nil {
		return err
	}
	if s.quartz, err = references.ComponentsToIQuartz(cc, 0); err != nil {
		return err
	}
	if err = s.IPIC6510.Setup(s, cfg, s.quartz); err != nil {
		return err
	}
	return nil
}
