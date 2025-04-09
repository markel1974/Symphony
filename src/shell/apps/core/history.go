/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"strconv"
	"strings"
)

func CreateHistory() *cli.Command {
	run := func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		if len(args) == 0 {
			r.History(interfaces.HistoryActionList, -1)
			return nil

		}
		if idx, err := strconv.Atoi(args[0]); err == nil {
			r.History(interfaces.HistoryActionExec, idx)
			return nil
		}
		name := strings.TrimSpace(strings.ToLower(args[0]))
		args = args[1:]
		switch name {
		case "clear":
			r.History(interfaces.HistoryActionClear, -1)
		case "exec":
			if len(args) > 0 {
				if idx, err := strconv.Atoi(args[0]); err == nil {
					r.History(interfaces.HistoryActionExec, idx)
				}
			}
		case "list":
			r.History(interfaces.HistoryActionList, -1)
		}
		return nil
	}
	root := cli.NewCommand("history", []string{"h"}, false, run)
	root.SetHelp("History", "History")

	return root
}
