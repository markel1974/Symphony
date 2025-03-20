package pic_6510

import (
	"github.com/markel1974/c64emu/src/references"
)

func Identifier() string {
	return "pic_6510"
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
	return NewPIC(parent, factory, label)
}
