package core

import (
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateLs() *cli.Command {
	run := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		r.WriteLn("")
		for _, c := range r.CWDChilds() {
			r.WriteLn(c)
		}
		return nil
	}
	root := cli.NewCommand("ls", []string{"dir"}, false, run)
	root.SetHelp("ls", "ls")

	return root
}
