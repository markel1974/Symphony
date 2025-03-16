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
	return references.IPic6510(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewPic(parent, factory, suffix)
}
