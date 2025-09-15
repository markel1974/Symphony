package stats

import (
	"fmt"
	"runtime"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// CreateCPUStatus returns a shell command that outputs various CPU-related metrics such as logical CPU count and goroutine count.
func CreateCPUStatus() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		process.Write(fmt.Sprintf("Number of logical CPUs: %d", runtime.NumCPU()), true)
		process.Write(fmt.Sprintf("Maximum number of CPUs that can be executing simultaneously: %d", runtime.GOMAXPROCS(0)), true)
		process.Write(fmt.Sprintf("Number of goroutines that currently exist: %d", runtime.NumGoroutine()), true)
		process.Write(fmt.Sprintf("Number of cgo calls made by the current process: %d", runtime.NumCgoCall()), true)
		return nil
	}
	root := process.NewCommand("cpu_status", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("CPU status", "CPU status")
	return root
}
