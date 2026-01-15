package stats

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

// CreateMemProfile creates a shell command for generating memory usage profiles.
// It triggers garbage collection and writes the heap profile to a specified file.
func CreateMemProfile() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		if len(args) <= 0 {
			process.Write("could not create mem profile: "+"missing filename", true)
			return nil
		}
		saneFilename := filepath.Base(args[0])
		f, err := os.Create(saneFilename)
		if err != nil {
			process.Write("could not create mem profile: "+err.Error(), true)
			return nil
		}
		defer f.Close()
		runtime.GC()
		if err = pprof.WriteHeapProfile(f); err != nil {
			process.Write("could not write mem profile: "+err.Error(), true)
		}
		process.Write("Memory Profiling started", true)
		return nil
	}
	root := process.NewCommand("memprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Memory profiling", "Memory profiling")
	return root
}
