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
	apps = append(apps, CreateProfileCPUStart())
	apps = append(apps, CreateProfileCPUStop())
	apps = append(apps, CreateProfileMemory())
	apps = append(apps, CreateMemoryStatus())
	apps = append(apps, CreateMemoryPlot())
	apps = append(apps, CreateCPUStatus())
	//apps = append(apps, CreateCPUUsage())
	for _, app := range apps {
		_ = root.AddCommand(app)
	}
	return root, apps
}
