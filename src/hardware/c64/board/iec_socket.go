package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type IECSocket struct {
	references.IIec // Incorpora l'interfaccia
}

func NewIECSocket() *IECSocket {
	return &IECSocket{
		IIec: nil,
	}
}

func (s *IECSocket) Connect(iec references.IIec, quartz references.IQuartz, cfg *config.Config) error {
	s.IIec = iec
	if err := s.IIec.Setup(quartz, cfg); err != nil {
		return err
	}
	return nil
}
