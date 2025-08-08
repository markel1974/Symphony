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
	"time"

	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

func CreatePs() interfaces.ICommand {
	run := func(process interfaces.IProcess, args []string) error {
		out := "Pid: Process"
		now := time.Now()
		pl := process.ProcessList()
		for _, v := range pl {
			diff := now.Sub(v.Time())
			hours := int(diff.Hours())
			minutes := int(diff.Minutes()) % 60
			seconds := int(diff.Seconds()) % 60
			out += fmt.Sprintf("\r\n%d: %s (%s) %s %02d:%02d:%02d", v.PID(), v.Name(), v.Line(), v.User(), hours, minutes, seconds)
		}
		process.Write(out, true)
		return nil
	}

	root := process.NewCommand("ps", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Processes", "Processes")

	return root
}
