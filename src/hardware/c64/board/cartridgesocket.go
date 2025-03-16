package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type CartridgeSocket struct {
	board *Board
	references.IExpansionSocketC64
}

func NewCartridgeSocket() *CartridgeSocket {
	return &CartridgeSocket{}
}

func (cs *CartridgeSocket) Connect(board *Board, m references.IExpansionSocketC64, b references.IExpansionC64, cfg *config.Config) error {
	cs.board = board
	cs.IExpansionSocketC64 = m
	if err := cs.IExpansionSocketC64.Setup(b, cfg); err != nil {
		return err
	}
	return nil
}
