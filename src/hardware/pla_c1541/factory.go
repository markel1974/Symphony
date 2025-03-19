package pla_c1541

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "pla_c1541"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*PLA)(nil)
	return references.IPLAc1541(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewPLA(parent, factory, label)
}
