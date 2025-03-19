package easyflash

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "easyFlash"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*CartridgeEasyFlash)(nil)
	return references.ICartridgeC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewEasyFlash(parent, factory, label)
}
