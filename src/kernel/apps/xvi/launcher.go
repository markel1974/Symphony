package xvi

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateXVI initializes and returns a new command for the XVI process setup and execution.
func CreateXVI() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		p := NewXVI()
		p.Setup(process, args)
		p.Start()
		return nil
	}
	root := process.NewCommand("xvi", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("XVI", "XVI")
	return root
}
