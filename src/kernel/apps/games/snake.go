package games

import (
	"github.com/markel1974/c64emu/src/kernel/apps/games/snake"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateSnake initializes a new command for the Snake game and sets up event handlers for input, painting, timers, and creation.
// It returns a pointer to a configured process.Command ready to run the Snake game.
func CreateSnake() *process.Command {
	onCreate := func(task interfaces.IProcess, args []string) error {
		w, h := task.GetScreenSize()
		s := snake.New()
		s.SetSize(h, w)
		task.SetContext(s)
		task.CreateTimer(0, 200, -1)
		return nil
	}
	onRead := func(task interfaces.IProcess, code int, key rune) {
		ctx := task.GetContext()
		s := ctx.(*snake.Snake)
		switch key {
		case 'a':
			if s.Direction != snake.Left {
				s.Direction = snake.Left
			}
		case 'd':
			if s.Direction != snake.Right {
				s.Direction = snake.Right
			}
		case 'w':
			if s.Direction != snake.Up {
				s.Direction = snake.Up
			}
		case 's':
			if s.Direction != snake.Down {
				s.Direction = snake.Down
			}
		case '1':
			s.Start()
		}
	}
	onTimer := func(task interfaces.IProcess, tid int, interval int) {
		ctx := task.GetContext()
		s, ok := ctx.(*snake.Snake)
		if !ok {
			return
		}
		s.Advance()
		task.PaintRequest()
	}
	onPaint := func(task interfaces.IProcess, surface interfaces.ISurface) {
		ctx := task.GetContext()
		s, ok := ctx.(*snake.Snake)
		if !ok {
			return
		}
		rows, columns := surface.GetSize()
		if s.Rows != rows || s.Columns != columns {
			s.SetSize(rows, columns)
		}
		s.Draw(surface)
	}
	root := process.NewCommand("snake", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Snake", "Snake")
	root.SetOnTimer(onTimer)
	root.SetOnRead(onRead)
	root.SetOnPaint(onPaint)

	return root
}
