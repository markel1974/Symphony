package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket is a type embedding the IROMLoaderC64 interface, designed for managing ROM loader functionalities.
type RomLoaderSocket struct {
	references.IROMLoaderC64
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewRomLoaderSocket initializes and returns a new instance of RomLoaderSocket structure.
func NewRomLoaderSocket(parent references.IComponent, label string) *RomLoaderSocket {
	s := &RomLoaderSocket{
		IROMLoaderC64: nil,
		parent:        parent,
		label:         label,
	}
	s.hwId = references.IdIROMLoaderC64(s.IROMLoaderC64, s.label, 0)
	return s
}

func (s *RomLoaderSocket) HardwareId() string {
	return s.hwId
}

// Mount initializes the RomLoaderSocket by setting up its IROMLoaderC64 interface and applying the provided configuration.
func (s *RomLoaderSocket) Mount() error {
	var err error
	s.component = s.parent.GetChildByHardwareId(s.HardwareId())
	if s.IROMLoaderC64, err = references.ComponentToIROMLoaderC64(s.component); err != nil {
		return err
	}
	if err = s.IROMLoaderC64.Bind(s); err != nil {
		return err
	}
	return nil
}
