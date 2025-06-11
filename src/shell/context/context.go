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

package context

import (
	"github.com/markel1974/c64emu/src/shell/adaptiveticker"
	"github.com/markel1974/c64emu/src/shell/apps"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/context/file_system"
	"github.com/markel1974/c64emu/src/shell/context/render"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/shell/shell"
	"io"
)

const (
	contextMaQueueLen = 1024
)

type Context struct {
	ticker   *adaptiveticker.AdaptiveTicker
	reader   io.Reader
	writer   io.Writer
	commands *cli.Command
	render   interfaces.IRender
	auth     interfaces.IAuthenticator
	sh       *shell.Shell
	enterKey rune
	kernel   *Kernel
	prompt   string
	autosave bool
}

func NewContext(ticker *adaptiveticker.AdaptiveTicker, reader io.Reader, writer io.Writer, auth interfaces.IAuthenticator, commands *cli.Command, prompt string, autosave bool) *Context {
	ctx := &Context{
		ticker:   ticker,
		reader:   reader,
		writer:   writer,
		auth:     auth,
		commands: commands,
		prompt:   prompt,
		kernel:   nil,
		autosave: autosave,
	}
	return ctx
}

func (c *Context) Setup(terminal interfaces.ITerminal) {
	c.render = render.NewRender(terminal)
	system := apps.NewRoot()
	systemCommands, commands := system.Build(c.commands)
	fs := file_system.NewCommandInteractor(commands, []interfaces.ICommand{systemCommands})
	ioAdapter := interfaces.IInputOutput(c)
	c.kernel = NewKernel(c.ticker, c.render, ioAdapter, fs)
	c.sh = shell.NewShell(c.auth, c.render, c, c.prompt, c.autosave)
}

func (c *Context) Exec() {
	c.render.WriteColor("Admin Console Ready", interfaces.ColorBlueDef, interfaces.ColorRedDef, interfaces.ModeNormal)
	c.sh.NextLine()
	c.kernel.Start()
}

func (c *Context) Type(kind interfaces.KeyType, key rune) {
	if kind == interfaces.KeyTypeCtrl {
		switch key {
		case 3:
			c.kernel.SetSelectionDisabled()
			c.kernel.KillForeground()
			c.sh.NextLine()
		case 4:
			c.kernel.ExecActivate()
		}
		return
	}
	if fgPid := c.kernel.GetForegroundPid(); fgPid != adaptiveticker.UnknownId {
		c.kernel.ExecRead(fgPid, int(kind), key)
		return
	}
	if quit := c.sh.KeyEvent(kind, key); quit {
		c.kernel.ExitRequested()
	}
}

func (c *Context) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *Context) Write(data []byte) (int, error) {
	return c.writer.Write(data)
}

func (c *Context) Close() {
}

func (c *Context) ExecCommand(line string) (bool, error) {
	return c.kernel.ExecCommand(line, nil)
}

func (c *Context) ExecSuggestion(in string, cursor int, count int) (int, bool) {
	ret := false
	data, suggestions, found := c.kernel.GetSuggestion(in, cursor)
	sLen := 0
	if found && len(suggestions) > 0 {
		sLen = len(suggestions)
		if idx := count % sLen; idx < sLen {
			if complete := suggestions[idx]; len(complete) > len(data) {
				tabLine := complete
				c.sh.Redraw(tabLine)
				c.sh.SetHistoryDefault(tabLine)
				ret = true
			}
		}
	}
	return sLen, ret
}

func (c *Context) SetScreenSize(w int, h int) {
	c.kernel.SetScreenSize(w, h)
}

func (c *Context) History(verb interfaces.HistoryAction, idx int) {
	switch verb {
	case interfaces.HistoryActionClear:
		c.sh.ClearHistory()
	case interfaces.HistoryActionExec:
		if arg, found := c.sh.GetHistoryAtPos(idx); found {
			_, _ = c.ExecCommand(arg)
		}
	case interfaces.HistoryActionList:
		c.render.Write(c.sh.GetHistory())
	}
}
