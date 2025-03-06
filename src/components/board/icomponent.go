package board

// IComponent represents an interface for components with identifiable and hierarchical structure.
type IComponent interface {
	IDumpable

	GetId() string

	GetParentId() string

	Emulate()

	Reset()
}

// Components represents a collection of IComponent instances managed in a container mapped by their unique IDs.
type Components struct {
	container map[string]IComponent
}

// NewComponents initializes and returns a pointer to a new Components instance with an empty container map.
func NewComponents() *Components {
	return &Components{
		container: make(map[string]IComponent),
	}
}

// Register adds the given IComponent instance to the container, identified by its unique ID obtained from GetId().
func (s *Components) Register(component IComponent) {
	id := component.GetId()
	s.container[id] = component
}
