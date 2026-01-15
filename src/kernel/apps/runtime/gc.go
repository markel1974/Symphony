package runtime

import (
	"runtime"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateGC initializes and returns a command that triggers garbage collection when executed.
func CreateGC() *process.Command {
	run := func(process interfaces.IUserProcess, args []string) error {
		runtime.GC()
		process.Write("GC Done", true)
		return nil
	}
	root := process.NewCommand("gc", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start Garbage", "Start Garbage")

	return root
}
