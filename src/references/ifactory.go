package references

type IFactory interface {
	Kind() interface{}

	Identifier() string

	Create(parent IComponent, factory IComponentFactory, suffix string) IComponent
}
