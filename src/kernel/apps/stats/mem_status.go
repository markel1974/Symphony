package stats

import (
	"fmt"
	"runtime"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateMemStatus initializes a shell command for monitoring runtime memory statistics and garbage collection cycles.
func CreateMemStatus() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		process.Write(fmt.Sprintf("Allocated memory in heap objects: %.3f MB", bToMb(m.Alloc)), true)
		process.Write(fmt.Sprintf("Total memory allocated for heap objects: %.3f MB", bToMb(m.TotalAlloc)), true)
		process.Write(fmt.Sprintf("Total memory obtained from the OS: %.3f MB", bToMb(m.Sys)), true)
		process.Write(fmt.Sprintf("Number of completed GC cycles: %d", m.NumGC), true)
		return nil
	}
	root := process.NewCommand("rt", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Runtime Status", "Runtime Status")
	return root
}
