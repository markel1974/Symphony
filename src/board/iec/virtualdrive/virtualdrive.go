package virtualdrive

type IVirtualDrive interface {
	Reset()
	Ready() bool
	GetPath() string
	Open(uint8, []byte) uint8
	Close(uint8) uint8
	Read(uint8, *uint8) uint8
	Write(uint8, uint8, bool) uint8
}
