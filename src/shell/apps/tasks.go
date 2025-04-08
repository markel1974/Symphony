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

package apps

import (
	"github.com/markel1974/c64emu/src/shell/apps/commandcreator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"strings"
)

func CreateTasks(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "task"
	root.Short = "Task"
	root.Long = "Task"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		if len(args) <= 0 {
			return
		}
		r := cmd.GetRootContext()
		kind := strings.TrimSpace(strings.ToLower(args[0]))
		args = args[1:]
		switch kind {
		case "list":
			r.WriteLn("")
			for _, task := range r.ListTasks() {
				r.WriteLn(task)
			}
		case "restore":
			if len(args) > 0 {
				r.RestoreTasks(args[0])
			}
		case "save":
			if len(args) > 0 {
				r.SaveTasks(args[0])
			}
		}
	}
	return root
}
