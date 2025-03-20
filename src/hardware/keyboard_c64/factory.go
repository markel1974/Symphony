package keyboard_c64

import (
	"github.com/markel1974/c64emu/src/references"
)

func Identifier() string {
	return "keyboard_c64"
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
	z := (*Keyboard)(nil)
	return references.IKeyboard(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewKeyboard(parent, factory, label)
}
