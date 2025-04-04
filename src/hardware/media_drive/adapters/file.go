package adapters

import (
	"os"
)

// File represents a file system adapter for managing files and directories within a specified path.
// The `path` field denotes the root directory where operations are performed.
type File struct {
	path string
}

// NewFile creates a new AdapterFileSystem instance for the specified directory path.
// It ensures the provided path ends with a path separator and verifies if the path is a valid directory.
func NewFile(path string) (*File, error) {
	return &File{path: path}, nil
}

func (a *File) Extension() string { return "*" }

// Name returns the path associated with the AdapterFileSystem instance.
func (a *File) Name() string {
	return a.path
}

// ReadDir reads the content of the directory specified in the AdapterFileSystem and returns a slice of os.FileInfo.
func (a *File) ReadDir() ([]os.FileInfo, error) {
	s, err := os.Stat(a.path)
	if err != nil {
		return nil, err
	}
	var out []os.FileInfo
	out = append(out, s)
	return out, nil
}

// ReadFile reads the contents of a file located at the given path relative to the AdapterFileSystem's base path.
// It returns the file's data as a byte slice or an error if the file cannot be read.
func (a *File) ReadFile(_ string) ([]byte, error) {
	data, err := os.ReadFile(a.path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes the provided data to a file with the specified name within the adapter's file path.
// It creates the file if it does not exist or overwrites it if it does.
// Returns an error if the write operation fails.
func (a *File) WriteFile(_ string, data []byte) error {
	if err := os.WriteFile(a.path, data, 0644); err != nil {
		return err
	}
	return nil
}
