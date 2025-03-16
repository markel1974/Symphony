package board

import (
	mos6510 "github.com/markel1974/c64emu/src/references"
)

type CPUSocket struct {
	board *Board
	pic   mos6510.IPIC6510
	banks mos6510.I6510Banks
}

func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{}
	return c
}

func (w *CPUSocket) Setup(board *Board) {
	w.board = board
	w.pic = board.pic
	w.banks = board.pla
}

func (w *CPUSocket) Reset() {
	//w.board.cpu.Reset()
}

func (w *CPUSocket) GetPic() mos6510.IPIC6510 {
	return w.pic
}

func (w *CPUSocket) GetBanks() mos6510.I6510Banks {
	return w.banks
}
