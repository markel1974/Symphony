package magicdesk

import (
	"github.com/markel1974/symphony/src/hardware/c64_cartridges_rev1/catalog"
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

// Identifier returns a string that uniquely identifies a Magic Desk cartridge within the system.
func Identifier() string {
	return "magic_desk"
}

// Factory represents a type responsible for creating specific instance components with defined identifiers and kinds.
type Factory struct {
}

// NewFactory creates and returns a new instance of Factory, which implements the IFactory interface.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier retrieves the unique string identifier associated with the Factory instance.
func (t *Factory) Identifier() string {
	return Identifier()
}

// GetType returns an integer identifier representing the type of the Magic Desk cartridge.
func GetType() int {
	return catalog.CartridgeMagicDesk
}

// Kind returns an instance implementing the IC64Cartridge interface associated with the CartridgeMagicDesk type.
func (t *Factory) Kind() interface{} {
	z := (*CartridgeMagicDesk)(nil)
	return references.IC64Cartridge(z)
}

// Create initializes a new MagicDesk component with the specified parent, factory, label, and instance number.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewMagicDesk(parent, factory, label, instance)
}

// init initializes the application by registering a new factory and associating a type with its creation function.
func init() {
	registry.RegisterComponentFactory(NewFactory())

	catalog.RegisterType(GetType(), New)
}
