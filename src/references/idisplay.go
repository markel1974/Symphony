package references

// IDisplayBuffer defines methods for interacting with a display buffer, allowing data manipulation at specified indices.
// Set sets a single byte of data at the given index.
// SetMulti8 sets a single byte of data and applies it across multiple relevant sections.
// Set8 sets an array of 8 bytes of data starting at the given index.
type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}
