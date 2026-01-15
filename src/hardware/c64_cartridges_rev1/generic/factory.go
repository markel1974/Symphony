package generic

import (
	"github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

func Identifier() string {
	return "generic"
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
	z := (*CartridgeGeneric)(nil)
	return references.IC64Cartridge(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewCartridgeGeneric(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())

	catalog.RegisterType(GetType(), New)
	catalog.RegisterSize(0x2000, New) //cartridge8k.New
	catalog.RegisterSize(0x4000, New) //cartridge16k.New
}
