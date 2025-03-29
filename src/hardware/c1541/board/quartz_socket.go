package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// QuartzSocket represents a wrapper structure that incorporates the IQuartz interface for managing clock functionalities.
type QuartzSocket struct {
	references.IQuartz
	label     string
	parent    references.IComponent
	component references.IComponent
}

// NewQuartzSocket creates and returns a new instance of QuartzSocket with its IQuartz interface initialized as nil.
func NewQuartzSocket(parent references.IComponent, label string) *QuartzSocket {
	return &QuartzSocket{
		IQuartz: nil,
		parent:  parent,
		label:   label,
	}
}

// Mount initializes the QuartzSocket by associating it with IQuartz and applying configuration settings.
func (s *QuartzSocket) Mount() error {
	var err error
	idQuartz := references.IdIQuartz(s.IQuartz, s.label, 0)
	if s.IQuartz, err = references.ComponentToIQuartz(s.parent.GetChildByHardwareId(idQuartz)); err != nil {
		return err
	}
	if err = s.IQuartz.Bind(s); err != nil {
		return err
	}
	return nil
}
