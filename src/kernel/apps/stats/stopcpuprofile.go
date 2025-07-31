package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"runtime/pprof"
)

// CreateProfileCPUStop creates and returns a shell command to stop CPU profiling.
func CreateProfileCPUStop() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		pprof.StopCPUProfile()
		task.WriteLn("Cpu Profiling stopped")
		return nil
	}
	root := shell.NewCommand("stopcpuprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Stop cpu profiling", "Stop cpu profiling")

	return root
}
