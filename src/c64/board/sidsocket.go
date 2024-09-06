package board

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
