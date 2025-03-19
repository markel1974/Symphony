package pic_6510

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "pic_6510"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Pic)(nil)
	return references.IPIC6510(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewPIC(parent, factory, label)
}
