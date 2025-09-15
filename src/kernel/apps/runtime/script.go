package runtime

import (
	"embed"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

//go:embed scripts/*.go
var scripts embed.FS

// CreateScript initializes and returns a command that triggers garbage collection when executed.
func CreateScript() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("script", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Script", "Script")
	data, _ := scripts.ReadFile("scripts/script.go")
	root.SetScript(string(data))
	return root
}
