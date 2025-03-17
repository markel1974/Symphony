package board

import (
	"github.com/markel1974/c64emu/src/references"
)

type QuartzSocket struct {
	references.IQuartz // Incorpora l'interfaccia
}

func NewQuartzSocket() *QuartzSocket {
	return &QuartzSocket{
		IQuartz: nil,
	}
}

func (s *QuartzSocket) Connect(q references.IQuartz) error {
	s.IQuartz = q
	if err := s.IQuartz.Setup(); err != nil {
		return err
	}
	return nil
}
