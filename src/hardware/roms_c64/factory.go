package roms_c64

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "roms_c64"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*RomLoader)(nil)
	return references.IROMLoaderC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewRomLoader(parent, factory, label)
}
