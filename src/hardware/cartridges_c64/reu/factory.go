package reu

import (
	"github.com/markel1974/c64emu/src/hardware/cartridges_c64/catalog"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

const baseId = "REU"
const identifier128K = "128K"
const identifier256K = "256K"
const identifier512K = "512K"
const identifier1M = "1M"
const identifier2M = "2M"
const identifier4M = "4M"
const identifier8M = "8M"
const identifier16M = "16M"

type Factory struct {
	kind string
	size int
}

func NewFactory(kind string) *Factory {
	f := &Factory{
		kind: baseId + kind,
	}
	switch kind {
	case identifier128K:
		f.size = size128K
	case identifier256K:
		f.size = size256K
	case identifier512K:
		f.size = size512K
	case identifier1M:
		f.size = size1M
	case identifier2M:
		f.size = size2M
	case identifier4M:
		f.size = size4M
	case identifier8M:
		f.size = size8M
	case identifier16M:
		f.size = size16M
	default:
		f.size = size512K
	}
	return f
}

func (t *Factory) Identifier() string {
	return t.kind
}

func (t *Factory) Kind() interface{} {
	z := (*REU)(nil)
	return references.ICartridgeC64(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return newReu(parent, factory, label, instance, t.kind, t.size)
}

func init() {
	registry.RegisterComponentFactory(NewFactory(identifier128K))
	registry.RegisterComponentFactory(NewFactory(identifier256K))
	registry.RegisterComponentFactory(NewFactory(identifier512K))
	registry.RegisterComponentFactory(NewFactory(identifier1M))
	registry.RegisterComponentFactory(NewFactory(identifier2M))
	registry.RegisterComponentFactory(NewFactory(identifier4M))
	registry.RegisterComponentFactory(NewFactory(identifier8M))
	registry.RegisterComponentFactory(NewFactory(identifier16M))

	catalog.RegisterHardware(Id128K, New128K)
	catalog.RegisterHardware(Id256K, New256K)
	catalog.RegisterHardware(Id512K, New512K)
	catalog.RegisterHardware(Id1M, New1M)
	catalog.RegisterHardware(Id2M, New2M)
	catalog.RegisterHardware(Id4M, New4M)
	catalog.RegisterHardware(Id8M, New8M)
	catalog.RegisterHardware(Id16M, New16M)
}
