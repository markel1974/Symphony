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
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
)

// Create initializes and returns the root command for runtime operations.
func Create() (*process.Command, []interfaces.ICommand) {
	run := func(process interfaces.IProcess, args []string) error {
		return nil
	}
	root := process.NewCommand("runtime", interfaces.CommandTypeDirectory, nil, false, run)
	root.SetHelp("Runtime", "Runtime")
	var apps []interfaces.ICommand
	apps = append(apps, CreateGC())
	for _, app := range apps {
		_ = root.AddCommand(app)
	}
	return root, apps
}
