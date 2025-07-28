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
	"github.com/markel1974/c64emu/src/kernel/shell"
)

func CreateTetris() *shell.Command {
	run := func(task interfaces.ITask, args []string) error {
		//r := cmd.GetRootContext()
		//w, h := r.GetScreenSize()
		tx := tetris.New(10, 18)
		//s.SetSize(h, w)
		//pid := task.PID()
		task.SetContext(tx)
		//r.SetContext(pid, tx)
		task.CreateTimer(0, 300, -1)
		//r.CreateTimer(pid, 0, 300, -1)
		return nil
	}
	readFn := func(task interfaces.ITask, code int, key rune) {
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
	TimerFn := func(task interfaces.ITask, tid int, interval int) {
		ctx := task.GetContext()
		tx, ok := ctx.(*tetris.Tetris)
		if !ok {
			return
		}
		tx.ApplyGravity()
		task.PaintRequest()
	}
	paintFn := func(task interfaces.ITask, surface interfaces.ISurface) {
		ctx := task.GetContext()
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

	root := shell.NewCommand("tetris", interfaces.CommandTypeFile, nil, true, run)
	root.SetHelp("Tetris", "Tetris")
	root.SetReadFn(readFn)
	root.SetTimerFn(TimerFn)
	root.SetPaintFn(paintFn)

	return root
}
