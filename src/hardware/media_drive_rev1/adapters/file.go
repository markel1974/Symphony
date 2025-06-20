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

// BlockRead attempts to read a block of data from a specified track and sector but always returns an error for a Void adapter.
func (a *File) BlockRead(ch IChannel, track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// BlockWrite writes a block to the specified track and sector and always returns an error for the Void adapter.
func (a *File) BlockWrite(ch IChannel, track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// MemoryRead reads a specified memory region defined by the address and length. Returns an error for unsupported adapters.
func (a *File) MemoryRead(address uint16, length int) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

// MemoryWrite writes a slice of bytes to the specified memory address and returns an error for void adapters.
func (a *File) MemoryWrite(address uint16, data []byte) error {
	return fmt.Errorf("unimplemented")
}

// MemoryExec executes a command at the specified memory address and always returns an error for the Void adapter.
func (a *File) MemoryExec(address uint16) error {
	return fmt.Errorf("unimplemented")
}

// Position sets the current position to the specified value, returning an error as the Void adapter is not functional.
func (a *File) Position(ch IChannel, position int) error {
	return fmt.Errorf("unimplemented")
}
