package mos6569

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

func Identifier() string {
	return "mos6569"
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return Identifier()
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewVIC(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
