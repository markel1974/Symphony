package iec

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "iec"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Dispatcher)(nil)
	return references.IIec(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewDispatcher(parent, factory, label)
}
