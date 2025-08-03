package interfaces

// CommandType definisce la natura strutturale di un nodo ICommand.
type CommandType int

const (
	// CommandTypeUnknown represents an undefined or unrecognized command type in the ICommand hierarchy.
	CommandTypeUnknown CommandType = iota

	// CommandTypeDirectory represents a command type corresponding to a directory node, allowing the addition of subcommands.
	CommandTypeDirectory

	// CommandTypeFile represents a command type corresponding to a file node.
	CommandTypeFile

	// CommandTypeLink represents a command type corresponding to a symbolic link node.
	CommandTypeLink
)

// ICommand defines an interface for handling hierarchical commands with support for events, execution, and navigation.
type ICommand interface {
	SetParent(parent ICommand)

	Type() CommandType

	PaintEvent() PaintFn

	ReadEvent() ReadFn

	ReadBroadcastEvent() ReadFn

	TimerEvent() TimerFn

	SetReadBroadcastFn(fn ReadFn)

	SetReadFn(fn ReadFn)

	SetTimerFn(fn TimerFn)

	SetPaintFn(fn PaintFn)

	Daemon() bool

	Background() bool

	SetHelp(short string, long string)

	Name() string

	Root() ICommand

	Parent() ICommand

	Childs() []ICommand

	DirectoryListing() []string

	Help() string

	FindChildren(name string) ICommand

	FindChildrenPrefix(prefix string) ICommand

	Find(args []string) (ICommand, []string, error)

	FindNext(next string) ICommand

	Traverse(args []string) ICommand

	SuggestionsFor(typedName string) []string

	Execute(task IProcess, arg []string) error

	Commands() []ICommand

	AddCommand(cx ICommand) error

	RemoveCommand(cx ICommand) error

	CommandPath() string

	Path() []string

	HasAlias(s string) bool

	NameAndAliases() string

	IsAdditionalHelpTopicCommand() bool

	HasHelpSubCommands() bool

	HasParent() bool
}
