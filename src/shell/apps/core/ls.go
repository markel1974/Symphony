package core

import (
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateLs() *cli.Command {
	run := func(task interfaces.ITask, args []string) error {
		task.WriteLn("")
		for _, c := range task.CWDChilds() {
			task.WriteLn(c)
		}
		return nil
	}
	root := cli.NewCommand("ls", interfaces.CommandTypeFile, []string{"dir"}, false, run)
	root.SetHelp("ls", "ls")

	return root
}
