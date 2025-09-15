package stats

import (
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateCPUProfileStart initializes and returns a shell command to start CPU profiling and save it to a specified file.
func CreateCPUProfileStart() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		if len(args) <= 0 {
			process.Write("could not create cpu profile: "+"missing filename", true)
			return nil
		}
		saneFilename := filepath.Base(args[0])
		f, err := os.Create(saneFilename)
		if err != nil {
			process.Write("could not create CPU profile: "+err.Error(), true)
			return nil
		}
		defer f.Close()
		if err = pprof.StartCPUProfile(f); err != nil {
			process.Write("could not start CPU profile: "+err.Error(), true)
			return nil
		}
		process.Write("CPU Profiling started", true)
		return nil
	}
	root := process.NewCommand("cpu_profile_start", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Start CPU profiling", "Start CPU profiling")
	return root
}
