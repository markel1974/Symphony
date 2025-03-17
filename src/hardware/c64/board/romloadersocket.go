package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type RomLoaderSocket struct {
	references.IROMLoaderC64 // Incorpora l'interfaccia
}

func NewRomLoaderSocket() *RomLoaderSocket {
	return &RomLoaderSocket{
		IROMLoaderC64: nil,
	}
}

func (s *RomLoaderSocket) Connect(rl references.IROMLoaderC64, cfg *config.Config) error {
	s.IROMLoaderC64 = rl
	if err := s.IROMLoaderC64.Setup(cfg); err != nil {
		return err
	}
	return nil
}
