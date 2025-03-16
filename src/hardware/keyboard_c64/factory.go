package keyboard_c64

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "keyboard_c64"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Keyboard)(nil)
	return references.IKeyboard(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, suffix string) references.IComponent {
	return NewKeyboard(parent, factory, suffix)
}
