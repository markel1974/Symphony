package fs_drive

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

func Identifier() string {
	return "fs_drive"
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return Identifier()
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewBoard(parent, factory, label)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
