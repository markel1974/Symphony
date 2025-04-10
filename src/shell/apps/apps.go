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

package apps

import (
	"github.com/markel1974/c64emu/src/shell/apps/core"
	"github.com/markel1974/c64emu/src/shell/apps/games"
	"github.com/markel1974/c64emu/src/shell/apps/runtime"
	"github.com/markel1974/c64emu/src/shell/apps/stats"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

type Root struct {
}

func NewRoot() *Root {
	return &Root{}
}

func (t *Root) Build(bin *cli.Command) (*cli.Command, *cli.Command) {
	sbinRun := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		return nil
	}
	sbin := cli.NewCommand("sbin", nil, false, sbinRun)
	sbin.SetHelp("SBin", "SBin")

	_ = sbin.AddCommand(stats.Create())
	_ = sbin.AddCommand(runtime.Create())

	rootRun := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		return nil
	}
	root := cli.NewCommand("", nil, false, rootRun)
	_ = root.AddCommand(sbin)
	_ = root.AddCommand(bin)
	_ = root.AddCommand(games.Create())

	coreCRun := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		return nil
	}
	coreC := cli.NewCommand("", nil, false, coreCRun)
	coreC.SetHelp("Core", "Core")
	_ = coreC.AddCommand(core.CreateExit())
	_ = coreC.AddCommand(core.CreateCD())
	_ = coreC.AddCommand(core.CreatePWD())
	_ = coreC.AddCommand(core.CreateActivate())
	_ = coreC.AddCommand(core.CreateKill())
	_ = coreC.AddCommand(core.CreateKillAll())
	_ = coreC.AddCommand(core.CreatePs())
	_ = coreC.AddCommand(core.CreateClear())
	_ = coreC.AddCommand(core.CreateFg())
	_ = coreC.AddCommand(core.CreateHistory())
	_ = coreC.AddCommand(core.CreateTasks())
	_ = coreC.AddCommand(core.CreateLs())
	_ = coreC.AddCommand(core.CreateHelp(root))

	return coreC, root
}
