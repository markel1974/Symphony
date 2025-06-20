package adapters

import (
	"fmt"
	"os"
)

func FileExtension() string { return "*" }

// File represents a file with an associated path.
type File struct {
	path string
}

// NewFile initializes and returns a new File instance for the given path or an error if the operation fails.
func NewFile(path string) (*File, error) {
	return &File{path: path}, nil
}

// Extension returns the file extension of the associated file path.
func (a *File) Extension() string { return FileExtension() }

// Name returns the path of the file as a string.
func (a *File) Name() string {
	return a.path
}

// ReadDir retrieves the file information of the file at the specified path and returns it in a slice.
func (a *File) ReadDir() ([]os.FileInfo, error) {
	s, err := os.Stat(a.path)
	if err != nil {
		return nil, err
	}
	out := []os.FileInfo{s}
	return out, nil
}

// ReadFile reads the content of the file specified by its path and returns the data as a byte slice or an error.
func (a *File) ReadFile(_ string) ([]byte, error) {
	data, err := os.ReadFile(a.path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes the provided data to the file at the specified path. Returns an error if the write operation fails.
func (a *File) WriteFile(_ string, data []byte) error {
	if err := os.WriteFile(a.path, data, 0644); err != nil {
		return err
	}
	return nil
}

// RenameFile renames a file from oldName to newName within the directory. Returns an error if the operation fails.
func (a *File) RenameFile(oldName string, newName string) error {
	return fmt.Errorf("unimplemented")
}

// Format applies a specific format operation to the directory, based on the provided name and id, and returns an error if unimplemented.
func (a *File) Format(name string, id string) error {
	return fmt.Errorf("unimplemented")
}

// Reset reinitializes the directory to a default or empty state, if supported. Returns an error if unimplemented.
func (a *File) Reset() error {
	return fmt.Errorf("unimplemented")
}

// Validate checks the state of the Directory instance and ensures it meets required conditions or constraints.
func (a *File) Validate() error {
	return fmt.Errorf("unimplemented")
}

// ScratchFile removes or marks a specified file for deletion within the directory, using the given command data.
func (a *File) ScratchFile(commandData string) error {
	return fmt.Errorf("unimplemented")
}
