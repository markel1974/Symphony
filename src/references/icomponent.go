package references

import "io"

// INode represents a hierarchical node structure allowing management of components, child nodes, and traversal operations.
// AddComponent adds a child component to the node and returns the new child node.
// GetComponent retrieves the component associated with the node.
// GetChildren returns a slice of all child nodes of the current node.
// FindNode searches for and returns a node based on a given path string.
type INode interface {
	AddComponent(component IComponent) INode
	RemoveComponent(component IComponent) bool
	GetComponent() IComponent
	GetChildren() []INode
	FindComponent(path string) IComponent
}

// IHardware defines an interface for hardware functionalities with a method to reset the hardware state.
type IHardware interface {
	Reset()
}

// INavigate is an interface for managing and navigating hierarchical components and properties.
// GetId retrieves the unique identifier of the component.
// GetNode retrieves the associated node.
// SetNode associates a node with the current component.
// GetChildren returns a slice of child components attached to the current component.
// GetComponentPath retrieves a component based on the provided path as a string.
// GetProperty fetches a property by its key, returning a value or an error.
// GetPropertyPath fetches a property using a component path and key, returning a value or an error.
// SetProperty sets a property by key with the provided value and may return an error.
// SetPropertyPath sets a property via a component path and key with the provided value and may return an error.
// Dump returns a map representation of the component's properties and configuration.
// DumpPath returns a map of properties and configuration for a component specified by a path.
// DumpAll returns a complete dump of all properties and configuration of the component and its children.
// Restore reconstructs the internal state using a map of properties and configuration.
// RestoreAll restores all configurations and properties for the current component and its children.
// RestorePath restores the configuration of a specific component using a path and a map of properties.
// CommandAdd adds a command to the component with a unique identifier and description.
// CommandExec executes a component's command by its identifier and returns the result or an error.
// CommandExecPath executes a command at a specific path of a component using an identifier, returning the result or error.
// CommandDocumentation populates a map with descriptions of available commands.
// Print outputs a formatted representation of the component to an io.Writer based on the provided format and flags.
type INavigate interface {
	GetId() string

	GetNode() INode

	SetNode(node INode)

	GetChild(id string) IComponent

	GetChildren() []IComponent

	GetComponentPath(string) IComponent

	GetProperty(string) (interface{}, error)

	GetPropertyPath(string, string) (interface{}, error)

	SetProperty(string, interface{}) error

	SetPropertyPath(string, string, interface{}) error

	Dump() (map[string]interface{}, error)

	DumpPath(string) (map[string]interface{}, error)

	DumpAll() (map[string]interface{}, error)

	Restore(map[string]interface{}) error

	RestoreAll(map[string]interface{}) error

	RestorePath(string, map[string]interface{}) error

	CommandAdd(id string, desc string, command interface{}) error

	CommandExec(string, ...interface{}) (interface{}, error)

	CommandExecPath(string, string, ...interface{}) (interface{}, error)

	CommandDocumentation(map[string]interface{})

	Print(io.Writer, string, bool)
}

// IComponent represents a composite interface that combines IHardware and INavigate capabilities.
type IComponent interface {
	IHardware
	INavigate
}

// IComponentFactory defines an interface for creating IComponent instances with hierarchical and contextual parameters.
type IComponentFactory interface {
	Create(parent IComponent, id string, instance int) (IComponent, error)

	CreateIQuartz(parent IComponent, id string, instance int) (IComponent, IQuartz, error)

	CreateIVIC(parent IComponent, id string, instance int) (IComponent, IVIC, error)

	CreateISID(parent IComponent, id string, instance int) (IComponent, ISID, error)

	CreateICIA(parent IComponent, id string, instance int) (IComponent, ICIA, error)

	CreateIVIA(parent IComponent, id string, instance int) (IComponent, IVIA, error)

	CreateI6510(parent IComponent, id string, instance int) (IComponent, I6510, error)

	CreateIPIC6510(parent IComponent, id string, instance int) (IComponent, IPIC6510, error)

	CreateIPLAc64(parent IComponent, id string, instance int) (IComponent, IPlaC64, error)

	CreateIPLAc1541(parent IComponent, id string, instance int) (IComponent, IPLAc1541, error)

	CreateIEC(parent IComponent, id string, instance int) (IComponent, IIec, error)

	CreateIROMLoaderC64(parent IComponent, id string, instance int) (IComponent, IROMLoaderC64, error)

	CreateIROMLoaderC1541(parent IComponent, id string, instance int) (IComponent, IROMLoaderC1541, error)

	CreateICartridgeManagerC64(parent IComponent, id string, instance int) (IComponent, ICartridgeManagerC64, error)

	CreateIThrottle(parent IComponent, id string, instance int) (IComponent, IThrottle, error)

	CreateIKeyboard(parent IComponent, id string, instance int) (IComponent, IKeyboard, error)

	CreateIJoystick(parent IComponent, id string, instance int) (IComponent, IJoystick, error)
}
