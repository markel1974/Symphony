package references

// IFactory is an interface for creating and identifying components in a structured and modular system.
// Kind returns the type or category of the factory's components.
// Identifier provides a unique string identifier for the factory.
// Create instantiates a new component with a parent, associated factory, and a specific suffix.
type IFactory interface {
	Identifier() string

	Create(parent IComponent, factory IComponentFactory, label string, instance int) IComponent
}
