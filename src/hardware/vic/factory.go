package mos6569

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "mos6569"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*VIC)(nil)
	return references.IVic(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewVIC(parent, factory, suffix)
}
