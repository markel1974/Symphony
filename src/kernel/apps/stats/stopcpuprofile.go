package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"runtime/pprof"
)

// CreateProfileCPUStop creates and returns a shell command to stop CPU profiling.
func CreateProfileCPUStop() *process.Command {
	run := func(task interfaces.IProcess, args []string) error {
		pprof.StopCPUProfile()
		task.WriteLn("Cpu Profiling stopped")
		return nil
	}
	root := process.NewCommand("stopcpuprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Stop cpu profiling", "Stop cpu profiling")

	return root
}
