package references

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"io"
	"log"
	"reflect"
	"strconv"
)

// INode defines a hierarchical structure for managing components and their relationships in a tree-like format.
// AddComponent adds a component as a child and returns the child node.
// RemoveComponent removes a specified component from the children, returning true if successful.
// GetComponent retrieves the component associated with the current node.
// GetChildren returns a slice of all child nodes.
// Traverse traverses the structure to locate a component by its path.
type INode interface {
	AddComponent(component IComponent) INode

	RemoveComponent(component IComponent) bool

	GetComponent() IComponent

	GetChildren() []INode

	Traverse(path string) IComponent
}

// IHardware defines an interface for hardware components, offering methods for emulation and reset operations.
// GetId retrieves the unique identifier of the component.
// EmulationRequired indicates whether emulation is necessary for the hardware component.
// Emulate starts or manages the emulation process for the hardware.
// Reset reinitializes or restores the hardware to a default state.
type IHardware interface {
	GetId() string

	GetName() string

	HardwareId() string

	Setup() error

	Connect() error

	EmulationRequired() bool

	Emulate()

	Internal() bool

	Reset()
}

// ICommand defines an interface for managing and executing commands in a system.
// GetCommand retrieves the underlying cli.Command object.
// CommandAdd registers a new command with a specified ID, description, and handler.
// CommandExec executes a registered command with the given arguments.
// CommandExecPath executes a command at a specific path with the given arguments.
// CommandDocumentation provides documentation for commands using a given map.
type ICommand interface {
	GetCommand() *shell.Command

	CommandAdd(id string, desc string, command interface{}) error

	CommandExec(string, ...string) (interface{}, error)

	CommandExecPath(string, string, ...string) (interface{}, error)

	CommandDocumentation(map[string]interface{})
}

// INavigate is an interface for managing and navigating hierarchical components and properties.
// GetNode retrieves the associated node.
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
// Print outputs a formatted representation of the component to an io.Writer based on the provided format and flags.
type INavigate interface {
	GetFactory() IComponentFactory

	GetNode() INode

	GetChildByHardwareId(id string) IComponent

	GetChildren() []IComponent

	GetComponentPath(string) IComponent

	PropertyGet(string) (interface{}, error)

	PropertyPath(string, string) (interface{}, error)

	PropertySet(string, interface{}) error

	PropertySetPath(string, string, interface{}) error

	Dump() (map[string]interface{}, error)

	DumpPath(string) (map[string]interface{}, error)

	DumpAll() (map[string]interface{}, error)

	Restore(map[string]interface{}) error

	RestoreAll(map[string]interface{}) error

	RestorePath(string, map[string]interface{}) error

	Print(io.Writer, string, bool)
}

// IComponent represents the core interface for components, combining hardware, navigation, and command functionalities.
type IComponent interface {
	IHardware

	INavigate

	ICommand
}

type ISocket interface {
	HardwareId() string

	Wire() error
}

// IComponentFactory defines methods for creating and managing various types of components in an emulation system.
type IComponentFactory interface {
	Create(parent IComponent, label string, id string, instance int) (IComponent, error)

	GetIDisplayBuffer() IDisplayBuffer

	GetIAudioRender() IAudioRender

	GetConfig() *config.Config
}

func IdInternalComponent(label string, instance int, id string) string {
	return label + ":" + id + ":" + strconv.Itoa(instance)
}

func InterfaceName(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Interface {
		return t.Elem().Name()
	}
	log.Fatal("not an interface")
	return ""
}
