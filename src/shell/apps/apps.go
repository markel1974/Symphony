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
	"github.com/markel1974/c64emu/src/shell/apps/games"
	"github.com/markel1974/c64emu/src/shell/apps/runtime"
	"github.com/markel1974/c64emu/src/shell/apps/stats"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"io"
)

type Apps struct {
	ctx    interfaces.IContext
	writer io.Writer
}

func NewApps(ctx interfaces.IContext, writer io.Writer) *Apps {
	return &Apps{
		ctx:    ctx,
		writer: writer,
	}
}

func (t *Apps) AddCommand(cmd *cli.Command, child *cli.Command) {
	_ = cmd.AddCommand(child)
}

func (t *Apps) CreateCommand() *cli.Command {
	cmd := cli.NewCommand()

	cmd.SetRootContext(t.ctx)
	cmd.SetOut(t.writer)
	cmd.SetErr(t.writer)
	return cmd
}

func (t *Apps) Build(bin *cli.Command) (*cli.Command, *cli.Command) {
	bin.Use = "bin"
	bin.Short = "Bin"

	system := t.CreateCommand()
	t.AddCommand(system, CreateExit(t))
	t.AddCommand(system, CreateCD(t))
	t.AddCommand(system, CreatePWD(t))
	t.AddCommand(system, CreateActivate(t))
	t.AddCommand(system, CreateKill(t))
	t.AddCommand(system, CreateKillAll(t))
	t.AddCommand(system, CreatePs(t))
	t.AddCommand(system, CreateClear(t))
	t.AddCommand(system, CreateFg(t))
	t.AddCommand(system, CreateHistory(t))
	t.AddCommand(system, CreateTasks(t))

	root := t.CreateCommand()
	sbin := t.CreateCommand()
	sbin.Use = "sbin"
	sbin.Short = "SBin"

	t.AddCommand(sbin, stats.Create(t))
	t.AddCommand(sbin, runtime.Create(t))

	t.AddCommand(root, sbin)
	t.AddCommand(root, bin)
	t.AddCommand(root, games.Create(t))
	return system, root
}
