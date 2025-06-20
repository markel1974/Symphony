package adapters

import (
	"fmt"
	"os"
)

// Void represents a placeholder adapter that does not perform any operations and always returns errors.
type Void struct {
}

// NewVoid initializes and returns a pointer to a new instance of an empty adapter (AdapterVoid).
func NewVoid() *Void {
	return &Void{}
}

// Name returns the string identifier "empty" for the AdapterVoid instance.
func (a *Void) Name() string {
	return "empty"
}

func (a *Void) Extension() string { return "" }

// ReadDir returns an error indicating that the adapter is invalid and does not support directory reading.
func (a *Void) ReadDir() ([]os.FileInfo, error) {
	return nil, fmt.Errorf("invalid void adapter")
}

// ReadFile attempts to read a file with the specified name but always returns an error indicating an invalid adapter.
func (a *Void) ReadFile(_ string) ([]byte, error) {
	return nil, fmt.Errorf("invalid void adapter")
}

// WriteFile attempts to write data to a file with the given name and returns an error indicating the operation is invalid.
func (a *Void) WriteFile(name string, data []byte) error {
	return fmt.Errorf("invalid void adapter")
}

// RenameFile renames a file from oldName to newName within the directory. Returns an error if the operation fails.
func (a *Void) RenameFile(oldName string, newName string) error {
	return fmt.Errorf("invalid void adapter")
}

// Format applies a specific format operation to the directory, based on the provided name and id, and returns an error if unimplemented.
func (a *Void) Format(name string, id string) error {
	return fmt.Errorf("invalid void adapter")
}

// Reset reinitializes the directory to a default or empty state, if supported. Returns an error if unimplemented.
func (a *Void) Reset() error {
	return fmt.Errorf("invalid void adapter")
}

// Validate checks the state of the Directory instance and ensures it meets required conditions or constraints.
func (a *Void) Validate() error {
	return fmt.Errorf("invalid void adapter")
}

// ScratchFile removes or marks a specified file for deletion within the directory, using the given command data.
func (a *Void) ScratchFile(commandData string) error {
	return fmt.Errorf("invalid void adapter")
}
