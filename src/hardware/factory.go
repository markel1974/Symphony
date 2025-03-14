package hardware

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	mos6526 "github.com/markel1974/c64emu/src/hardware/cia"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(parent component.IComponent, id string, suffix string) (component.IComponent, error) {
	switch id {
	case "cia":
		return mos6526.NewCIA(parent, suffix), nil
	default:
		return nil, fmt.Errorf("unknown component %s", id)
	}
}
