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

package games

import (
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func Create() *cli.Command {
	root := cli.NewCommand("games", nil, false)
	root.SetHelp("Games", "Games")
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		return nil
	}
	_ = root.AddCommand(CreateSnake())
	_ = root.AddCommand(CreateTetris())
	_ = root.AddCommand(CreateInvaders())
	return root
}
