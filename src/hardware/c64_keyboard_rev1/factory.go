package c64_keyboard_rev1

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

func Identifier() string {
	return "c64_keyboard"
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
	return NewKeyboard(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
