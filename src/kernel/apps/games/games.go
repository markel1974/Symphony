package games

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// Create initializes and returns the root command for the games directory with its subcommands.
func Create() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("games", interfaces.CommandTypeDirectory, nil, false, run)
	root.SetHelp("Games", "Games")

	_ = root.AddCommand(CreateSnake())
	_ = root.AddCommand(CreateTetris())
	_ = root.AddCommand(CreateInvaders())
	return root
}
