package system

/*
import (
	"github.com/markel1974/symphony/src/kernel/interfaces"
	"github.com/markel1974/symphony/src/kernel/process"
	"strconv"
	"strings"
)

func CreateHistory() *process.Command {
	run := func(task interfaces.IUserProcess, args []string) error {
		if len(args) == 0 {
			task.History(interfaces.HistoryActionList, -1)
			return nil

		}
		if idx, err := strconv.Atoi(args[0]); err == nil {
			task.History(interfaces.HistoryActionExec, idx)
			return nil
		}
		name := strings.TrimSpace(strings.ToLower(args[0]))
		args = args[1:]
		switch name {
		case "clear":
			task.History(interfaces.HistoryActionClear, -1)
		case "exec":
			if len(args) > 0 {
				if idx, err := strconv.Atoi(args[0]); err == nil {
					task.History(interfaces.HistoryActionExec, idx)
				}
			}
		case "list":
			task.History(interfaces.HistoryActionList, -1)
		}
		return nil
	}
	root := process.NewCommand("history", interfaces.CommandTypeFile, []string{"h"}, false, run)
	root.SetHelp("History", "History")

	return root
}

*/
