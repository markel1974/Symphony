package board

import (
	"github.com/markel1974/c64emu/src/references"
)

func Identifier() string {
	return "vic20"
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return Identifier()
}

func (t *Factory) Kind() interface{} {
	z := (*Board)(nil)
	return references.IBoard(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewBoard(parent, factory, label)
}
