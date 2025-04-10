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
)

func CreateActivate() *cli.Command {
	run := func(task interfaces.ITask, args []string) error {
		targetPid := -1
		if len(args) > 0 {
			targetPid, _ = strconv.Atoi(args[0])
		}
		task.SetSelectionMode(targetPid)
		return nil
	}
	readFn := func(task interfaces.ITask, code int, key rune) {
		if code == 1 {
			switch interfaces.CursorCodeDef(key) {
			case interfaces.CursorUpDef:
				task.SetSelectionOptions('y', -1)
			case interfaces.CursorDownDef:
				task.SetSelectionOptions('y', 1)
			case interfaces.CursorLeftDef:
				task.SetSelectionOptions('x', -1)
			case interfaces.CursorRightDef:
				task.SetSelectionOptions('x', 1)
			}
		} else {
			switch key {
			case 'w':
				task.SetSelectionOptions('y', -1)
			case 's':
				task.SetSelectionOptions('y', 1)
			case 'a':
				task.SetSelectionOptions('x', -1)
			case 'd':
				task.SetSelectionOptions('x', 1)
			case '+':
				task.SetSelectionOptions('z', 0.1)
			case '-':
				task.SetSelectionOptions('z', -0.1)
			case '\t':
				task.SetSelectionModeNext()
			case 'q':
				task.SetSelectionModePrevious()
			}
		}
	}

	root := cli.NewCommand("activate", nil, true, run)
	root.SetHelp("Activate", "Activate")
	root.SetReadFn(readFn)

	return root
}
