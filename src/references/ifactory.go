package references

import "github.com/markel1974/c64emu/src/component"

type IComponentFactory interface {
	Create(component.IComponent, string, string) (component.IComponent, error)
}
