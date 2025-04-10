package interfaces

// PathSeparator defines the standard character used to separate elements in a file system path.
const PathSeparator = "/"

// RunFn defines a function type executed with a given context, command, process ID, and arguments, returning an error if any.
type RunFn func(r IContext, cmd ICommand, pid int, args []string) error

// TimerFn defines a function type for handling timer-based events within a context, providing arguments for event details.
type TimerFn func(r IContext, cmd ICommand, pid int, tid int, ctx interface{}, interval int)

// ReadFn defines a function type for handling read events, providing context, command, process ID, and input details.
type ReadFn func(r IContext, cmd ICommand, pid int, ctx interface{}, code int, key rune)

// PaintFn defines a function type responsible for handling paint events given the context, command, process ID, user context, and surface.
type PaintFn func(r IContext, cmd ICommand, pid int, ctx interface{}, surface ISurface)

// ICommand represents an interface for building hierarchical command structures with event handling and utility methods.
type ICommand interface {
	PaintEvent() PaintFn

	ReadEvent() ReadFn

	TimerEvent() TimerFn

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

	Help() string

	FindChildren(name string) ICommand

	FindChildrenPrefix(prefix string) ICommand

	Find(args []string) (ICommand, []string, error)

	Traverse(args []string) (ICommand, []string, error)

	SuggestionsFor(typedName string) []string

	FindSuggestions(arg string) []string

	VisitParents(fn func(ICommand))

	Execute(r IContext, arg []string, pid int) error

	Commands() []ICommand

	//AddCommand(cx ...ICommand) error

	//RemoveCommand(cx ...ICommand)

	CommandPath() string

	Path() []string

	HasAlias(s string) bool

	NameAndAliases() string

	HasSubCommands() bool

	IsAdditionalHelpTopicCommand() bool

	HasHelpSubCommands() bool

	HasParent() bool
}
