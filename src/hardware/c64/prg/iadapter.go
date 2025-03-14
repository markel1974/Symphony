package prg

// IAdapter defines an interface for reading and writing data to a device or memory using an address and data value.
// Read fetches a byte of data from the specified address.
// Write stores a byte of data at the specified address.
type IAdapter interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}
