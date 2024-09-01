package mos6522

type IWiring interface {
	Reset()
	ReadPRA(uint8, uint8) uint8
	ReadPRB(uint8, uint8) uint8
	WritePRA(uint8, uint8)
	WritePRB(uint8, uint8)
	WriteDDRA(uint8, uint8)
	WriteDDRB(uint8, uint8)
}
