package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket represents a ROM loader socket for facilitating loading and management of ROMs in the system.
// It implements the IROMLoaderC1541 interface to configure and load ROM data for a C1541 drive simulation.
// This type integrates ROM-specific behaviors and configuration through the embedded interface.
type RomLoaderSocket struct {
	references.IROMLoaderC1541
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRomLoaderSocket creates and initializes a new instance of RomLoaderSocket with its IROMLoaderC1541 interface set to nil.
func NewRomLoaderSocket(parent references.IComponent, label string) *RomLoaderSocket {
	s := &RomLoaderSocket{
		IROMLoaderC1541: nil,
		parent:          parent,
		label:           label,
	}
	s.hwId = references.IdIROMLoaderC1541(s.IROMLoaderC1541, s.label, 0)
	return s
}

func (s *RomLoaderSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the RomLoaderSocket by mapping components and configuring the associated IROMLoaderC1541.
func (s *RomLoaderSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IROMLoaderC1541, err = references.ComponentToIROMLoaderC1541(s.component); err != nil {
		return err
	}
	if err = s.IROMLoaderC1541.Bind(s); err != nil {
		return err
	}
	return nil
}
