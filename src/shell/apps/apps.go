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
)

type Root struct {
}

func NewRoot() *Root {
	return &Root{}
}

func (t *Root) AddCommand(cmd *cli.Command, child *cli.Command) {
	_ = cmd.AddCommand(child)
}

func (t *Root) CreateCommand() *cli.Command {
	cmd := cli.NewCommand()
	return cmd
}

func (t *Root) Build(bin *cli.Command) (*cli.Command, *cli.Command) {
	bin.Use = "bin"
	bin.Short = "Bin"

	coreC := t.CreateCommand()
	t.AddCommand(coreC, core.CreateExit(t))
	t.AddCommand(coreC, core.CreateCD(t))
	t.AddCommand(coreC, core.CreatePWD(t))
	t.AddCommand(coreC, core.CreateActivate(t))
	t.AddCommand(coreC, core.CreateKill(t))
	t.AddCommand(coreC, core.CreateKillAll(t))
	t.AddCommand(coreC, core.CreatePs(t))
	t.AddCommand(coreC, core.CreateClear(t))
	t.AddCommand(coreC, core.CreateFg(t))
	t.AddCommand(coreC, core.CreateHistory(t))
	t.AddCommand(coreC, core.CreateTasks(t))
	t.AddCommand(coreC, core.CreateLs(t))

	root := t.CreateCommand()
	sbin := t.CreateCommand()
	sbin.Use = "sbin"
	sbin.Short = "SBin"

	t.AddCommand(sbin, stats.Create(t))
	t.AddCommand(sbin, runtime.Create(t))

	t.AddCommand(root, sbin)
	t.AddCommand(root, bin)
	t.AddCommand(root, games.Create(t))
	return coreC, root
}
