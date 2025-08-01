package shell

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
)

func CreateXShell() *shell.Command {
	onCreate := func(process interfaces.IProcess, args []string) error {
		s := NewXShell("%", true)
		process.SetContext(s)
		process.WriteHighlights("Admin Console Ready")
		s.SetPromptPrefix(process.CWDName())
		s.NextLine(process, true)
		return nil
	}
	onRead := func(process interfaces.IProcess, code int, key rune) {
		ctx := process.GetContext()
		s, ok := ctx.(*XShell)
		if !ok {
			return
		}
		s.KeyHandler(process, interfaces.KeyTypeKey, key)
	}
	onTimer := func(process interfaces.IProcess, tid int, interval int) {
	}
	root := shell.NewCommand("xshell", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("XShell", "XShell")
	root.SetTimerFn(onTimer)
	root.SetReadFn(onRead)

	return root
}
