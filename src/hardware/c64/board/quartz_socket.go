package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// QuartzSocket represents a wrapper structure that incorporates the IQuartz interface for managing clock functionalities.
type QuartzSocket struct {
	references.IQuartz // Incorpora l'interfaccia
}

// NewQuartzSocket creates and returns a new instance of QuartzSocket with its IQuartz interface initialized as nil.
func NewQuartzSocket() *QuartzSocket {
	return &QuartzSocket{
		IQuartz: nil,
	}
}

// Setup initializes the QuartzSocket by associating it with IQuartz and applying configuration settings.
func (s *QuartzSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if s.IQuartz, err = references.ComponentsToIQuartz(c, 0); err != nil {
		return err
	}
	if err = s.IQuartz.Setup(s, cfg); err != nil {
		return err
	}
	return nil
}
