package interfaces

type IFileSystem interface {
	IServer

	CallAddSearchPath(router IRouter, sp ICommand)

	CallCWDName(router IRouter) string

	CallCWDSet(router IRouter, arg string) bool

	CallCWDCommandPath(router IRouter) string

	CallCWDPath(router IRouter) []string

	CallCWDDirectoryListing(router IRouter) []string

	CallFind(router IRouter, line string) (ICommand, []string, error)

	CallHelp(router IRouter, path string) (string, error)

	CallSuggestion(router IRouter, in string, cursor int) (string, []string, bool)
}
