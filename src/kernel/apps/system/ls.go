package system

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
)

func CreateLs() *shell.Command {
	run := func(task interfaces.IProcess, args []string) error {
		for _, c := range task.CWDDirectoryListing() {
			task.WriteLn(c)
		}
		return nil
	}
	root := shell.NewCommand("ls", interfaces.CommandTypeFile, []string{"dir"}, false, run)
	root.SetHelp("ls", "ls")

	return root
}
