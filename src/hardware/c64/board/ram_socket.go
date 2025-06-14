package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RamSocket represents a memory socket for the C64 system, managing its connection, binding, and hardware ID functionality.
type RamSocket struct {
	references.IRamC64
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRamSocket initializes a new RamSocket instance with a parent component and label, generating its hardware ID.
func NewRamSocket(parent references.IComponent, label string) *RamSocket {
	s := &RamSocket{
		IRamC64: nil,
		parent:  parent,
		label:   label,
	}
	s.hwId = references.IdIRamC64(s.IRamC64, s.label, 0)
	return s
}

// HardwareId returns the unique hardware identifier (hwId) of the RamSocket instance.
func (s *RamSocket) HardwareId() string {
	return s.hwId
}

// Mount associates the RamSocket with its parent component and binds it to an IRamC64 instance, returning an error if unsuccessful.
func (s *RamSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IRamC64, err = references.ComponentToIRamC64(s.component); err != nil {
		return err
	}
	if err = s.IRamC64.Bind(s); err != nil {
		return err
	}
	return nil
}
