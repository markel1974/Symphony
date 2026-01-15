package stats

import (
	"runtime/pprof"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateCPUProfileStop creates and returns a shell command to stop CPU profiling.
func CreateCPUProfileStop() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		pprof.StopCPUProfile()
		process.Write("CPU Profiling stopped", true)
		return nil
	}
	root := process.NewCommand("cpu_profile_stop", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Stop CPU profiling", "Stop CPU profiling")

	return root
}
