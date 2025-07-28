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
	"github.com/markel1974/c64emu/src/kernel/apps/games"
	"github.com/markel1974/c64emu/src/kernel/apps/runtime"
	"github.com/markel1974/c64emu/src/kernel/apps/stats"
	"github.com/markel1974/c64emu/src/kernel/apps/system"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/shell"
)

// Root represents the top-level structure used to build and organize CLI command hierarchies.
type Root struct {
}

// NewRoot creates and returns a new instance of the Root structure.
func NewRoot() *Root {
	return &Root{}
}

// Build constructs and returns two command trees, coreC and root, initialized with their respective subcommands and functionality.
func (t *Root) Build(bin *shell.Command) (*shell.Command, *shell.Command) {
	sbinRun := func(task interfaces.ITask, args []string) error {
		return nil
	}
	sbin := shell.NewCommand("sbin", interfaces.CommandTypeDirectory, nil, false, sbinRun)
	sbin.SetHelp("SBin", "SBin")

	_ = sbin.AddCommand(stats.Create())
	_ = sbin.AddCommand(runtime.Create())

	rootRun := func(task interfaces.ITask, args []string) error {
		return nil
	}
	root := shell.NewCommand("", interfaces.CommandTypeDirectory, nil, false, rootRun)
	_ = root.AddCommand(sbin)
	_ = root.AddCommand(bin)
	_ = root.AddCommand(games.Create())

	coreCRun := func(task interfaces.ITask, args []string) error {
		return nil
	}
	coreC := shell.NewCommand("", interfaces.CommandTypeDirectory, nil, false, coreCRun)
	coreC.SetHelp("Core", "Core")
	_ = coreC.AddCommand(system.CreateExit())
	_ = coreC.AddCommand(system.CreateCD())
	_ = coreC.AddCommand(system.CreatePWD())
	_ = coreC.AddCommand(system.CreateActivate())
	_ = coreC.AddCommand(system.CreateKill())
	_ = coreC.AddCommand(system.CreateKillAll())
	_ = coreC.AddCommand(system.CreatePs())
	_ = coreC.AddCommand(system.CreateClear())
	_ = coreC.AddCommand(system.CreateFg())
	_ = coreC.AddCommand(system.CreateHistory())
	_ = coreC.AddCommand(system.CreateTasks())
	_ = coreC.AddCommand(system.CreateLs())
	_ = coreC.AddCommand(system.CreateHelp())

	return coreC, root
}
