package board

// IComponent represents an interface for components with identifiable and hierarchical structure.
type IComponent interface {
	GetId() string

	GetParentId() string

	GetProperties() *Properties

	Emulate()

	Reset()
}
