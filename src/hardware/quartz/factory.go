package quartz

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "quartz"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Quartz)(nil)
	return references.IQuartz(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewQuartz(parent, factory, label)
}
