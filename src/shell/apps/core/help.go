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
	"github.com/markel1974/c64emu/src/shell/apps/commandcreator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"strings"
)

func CreateHelp(t commandcreator.ICreator, root *cli.Command) *cli.Command {
	help := t.CreateCommand()
	help.SetName("help", nil)
	help.ShortHelp = "Help about any command"
	help.LongHelp = `Help provides help for any command in the application. Simply type ` + help.Name() + ` help [path to command] for full details.`
	help.Run = func(r interfaces.IContext, c *cli.Command, pid int, args []string) error {
		if len(args) == 0 {
			return nil
		}
		path := r.CWDPath()
		path = append(path, args[0])
		cmd, _, err := root.Traverse(r.GetWriter(), path)
		if cmd == nil || err != nil || cmd == root {
			r.WriteLn("")
			r.WriteLn("unknown help topic: " + strings.Join(args, " "))
			c.Root().Usage(r.GetWriter())
		} else {
			r.WriteLn("")
			cmd.InitDefaultHelpFlag(r.GetWriter())
			cmd.Help(r.GetWriter())
		}
		return nil
	}
	return help
}
