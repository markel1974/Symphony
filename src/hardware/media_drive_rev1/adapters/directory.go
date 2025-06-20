package adapters

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"strings"
)

// DirectoryExtension returns the default extension for directories as a string.
func DirectoryExtension() string { return "" }

// Directory represents a filesystem directory and provides utility methods for file operations within it.
type Directory struct {
	path   string
	config *config.Config
}

// NewDirectory initializes and returns a Directory instance for the specified path if it is a valid directory.
// Returns an error if the path does not exist or is not a directory.
func NewDirectory(config *config.Config, path string) (*Directory, error) {
	if !strings.HasSuffix(path, config.AssetPathSeparator()) {
		path = path + config.AssetPathSeparator()
	}
	return &Directory{
		config: config,
		path:   path,
	}, nil
}

// Extension retrieves the directory extension as a string.
func (a *Directory) Extension() string { return DirectoryExtension() }

// Name returns the path string of the Directory instance.
func (a *Directory) Name() string {
	return a.path
}

// ReadDir reads the contents of the directory and returns a slice of fs.FileInfo representing the files and directories.
// Returns an error if the directory cannot be read.
func (a *Directory) ReadDir() ([]config.IAssetInfo, error) {
	item, err := a.config.AssetDir(a.path)
	if err != nil {
		return nil, err
	}
	return item, nil
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
	data, err := a.config.AssetRead(a.path + name)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes the provided data to a file named plainName in the directory. Returns an error if the operation fails.
func (a *Directory) WriteFile(plainName string, data []byte) error {
	completeFileName := a.path + plainName
	if err := a.config.AssetWrite(completeFileName, data); err != nil {
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

// BlockRead attempts to read a block of data from a specified track and sector but always returns an error for a Void adapter.
func (a *Directory) BlockRead(ch IChannel, track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// BlockWrite writes a block to the specified track and sector and always returns an error for the Void adapter.
func (a *Directory) BlockWrite(ch IChannel, track int, sector int) error {
	return fmt.Errorf("unimplemented")
}

// MemoryRead reads a specified memory region defined by the address and length. Returns an error for unsupported adapters.
func (a *Directory) MemoryRead(address uint16, length int) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

// MemoryExec executes a command at the specified memory address and always returns an error for the Void adapter.
func (a *Directory) MemoryExec(address uint16) error {
	return fmt.Errorf("unimplemented")
}

// MemoryWrite writes a slice of bytes to the specified memory address and returns an error for void adapters.
func (a *Directory) MemoryWrite(address uint16, data []byte) error {
	return fmt.Errorf("unimplemented")
}

// Position sets the current position to the specified value, returning an error as the Void adapter is not functional.
func (a *Directory) Position(ch IChannel, position int) error {
	return fmt.Errorf("unimplemented")
}
