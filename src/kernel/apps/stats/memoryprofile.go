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

package stats

import (
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/shell"
	"os"
	"runtime"
	"runtime/pprof"
)

func CreateProfileMemory() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		//r := cmd.GetRootContext()
		task.WriteLn("")
		if len(args) <= 0 {
			task.WriteLn("could not create mem profile: " + "missing filename")
			return nil
		}

		f, err := os.Create(args[0])
		if err != nil {
			task.WriteLn("could not create mem profile: " + err.Error())
			return nil
		}
		defer f.Close()

		runtime.GC()
		if err = pprof.WriteHeapProfile(f); err != nil {
			task.WriteLn("could not write mem profile: " + err.Error())
		}

		task.WriteLn("Cpu Profiling started")

		return nil
	}
	root := shell.NewCommand("memprofile", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Memory profiling", "Memory profiling")

	return root
}
