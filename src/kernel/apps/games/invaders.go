package games

import (
	"github.com/markel1974/c64emu/src/kernel/apps/games/invaders"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateInvaders initializes and returns a new process.Command for the "Invaders" game.
// It sets up handlers for creation, input, timers, and rendering.
func CreateInvaders() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		w, h := process.GetScreenSize()
		g := invaders.NewGame(w, h)
		g.SetMenuState()
		process.SetContext(g)
		process.CreateTimer(0, 100, -1)
		return nil
	}
	onRead := func(process interfaces.IUserProcess, code int, key rune) {
		ctx := process.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.HandleKey(key)
	}
	onTimer := func(process interfaces.IUserProcess, tid int, interval int) {
		ctx := process.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.Update()
		process.PaintRequest()
	}
	onPaint := func(process interfaces.IUserProcess, surface interfaces.ISurface) {
		ctx := process.GetContext()
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
