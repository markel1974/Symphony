package final_cartridge_iii

import (
	"github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

func Identifier() string {
	return "final_cartridge_iii"
}

func GetType() int {
	return catalog.CartridgeFinalIII
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
	z := (*CartridgeFinalCartridgeIII)(nil)
	return references.IC64Cartridge(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewCartridgeFinalCartridgeIII(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())

	catalog.RegisterType(GetType(), New)
}
