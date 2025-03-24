package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type PICSocket struct {
	references.IPIC6510
	quartz references.IQuartz
}

func NewPICSocket() *PICSocket {
	return &PICSocket{
		IPIC6510: nil,
	}
}

func (s *PICSocket) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if s.IPIC6510, err = references.ComponentsToIPIC6510(cc, 0); err != nil {
		return err
	}
	if s.quartz, err = references.ComponentsToIQuartz(cc, 0); err != nil {
		return err
	}
	if err = s.IPIC6510.Setup(s, cfg); err != nil {
		return err
	}
	return nil
}

func (s *PICSocket) Connect() error {
	if err := s.IPIC6510.Connect(s.quartz); err != nil {
		return err
	}
	return nil
}
