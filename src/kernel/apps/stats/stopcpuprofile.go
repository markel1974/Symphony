package stats

import (
	"runtime/pprof"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateProfileCPUStop creates and returns a shell command to stop CPU profiling.
func CreateProfileCPUStop() interfaces.ICommand {
	run := func(process interfaces.IProcess, args []string) error {
		pprof.StopCPUProfile()
		process.Write("Cpu Profiling stopped", true)
		return nil
	}
	root := process.NewCommand("stopcpuprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Stop cpu profiling", "Stop cpu profiling")

	return root
}
