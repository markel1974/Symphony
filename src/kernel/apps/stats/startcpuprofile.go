package stats

import (
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateProfileCPUStart initializes and returns a shell command to start CPU profiling and save it to a specified file.
func CreateProfileCPUStart() interfaces.ICommand {
	run := func(task interfaces.IProcess, args []string) error {
		if len(args) <= 0 {
			task.Write("could not create cpu profile: "+"missing filename", true)
			return nil
		}
		saneFilename := filepath.Base(args[0])
		f, err := os.Create(saneFilename)
		if err != nil {
			task.Write("could not create CPU profile: "+err.Error(), true)
			return nil
		}
		defer f.Close()
		if err = pprof.StartCPUProfile(f); err != nil {
			task.Write("could not start CPU profile: "+err.Error(), true)
			return nil
		}
		task.Write("Cpu Profiling started", true)
		return nil
	}
	root := process.NewCommand("startcpuprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start cpu profiling", "Start cpu profiling")
	return root
}
