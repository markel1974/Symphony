package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type CartridgeSocket struct {
	references.ICartridgeManagerC64
	cfg *config.Config
}

func NewCartridgeSocket() *CartridgeSocket {
	return &CartridgeSocket{}
}

func (cs *CartridgeSocket) Connect(cartMan references.ICartridgeManagerC64, exp references.IExpansionC64, cfg *config.Config) error {
	cs.ICartridgeManagerC64 = cartMan
	cs.cfg = cfg
	if err := cs.ICartridgeManagerC64.Setup(exp, cfg); err != nil {
		return err
	}
	return nil
}
