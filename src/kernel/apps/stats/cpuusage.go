package stats

import (
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
	"time"
)

// CreateCPUUsage generates a shell command to compute and display CPU usage over a 3-second interval.
func CreateCPUUsage() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		task.WriteLn("Computing cpu usage")
		idle0, total0 := getCPUSample()
		time.Sleep(3 * time.Second)
		idle1, total1 := getCPUSample()
		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks
		task.WriteLn(fmt.Sprintf("CPU Usage: %f", cpuUsage))
		return nil
	}
	root := shell.NewCommand("usage", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("CPU usage", "CPU usage")

	return root
}
