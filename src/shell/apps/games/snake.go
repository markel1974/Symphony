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
	"github.com/markel1974/c64emu/src/shell/apps/games/snake"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateSnake() *cli.Command {
	root := cli.NewCommand("snake", nil, true)
	root.SetHelp("Snake", "Snake")
	root.Run = func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		w, h := r.GetScreenSize()
		s := snake.New()
		s.SetSize(h, w)
		r.SetContext(pid, s)
		r.CreateTimer(pid, 0, 200, -1)
		return nil
	}
	root.ReadEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
		s := ctx.(*snake.Snake)
		switch key {
		case 'a':
			if s.Direction != snake.Left {
				s.Direction = snake.Left
			}
		case 'd':
			if s.Direction != snake.Right {
				s.Direction = snake.Right
			}
		case 'w':
			if s.Direction != snake.Up {
				s.Direction = snake.Up
			}
		case 's':
			if s.Direction != snake.Down {
				s.Direction = snake.Down
			}
		case '1':
			s.Start()
		}
	}
	root.TimerEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		s, ok := ctx.(*snake.Snake)
		if !ok {
			return
		}
		s.Advance()
		r.PaintRequest(pid)
	}
	root.PaintEvent = func(r interfaces.IContext, cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		s, ok := ctx.(*snake.Snake)
		if !ok {
			return
		}
		rows, columns := surface.GetSize()
		if s.Rows != rows || s.Columns != columns {
			s.SetSize(rows, columns)
		}
		s.Draw(surface)
	}
	return root
}
