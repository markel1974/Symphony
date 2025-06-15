package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket represents a ROM loader socket for facilitating loading and management of ROMs in the system.
// It implements the IC1541Roms interface to configure and load ROM data for a C1541 drive simulation.
// This type integrates ROM-specific behaviors and configuration through the embedded interface.
type RomLoaderSocket struct {
	references.IC1541Roms
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRomLoaderSocket creates and initializes a new instance of RomLoaderSocket with its IC1541Roms interface set to nil.
func NewRomLoaderSocket(parent references.IComponent, label string) *RomLoaderSocket {
	s := &RomLoaderSocket{
		IC1541Roms: nil,
		parent:     parent,
		label:      label,
	}
	s.hwId = references.IdIC1541Roms(s.IC1541Roms, s.label, 0)
	return s
}

func (s *RomLoaderSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the RomLoaderSocket by mapping components and configuring the associated IC1541Roms.
func (s *RomLoaderSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IC1541Roms, err = references.ComponentToIC1541Roms(s.component); err != nil {
		return err
	}
	if err = s.IC1541Roms.Bind(s); err != nil {
		return err
	}
	return nil
}
