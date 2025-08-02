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

package system

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"strconv"
)

func CreateActivate() *process.Command {
	run := func(task interfaces.IProcess, args []string) error {
		targetPid := -1
		if len(args) > 0 {
			targetPid, _ = strconv.Atoi(args[0])
		}
		task.ProcessSelection(targetPid)
		return nil
	}
	readFn := func(task interfaces.IProcess, code int, key rune) {
		if code == 1 {
			switch interfaces.CursorCodeDef(key) {
			case interfaces.CursorUpDef:
				task.ProcessSelectionOptions('y', -1)
			case interfaces.CursorDownDef:
				task.ProcessSelectionOptions('y', 1)
			case interfaces.CursorLeftDef:
				task.ProcessSelectionOptions('x', -1)
			case interfaces.CursorRightDef:
				task.ProcessSelectionOptions('x', 1)
			}
		} else {
			switch key {
			case 'w':
				task.ProcessSelectionOptions('y', -1)
			case 's':
				task.ProcessSelectionOptions('y', 1)
			case 'a':
				task.ProcessSelectionOptions('x', -1)
			case 'd':
				task.ProcessSelectionOptions('x', 1)
			case '+':
				task.ProcessSelectionOptions('z', 0.1)
			case '-':
				task.ProcessSelectionOptions('z', -0.1)
			case '\t':
				task.ProcessSelectionNext()
			case 'q':
				task.ProcessSelectionPrevious()
			}
		}
	}

	root := process.NewCommand("activate", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Activate", "Activate")
	root.SetReadFn(readFn)

	return root
}
