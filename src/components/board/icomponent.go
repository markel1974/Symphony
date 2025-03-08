package board

// IComponent represents a generic interface for a component in a hierarchical structure.
type IComponent interface {
	GetId() string

	SetNode(*Node)

	GetNode() *Node

	GetProperty(prop string) (interface{}, error)

	SetProperty(prop string, val interface{}) error

	Dump() (map[string]interface{}, error)

	Restore(state map[string]interface{}) error

	RunCommand(cmd string, args []string) (map[string]interface{}, error)

	Reset()
}

// BaseComponent represents a fundamental component with an ID, hierarchical path, and associated properties.
type BaseComponent struct {
	id         string
	node       *Node
	properties *Properties
}

// NewBaseComponent creates and returns a new instance of the BaseComponent.
func NewBaseComponent(name string, suffix string, runFn RunFn) *BaseComponent {
	id := name
	if len(suffix) > 0 {
		id += "_" + suffix
	}
	bc := &BaseComponent{
		id:         id,
		node:       nil,
		properties: NewProperties(runFn),
	}
	return bc
}

// GetId returns the unique identifier (id) of the BaseComponent.
func (b *BaseComponent) GetId() string {
	return b.id
}

func (b *BaseComponent) SetNode(n *Node) {
	b.node = n
}

// GetNode returns the path of the BaseComponent as a slice of strings.
func (b *BaseComponent) GetNode() *Node {
	return b.node
}

// AddProperty registers a new property to the BaseComponent with the given id, description, read-only flag, getter, and setter.
func (b *BaseComponent) AddProperty(id string, desc string, ro bool, get interface{}, set interface{}) {
	p := CreatePropertyInfo(id, desc, ro, get, set)
	b.properties.Add(p)
}

// GetProperty retrieves the value of a specified property by its name and returns it or an error if not found or inaccessible.
func (b *BaseComponent) GetProperty(prop string) (interface{}, error) {
	v, err := b.properties.GetProperty(prop)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// SetProperty updates the value of a specified property if it exists and is valid, returning an error otherwise.
func (b *BaseComponent) SetProperty(prop string, val interface{}) error {
	if err := b.properties.SetProperty(prop, val); err != nil {
		return err
	}
	return nil
}

// Dump returns a map representation of the component's current state by invoking the Dump method of its properties.
// Returns any error encountered during retrieval of property values.
func (b *BaseComponent) Dump() (map[string]interface{}, error) {
	state, err := b.properties.Dump()
	if err != nil {
		return nil, err
	}
	return state, nil
}

// Restore restores the component's state using the provided state map and returns an error if restoration fails.
func (b *BaseComponent) Restore(state map[string]interface{}) error {
	if err := b.properties.Restore(state); err != nil {
		return err
	}
	return nil
}

// RunCommand executes a command with specified arguments using the component's properties and returns the result or an error.
func (b *BaseComponent) RunCommand(cmd string, args []string) (map[string]interface{}, error) {
	d, err := b.properties.Run(cmd, args)
	if err != nil {
		return nil, err
	}
	return d, nil
}
