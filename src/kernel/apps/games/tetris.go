package games

import (
	"github.com/markel1974/symphony/src/kernel/apps/games/tetris"
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateTetris initializes and returns a Tetris game command process with input, timer, and paint event handlers.
func CreateTetris() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		t := tetris.New()
		t.Setup(process)
		t.Start()
		return nil
	}
	root := process.NewCommand("tetris", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Tetris", "Tetris")
	return root
}
