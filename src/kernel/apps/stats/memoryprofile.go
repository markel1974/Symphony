package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"os"
	"runtime"
	"runtime/pprof"
)

// CreateProfileMemory creates a shell command for generating memory usage profiles.
// It triggers garbage collection and writes the heap profile to a specified file.
func CreateProfileMemory() *shell.Command {
	run := func(task interfaces.IProcess, args []string) error {
		//r := cmd.GetRootContext()
		if len(args) <= 0 {
			task.WriteLn("could not create mem profile: " + "missing filename")
			return nil
		}

		f, err := os.Create(args[0])
		if err != nil {
			task.WriteLn("could not create mem profile: " + err.Error())
			return nil
		}
		defer f.Close()

		runtime.GC()
		if err = pprof.WriteHeapProfile(f); err != nil {
			task.WriteLn("could not write mem profile: " + err.Error())
		}

		task.WriteLn("Cpu Profiling started")

		return nil
	}
	root := shell.NewCommand("memprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Memory profiling", "Memory profiling")

	return root
}
