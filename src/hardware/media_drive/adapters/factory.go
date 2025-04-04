package adapters

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// Factory represents a factory for creating and managing IAdapter instances.
type Factory struct {
	void *Void
}

// NewFactory creates and initializes a new Factory instance with a Void adapter.
func NewFactory() *Factory {
	f := &Factory{
		void: NewVoid(),
	}
	return f
}

// Void retrieves the void adapter, a placeholder adapter that does not perform any operations and always returns errors.
func (f *Factory) Void() IAdapter {
	return f.void
}

// Create initializes and returns a new IAdapter instance for a specified path or an error if the operation fails.
func (f *Factory) Create(p string) (IAdapter, error) {
	d, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if d.IsDir() {
		return NewFileSystem(p)
	}
	ext := strings.TrimSpace(strings.ToLower(path.Ext(p)))
	if ext == ".zip" {
		return NewZip(p)
	}
	return nil, fmt.Errorf("unsupported file: %s", p)
}
