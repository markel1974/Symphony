package stats

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"runtime"
)

// CreateCPUStatus returns a shell command that outputs various CPU-related metrics such as logical CPU count and goroutine count.
func CreateCPUStatus() interfaces.ICommand {
	run := func(task interfaces.IProcess, args []string) error {
		task.WriteLn(fmt.Sprintf("Number of logical CPUs: %d", runtime.NumCPU()))
		task.WriteLn(fmt.Sprintf("Maximum number of CPUs that can be executing simultaneously: %d", runtime.GOMAXPROCS(0)))
		task.WriteLn(fmt.Sprintf("Number of goroutines that currently exist: %d", runtime.NumGoroutine()))
		task.WriteLn(fmt.Sprintf("Number of cgo calls made by the current process: %d", runtime.NumCgoCall()))
		return nil
	}
	root := process.NewCommand("cpu", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("CPUs status", "CPUs status")

	return root
}
