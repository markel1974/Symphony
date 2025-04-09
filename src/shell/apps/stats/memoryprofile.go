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
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"os"
	"runtime"
	"runtime/pprof"
)

func CreateProfileMemory() *cli.Command {
	root := cli.NewCommand("memprofile", nil, false)
	root.SetHelp("Memory profiling", "Memory profiling")
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		//r := cmd.GetRootContext()
		r.WriteLn("")
		if len(args) <= 0 {
			r.WriteLn("could not create mem profile: " + "missing filename")
			return nil
		}

		f, err := os.Create(args[0])
		if err != nil {
			r.WriteLn("could not create mem profile: " + err.Error())
			return nil
		}
		defer f.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			r.WriteLn("could not write mem profile: " + err.Error())
		}

		r.WriteLn("Cpu Profiling started")

		return nil
	}
	return root
}
