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
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreateKill() interfaces.ICommand {
	run := func(task interfaces.IProcess, args []string) error {
		if len(args) <= 0 {
			task.WriteLn("Empty argument")
			return nil
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			task.WriteLn("Invalid argument: " + args[0])
			return nil
		}
		if !task.IsActive(pid) {
			task.WriteLn("Unknown Task: " + args[0])
			return nil
		}
		task.Kill(pid)
		return nil
	}
	root := process.NewCommand("kill", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Kill", "Kill")

	return root
}
