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
	"github.com/markel1974/c64emu/src/kernel/apps/games/tetris"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/servers/shell"
)

func CreateTetris() *shell.Command {
	onCreate := func(task interfaces.IProcess, args []string) error {
		task.SetContext(tetris.New(10, 18))
		task.CreateTimer(0, 300, -1)
		return nil
	}
	onRead := func(task interfaces.IProcess, code int, key rune) {
		ctx := task.GetContext()
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
	onTimer := func(task interfaces.IProcess, tid int, interval int) {
		ctx := task.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.ApplyGravity()
		task.PaintRequest()
	}
	onPaint := func(task interfaces.IProcess, surface interfaces.ISurface) {
		ctx := task.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.Draw(surface)
	}
	root := shell.NewCommand("tetris", interfaces.CommandTypeFile, nil, true, onCreate)
	root.SetHelp("Tetris", "Tetris")
	root.SetReadFn(onRead)
	root.SetTimerFn(onTimer)
	root.SetPaintFn(onPaint)
	return root
}
