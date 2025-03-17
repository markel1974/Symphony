package board

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
	"log"
	"os"
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

func (cs *CartridgeSocket) Initialize() error {
	for _, cartName := range cs.cfg.GetCartridges() {
		var data []uint8
		if len(cartName.Path) > 0 {
			var err error
			if data, err = os.ReadFile(cartName.Path); err != nil {
				log.Printf("can't add cartridge: %s", err.Error())
				return err
			}
		}
		cartId, err := cs.Add(cartName.Kind, cartName.Path, data)
		if err != nil {
			return fmt.Errorf("can't add cartridge: %s", err.Error())
		}
		log.Printf("cartridge: %s [%s] successfully added", cartName, cartId)
	}
	return nil
}
