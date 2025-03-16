package mos6510

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "mos6510"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*CPU)(nil)
	return references.I6510(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewCPU(parent, factory, suffix)
}
