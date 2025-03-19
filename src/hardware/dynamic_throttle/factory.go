package dynamic_throttle

import (
	"github.com/markel1974/c64emu/src/references"
)

const componentId = "dynamic_throttle"

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (t *Factory) Identifier() string {
	return componentId
}

func (t *Factory) Kind() interface{} {
	z := (*DynamicThrottle)(nil)
	return references.IThrottle(z)
}

func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewDynamicThrottle(parent, factory, label)
}
