package mos6581

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "mos6581"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*SID)(nil)
	return references.ISID(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewSID(parent, factory, label)
}
