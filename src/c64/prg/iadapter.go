package prg

type IAdapter interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}
