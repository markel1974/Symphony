package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RomSocket represents a component that interfaces with IRomsC64 for managing and loading C64 ROM data.
// It attaches to a parent IComponent and provides labeled identification for hierarchical navigation and management.
// RomSocket includes methods for hardware identification and initialization by mounting and binding to IRomsC64.
type RomSocket struct {
	references.IRomsC64
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRomSocket creates a new instance of RomSocket, initializes it with a parent component and label, and assigns a hardware ID.
func NewRomSocket(parent references.IComponent, label string) *RomSocket {
	s := &RomSocket{
		IRomsC64: nil,
		parent:   parent,
		label:    label,
	}
	s.hwId = references.IdIRomsC64(s.IRomsC64, s.label, 0)
	return s
}

// HardwareId retrieves the hardware identifier of the RomSocket instance.
func (s *RomSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes and binds the RomSocket to the appropriate IRomsC64 component by its hardware ID and returns an error if binding fails.
func (s *RomSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IRomsC64, err = references.ComponentToIRomsC64(s.component); err != nil {
		return err
	}
	if err = s.IRomsC64.Bind(s); err != nil {
		return err
	}
	return nil
}
