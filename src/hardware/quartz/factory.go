package quartz

import (
	"github.com/markel1974/c64emu/src/references"
)

type Factory struct {
}

func Identifier() string {
	return "quartz"
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return Identifier()
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewQuartz(parent, factory, label)
}
