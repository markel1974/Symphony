package board

import (
	"github.com/markel1974/c64emu/src/config"
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

func (s *QuartzSocket) Setup(c map[string]references.IComponent, _ *config.Config) error {
	var err error
	if s.IQuartz, err = references.ComponentsToIQuartz(c, 0); err != nil {
		return err
	}
	if err = s.IQuartz.Setup(); err != nil {
		return err
	}
	return nil
}

func (s *QuartzSocket) Connect() error {
	return nil
}
