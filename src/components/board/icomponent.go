package board

// IComponent represents an interface for components with identifiable and hierarchical structure.
type IComponent interface {
	IDumpable

	GetId() string

	GetParentId() string

	Emulate()

	Reset()
}
