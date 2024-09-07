package board

import "github.com/markel1974/c64emu/src/components/sid"

type SidSocket struct {
	board *Board
}

func NewSidSocket() *SidSocket {
	c := &SidSocket{}
	return c
}

func (w *SidSocket) Setup(board *Board) {
	w.board = board
}

func (w *SidSocket) Reset() {
}

func (w *SidSocket) GetPlayer() mos6581.IPlayer {
	return w.board.player
}
