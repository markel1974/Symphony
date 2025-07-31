package runtime

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"runtime"
)

// CreateGC initializes and returns a command that triggers garbage collection when executed.
func CreateGC() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		//r := cmd.GetRootContext()
		runtime.GC()
		task.WriteLn("GC Done")
		return nil
	}
	root := shell.NewCommand("gc", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start Garbage", "Start Garbage")

	return root
}
