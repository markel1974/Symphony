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
	"fmt"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/shell"
	"runtime"
)

func CreateMemoryStatus() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		task.WriteLn(fmt.Sprintf("Allocated memory in heap objects: %.3f MB", bToMb(m.Alloc)))
		task.WriteLn(fmt.Sprintf("Total memory allocated for heap objects: %.3f MB", bToMb(m.TotalAlloc)))
		task.WriteLn(fmt.Sprintf("Total memory obtained from the OS: %.3f MB", bToMb(m.Sys)))
		task.WriteLn(fmt.Sprintf("Number of completed GC cycles: %d", m.NumGC))
		return nil
	}
	root := shell.NewCommand("rt", interfaces.CommandTypeFile, nil, false, run)
	root.SetHelp("Runtime Status", "Runtime Status")

	return root
}
