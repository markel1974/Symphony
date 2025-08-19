package xshell

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreateXShell() *process.Command {
	onCreate := func(process interfaces.IUserProcess, args []string) error {
		s := NewXShell("% ", true)
		s.Setup(process)
		return nil
	}
	root := process.NewCommand("xsh", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("XShell", "XShell")
	root.SetReadFilter(true)
	return root
}
