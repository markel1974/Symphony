package xshell

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreateXShell() *process.Command {

	onCreate := func(process interfaces.IUserProcess, args []string) error {
		s := NewXShell("% ", true)
		process.SetContext(s)
		s.Start(process)
		return nil
	}
	onError := func(process interfaces.IUserProcess, err error) {
		ctx := process.GetContext()
		s, ok := ctx.(*XShell)
		if !ok {
			return
		}
		s.ErrorHandler(process, err)
	}
	onRead := func(process interfaces.IUserProcess, code int, key rune) {
		ctx := process.GetContext()
		s, ok := ctx.(*XShell)
		if !ok {
			return
		}
		s.KeyHandler(process, code, key)
	}
	onReadBroast := func(process interfaces.IUserProcess, code int, key rune) {
		ctx := process.GetContext()
		s, _ := ctx.(*XShell)
		if s == nil {
			return
		}
		s.BroadcastKeyHandler(process, code, key)
	}
	onActivate := func(process interfaces.IUserProcess) {
		ctx := process.GetContext()
		s, _ := ctx.(*XShell)
		if s == nil {
			return
		}
		s.ActivateHandler(process)
	}
	root := process.NewCommand("xsh", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("XShell", "XShell")
	root.SetOnError(onError)
	root.SetOnRead(onRead)
	root.SetOnReadBroadcast(onReadBroast)
	root.SetOnActivate(onActivate)

	return root
}
