package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RomSocket represents a component that interfaces with IC64Roms for managing and loading C64 ROM data.
// It attaches to a parent IComponent and provides labeled identification for hierarchical navigation and management.
// RomSocket includes methods for hardware identification and initialization by mounting and binding to IC64Roms.
type RomSocket struct {
	references.IC64Roms
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRomSocket creates a new instance of RomSocket, initializes it with a parent component and label, and assigns a hardware ID.
func NewRomSocket(parent references.IComponent, label string) *RomSocket {
	s := &RomSocket{
		IC64Roms: nil,
		parent:   parent,
		label:    label,
	}
	s.hwId = references.IdIC64Roms(s.IC64Roms, s.label, 0)
	return s
}

// HardwareId retrieves the hardware identifier of the RomSocket instance.
func (s *RomSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes and binds the RomSocket to the appropriate IC64Roms component by its hardware ID and returns an error if binding fails.
func (s *RomSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC64Roms, err = references.ComponentToIC64Roms(s.component); err != nil {
		return err
	}
	if err = s.IC64Roms.Bind(s); err != nil {
		return err
	}
	return nil
}
