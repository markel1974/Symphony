package board

// IComponent defines an interface for components with identifiable parent-child relationships and state dumping capabilities.
// GetId returns the unique identifier of the component.
// GetParentId returns the identifier of the component's parent, if applicable.
type IComponent interface {
	IDumpable

	GetId() string

	GetParentId() string
}
