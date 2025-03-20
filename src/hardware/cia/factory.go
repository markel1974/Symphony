package mos6526

import (
	"github.com/markel1974/c64emu/src/references"
)

func Identifier() string {
	return "mos6526"
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
	return NewCIA(parent, factory, label)
}
