package mos6581

import "github.com/markel1974/c64emu/src/components/board"

type ISocket interface {
	GetPlayer() board.IPlayer
}
