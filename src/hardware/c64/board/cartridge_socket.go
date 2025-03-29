package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// CartridgeSocket represents a socket for handling interactions with a cartridge, expansion, and associated configuration.
type CartridgeSocket struct {
	references.ICartridgeManagerC64
	label     string
	parent    references.IComponent
	component references.IComponent
	expansion references.IExpansionC64
	hwId      string
}

// NewCartridgeSocket creates and returns a new CartridgeSocket initialized with the provided expansion interface.
func NewCartridgeSocket(parent references.IComponent, label string, expansion references.IExpansionC64) *CartridgeSocket {
	cs := &CartridgeSocket{
		parent:    parent,
		label:     label,
		expansion: expansion,
	}
	cs.hwId = references.IdICartridgeManagerC64(cs.ICartridgeManagerC64, cs.label, 0)
	return cs
}

func (cs *CartridgeSocket) HardwareId() string {
	return cs.hwId
}

// Mount initializes the CartridgeSocket and its associated ICartridgeManagerC64 instance with provided components and config.
func (cs *CartridgeSocket) Mount() error {
	var err error
	cs.component = cs.parent.GetChildByHardwareId(cs.HardwareId())
	if cs.ICartridgeManagerC64, err = references.ComponentToICartridgeManagerC64(cs.component); err != nil {
		return err
	}
	//TODO expansion an IComponent
	if err = cs.ICartridgeManagerC64.Bind(cs, cs.expansion); err != nil {
		return err
	}
	return nil
}
