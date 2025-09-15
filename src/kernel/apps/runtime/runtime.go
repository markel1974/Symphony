package runtime

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// Create initializes and returns the root command for runtime operations.
func Create() (*process.Command, []interfaces.ICommand) {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("runtime", interfaces.CommandTypeDirectory, nil, false, run)
	root.SetHelp("Runtime", "Runtime")
	var apps []interfaces.ICommand
	apps = append(apps, CreateGC())
	apps = append(apps, CreateExample())
	for _, app := range apps {
		_ = root.AddCommand(app)
	}
	return root, apps
}
