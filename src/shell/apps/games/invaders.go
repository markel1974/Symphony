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
	run := func(task interfaces.ITask, args []string) error {
		w, h := task.GetScreenSize()
		g := invaders.NewGame(w, h)
		g.SetMenuState()
		task.SetContext(g)
		task.CreateTimer(0, 100, -1)
		return nil
	}
	readFn := func(task interfaces.ITask, code int, key rune) {
		ctx := task.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.HandleKey(key)
	}
	timerFn := func(task interfaces.ITask, tid int, interval int) {
		ctx := task.GetContext()
		g, ok := ctx.(*invaders.Invaders)
		if !ok {
			return
		}
		g.Update()
		task.PaintRequest()
	}
	paintFn := func(task interfaces.ITask, surface interfaces.ISurface) {
		ctx := task.GetContext()
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

	root := cli.NewCommand("invaders", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Invaders", "Invaders")
	root.SetTimerFn(timerFn)
	root.SetPaintFn(paintFn)
	root.SetReadFn(readFn)

	return root
}
