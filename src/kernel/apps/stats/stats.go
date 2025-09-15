package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// Create initializes and returns the root command for system statistics operations with multiple subcommands attached.
func Create() (*process.Command, []interfaces.ICommand) {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("stats", interfaces.CommandTypeDirectory, nil, false, run)
	root.SetHelp("System Stats", "System Stats")

	var apps []interfaces.ICommand
	apps = append(apps, CreateCPUProfileStart())
	apps = append(apps, CreateCPUProfileStop())
	apps = append(apps, CreateMemProfile())
	apps = append(apps, CreateMemStatus())
	apps = append(apps, CreateMemPlot())
	apps = append(apps, CreateCPUStatus())
	for _, app := range apps {
		_ = root.AddCommand(app)
	}
	return root, apps
}
