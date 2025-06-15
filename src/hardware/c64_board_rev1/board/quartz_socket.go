package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// QuartzSocket represents a wrapper structure that incorporates the IQuartz interface for managing clock functionalities.
type QuartzSocket struct {
	references.IQuartz // Incorpora l'interfaccia
	label              string
	parent             references.IComponent
	component          references.IComponent
	hwId               string
}

// NewQuartzSocket creates and returns a new instance of QuartzSocket with its IQuartz interface initialized as nil.
func NewQuartzSocket(parent references.IComponent, label string) *QuartzSocket {
	s := &QuartzSocket{
		IQuartz: nil,
		parent:  parent,
		label:   label,
	}
	s.hwId = references.IdIQuartz(s.IQuartz, s.label, 0)
	return s
}

func (s *QuartzSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the QuartzSocket by associating it with IQuartz and applying configuration settings.
func (s *QuartzSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IQuartz, err = references.ComponentToIQuartz(s.component); err != nil {
		return err
	}
	if err = s.IQuartz.Bind(s, references.IQuartz1Mhz); err != nil {
		return err
	}
	return nil
}
