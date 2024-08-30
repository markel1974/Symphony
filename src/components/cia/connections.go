package cia

const intrCia1Id = 4

type IWiring interface {
	Reset()
	ReadPortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8
	ReadPortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8) uint8
	WritePortA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WritePortB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WriteDdrA(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
	WriteDdrB(prA uint8, ddrA uint8, prB uint8, ddrB uint8)
}
