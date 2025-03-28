package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket represents a ROM loader socket for facilitating loading and management of ROMs in the system.
// It implements the IROMLoaderC1541 interface to configure and load ROM data for a C1541 drive simulation.
// This type integrates ROM-specific behaviors and configuration through the embedded interface.
type RomLoaderSocket struct {
	references.IROMLoaderC1541 // Incorpora l'interfaccia
}

// NewRomLoaderSocket creates and initializes a new instance of RomLoaderSocket with its IROMLoaderC1541 interface set to nil.
func NewRomLoaderSocket() *RomLoaderSocket {
	return &RomLoaderSocket{
		IROMLoaderC1541: nil,
	}
}

// Mount initializes the RomLoaderSocket by mapping components and configuring the associated IROMLoaderC1541.
func (s *RomLoaderSocket) Mount(c map[string]references.IComponent, _ *config.Config, label string) error {
	var err error
	s.IROMLoaderC1541, err = references.ComponentsToIROMLoaderC1541(c, label, 0)
	if err != nil {
		return err
	}
	if err = s.IROMLoaderC1541.Bind(s); err != nil {
		return err
	}
	return nil
}
