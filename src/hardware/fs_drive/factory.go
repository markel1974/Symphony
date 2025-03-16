package fs_drive

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "fs_drive"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*FSDrive)(nil)
	return references.IIecDevice(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewBoard(parent, factory, suffix)
}
