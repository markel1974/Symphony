package adapters

import (
	"fmt"
	"os"
)

// Void represents a type used primarily as a no-op or placeholder for unsupported operations or empty implementations.
type Void struct {
}

// NewVoid initializes and returns a new instance of the Void adapter, which serves as a no-operation placeholder.
func NewVoid() *Void {
	return &Void{}
}

// Name returns the name of the Void type, always as the string "empty".
func (a *Void) Name() string {
	return "empty"
}

// Extension returns an empty string as the file extension for the Void type.
func (a *Void) Extension() string { return "" }

// ReadDir returns an error indicating the void adapter is not able to read directory contents.
func (a *Void) ReadDir() ([]os.FileInfo, error) {
	return nil, fmt.Errorf("invalid void adapter")
}

// ReadFile attempts to read a file but always returns an error indicating the operation is invalid for the Void adapter.
func (a *Void) ReadFile(_ string) ([]byte, error) {
	return nil, fmt.Errorf("invalid void adapter")
}

// WriteFile attempts to write data to the specified file but always returns an error as it is not supported by the Void adapter.
func (a *Void) WriteFile(name string, data []byte) error {
	return fmt.Errorf("invalid void adapter")
}

// RenameFile attempts to rename a file from oldName to newName but always returns an error as this is a void adapter.
func (a *Void) RenameFile(oldName string, newName string) error {
	return fmt.Errorf("invalid void adapter")
}

// Format attempts to format the void adapter but always returns an error indicating the adapter is invalid.
func (a *Void) Format(name string, id string) error {
	return fmt.Errorf("invalid void adapter")
}

// Reset returns an error indicating that the Reset operation is not supported by the Void adapter.
func (a *Void) Reset() error {
	return fmt.Errorf("invalid void adapter")
}

// Validate checks the Void adapter for validity and always returns an error indicating it is invalid.
func (a *Void) Validate() error {
	return fmt.Errorf("invalid void adapter")
}

// ScratchFile attempts to execute a scratch operation but always returns an error indicating an invalid void adapter.
func (a *Void) ScratchFile(commandData string) error {
	return fmt.Errorf("invalid void adapter")
}

// BlockRead attempts to read a block of data from a specified track and sector but always returns an error for a Void adapter.
func (a *Void) BlockRead(track int, sector int) error {
	return fmt.Errorf("invalid void adapter")
}

// BlockWrite writes a block to the specified track and sector and always returns an error for the Void adapter.
func (a *Void) BlockWrite(track int, sector int) error {
	return fmt.Errorf("invalid void adapter")
}

// MemoryRead reads a specified memory region defined by the address and length. Returns an error for unsupported adapters.
func (a *Void) MemoryRead(address uint16, length int) ([]byte, error) {
	return nil, fmt.Errorf("invalid void adapter")
}

// MemoryWrite writes a slice of bytes to the specified memory address and returns an error for void adapters.
func (a *Void) MemoryWrite(address uint16, data []byte) error {
	return fmt.Errorf("invalid void adapter")
}

// Position sets the current position to the specified value, returning an error as the Void adapter is not functional.
func (a *Void) Position(position int) error {
	return fmt.Errorf("invalid void adapter")
}
