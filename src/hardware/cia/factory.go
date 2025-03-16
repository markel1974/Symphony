package mos6526

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "mos6526"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*CIA)(nil)
	return references.ICIA(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewCIA(parent, factory, suffix)
}
