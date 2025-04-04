package adapters

import "os"

// IAdapter defines an interface for interacting with a file-like storage system or adapter.
// Name returns the name of the adapter as a string.
// ReadDir retrieves a list of files and directories in the current directory as os.FileInfo slices.
// ReadFile reads and returns the content of a specified file as a byte slice.
// WriteFile writes data to a specified file, creating or overwriting the file.
type IAdapter interface {
	Name() string

	ReadDir() ([]os.FileInfo, error)

	ReadFile(name string) ([]byte, error)

	WriteFile(name string, data []byte) error
}
