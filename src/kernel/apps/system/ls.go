package system

import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

func CreateLs() interfaces.ICommand {
	run := func(task interfaces.IUserProcess, args []string) error {
		for _, c := range task.CWDDirectoryListing() {
			task.Write(c, true)
		}
		return nil
	}
	root := process.NewCommand("ls", interfaces.CommandTypeFile, []string{"dir"}, false, run)
	root.SetHelp("ls", "ls")

	return root
}
