package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeSocket represents a socket for handling interactions with a cartridge, expansion, and associated configuration.
type CartridgeSocket struct {
	references.ICartridgeManagerC64
	expansion references.IExpansionC64
	cfg       *config.Config
}

// NewCartridgeSocket creates and returns a new CartridgeSocket initialized with the provided expansion interface.
func NewCartridgeSocket(expansion references.IExpansionC64) *CartridgeSocket {
	return &CartridgeSocket{
		expansion: expansion,
	}
}

// Setup initializes the CartridgeSocket and its associated ICartridgeManagerC64 instance with provided components and config.
func (cs *CartridgeSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if cs.ICartridgeManagerC64, err = references.ComponentsToICartridgeManagerC64(c, 0); err != nil {
		return err
	}
	cs.cfg = cfg
	if err = cs.ICartridgeManagerC64.Setup(cs, cfg, cs.expansion); err != nil {
		return err
	}
	return nil
}
