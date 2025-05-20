package config

// Cartridge represents a data structure containing metadata and binary data of a cartridge file.
type Cartridge struct {
	kind string
	name string
	path string
	data []byte
}

// NewCartridge creates a new Cartridge instance using the specified kind and file path, loading its data from the file.
// Returns a pointer to the Cartridge and an error if the file cannot be read or loaded.
func NewCartridge(kind string, filePath string, name string, data []byte) (*Cartridge, error) {
	c := &Cartridge{
		kind: kind,
		name: name,
		path: filePath,
		data: data,
	}
	return c, nil
}

// GetKind returns the kind of the cartridge as a string.
func (c *Cartridge) GetKind() string {
	return c.kind
}

// GetName returns the name of the Cartridge instance.
func (c *Cartridge) GetName() string {
	return c.name
}

// GetPath returns the file path of the Cartridge instance as a string.
func (c *Cartridge) GetPath() string {
	return c.path
}

// GetData retrieves the data stored in the Cartridge as a byte slice.
func (c *Cartridge) GetData() []byte {
	return c.data
}
