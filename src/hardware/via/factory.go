package mos6522

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "mos6522"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*VIA)(nil)
	return references.IVIA(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewVIA(parent, factory, suffix)
}
