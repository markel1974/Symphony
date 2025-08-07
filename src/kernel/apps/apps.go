package apps

import (
	"github.com/markel1974/c64emu/src/kernel/apps/games"
	"github.com/markel1974/c64emu/src/kernel/apps/runtime"
	"github.com/markel1974/c64emu/src/kernel/apps/stats"
	"github.com/markel1974/c64emu/src/kernel/apps/system"
	"github.com/markel1974/c64emu/src/kernel/apps/xshell"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// Root represents the top-level structure used to build and organize CLI command hierarchies.
type Root struct {
}

// NewRoot creates and returns a new instance of the Root structure.
func NewRoot() *Root {
	return &Root{}
}

// Build constructs and returns two command trees, coreC and root, initialized with their respective subcommands and functionality.
func (t *Root) Build(bin interfaces.ICommand) (string, interfaces.ICommand, interfaces.ICommand) {
	var aliases []interfaces.ICommand
	sbin := process.NewCommand("sbin", interfaces.CommandTypeDirectory, nil, false, func(process interfaces.IProcess, args []string) error {
		return nil
	})
	sbin.SetHelp("SBin", "SBin")

	var sbinCommands []interfaces.ICommand

	xsh := xshell.CreateXShell()
	sbinCommands = append(sbinCommands, xsh)
	sbinCommands = append(sbinCommands, system.CreateExit())
	sbinCommands = append(sbinCommands, system.CreateCD())
	sbinCommands = append(sbinCommands, system.CreatePWD())
	//sbinCommands = append(sbinCommands, system.CreateActivate())
	sbinCommands = append(sbinCommands, system.CreateKill())
	sbinCommands = append(sbinCommands, system.CreateKillAll())
	sbinCommands = append(sbinCommands, system.CreatePs())
	sbinCommands = append(sbinCommands, system.CreateClear())
	sbinCommands = append(sbinCommands, system.CreateFg())
	//sbinCommands = append(sbinCommands, system.CreateHistory())
	//sbinCommands = append(sbinCommands, system.CreateTasks())
	sbinCommands = append(sbinCommands, system.CreateLs())
	sbinCommands = append(sbinCommands, system.CreateHelp())

	for _, app := range sbinCommands {
		_ = sbin.AddCommand(app)
	}

	aliases = append(aliases, sbinCommands...)

	if statsRoot, statsApp := stats.Create(); statsRoot != nil {
		_ = sbin.AddCommand(statsRoot)
		aliases = append(aliases, statsApp...)
	}

	if runtimeRoot, runtimeApp := runtime.Create(); runtimeRoot != nil {
		_ = sbin.AddCommand(runtimeRoot)
		aliases = append(aliases, runtimeApp...)
	}

	root := process.NewCommand("/", interfaces.CommandTypeDirectory, nil, false, func(process interfaces.IProcess, args []string) error {
		return nil
	})
	_ = root.AddCommand(sbin)
	_ = root.AddCommand(bin)
	_ = root.AddCommand(games.Create())

	aliasesCommand := process.NewCommand("", interfaces.CommandTypeDirectory, nil, false, func(process interfaces.IProcess, args []string) error {
		return nil
	})
	aliasesCommand.SetHelp("Aliases", "Aliases")
	//aliases
	for _, app := range aliases {
		_ = aliasesCommand.AddCommand(app)
	}
	return xsh.Name(), aliasesCommand, root
}
