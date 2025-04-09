package core

import (
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateLs() *cli.Command {
	root := cli.NewCommand("ls", []string{"dir"}, false)
	root.SetHelp("ls", "ls")
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		r.WriteLn("")
		for _, c := range r.CWDChilds() {
			r.WriteLn(c)
		}
		return nil
	}
	return root
}
