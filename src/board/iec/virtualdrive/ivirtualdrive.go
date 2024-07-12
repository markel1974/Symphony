package virtualdrive

type IVirtualDrive interface {
	Reset()
	Ready() bool
	GetPath() string
	Open(uint8, []uint8) uint8
	Close(uint8) uint8
	Read(uint8) (uint8, uint8)
	Write(uint8, uint8, bool) uint8
}

type VirtualDrive struct {
	deviceNumber uint8
	errorData    []byte
	errorIdx     int
	cmdBuf       []uint8 // Buffer for incoming command strings
	cmdLen       int     // Length of received command
}

func NewVirtualDrive(deviceNumber uint8) *VirtualDrive {
	return &VirtualDrive{
		deviceNumber: deviceNumber,
		errorData:    nil,
		errorIdx:     0,
		cmdBuf:       make([]uint8, 64),
		cmdLen:       0,
	}
}

func (vd *VirtualDrive) SetError(e int) {
	vd.errorData = Errors[e]
	vd.errorIdx = 0
}

func (vd *VirtualDrive) CommandSet(data uint8) bool {
	if vd.cmdLen >= 58 {
		return false
	}
	vd.cmdBuf[vd.cmdLen] = data
	vd.cmdLen++
	return true
}

func (vd *VirtualDrive) CommandClear() {
	vd.cmdLen = 0
}

func (vd *VirtualDrive) ParseFileName(name string, convertCharset bool) (string, int, int, int) {
	//TODO IMPLEMENT
	mode := FMODE_READ
	kind := FTYPE_PRG
	//dest uint8* , dest_len int& , mode int& , kind int, rec_len int& ,
	return "", mode, kind, 0
}

func (vd *VirtualDrive) CommandExecBuf() {
	//TODO IMPLEMENT
	vd.CommandExec(vd.cmdBuf)
	vd.cmdLen = 0
}

func (vd *VirtualDrive) CommandExec(data []uint8) {
	//TODO IMPLEMENT
	//vd.executeCmd(vd.cmdBuf, vd.cmdLen)
	vd.cmdLen = 0
}

func (vd *VirtualDrive) LedTurnOn() {
}

func (vd *VirtualDrive) LedTurnOff() {
}

func (vd *VirtualDrive) RetrieveError() uint8 {
	if vd.errorIdx < len(vd.errorData) {
		v := vd.errorData[vd.errorIdx]
		vd.errorIdx++
		return v
	}
	return '\r'
}
