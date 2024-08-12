package virtualdrive

type IIec interface {
	PeripheralAtnResponse(deviceNumber uint8, data uint8)
	PeripheralRead() uint8
	PeripheralWrite(deviceNumber uint8, data uint8)
}

type IVirtualDrive interface {
	Reset()
	Emulate()
	Ready() bool
	GetDeviceNumber() uint8
	AtnStateChanged(bool, bool)
	BusStateChanged(uint8)

	//GetPath() string
	//Open(uint8, []uint8) uint8
	//Close(uint8) uint8
	//Read(uint8) (uint8, uint8)
	//Write(uint8, uint8, bool) uint8
}

func ParseFileName(name string, convertCharset bool) (string, int, int, int) {
	//TODO IMPLEMENT
	mode := FMODE_READ
	kind := FTYPE_PRG
	//dest uint8* , dest_len int& , mode int& , kind int, rec_len int& ,
	return "", mode, kind, 0
}
