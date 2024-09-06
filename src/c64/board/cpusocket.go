package board

import mos6510 "github.com/markel1974/c64emu/src/components/6510"

type CPUSocket struct {
	pic   mos6510.IPic
	banks mos6510.IBanks
}

func NewCPUSocket() *CPUSocket {
	c := &CPUSocket{}
	return c
}

func (w *CPUSocket) Setup(board *Board) {
	w.pic = board.pic
	w.banks = board.banks
}

func (w *CPUSocket) GetPic() mos6510.IPic {
	return w.pic
}

func (w *CPUSocket) GetBanks() mos6510.IBanks {
	return w.banks
}
