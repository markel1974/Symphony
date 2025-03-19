package board

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "c1541"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Board)(nil)
	return references.IIecDevice(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewBoard(parent, factory, label)
}
