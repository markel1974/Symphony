package magicdesk

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

func Identifier() string {
	return "magic_desk"
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
	z := (*CartridgeMagicDesk)(nil)
	return references.ICartridgeC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewMagicDesk(parent, factory, label)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
