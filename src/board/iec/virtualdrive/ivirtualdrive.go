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

type VirtualDrive struct {
	deviceNumber uint8
	error        int
}

func NewVirtualDrive(deviceNumber uint8) *VirtualDrive {
	return &VirtualDrive{deviceNumber: deviceNumber, error: 0}
}

func (vd *VirtualDrive) SetError(e int) {
	vd.error = e
}

func (v *VirtualDrive) LedTurnOn() {
}

func (v *VirtualDrive) LedTurnOff() {
}
