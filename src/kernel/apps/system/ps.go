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
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreatePs() *process.Command {
	run := func(process interfaces.IProcess, args []string) error {
		out := "\r\nPid: Process"
		pl := process.ProcessList()
		for _, v := range pl {
			out += fmt.Sprintf("\r\n%d: %s", v.Pid, v.Name)
		}
		process.WriteLn(out)
		return nil
	}

	root := process.NewCommand("ps", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Processes", "Processes")

	return root
}
