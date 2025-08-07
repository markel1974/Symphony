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

/*
import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"strconv"
)


import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"strconv"
)

func CreateActivate() *process.Command {
	run := func(process interfaces.IProcess, args []string) error {
		targetPid := -1
		if len(args) > 0 {
			targetPid, _ = strconv.Atoi(args[0])
		}
		process.WindowsSelectionb(targetPid)
		return nil
	}
	readFn := func(task interfaces.IProcess, code int, key rune) {
		if code == 1 {
			switch interfaces.CursorCodeDef(key) {
			case interfaces.CursorUpDef:
				process.ProcessSelectionOptions('y', -1)
			case interfaces.CursorDownDef:
				process.ProcessSelectionOptions('y', 1)
			case interfaces.CursorLeftDef:
				process.ProcessSelectionOptions('x', -1)
			case interfaces.CursorRightDef:
				process.ProcessSelectionOptions('x', 1)
			}
		} else {
			switch key {
			case 'w':
				process.ProcessSelectionOptions('y', -1)
			case 's':
				process.ProcessSelectionOptions('y', 1)
			case 'a':
				process.ProcessSelectionOptions('x', -1)
			case 'd':
				process.ProcessSelectionOptions('x', 1)
			case '+':
				process.ProcessSelectionOptions('z', 0.1)
			case '-':
				process.ProcessSelectionOptions('z', -0.1)
			case '\t':
				process.ProcessSelectionNext()
			case 'q':
				process.ProcessSelectionPrevious()
			}
		}
	}

	root := process.NewCommand("activate", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Activate", "Activate")
	root.SetOnRead(readFn)

	return root
}


*/
