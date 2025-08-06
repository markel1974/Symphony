package runtime

import (
	"runtime"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateGC initializes and returns a command that triggers garbage collection when executed.
func CreateGC() *process.Command {
	run := func(task interfaces.IProcess, args []string) error {
		//r := cmd.GetRootContext()
		runtime.GC()
		task.Write("GC Done", true)
		return nil
	}
	root := process.NewCommand("gc", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start Garbage", "Start Garbage")

	return root
}
