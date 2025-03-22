package config

// Drive represents a virtual drive with a specific kind, data, ID, and write protection status.
type Drive struct {
	kind           string
	data           []byte
	id             string
	writeProtected bool
}

// NewDrive creates a new Drive instance with the specified kind and file path, initializing its data and write protection status.
// It returns the created Drive pointer and an error if the file cannot be accessed or read.
func NewDrive(kind string, path string) (*Drive, error) {
	data, wp, err := ImageFromFile(path)
	if err != nil {
		return nil, err
	}
	if len(kind) == 0 {
		kind = "c1541"
	}
	d := &Drive{
		kind:           kind,
		id:             path,
		data:           data,
		writeProtected: wp,
	}
	return d, nil
}

// GetKind returns the type of the drive as a string.
func (d *Drive) GetKind() string {
	return d.kind
}

// GetData retrieves the raw data stored in the Drive as a byte slice.
func (d *Drive) GetData() []byte {
	return d.data
}

// GetId returns the unique identifier of the Drive instance as a string.
func (d *Drive) GetId() string {
	return d.id
}

// IsWriteProtected checks if the drive is set to write-protected mode and returns true if it is; otherwise, false.
func (d *Drive) IsWriteProtected() bool {
	return d.writeProtected
}
