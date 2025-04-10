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
	run := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		targetPid := -1
		if len(args) > 0 {
			targetPid, _ = strconv.Atoi(args[0])
		}
		r.SetSelectionMode(targetPid)
		return nil
	}
	readFn := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, ctx interface{}, code int, key rune) {
		if code == 1 {
			switch interfaces.CursorCodeDef(key) {
			case interfaces.CursorUpDef:
				r.SetSelectionOptions('y', -1)
			case interfaces.CursorDownDef:
				r.SetSelectionOptions('y', 1)
			case interfaces.CursorLeftDef:
				r.SetSelectionOptions('x', -1)
			case interfaces.CursorRightDef:
				r.SetSelectionOptions('x', 1)
			}
		} else {
			switch key {
			case 'w':
				r.SetSelectionOptions('y', -1)
			case 's':
				r.SetSelectionOptions('y', 1)
			case 'a':
				r.SetSelectionOptions('x', -1)
			case 'd':
				r.SetSelectionOptions('x', 1)
			case '+':
				r.SetSelectionOptions('z', 0.1)
			case '-':
				r.SetSelectionOptions('z', -0.1)
			case '\t':
				r.SetSelectionModeNext()
			case 'q':
				r.SetSelectionModePrevious()
			}
		}
	}

	root := cli.NewCommand("activate", nil, true, run)
	root.SetHelp("Activate", "Activate")
	root.SetReadFn(readFn)

	return root
}
