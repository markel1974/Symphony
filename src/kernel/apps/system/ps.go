package system

import (
	"fmt"
	"time"

	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
)

func CreatePs() interfaces.ICommand {
	run := func(process interfaces.IUserProcess, args []string) error {
		out := "Pid: Process"
		now := time.Now()
		pl := process.ProcessList()
		for _, v := range pl {
			diff := now.Sub(v.Time())
			hours := int(diff.Hours())
			minutes := int(diff.Minutes()) % 60
			seconds := int(diff.Seconds()) % 60
			out += fmt.Sprintf("\r\n%d: %s (%s) %s %02d:%02d:%02d", v.PID(), v.Name(), v.Line(), v.User(), hours, minutes, seconds)
		}
		process.Write(out, true)
		return nil
	}

	root := process.NewCommand("ps", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Processes", "Processes")

	return root
}
