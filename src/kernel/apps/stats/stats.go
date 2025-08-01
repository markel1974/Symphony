package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
)

// Create initializes and returns the root command for system statistics operations with multiple subcommands attached.
func Create() *shell.Command {
	run := func(task interfaces.IProcess, args []string) error {
		return nil
	}
	root := shell.NewCommand("stats", interfaces.CommandTypeDirectory, nil, false, run)
	root.SetHelp("System Stats", "System Stats")

	_ = root.AddCommand(CreateProfileCPUStart())
	_ = root.AddCommand(CreateProfileCPUStop())
	_ = root.AddCommand(CreateProfileMemory())
	_ = root.AddCommand(CreateMemoryStatus())
	_ = root.AddCommand(CreateMemoryPlot())
	_ = root.AddCommand(CreateCPUStatus())
	_ = root.AddCommand(CreateCPUUsage())

	return root
}
