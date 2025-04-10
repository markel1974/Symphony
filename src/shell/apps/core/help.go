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
	"strings"
)

func CreateHelp(root *cli.Command) *cli.Command {
	run := func(task interfaces.ITask, args []string) error {
		if len(args) == 0 {
			return nil
		}
		path := task.CWDPath()
		path = append(path, args[0])
		cmd, _, err := root.Traverse(path)
		if cmd == nil || err != nil || cmd == root {
			task.WriteLn("")
			task.WriteLn("unknown help topic: " + strings.Join(args, " "))
			task.WriteLn(task.GetCommand().Root().Help())
		} else {
			task.WriteLn("")
			task.WriteLn(cmd.Help())
		}
		return nil
	}

	help := cli.NewCommand("help", nil, false, run)
	help.SetHelp("Help about any command", `Help provides help for any command in the application. Simply type `+help.Name()+` help [path to command] for full details.`)

	return help
}
