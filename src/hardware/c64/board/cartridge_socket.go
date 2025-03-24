package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type CartridgeSocket struct {
	references.ICartridgeManagerC64
	exp references.IExpansionC64
	cfg *config.Config
}

func NewCartridgeSocket(exp references.IExpansionC64) *CartridgeSocket {
	return &CartridgeSocket{
		exp: exp,
	}
}

func (cs *CartridgeSocket) Setup(c map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if cs.ICartridgeManagerC64, err = references.ComponentsToICartridgeManagerC64(c, 0); err != nil {
		return err
	}
	cs.cfg = cfg
	if err = cs.ICartridgeManagerC64.Setup(cs.exp, cfg); err != nil {
		return err
	}
	return nil
}

func (cs *CartridgeSocket) Connect() error {
	return nil
}
