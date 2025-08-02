package xshell

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreateXShell() *process.Command {
	onCreate := func(process interfaces.IProcess, args []string) error {
		s := NewXShell("% ", true)
		process.SetContext(s)
		s.Start(process)
		return nil
	}
	onRead := func(process interfaces.IProcess, code int, key rune) {
		ctx := process.GetContext()
		s, ok := ctx.(*XShell)
		if !ok {
			return
		}
		s.KeyHandler(process, code, key)
	}

	orReadBroast := func(process interfaces.IProcess, code int, key rune) {
		ctx := process.GetContext()
		s, _ := ctx.(*XShell)
		if s == nil {
			return
		}
		s.BroadcastKeyHandler(process, code, key)
	}
	onTimer := func(process interfaces.IProcess, tid int, interval int) {
	}
	root := process.NewCommand("xsh", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("XShell", "XShell")
	root.SetTimerFn(onTimer)
	root.SetReadFn(onRead)
	root.SetReadBroadcastFn(orReadBroast)

	return root
}
