package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RamSocket represents a memory socket for the C64 system, managing its connection, binding, and hardware Id functionality.
type RamSocket struct {
	references.IC64Ram
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRamSocket initializes a new RamSocket instance with a parent component and label, generating its hardware Id.
func NewRamSocket(parent references.IComponent, label string) *RamSocket {
	s := &RamSocket{
		IC64Ram: nil,
		parent:  parent,
		label:   label,
	}
	s.hwId = references.IdIC64Ram(s.IC64Ram, s.label, 0)
	return s
}

// HardwareId returns the unique hardware identifier (hwId) of the RamSocket instance.
func (s *RamSocket) HardwareId() string {
	return s.hwId
}

// Wire associates the RamSocket with its parent component and binds it to an IC64Ram instance, returning an error if unsuccessful.
func (s *RamSocket) Wire() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC64Ram, err = references.ComponentToIC64Ram(s.component); err != nil {
		return err
	}
	if err = s.IC64Ram.Bind(s); err != nil {
		return err
	}
	return nil
}
