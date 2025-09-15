package runtime

import (
	"embed"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

//go:embed scripts/*.go
var scripts embed.FS

// CreateExample initializes and returns a command that triggers garbage collection when executed.
func CreateExample() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("example", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Example", "Example")
	data, _ := scripts.ReadFile("scripts/example.go")
	root.SetScript(string(data))
	return root
}
