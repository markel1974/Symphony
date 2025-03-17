package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// RomLoaderSocket is a struct that embeds the IROMLoaderC1541 interface for managing ROM loading functionalities.
// It connects to and configures a ROM loader using the provided configuration via the Connect method.
type RomLoaderSocket struct {
	references.IROMLoaderC1541 // Incorpora l'interfaccia
}

// NewRomLoaderSocket returns a new instance of RomLoaderSocket with IROMLoaderC1541 interface initialized to nil.
func NewRomLoaderSocket() *RomLoaderSocket {
	return &RomLoaderSocket{
		IROMLoaderC1541: nil,
	}
}

// Connect initializes the RomLoaderSocket by associating it with an IROMLoaderC1541 instance and configuring it using cfg.
// Returns an error if the setup of the IROMLoaderC1541 fails.
func (s *RomLoaderSocket) Connect(rl references.IROMLoaderC1541, cfg *config.Config) error {
	s.IROMLoaderC1541 = rl
	if err := s.IROMLoaderC1541.Setup(cfg); err != nil {
		return err
	}
	return nil
}
