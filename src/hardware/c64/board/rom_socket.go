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

func (s *RomLoaderSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if s.IROMLoaderC64, err = references.ComponentsToIROMLoaderC64(c, 0); err != nil {
		return err
	}
	if err = s.IROMLoaderC64.Setup(s, cfg); err != nil {
		return err
	}
	return nil
}

func (s *RomLoaderSocket) Connect() error {
	if err := s.IROMLoaderC64.Connect(); err != nil {
		return err
	}
	return nil
}
