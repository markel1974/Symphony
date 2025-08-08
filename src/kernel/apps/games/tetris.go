package games

import (
	"github.com/markel1974/c64emu/src/kernel/apps/games/tetris"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateTetris initializes and returns a Tetris game command process with input, timer, and paint event handlers.
func CreateTetris() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		process.SetContext(tetris.New(10, 18))
		process.CreateTimer(0, 300, -1)
		return nil
	}
	onRead := func(process interfaces.IUserProcess, code int, key rune) {
		ctx := process.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		switch key {
		case 'a':
			tx.MoveLeft()
		case 'd':
			tx.MoveRight()
		case 'w':
			tx.RotateRight()
		case 's':
			tx.MoveDown()
		case ' ':
			tx.Drop()
		case '1':
			//r := cmd.GetRootContext()
			//w, h := r.GetScreenSize()
			tx.Init(10, 18)
		}
	}
	onTimer := func(process interfaces.IUserProcess, tid int, interval int) {
		ctx := process.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.ApplyGravity()
		process.PaintRequest()
	}
	onPaint := func(process interfaces.IUserProcess, surface interfaces.ISurface) {
		ctx := process.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.Draw(surface)
	}
	root := process.NewCommand("tetris", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Tetris", "Tetris")
	root.SetOnRead(onRead)
	root.SetOnTimer(onTimer)
	root.SetOnPaint(onPaint)
	return root
}
