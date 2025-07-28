package interfaces

type IFileSystem interface {
	AddSearchPath(sp ICommand)

	CWD() ICommand

	CWDSet(arg string) bool

	Find(line string) (ICommand, []string, error)

	Help(path string) (string, error)

	Suggestion(in string, cursor int) (string, []string, bool)
}
