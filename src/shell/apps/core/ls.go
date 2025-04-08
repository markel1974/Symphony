package core

import (
	"github.com/markel1974/c64emu/src/shell/apps/commandcreator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateLs(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "ls"
	root.Aliases = []string{"dir"}
	root.Short = "ls"
	root.Long = "ls"
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		r.WriteLn("")
		for _, c := range r.Childs() {
			r.WriteLn(c)
		}
		return nil
	}
	return root
}
