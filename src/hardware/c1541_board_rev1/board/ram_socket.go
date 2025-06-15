package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RamSocket provides a structure for managing a socket interface bound to a RAM module compliant with IC1541Ram.
// It includes properties for identification, hierarchical parent/child relationships, and component association.
type RamSocket struct {
	references.IC1541Ram
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRamSocket creates and initializes a new RamSocket instance with a parent component and a label.
func NewRamSocket(parent references.IComponent, label string) *RamSocket {
	s := &RamSocket{
		IC1541Ram: nil,
		parent:    parent,
		label:     label,
	}
	s.hwId = references.IdIC1541Ram(s.IC1541Ram, s.label, 0)
	return s
}

// HardwareId returns the unique hardware ID of the RamSocket as a string.
func (s *RamSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes and binds the RAM socket to its corresponding RAM module, returning an error if the process fails.
func (s *RamSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC1541Ram, err = references.ComponentToIC1541Ram(s.component); err != nil {
		return err
	}
	if err = s.IC1541Ram.Bind(s); err != nil {
		return err
	}
	return nil
}
