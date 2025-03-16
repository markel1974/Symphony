package references

type IFactory interface {
	Create(parent IComponent, factory IComponentFactory, suffix string) IComponent

	Kind() interface{}

	Identifier() string
}
