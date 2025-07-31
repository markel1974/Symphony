package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"os"
	"runtime/pprof"
)

// CreateProfileCPUStart initializes and returns a shell command to start CPU profiling and save it to a specified file.
func CreateProfileCPUStart() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		if len(args) <= 0 {
			task.WriteLn("could not create cpu profile: " + "missing filename")
			return nil
		}
		f, err := os.Create(args[0])
		if err != nil {
			task.WriteLn("could not create CPU profile: " + err.Error())
			return nil
		}
		defer f.Close()
		if err = pprof.StartCPUProfile(f); err != nil {
			task.WriteLn("could not start CPU profile: " + err.Error())
			return nil
		}
		task.WriteLn("Cpu Profiling started")
		return nil
	}
	root := shell.NewCommand("startcpuprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start cpu profiling", "Start cpu profiling")
	return root
}
