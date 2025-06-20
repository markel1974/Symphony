package adapters

import (
	"fmt"
	"os"
	"strings"
)

// DirectoryExtension returns the default extension for directories as a string.
func DirectoryExtension() string { return "" }

// Directory represents a filesystem directory and provides utility methods for file operations within it.
type Directory struct {
	path string
}

// NewDirectory initializes and returns a Directory instance for the specified path if it is a valid directory.
// Returns an error if the path does not exist or is not a directory.
func NewDirectory(path string) (*Directory, error) {
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
	return &Directory{path: path}, nil
}

// Extension retrieves the directory extension as a string.
func (a *Directory) Extension() string { return DirectoryExtension() }

// Name returns the path string of the Directory instance.
func (a *Directory) Name() string {
	return a.path
}

// ReadDir reads the contents of the directory and returns a slice of os.FileInfo representing the files and directories.
// Returns an error if the directory cannot be read.
func (a *Directory) ReadDir() ([]os.FileInfo, error) {
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

// ReadFile retrieves the content of a file with the given plain name in the directory, returning the data or an error.
func (a *Directory) ReadFile(plainName string) ([]byte, error) {
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

// WriteFile writes the provided data to a file named plainName in the directory. Returns an error if the operation fails.
func (a *Directory) WriteFile(plainName string, data []byte) error {
	completeFileName := a.path + plainName
	if err := os.WriteFile(completeFileName, data, 0644); err != nil {
		return err
	}
	return nil
}

// RenameFile renames a file from oldName to newName within the directory. Returns an error if the operation fails.
func (a *Directory) RenameFile(oldName string, newName string) error {
	return fmt.Errorf("unimplemented")
}

// Format applies a specific format operation to the directory, based on the provided name and id, and returns an error if unimplemented.
func (a *Directory) Format(name string, id string) error {
	return fmt.Errorf("unimplemented")
}

// Reset reinitializes the directory to a default or empty state, if supported. Returns an error if unimplemented.
func (a *Directory) Reset() error {
	return fmt.Errorf("unimplemented")
}

// Validate checks the state of the Directory instance and ensures it meets required conditions or constraints.
func (a *Directory) Validate() error {
	return fmt.Errorf("unimplemented")
}

// ScratchFile removes or marks a specified file for deletion within the directory, using the given command data.
func (a *Directory) ScratchFile(commandData string) error {
	return fmt.Errorf("unimplemented")
}
