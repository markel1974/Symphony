package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket is a type embedding the IROMLoaderC64 interface, designed for managing ROM loader functionalities.
type RomLoaderSocket struct {
	references.IROMLoaderC64 // Incorpora l'interfaccia
}

// NewRomLoaderSocket initializes and returns a new instance of RomLoaderSocket structure.
func NewRomLoaderSocket() *RomLoaderSocket {
	return &RomLoaderSocket{
		IROMLoaderC64: nil,
	}
}

func (s *RomLoaderSocket) Bind() error {
	if err := s.IROMLoaderC64.Bind(s); err != nil {
		return err
	}
	return nil
}

// Mount initializes the RomLoaderSocket by setting up its IROMLoaderC64 interface and applying the provided configuration.
func (s *RomLoaderSocket) Mount(cc map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	if s.IROMLoaderC64, err = references.ComponentsToIROMLoaderC64(cc, label, 0); err != nil {
		return err
	}
	if err = s.IROMLoaderC64.Bind(s); err != nil {
		return err
	}
	return nil
}
