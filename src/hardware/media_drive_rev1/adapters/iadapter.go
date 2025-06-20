package adapters

import "os"

type IChannel interface {
}

// IAdapter defines an interface for file system and archival operations such as reading, writing, and directory listing.
// Name returns the name of the adapter.
// Extension retrieves the default file extension associated with the adapter.
// ReadDir lists the contents of a directory, returning a slice of file information or an error.
// ReadFile reads and returns the content of a specified file or an error.
// WriteFile writes data to a specified file and returns an error if the operation fails.
type IAdapter interface {
	Name() string

	Extension() string

	ReadDir() ([]os.FileInfo, error)

	ReadFile(name string) ([]byte, error)

	WriteFile(name string, data []byte) error

	RenameFile(oldName string, newName string) error

	Format(name string, id string) error

	Reset() error

	Validate() error

	ScratchFile(commandData string) error

	BlockRead(c IChannel, track int, sector int) error

	BlockWrite(c IChannel, track int, sector int) error

	MemoryRead(address uint16, length int) ([]byte, error)

	MemoryWrite(address uint16, data []byte) error

	MemoryExec(address uint16) error

	Position(c IChannel, position int) error
}
