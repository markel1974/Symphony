package stats

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"runtime"
)

// CreateMemoryStatus initializes a shell command for monitoring runtime memory statistics and garbage collection cycles.
func CreateMemoryStatus() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		task.WriteLn(fmt.Sprintf("Allocated memory in heap objects: %.3f MB", bToMb(m.Alloc)))
		task.WriteLn(fmt.Sprintf("Total memory allocated for heap objects: %.3f MB", bToMb(m.TotalAlloc)))
		task.WriteLn(fmt.Sprintf("Total memory obtained from the OS: %.3f MB", bToMb(m.Sys)))
		task.WriteLn(fmt.Sprintf("Number of completed GC cycles: %d", m.NumGC))
		return nil
	}
	root := shell.NewCommand("rt", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Runtime Status", "Runtime Status")

	return root
}
