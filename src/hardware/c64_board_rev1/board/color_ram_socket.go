package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// ColorRamSocket represents a hardware socket for the Color RAM in the C64 system, capable of interfacing with IC64ColorRam.
// It serves as a bridge between components, managing hardware identification, hierarchy, and interactions.
type ColorRamSocket struct {
	references.IC64ColorRam
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewColorRamSocket creates and initializes a new RamSocket instance with the provided parent component and label.
func NewColorRamSocket(parent references.IComponent, label string) *ColorRamSocket {
	s := &ColorRamSocket{
		IC64ColorRam: nil,
		parent:       parent,
		label:        label,
	}
	s.hwId = references.IdIC64ColorRam(s.IC64ColorRam, s.label, 0)
	return s
}

// HardwareId retrieves the hardware identifier for the ColorRamSocket instance.
func (s *ColorRamSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes and binds the component of the ColorRamSocket to a hardware-defined child component, returning any error.
func (s *ColorRamSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC64ColorRam, err = references.ComponentToIC64ColorRam(s.component); err != nil {
		return err
	}
	if err = s.IC64ColorRam.Bind(s); err != nil {
		return err
	}
	return nil
}
