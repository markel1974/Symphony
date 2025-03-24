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

// Setup initializes the RomLoaderSocket by mapping components and configuring the associated IROMLoaderC1541.
func (s *RomLoaderSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	s.IROMLoaderC1541, err = references.ComponentsToIROMLoaderC1541(c, 0)
	if err != nil {
		return err
	}
	if err = s.IROMLoaderC1541.Setup(s, cfg); err != nil {
		return err
	}
	return nil
}

// Connect establishes a connection by invoking the Connect method of the IROMLoaderC1541 interface. Returns an error if unsuccessful.
func (s *RomLoaderSocket) Connect() error {
	if err := s.IROMLoaderC1541.Connect(); err != nil {
		return err
	}
	return nil
}
