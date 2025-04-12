package interfaces

// ICommand defines an interface for handling hierarchical commands with support for events, execution, and navigation.
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

	Traverse(args []string) ICommand

	//TODO REMOVE
	SuggestionsFor(typedName string) []string

	SuggestionsFor_NEW(typedName string) []string

	FindSuggestions(arg string) []string

	VisitParents(fn func(ICommand))

	Execute(task ITask, arg []string) error

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
