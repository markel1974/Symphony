package adapters

import (
	"github.com/markel1974/symphony/src/config"
	"path"
	"strings"
)

// Factory represents a factory for creating and managing IAdapter instances.
type Factory struct {
}

// NewFactory creates and initializes a new Factory instance with a Void adapter.
func NewFactory() *Factory {
	f := &Factory{}
	return f
}

// Create initializes and returns a new IAdapter instance for a specified path or an error if the operation fails.
func (f *Factory) Create(config *config.Config, asset string) (IAdapter, error) {
	info, err := config.AssetInfo(asset)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return NewDirectory(config, asset)
	}
	ext := strings.TrimSpace(strings.ToLower(path.Ext(asset)))
	switch ext {
	case ZipExtension():
		return NewZip(config, asset)
	default:
		return NewFile(config, asset)
	}
}
