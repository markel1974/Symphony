package games

import (
	"github.com/markel1974/symphony/src/kernel/apps/games/invaders"
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateInvaders initializes and returns a new process.Command for the "Invaders" game.
// It sets up handlers for creation, input, timers, and rendering.
func CreateInvaders() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		g := invaders.New()
		g.Setup(process)
		g.Start()
		return nil
	}
	root := process.NewCommand("invaders", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Invaders", "Invaders")
	return root
}
