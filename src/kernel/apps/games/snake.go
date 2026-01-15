package games

import (
	"github.com/markel1974/symphony/src/kernel/apps/games/snake"
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateSnake initializes a new command for the Snake game and sets up event handlers for input, painting, timers, and creation.
// It returns a pointer to a configured process.Command ready to run the Snake game.
func CreateSnake() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		s := snake.New()
		s.Setup(process)
		s.Start()
		return nil
	}
	root := process.NewCommand("snake", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Snake", "Snake")
	return root
}
