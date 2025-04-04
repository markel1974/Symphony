package adapters

import (
	"fmt"
	"os"
	"strings"
)

// FileSystem represents a file system adapter for managing files and directories within a specified path.
// The `path` field denotes the root directory where operations are performed.
type FileSystem struct {
	path string
}

// NewFileSystem creates a new AdapterFileSystem instance for the specified directory path.
// It ensures the provided path ends with a path separator and verifies if the path is a valid directory.
func NewFileSystem(path string) (*FileSystem, error) {
	if !strings.HasSuffix(path, string(os.PathSeparator)) {
		path = path + string(os.PathSeparator)
	}
	d, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !d.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", path)
	}
	return &FileSystem{path: path}, nil
}

func (a *FileSystem) Extension() string { return "" }

// Name returns the path associated with the AdapterFileSystem instance.
func (a *FileSystem) Name() string {
	return a.path
}

// ReadDir reads the content of the directory specified in the AdapterFileSystem and returns a slice of os.FileInfo.
func (a *FileSystem) ReadDir() ([]os.FileInfo, error) {
	items, err := os.ReadDir(a.path)
	if err != nil {
		return nil, err
	}
	var out []os.FileInfo
	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, info)
	}
	return out, nil
}

// ReadFile reads the contents of a file located at the given path relative to the AdapterFileSystem's base path.
// It returns the file's data as a byte slice or an error if the file cannot be read.
func (a *FileSystem) ReadFile(plainName string) ([]byte, error) {
	f, err := a.ReadDir()
	if err != nil {
		return nil, err
	}
	var name = ""
	for _, item := range f {
		if string(CreateFileName(item.Name())) == plainName {
			name = item.Name()
			break
		}
	}
	if len(name) == 0 {
		return nil, fmt.Errorf("file not found %s", plainName)
	}
	data, err := os.ReadFile(a.path + name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes the provided data to a file with the specified name within the adapter's file path.
// It creates the file if it does not exist or overwrites it if it does.
// Returns an error if the write operation fails.
func (a *FileSystem) WriteFile(plainName string, data []byte) error {
	completeFileName := a.path + plainName
	if err := os.WriteFile(completeFileName, data, 0644); err != nil {
		return err
	}
	return nil
}
