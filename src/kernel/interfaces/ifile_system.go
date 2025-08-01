package interfaces

type IFileSystem interface {
	AddSearchPath(sp ICommand)

	CWDName() string

	CWDSet(arg string) bool

	CWDCommandPath() string

	CWDPath() []string

	CWDDirectoryListing() []string

	Find(line string) (ICommand, []string, error)

	Help(path string) (string, error)

	Suggestion(in string, cursor int) (string, []string, bool)
}
