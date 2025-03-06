package board

// IDumpable represents an interface for serializing an object's state into a map and restoring it from the same.
// Dump exports the object's state into a provided map.
// Restore imports the object's state from a provided map.
// GetProperties returns metadata about the object's properties.
type IDumpable interface {
	Dump(d *Dumper) error

	Restore(d *Dumper) error

	GetProperties() map[string]*PropertyInfo
}
