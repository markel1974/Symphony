package cartridges_c64

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "cartridges_c64"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*Manager)(nil)
	return references.ICartridgeManagerC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewManager(parent, factory, label)
}
