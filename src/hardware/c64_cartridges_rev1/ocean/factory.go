package ocean

import (
	"github.com/markel1974/c64emu/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

func Identifier() string {
	return "ocean"
}

// GetType returns the type identifier of the Ocean cartridge as an integer constant.
func GetType() int {
	return catalog.CartridgeOcean
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
	z := (*CartridgeOcean)(nil)
	return references.ICartridgeC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewCartridgeOcean(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())

	catalog.RegisterType(GetType(), New)
	catalog.RegisterSizeDefault(New)
}
