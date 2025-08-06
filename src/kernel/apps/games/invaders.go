package games

import (
	"github.com/markel1974/c64emu/src/kernel/apps/games/invaders"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateInvaders initializes and returns a new process.Command for the "Invaders" game.
// It sets up handlers for creation, input, timers, and rendering.
func CreateInvaders() *process.Command {
	onCreate := func(task interfaces.IProcess, args []string) error {
		w, h := task.GetScreenSize()
		g := invaders.NewGame(w, h)
		g.SetMenuState()
		task.SetContext(g)
		task.CreateTimer(0, 100, -1)
		return nil
	}
	onRead := func(task interfaces.IProcess, code int, key rune) {
		ctx := task.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.HandleKey(key)
	}
	onTimer := func(task interfaces.IProcess, tid int, interval int) {
		ctx := task.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.Update()
		task.PaintRequest()
	}
	onPaint := func(task interfaces.IProcess, surface interfaces.ISurface) {
		ctx := task.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		rows, columns := surface.GetSize()
		w, h := g.GetSize()
		if h != rows || w != columns {
			g.SetSize(columns, rows)
			g.SetMenuState()
		}
		g.Draw(surface)
	}
	root := process.NewCommand("invaders", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetOnRead(onRead)
	root.SetOnTimer(onTimer)
	root.SetOnPaint(onPaint)
	root.SetHelp("Invaders", "Invaders")
	return root
}
