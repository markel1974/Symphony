package references

import "io"

type INode interface {
	AddComponent(component IComponent) INode
	GetComponent() IComponent
	GetChildren() []INode
	FindNode(path string) INode
}

// IHardware defines an interface for hardware components with a Reset method.
type IHardware interface {
	Reset()
}

// INavigate provides an interface for hierarchical navigation and manipulation of node structures and properties.
type INavigate interface {
	GetId() string

	GetNode() INode

	SetNode(node INode)

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
