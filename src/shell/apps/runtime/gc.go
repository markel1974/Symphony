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

package runtime

import (
	"github.com/markel1974/c64emu/src/shell/apps/commandcreator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"runtime"
)

func CreateGC(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.SetName("gc", nil)
	root.ShortHelp = "Start Garbage"
	root.LongHelp = "Start Garbage"
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		//r := cmd.GetRootContext()
		r.WriteLn("")
		runtime.GC()
		r.WriteLn("GC Done")
		return nil
	}
	return root
}
