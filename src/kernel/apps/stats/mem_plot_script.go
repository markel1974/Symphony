package stats

import (
	"embed"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

//go:embed scripts/*.go
var scripts embed.FS

// CreateMemPlotScript generates a process command for a runtime memory plot script.
// It sets the script content from a preloaded file and provides help information.
func CreateMemPlotScript() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("rtplot_script", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Runtime Plot Script", "Runtime Plot Script")
	data, _ := scripts.ReadFile("scripts/mem_plot.go")
	root.SetScript(string(data))
	return root
}
