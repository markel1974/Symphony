package virtualdrive

import "github.com/markel1974/c64emu/src/config"

type IIec interface {
	PeripheralRead() uint8
	PeripheralWrite(deviceNumber uint8, data uint8)
}

type IVirtualDrive interface {
	Setup(*config.Config)
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
