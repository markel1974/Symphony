package interfaces

type IFileSystem interface {
	IServer

	CallCWDPath(router IRouter) string

	CallCWDDirectoryListing(router IRouter) []string

	CallFind(router IRouter, line string) (ICommand, []string, error)

	CallHelp(router IRouter, path string) (string, error)
}
