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
	"github.com/markel1974/c64emu/src/shell/apps/games/invaders"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateInvaders() *cli.Command {
	run := func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		w, h := r.GetScreenSize()
		g := invaders.NewGame(w, h)
		g.SetMenuState()
		r.SetContext(pid, g)
		r.CreateTimer(pid, 0, 100, -1)
		return nil
	}
	root := cli.NewCommand("invaders", nil, true, run)
	root.SetHelp("Invaders", "Invaders")

	root.ReadEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.HandleKey(key)
	}
	root.TimerEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.Update()
		r.PaintRequest(pid)
	}
	root.PaintEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		rows, columns := surface.GetSize()
		w, h := g.GetSize()
		if h != rows || w != columns {
			g.SetSize(columns, rows)
			g.SetMenuState()
		}

		g.Draw(surface)
	}

	return root
}
