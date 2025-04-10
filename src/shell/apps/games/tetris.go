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
	"github.com/markel1974/c64emu/src/shell/apps/games/tetris"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
)

func CreateTetris() *cli.Command {
	run := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, args []string) error {
		//r := cmd.GetRootContext()
		//w, h := r.GetScreenSize()
		tx := tetris.New(10, 18)
		//s.SetSize(h, w)
		r.SetContext(pid, tx)
		r.CreateTimer(pid, 0, 300, -1)
		return nil
	}
	readFn := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, ctx interface{}, code int, key rune) {
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		switch key {
		case 'a':
			tx.MoveLeft()
		case 'd':
			tx.MoveRight()
		case 'w':
			tx.RotateRight()
		case 's':
			tx.MoveDown()
		case ' ':
			tx.Drop()
		case '1':
			//r := cmd.GetRootContext()
			//w, h := r.GetScreenSize()
			tx.Init(10, 18)
		}
	}
	TimerFn := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, tid int, ctx interface{}, interval int) {
		//r := cmd.GetRootContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.ApplyGravity()
		r.PaintRequest(pid)
	}
	paintFn := func(r interfaces.IContext, cmd interfaces.ICommand, pid int, ctx interface{}, surface interfaces.ISurface) {
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		//rows, columns := surface.GetSize()
		//h, w := t.GetSize()
		//if h != rows || w != columns {
		//	fmt.Println(h, w)
		//	t.Init(rows, columns)
		//}
		tx.Draw(surface)
	}

	root := cli.NewCommand("tetris", nil, true, run)
	root.SetHelp("Tetris", "Tetris")
	root.SetReadFn(readFn)
	root.SetTimerFn(TimerFn)
	root.SetPaintFn(paintFn)

	return root
}
